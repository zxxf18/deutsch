// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"net/http"

	"deutsch/internal/logic/question"
	"deutsch/internal/pkg/httpparse"
	"deutsch/internal/svc"
	"deutsch/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &types.GetQuestionRequest{QuestionId: httpparse.PathVar(r, "question_id")}

		l := question.NewGetQuestionLogic(r.Context(), svcCtx)
		resp, err := l.GetQuestion(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
