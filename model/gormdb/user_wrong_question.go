package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserWrongQuestion 用户错题本
type UserWrongQuestion struct {
	ID         string    `gorm:"primaryKey;type:char(36);not null" json:"id"`
	UserID     string    `gorm:"uniqueIndex:idx_user_wrong;type:char(36);not null" json:"userId"`
	QuestionID string    `gorm:"uniqueIndex:idx_user_wrong;type:char(36);not null" json:"questionId"`
	AddedAt    time.Time `gorm:"autoCreateTime" json:"addedAt"`
}

func (UserWrongQuestion) TableName() string {
	return "user_wrong_questions"
}

func (u *UserWrongQuestion) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
