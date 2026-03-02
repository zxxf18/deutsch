package repository

import (
	"context"

	"deutsch/model/gormdb"

	"gorm.io/gorm"
)

// ConfigGormRepo 配置仓库 GORM 实现
type ConfigGormRepo struct {
	DB *gorm.DB
}

// NewConfigGormRepo 创建配置仓库
func NewConfigGormRepo(db *gorm.DB) ConfigRepository {
	return &ConfigGormRepo{DB: db}
}

// ListStates 联邦州列表，按 sort_order 排序
func (r *ConfigGormRepo) ListStates(ctx context.Context) ([]*gormdb.GermanState, error) {
	var states []*gormdb.GermanState
	err := r.DB.WithContext(ctx).Order("sort_order ASC, slug ASC").Find(&states).Error
	return states, err
}
