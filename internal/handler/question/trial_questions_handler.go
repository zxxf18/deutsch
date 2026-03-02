// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"net/http"

	"deutsch/internal/logic/question"
	"deutsch/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TrialQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := question.NewTrialQuestionsLogic(r.Context(), svcCtx)
		resp, err := l.TrialQuestions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
