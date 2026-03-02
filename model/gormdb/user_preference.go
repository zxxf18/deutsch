package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserPreference 用户偏好（考试目标州、主题等）
type UserPreference struct {
	ID                   string    `gorm:"primaryKey;type:char(36);not null" json:"id"`
	UserID               string    `gorm:"uniqueIndex;type:char(36);not null" json:"userId"`
	PreferredExamStateID *string   `gorm:"type:char(36)" json:"preferredExamStateId"` // 模拟考试默认目标州，空则随机
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}

func (u *UserPreference) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
