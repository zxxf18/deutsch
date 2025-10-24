package gorm

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User GORM 用户模型
type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"` // UUID
	Phone       string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Password    string         `gorm:"type:varchar(64);not null" json:"-"` // SHA256
	Username    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Nickname    string         `gorm:"type:varchar(50)" json:"nickname"`
	IsEnabled   bool           `gorm:"default:true" json:"is_enabled"`
	Role        string         `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	Description string         `gorm:"type:text" json:"description"`
	InviteCode  string         `gorm:"type:varchar(12)" json:"invite_code"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	// 自定义方法
	gorm.Model
}

// HashPassword 生成 SHA256 密码
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// BeforeCreate 钩子：自动生成 UserId (UUID)
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserId == "" {
		u.UserId = uuid.New().String()
	}
	return nil
}

// FindByPhone 查询用户（忽略已删除）
func (u *User) FindByPhone(db *gorm.DB, phone string) (*User, error) {
	var user User
	err := db.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error
	return &user, err
}

// Insert 创建用户
func (u *User) Insert(db *gorm.DB) error {
	return db.Create(u).Error
}
