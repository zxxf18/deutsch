// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package question

import (
	"net/http"

	"deutsch/internal/logic/question"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TrialCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TrialCheckRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question.NewTrialCheckLogic(r.Context(), svcCtx)
		resp, err := l.TrialCheck(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
