package gormdb

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

// User 用户模型（UUID主键 + 软删除）
type User struct {
	ID           string         `gorm:"primaryKey;type:char(36);not null" json:"id"` // UUID string
	Username     string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Phone        string         `gorm:"uniqueIndex;type:varchar(20)" json:"phone,omitempty"` // 支持GetByPhone
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`                 // 隐藏密码
	Role         string         `gorm:"type:enum('user','admin','guest');default:user" json:"role"`
	Nickname     string         `gorm:"type:varchar(50)" json:"nickname"`
	IsEnabled    bool           `gorm:"default:true" json:"is_enabled"`
	Description  string         `gorm:"type:text" json:"description"`
	InviteCode   string         `gorm:"type:varchar(32)" json:"invite_code"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段（隐藏）
}

// BeforeCreate 生成UUID主键
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// BeforeSave 使用SHA256哈希密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.PasswordHash != "" {
		hash := sha256.Sum256([]byte(u.PasswordHash))
		u.PasswordHash = hex.EncodeToString(hash[:])
	}
	return nil
}

// TableName 指定表名
func (u *User) TableName() string {
	return "users"
}
