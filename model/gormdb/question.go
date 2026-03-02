package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Question 题目
type Question struct {
	ID          string    `gorm:"primaryKey;type:char(36);not null" json:"id"`
	StateID     *string   `gorm:"index;type:char(36)" json:"stateId,omitempty"`
	QuestionDe  string    `gorm:"column:question_de;type:text;not null" json:"questionDe"`
	QuestionCn  string    `gorm:"column:question_cn;type:text" json:"questionCn"`
	Explanation string    `gorm:"type:text" json:"explanation"`
	HasImage    bool      `gorm:"column:has_image;default:false" json:"hasImage"`
	SortOrder   int       `gorm:"column:sort_order;default:0" json:"sortOrder"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}

// BeforeCreate 生成 UUID
func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	return nil
}

// QuestionOption 题目选项
type QuestionOption struct {
	ID          string `gorm:"primaryKey;type:char(36);not null" json:"id"`
	QuestionID  string `gorm:"column:question_id;index;type:char(36);not null" json:"questionId"`
	OptionIndex int    `gorm:"column:option_index;not null" json:"optionIndex"`
	OptionDe    string `gorm:"column:option_de;type:varchar(500);not null" json:"optionDe"`
	OptionCn    string `gorm:"column:option_cn;type:varchar(500)" json:"optionCn"`
	IsCorrect   bool   `gorm:"column:is_correct;default:false" json:"isCorrect"`
}

// TableName 指定表名
func (QuestionOption) TableName() string {
	return "question_options"
}

// BeforeCreate 生成 UUID
func (o *QuestionOption) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}
