package repository

import (
	"context"

	"deutsch/model/gormdb"
)

// ConfigRepository 配置与元数据
type ConfigRepository interface {
	ListStates(ctx context.Context) ([]*gormdb.GermanState, error)
}
