package gormdb

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InviteCode 邀请码模型
type InviteCode struct {
	ID        string    `gorm:"primaryKey;type:char(36);not null" json:"id"` // UUID string
	Code      string    `gorm:"uniqueIndex;type:varchar(20);not null" json:"code"`
	UsedBy    string    `gorm:"index;type:char(36)" json:"usedBy,omitempty"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
	IsEnabled bool      `gorm:"default:true" json:"is_enabled"`
	CreatedBy string    `gorm:"index;type:char(36)" json:"createdBy,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// BeforeCreate 生成UUID主键和默认过期时间
func (ic *InviteCode) BeforeCreate(tx *gorm.DB) error {
	if ic.ID == "" {
		ic.ID = uuid.New().String()
	}
	if ic.ExpiresAt.IsZero() {
		ic.ExpiresAt = time.Now().Add(30 * 24 * time.Hour) // 30天过期
	}
	if ic.Code == "" {
		ic.Code = GenerateInviteCode()
	}
	return nil
}

func GenerateInviteCode() string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// TableName 指定表名
func (ic *InviteCode) TableName() string {
	return "invite_codes"
}
