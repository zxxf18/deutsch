package question

import (
	"context"
	"math/rand"
	"sort"

	"deutsch/internal/types"
	"deutsch/model/gormdb"
	"deutsch/model/repository"
)

// BuildStateIDToSlug 构建州 id -> slug 映射，供 progress 等包复用
func BuildStateIDToSlug(states []*gormdb.GermanState) map[string]string {
	m := make(map[string]string, len(states))
	for _, s := range states {
		m[s.ID] = s.Slug
	}
	return m
}

// randShuffle 打乱切片
func randShuffle[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}

// stateResolver 根据题目 state_id 返回 "general" 或州 slug
func resolveState(q *gormdb.Question, idToSlug map[string]string) string {
	if q.StateID == nil || *q.StateID == "" {
		return "general"
	}
	if slug, ok := idToSlug[*q.StateID]; ok {
		return slug
	}
	return *q.StateID // 降级：未找到则用 id
}

// BuildQuestionItems 批量构建 QuestionItem（含选项与 state 字段）
// idToSlug: 州 id -> slug 映射，nil 时 State 留空
func BuildQuestionItems(ctx context.Context, repo repository.QuestionRepository, questions []*gormdb.Question, idToSlug map[string]string) ([]types.QuestionItem, error) {
	if len(questions) == 0 {
		return nil, nil
	}
	ids := make([]string, 0, len(questions))
	for _, q := range questions {
		ids = append(ids, q.ID)
	}
	optsMap, err := repo.GetOptionsByQuestionIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	items := make([]types.QuestionItem, 0, len(questions))
	for _, q := range questions {
		opts := optsMap[q.ID]
		state := ""
		if idToSlug != nil {
			state = resolveState(q, idToSlug)
		}
		items = append(items, toQuestionItem(q, opts, state))
	}
	return items, nil
}

// toQuestionItem 将 Question + 选项转为 QuestionItem
func toQuestionItem(q *gormdb.Question, opts []*gormdb.QuestionOption, state string) types.QuestionItem {
	item := types.QuestionItem{
		Id:          q.ID,
		QuestionDe:  q.QuestionDe,
		QuestionCn:  q.QuestionCn,
		Explanation: q.Explanation,
		HasImage:    q.HasImage,
		State:       state,
	}
	if len(opts) == 0 {
		return item
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].OptionIndex < opts[j].OptionIndex })
	optionsDe := make([]string, 0, len(opts))
	optionsCn := make([]string, 0, len(opts))
	for _, o := range opts {
		optionsDe = append(optionsDe, o.OptionDe)
		optionsCn = append(optionsCn, o.OptionCn)
		if o.IsCorrect {
			item.CorrectAnswer = o.OptionIndex
		}
	}
	item.OptionsDe = optionsDe
	item.OptionsCn = optionsCn
	return item
}

// BuildTrialQuestionItems 批量构建 TrialQuestionItem（不含 correctAnswer、explanation）
func BuildTrialQuestionItems(ctx context.Context, repo repository.QuestionRepository, questions []*gormdb.Question, idToSlug map[string]string) ([]types.TrialQuestionItem, error) {
	if len(questions) == 0 {
		return nil, nil
	}
	ids := make([]string, 0, len(questions))
	for _, q := range questions {
		ids = append(ids, q.ID)
	}
	optsMap, err := repo.GetOptionsByQuestionIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	items := make([]types.TrialQuestionItem, 0, len(questions))
	for _, q := range questions {
		opts := optsMap[q.ID]
		state := ""
		if idToSlug != nil {
			state = resolveState(q, idToSlug)
		}
		items = append(items, toTrialQuestionItem(q, opts, state))
	}
	return items, nil
}

func toTrialQuestionItem(q *gormdb.Question, opts []*gormdb.QuestionOption, state string) types.TrialQuestionItem {
	item := types.TrialQuestionItem{
		Id:         q.ID,
		QuestionDe: q.QuestionDe,
		QuestionCn: q.QuestionCn,
		HasImage:   q.HasImage,
		State:      state,
	}
	if len(opts) == 0 {
		return item
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].OptionIndex < opts[j].OptionIndex })
	optionsDe := make([]string, 0, len(opts))
	optionsCn := make([]string, 0, len(opts))
	for _, o := range opts {
		optionsDe = append(optionsDe, o.OptionDe)
		optionsCn = append(optionsCn, o.OptionCn)
	}
	item.OptionsDe = optionsDe
	item.OptionsCn = optionsCn
	return item
}
