package repository

import (
	"context"

	"deutsch/model/gormdb"
)

// ProgressRepository 用户进度仓储
type ProgressRepository interface {
	// UserPreference
	GetPreference(ctx context.Context, userID string) (*gormdb.UserPreference, error)
	UpsertPreference(ctx context.Context, pref *gormdb.UserPreference) error

	// Learning progress
	UpsertQuestionProgress(ctx context.Context, userID, questionID string, correct bool) error
	GetProgressByUser(ctx context.Context, userID string) ([]*gormdb.UserQuestionProgress, error)

	// Exam records
	CreateExamRecord(ctx context.Context, record *gormdb.ExamRecord) error
	GetExamRecordsByUser(ctx context.Context, userID string, offset, limit int) ([]*gormdb.ExamRecord, int64, error)
	GetExamRecordByID(ctx context.Context, id, userID string) (*gormdb.ExamRecord, error)

	// Wrong questions
	AddWrongQuestion(ctx context.Context, userID, questionID string) error
	RemoveWrongQuestion(ctx context.Context, userID, questionID string) error
	GetWrongQuestionIDs(ctx context.Context, userID string, offset, limit int) ([]string, int64, error)
}
