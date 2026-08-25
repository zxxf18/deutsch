package repository

import (
	"context"
	"time"

	"deutsch/model/gormdb"

	"gorm.io/gorm"
)

// ProgressGormRepo 进度仓储 GORM 实现
type ProgressGormRepo struct {
	DB *gorm.DB
}

// NewProgressGormRepo 创建进度仓储
func NewProgressGormRepo(db *gorm.DB) ProgressRepository {
	return &ProgressGormRepo{DB: db}
}

func (r *ProgressGormRepo) GetPreference(ctx context.Context, userID string) (*gormdb.UserPreference, error) {
	var pref gormdb.UserPreference
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pref, err
}

func (r *ProgressGormRepo) UpsertPreference(ctx context.Context, pref *gormdb.UserPreference) error {
	existing, err := r.GetPreference(ctx, pref.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		return r.DB.WithContext(ctx).Create(pref).Error
	}
	return r.DB.WithContext(ctx).Model(existing).Updates(map[string]interface{}{
		"preferred_exam_state_id": pref.PreferredExamStateID,
		"updated_at":              time.Now(),
	}).Error
}

func (r *ProgressGormRepo) UpsertQuestionProgress(ctx context.Context, userID, questionID string, correct bool) error {
	now := time.Now()
	var p gormdb.UserQuestionProgress
	err := r.DB.WithContext(ctx).Where("user_id = ? AND question_id = ?", userID, questionID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		p = gormdb.UserQuestionProgress{
			UserID:          userID,
			QuestionID:      questionID,
			CorrectCount:    0,
			WrongCount:      0,
			LastPracticedAt: now,
		}
		if correct {
			p.CorrectCount = 1
		} else {
			p.WrongCount = 1
		}
		return r.DB.WithContext(ctx).Create(&p).Error
	}
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"last_practiced_at": now,
	}
	if correct {
		upd["correct_count"] = p.CorrectCount + 1
	} else {
		upd["wrong_count"] = p.WrongCount + 1
	}
	return r.DB.WithContext(ctx).Model(&p).Updates(upd).Error
}

func (r *ProgressGormRepo) RecordPractice(ctx context.Context, userID, questionID string, correct bool) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &ProgressGormRepo{DB: tx}
		if err := txRepo.UpsertQuestionProgress(ctx, userID, questionID, correct); err != nil {
			return err
		}
		if !correct {
			return txRepo.AddWrongQuestion(ctx, userID, questionID)
		}
		return nil
	})
}

func (r *ProgressGormRepo) GetProgressByUser(ctx context.Context, userID string) ([]*gormdb.UserQuestionProgress, error) {
	var list []*gormdb.UserQuestionProgress
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *ProgressGormRepo) CreateExamRecord(ctx context.Context, record *gormdb.ExamRecord) error {
	return r.DB.WithContext(ctx).Create(record).Error
}

func (r *ProgressGormRepo) CreateExamRecordWithWrongQuestions(ctx context.Context, record *gormdb.ExamRecord, wrongQuestionIDs []string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &ProgressGormRepo{DB: tx}
		if err := txRepo.CreateExamRecord(ctx, record); err != nil {
			return err
		}
		for _, questionID := range wrongQuestionIDs {
			if err := txRepo.AddWrongQuestion(ctx, record.UserID, questionID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProgressGormRepo) GetExamRecordsByUser(ctx context.Context, userID string, offset, limit int) ([]*gormdb.ExamRecord, int64, error) {
	var list []*gormdb.ExamRecord
	var total int64
	db := r.DB.WithContext(ctx).Model(&gormdb.ExamRecord{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *ProgressGormRepo) GetExamRecordByID(ctx context.Context, id, userID string) (*gormdb.ExamRecord, error) {
	var rec gormdb.ExamRecord
	err := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&rec).Error
	return &rec, err
}

func (r *ProgressGormRepo) AddWrongQuestion(ctx context.Context, userID, questionID string) error {
	wq := &gormdb.UserWrongQuestion{
		UserID:     userID,
		QuestionID: questionID,
	}
	return r.DB.WithContext(ctx).Where("user_id = ? AND question_id = ?", userID, questionID).
		FirstOrCreate(wq).Error
}

func (r *ProgressGormRepo) RemoveWrongQuestion(ctx context.Context, userID, questionID string) error {
	return r.DB.WithContext(ctx).Where("user_id = ? AND question_id = ?", userID, questionID).
		Delete(&gormdb.UserWrongQuestion{}).Error
}

func (r *ProgressGormRepo) GetWrongQuestionIDs(ctx context.Context, userID string, offset, limit int) ([]string, int64, error) {
	var list []*gormdb.UserWrongQuestion
	var total int64
	db := r.DB.WithContext(ctx).Model(&gormdb.UserWrongQuestion{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("added_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	ids := make([]string, 0, len(list))
	for _, w := range list {
		ids = append(ids, w.QuestionID)
	}
	return ids, total, nil
}
