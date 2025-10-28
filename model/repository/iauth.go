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

// UserRepository 定义用户数据操作（UUID string键）
type UserRepository interface {
	Create(ctx context.Context, user *gormdb.User) error
	GetByEmail(ctx context.Context, email string) (*gormdb.User, error)
	GetByPhone(ctx context.Context, phone string) (*gormdb.User, error)
	GetByUserID(ctx context.Context, userID string) (*gormdb.User, error)
	Update(ctx context.Context, user *gormdb.User) error
	List(ctx context.Context, filter *Filter) ([]*gormdb.User, int64, error) // 分页，admin用
}

// InviteCodeRepository 定义邀请码操作（UUID string键）
type InviteCodeRepository interface {
	CreateBatch(ctx context.Context, count int, createdBy string) ([]string, error) // 返回生成的codes
	Validate(ctx context.Context, code string) (bool, *time.Time, error)            // valid, expiresAt
	MarkUsed(ctx context.Context, code string, usedBy string) error                 // usedBy string (UUID)
	List(ctx context.Context, filter *Filter) ([]*gormdb.InviteCode, int64, error)
	ListByCreator(ctx context.Context, createdBy string, filter *Filter) ([]*gormdb.InviteCode, int64, error) // createdBy string (UUID)
	GetByCode(ctx context.Context, code string) (*gormdb.InviteCode, error)
}
