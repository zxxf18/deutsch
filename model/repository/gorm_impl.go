package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"deutsch/model/gormdb"
)

// UserGormRepo GORM实现
type UserGormRepo struct {
	DB *gorm.DB
}

func NewUserGormRepo(db *gorm.DB) UserRepository {
	return &UserGormRepo{DB: db}
}

func (r *UserGormRepo) Create(ctx context.Context, user *gormdb.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserGormRepo) GetByEmail(ctx context.Context, email string) (*gormdb.User, error) {
	var user gormdb.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserGormRepo) GetByEmailIncludingDeleted(ctx context.Context, email string) (*gormdb.User, error) {
	var user gormdb.User
	err := r.DB.WithContext(ctx).Unscoped().Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserGormRepo) Restore(ctx context.Context, user *gormdb.User) error {
	return r.DB.WithContext(ctx).Unscoped().Model(user).Update("deleted_at", nil).Error
}

func (r *UserGormRepo) GetByPhone(ctx context.Context, phone string) (*gormdb.User, error) {
	var user gormdb.User
	err := r.DB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *UserGormRepo) GetByUserID(ctx context.Context, userID string) (*gormdb.User, error) {
	var user gormdb.User
	err := r.DB.WithContext(ctx).First(&user, "id = ?", userID).Error
	return &user, err
}

func (r *UserGormRepo) Update(ctx context.Context, user *gormdb.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserGormRepo) Delete(ctx context.Context, user *gormdb.User) error {
	return r.DB.WithContext(ctx).Delete(user).Error
}

func (r *UserGormRepo) List(ctx context.Context, filter *Filter) ([]*gormdb.User, int64, error) {
	var users []*gormdb.User
	var total int64

	query := r.DB.WithContext(ctx).Model(&gormdb.User{})

	// 分页计算
	pageNo := 1
	pageSize := 10
	if filter != nil {
		if filter.PageNo > 0 {
			pageNo = filter.PageNo
		}
		if filter.PageSize > 0 {
			pageSize = filter.PageSize
		}
	}
	offset := (pageNo - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(&users).Offset(-1).Count(&total).Error // Count需重置Offset
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// InviteGormRepo GORM实现
type InviteGormRepo struct {
	DB *gorm.DB
}

func NewInviteGormRepo(db *gorm.DB) InviteCodeRepository {
	return &InviteGormRepo{DB: db}
}

func (r *InviteGormRepo) CreateBatch(ctx context.Context, count int, createdBy string) ([]*gormdb.InviteCode, error) {
	result := make([]*gormdb.InviteCode, 0, count)
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err := recover(); err != nil {
			tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	for i := 0; i < count; i++ {
		ic := &gormdb.InviteCode{CreatedBy: createdBy, IsEnabled: true} // BeforeCreate生成ID/Code/Expiries
		if err := tx.Create(ic).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		result = append(result, ic)
	}
	return result, nil
}

func (r *InviteGormRepo) Validate(ctx context.Context, code string) (bool, *time.Time, error) {
	var ic gormdb.InviteCode
	err := r.DB.WithContext(ctx).
		Where("code = ? AND (used_by IS NULL OR used_by = '') AND expires_at > NOW() AND is_enabled = true", code).
		First(&ic).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &ic.ExpiresAt, nil
}

func (r *InviteGormRepo) MarkUsed(ctx context.Context, code string, usedBy string) error {
	return r.DB.WithContext(ctx).
		Model(&gormdb.InviteCode{}).
		Where("code = ?", code).
		Update("used_by", usedBy).Error
}

func (r *InviteGormRepo) List(ctx context.Context, filter *InviteCodeListFilter) ([]*gormdb.InviteCode, int64, error) {
	var invites []*gormdb.InviteCode
	var total int64

	query := r.DB.WithContext(ctx).Model(&gormdb.InviteCode{})

	// 仅查看可用的邀请码：未使用、未过期、已启用
	if filter != nil && filter.AvailableOnly {
		query = query.Where("(used_by IS NULL OR used_by = '') AND expires_at > NOW() AND is_enabled = true")
	}

	// 分页计算
	pageNo := 1
	pageSize := 10
	if filter != nil {
		if filter.PageNo > 0 {
			pageNo = filter.PageNo
		}
		if filter.PageSize > 0 {
			pageSize = filter.PageSize
		}
	}
	offset := (pageNo - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(&invites).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return invites, total, nil
}

func (r *InviteGormRepo) ListByCreator(ctx context.Context, createdBy string, filter *Filter) ([]*gormdb.InviteCode, int64, error) {
	var invites []*gormdb.InviteCode
	var total int64

	query := r.DB.WithContext(ctx).Model(&gormdb.InviteCode{}).Where("created_by = ?", createdBy)

	// 分页计算（同上）
	pageNo := 1
	pageSize := 10
	if filter != nil {
		if filter.PageNo > 0 {
			pageNo = filter.PageNo
		}
		if filter.PageSize > 0 {
			pageSize = filter.PageSize
		}
	}
	offset := (pageNo - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(&invites).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return invites, total, nil
}

func (r *InviteGormRepo) GetByID(ctx context.Context, id string) (*gormdb.InviteCode, error) {
	var ic gormdb.InviteCode
	err := r.DB.WithContext(ctx).First(&ic, "id = ?", id).Error
	return &ic, err
}

func (r *InviteGormRepo) GetByCode(ctx context.Context, code string) (*gormdb.InviteCode, error) {
	var ic gormdb.InviteCode
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&ic).Error
	return &ic, err
}

func (r *InviteGormRepo) Update(ctx context.Context, ic *gormdb.InviteCode) error {
	return r.DB.WithContext(ctx).Save(ic).Error
}

func (r *InviteGormRepo) Delete(ctx context.Context, ic *gormdb.InviteCode) error {
	return r.DB.WithContext(ctx).Delete(ic).Error
}
