package repository

import (
	"context"

	"deutsch/model/gormdb"
)

// QuestionRepository 题目仓库
type QuestionRepository interface {
	GetGeneral(ctx context.Context) ([]*gormdb.Question, error)
	GetAll(ctx context.Context) ([]*gormdb.Question, error)
	GetByStateID(ctx context.Context, stateID string) ([]*gormdb.Question, error)
	GetByIDs(ctx context.Context, ids []string) ([]*gormdb.Question, error)
	GetByID(ctx context.Context, id string) (*gormdb.Question, error)
	GetOptionsByQuestionIDs(ctx context.Context, questionIDs []string) (map[string][]*gormdb.QuestionOption, error)
	GetRandomGeneral(ctx context.Context, n int) ([]*gormdb.Question, error)
	GetRandomByStateID(ctx context.Context, stateID string, n int) ([]*gormdb.Question, error)
	CountByStateID(ctx context.Context, stateID string) (int64, error)
}
