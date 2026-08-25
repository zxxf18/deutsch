package repository

import (
	"context"
	"time"

	"deutsch/model/gormdb"
)

// Filter 分页过滤器
type Filter struct {
	PageNo   int `form:"pageNo" json:"pageNo,omitempty"`
	PageSize int `form:"pageSize" json:"pageSize,omitempty"`
}

// InviteCodeListFilter 邀请码列表过滤（含分页与可用性筛选）
type InviteCodeListFilter struct {
	Filter
	AvailableOnly bool // 仅查看可用的邀请码（未使用、未过期、已启用）
}

// UserRepository 定义用户数据操作（UUID string键）
type UserRepository interface {
	Create(ctx context.Context, user *gormdb.User) error
	GetByEmail(ctx context.Context, email string) (*gormdb.User, error)
	GetByEmailIncludingDeleted(ctx context.Context, email string) (*gormdb.User, error) // 含软删除，用于注册时恢复
	GetByUsername(ctx context.Context, username string) (*gormdb.User, error)
	GetByUsernameIncludingDeleted(ctx context.Context, username string) (*gormdb.User, error)
	Restore(ctx context.Context, user *gormdb.User) error // 恢复软删除
	GetByPhone(ctx context.Context, phone string) (*gormdb.User, error)
	GetByUserID(ctx context.Context, userID string) (*gormdb.User, error)
	Update(ctx context.Context, user *gormdb.User) error
	Delete(ctx context.Context, user *gormdb.User) error
	List(ctx context.Context, filter *Filter) ([]*gormdb.User, int64, error) // 分页，admin用
}

// InviteCodeRepository 定义邀请码操作（UUID string键）
type InviteCodeRepository interface {
	CreateBatch(ctx context.Context, count int, createdBy string) ([]*gormdb.InviteCode, error)
	Validate(ctx context.Context, code string) (bool, *time.Time, error) // valid, expiresAt
	MarkUsed(ctx context.Context, code string, usedBy string) error      // usedBy string (UUID)
	List(ctx context.Context, filter *InviteCodeListFilter) ([]*gormdb.InviteCode, int64, error)
	ListByCreator(ctx context.Context, createdBy string, filter *Filter) ([]*gormdb.InviteCode, int64, error) // createdBy string (UUID)
	GetByID(ctx context.Context, id string) (*gormdb.InviteCode, error)
	GetByCode(ctx context.Context, code string) (*gormdb.InviteCode, error)
	Update(ctx context.Context, ic *gormdb.InviteCode) error
	Delete(ctx context.Context, ic *gormdb.InviteCode) error
}
