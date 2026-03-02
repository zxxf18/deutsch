package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserQuestionProgress 用户题目练习进度
type UserQuestionProgress struct {
	ID              string    `gorm:"primaryKey;type:char(36);not null" json:"id"`
	UserID          string    `gorm:"uniqueIndex:idx_user_question;type:char(36);not null" json:"userId"`
	QuestionID      string    `gorm:"uniqueIndex:idx_user_question;type:char(36);not null" json:"questionId"`
	CorrectCount    int       `gorm:"default:0" json:"correctCount"`
	WrongCount      int       `gorm:"default:0" json:"wrongCount"`
	LastPracticedAt time.Time `json:"lastPracticedAt"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (UserQuestionProgress) TableName() string {
	return "user_question_progress"
}

func (u *UserQuestionProgress) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
