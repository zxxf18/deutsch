package middleware

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"deutsch/internal/sso"
	"deutsch/model/gormdb"
	"deutsch/model/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrUserDisabled = errors.New("business account is disabled")

// EnsureSSOUser provisions only business profile data. With the pinned Casdoor
// issuer, its immutable UUID subject is also the existing progress owner ID.
// Never link by mutable email/username, restore deleted users, or store passwords.
func EnsureSSOUser(ctx context.Context, repo repository.UserRepository, identity sso.Identity) (*gormdb.User, error) {
	if _, err := uuid.Parse(identity.Subject); err != nil || len(identity.Subject) != 36 {
		return nil, errors.New("invalid Casdoor subject")
	}
	user, err := repo.GetByUserID(ctx, identity.Subject)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if strings.TrimSpace(identity.Email) == "" || utf8.RuneCountInString(identity.Email) > 100 || strings.TrimSpace(identity.Username) == "" || utf8.RuneCountInString(identity.Username) > 50 {
			return nil, errors.New("identity does not fit business profile")
		}
		nickname := []rune(identity.DisplayName)
		if len(nickname) > 50 {
			nickname = nickname[:50]
		}
		role := "user"
		if identity.Role == "admin" {
			role = "admin"
		}
		user = &gormdb.User{ID: identity.Subject, Username: identity.Username, Email: identity.Email, Nickname: string(nickname), Role: role, IsEnabled: true}
		if err = repo.Create(ctx, user); err != nil {
			createErr := err
			// A concurrent first request may have won the insert, or the business
			// profile may predate SSO and already own this verified email. Re-read
			// the immutable subject first, then use the verified email as the
			// migration key to preserve the existing business ID and progress.
			var subjectErr error
			user, subjectErr = repo.GetByUserID(ctx, identity.Subject)
			if subjectErr == nil {
				err = nil
			} else if !errors.Is(subjectErr, gorm.ErrRecordNotFound) {
				err = subjectErr
			} else {
				var emailErr error
				user, emailErr = repo.GetByEmail(ctx, identity.Email)
				if emailErr == nil {
					err = nil
				} else {
					// Keep the original create error for a missing email match. This
					// avoids hiding the actual database constraint or outage.
					err = createErr
				}
			}
		}
	}
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsEnabled || user.DeletedAt.Valid {
		return nil, ErrUserDisabled
	}
	return user, nil
}
