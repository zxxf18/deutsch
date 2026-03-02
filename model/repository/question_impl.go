package repository

import (
	"context"

	"deutsch/model/gormdb"

	"gorm.io/gorm"
)

// QuestionGormRepo 题目仓库 GORM 实现
type QuestionGormRepo struct {
	DB *gorm.DB
}

// NewQuestionGormRepo 创建题目仓库
func NewQuestionGormRepo(db *gorm.DB) QuestionRepository {
	return &QuestionGormRepo{DB: db}
}

// GetGeneral 通用题（state_id = 通用州 ID），按 sort_order 排序
func (r *QuestionGormRepo) GetGeneral(ctx context.Context) ([]*gormdb.Question, error) {
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Where("state_id = ?", gormdb.GeneralStateID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

// GetAll 所有题目（含通用与各州），按 sort_order 排序
func (r *QuestionGormRepo) GetAll(ctx context.Context) ([]*gormdb.Question, error) {
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Order("sort_order ASC").Find(&list).Error
	return list, err
}

// GetByStateID 某州题目，按 sort_order 排序
func (r *QuestionGormRepo) GetByStateID(ctx context.Context, stateID string) ([]*gormdb.Question, error) {
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Where("state_id = ?", stateID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

// GetByIDs 按 ID 列表查询
func (r *QuestionGormRepo) GetByIDs(ctx context.Context, ids []string) ([]*gormdb.Question, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&list).Error
	return list, err
}

// GetByID 单题
func (r *QuestionGormRepo) GetByID(ctx context.Context, id string) (*gormdb.Question, error) {
	var q gormdb.Question
	err := r.DB.WithContext(ctx).First(&q, "id = ?", id).Error
	return &q, err
}

// GetOptionsByQuestionIDs 批量查询选项，返回 map[questionID][]Option
func (r *QuestionGormRepo) GetOptionsByQuestionIDs(ctx context.Context, questionIDs []string) (map[string][]*gormdb.QuestionOption, error) {
	if len(questionIDs) == 0 {
		return nil, nil
	}
	var list []*gormdb.QuestionOption
	err := r.DB.WithContext(ctx).Where("question_id IN ?", questionIDs).Order("question_id, option_index").Find(&list).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string][]*gormdb.QuestionOption)
	for _, o := range list {
		m[o.QuestionID] = append(m[o.QuestionID], o)
	}
	return m, nil
}

// GetRandomGeneral 随机取 n 道通用题
func (r *QuestionGormRepo) GetRandomGeneral(ctx context.Context, n int) ([]*gormdb.Question, error) {
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Where("state_id = ?", gormdb.GeneralStateID).Order("RAND()").Limit(n).Find(&list).Error
	return list, err
}

// GetRandomByStateID 随机取 n 道该州题目
func (r *QuestionGormRepo) GetRandomByStateID(ctx context.Context, stateID string, n int) ([]*gormdb.Question, error) {
	var list []*gormdb.Question
	err := r.DB.WithContext(ctx).Where("state_id = ?", stateID).Order("RAND()").Limit(n).Find(&list).Error
	return list, err
}
