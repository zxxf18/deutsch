// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"net/http"

	"deutsch/internal/logic/progress"
	"deutsch/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListWrongQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseFilterFromQuery(r)
		l := progress.NewListWrongQuestionsLogic(r.Context(), svcCtx)
		resp, err := l.ListWrongQuestions(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
