package gormdb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExamAnswers 考试答案：question_id -> 选项下标 (0-based)
type ExamAnswers map[string]int

func (a ExamAnswers) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *ExamAnswers) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for ExamAnswers")
	}
	return json.Unmarshal(b, a)
}

// ExamRecord 考试记录
type ExamRecord struct {
	ID        string      `gorm:"primaryKey;type:char(36);not null" json:"id"`
	UserID    string      `gorm:"index;type:char(36);not null" json:"userId"`
	StateID   *string     `gorm:"type:char(36)" json:"stateId"` // 目标州，空为全德通用
	Total     int         `gorm:"not null" json:"total"`
	Score     int         `gorm:"not null" json:"score"`
	Passed    bool        `gorm:"not null" json:"passed"`
	Answers   ExamAnswers `gorm:"type:json" json:"answers"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"createdAt"`
}

func (ExamRecord) TableName() string {
	return "exam_records"
}

func (e *ExamRecord) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}
