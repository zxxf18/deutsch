package gormdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GermanState 德国联邦州
type GermanState struct {
	ID        string    `gorm:"primaryKey;type:char(36);not null" json:"id"`
	Slug      string    `gorm:"uniqueIndex;type:varchar(50);not null" json:"slug"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	NameCn    string    `gorm:"column:name_cn;type:varchar(100)" json:"nameCn"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName 指定表名
func (GermanState) TableName() string {
	return "german_states"
}

// BeforeCreate 生成 UUID
func (s *GermanState) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
