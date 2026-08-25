package auth

import (
	"context"
	"strings"

	"deutsch/internal/code"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func genUsernameFromEmail(email string) string {
	idx := strings.Index(email, "@")
	if idx > 0 {
		return email[:idx]
	}
	return email
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	resp = &types.RegisterResponse{}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" {
		return nil, code.NewCodeError(code.CodeValidationError)
	}
	if len(req.Password) < 6 {
		return nil, code.NewCodeError(code.CodePasswordTooShort)
	}
	if req.InviteCode == "" {
		return nil, code.NewCodeError(code.CodeInviteCodeRequired)
	}

	username := req.Username
	if username == "" {
		username = genUsernameFromEmail(req.Email)
		for len(username) < 6 {
			username += "0"
		}
		if len(username) > 50 {
			username = username[:50]
		}
	}
	if strings.Contains(username, "@") || len(username) < 6 || len(username) > 50 {
		return nil, code.NewCodeError(code.CodeUsernameInvalid)
	}

	valid, _, err := l.svcCtx.InviteCodeRepo.Validate(l.ctx, req.InviteCode)
	if err != nil {
		l.Errorf("failed to validate invite code: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	if !valid {
		return nil, code.NewCodeError(code.CodeInvalidInviteCode)
	}

	existed, err := l.svcCtx.UserRepo.GetByEmailIncludingDeleted(l.ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		l.Errorf("failed to get user by email: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}

	usernameOwner, usernameErr := l.svcCtx.UserRepo.GetByUsernameIncludingDeleted(l.ctx, username)
	if usernameErr != nil && usernameErr != gorm.ErrRecordNotFound {
		l.Errorf("failed to get user by username: %+v", usernameErr)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	if usernameOwner != nil && (existed == nil || usernameOwner.ID != existed.ID) {
		return nil, code.NewCodeError(code.CodeUsernameExists)
	}

	passwordEncrypted, err := l.svcCtx.PasswordCipher.Encrypt(req.Password)
	if err != nil {
		l.Errorf("failed to encrypt password: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	var user *gormdb.User
	if existed != nil {
		if existed.DeletedAt.Valid {
			// 方案B：已删除用户用同邮箱重新注册 -> 恢复账号
			existed.DeletedAt = gorm.DeletedAt{}
			existed.Username = username
			existed.PasswordEncrypted = passwordEncrypted
			existed.Phone = phone
			existed.Nickname = req.Nickname
			existed.IsEnabled = true
			existed.Description = ""
			if err := l.svcCtx.UserRepo.Restore(l.ctx, existed); err != nil {
				l.Errorf("failed to restore user: %+v", err)
				return nil, code.NewCodeError(code.CodeDatabaseError)
			}
			if err := l.svcCtx.UserRepo.Update(l.ctx, existed); err != nil {
				l.Errorf("failed to update restored user: %+v", err)
				return nil, code.NewCodeError(code.CodeDatabaseError)
			}
			user = existed
		} else {
			return nil, code.NewCodeError(code.CodeEmailExists)
		}
	} else {
		user = &gormdb.User{
			Username:          username,
			Email:             req.Email,
			Phone:             phone,
			PasswordEncrypted: passwordEncrypted,
			Role:              "user",
			Nickname:          req.Nickname,
			IsEnabled:         true,
		}
		if err := l.svcCtx.UserRepo.Create(l.ctx, user); err != nil {
			l.Errorf("failed to create user: %+v", err)
			return nil, code.NewCodeError(code.CodeDatabaseError)
		}
	}

	if err := l.svcCtx.InviteCodeRepo.MarkUsed(l.ctx, req.InviteCode, user.ID); err != nil {
		l.Errorf("failed to mark invite code used: %+v", err)
		// 用户已创建，仅记录日志
	}

	jwtImpl, _ := jwt.New(l.ctx, l.svcCtx.Config.JWTAuth.AccessSecret, l.svcCtx.Config.JWTAuth.AccessExpire)
	jwtData := map[string]string{"userID": user.ID, "role": user.Role}
	jwtResp, err := jwtImpl.GenerateJWT(jwtData)
	if err != nil {
		l.Errorf("failed to generate jwt: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.ID = user.ID
	resp.Data.Username = user.Username
	resp.Data.Email = user.Email
	resp.Data.Role = user.Role
	resp.Data.Nickname = user.Nickname
	resp.Data.JwtToken = jwtResp.Token
	resp.Data.Expires = jwtResp.Expire.UnixMilli()
	resp.Data.MaxRefresh = jwtResp.MaxRefresh.UnixMilli()
	resp.Data.CreatedAt = user.CreatedAt.UnixMilli()
	resp.Data.UpdatedAt = user.UpdatedAt.UnixMilli()
	return
}
