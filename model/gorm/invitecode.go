package gorm

import (
	"database/sql"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// InviteCode GORM 邀请码模型
type InviteCode struct {
	ID        uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string        `gorm:"type:varchar(12);uniqueIndex;not null" json:"code"`
	CreatorId uint          `gorm:"not null" json:"creator_id"` // users.id
	UsedBy    sql.NullInt64 `gorm:"-" json:"used_by"`           // users.id，可空
	IsUsed    bool          `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

	gorm.Model
}

// FindByCode 查询未使用邀请码
func (ic *InviteCode) FindByCode(db *gorm.DB, code string) (*InviteCode, error) {
	var invite InviteCode
	err := db.Where("code = ? AND is_used = false", code).First(&invite).Error
	return &invite, err
}

// Insert 创建邀请码
func (ic *InviteCode) Insert(db *gorm.DB) error {
	return db.Create(ic).Error
}

// UpdateUsed 更新为已使用
func (ic *InviteCode) UpdateUsed(db *gorm.DB, usedBy uint) error {
	return db.Model(ic).Updates(map[string]interface{}{
		"is_used": true,
		"used_by": usedBy,
	}).Error
}

// GenerateInviteCode 生成 10 位随机邀请码
func GenerateInviteCode() string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
