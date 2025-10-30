// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package invitecode

import (
	"net/http"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func EnableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnableInviteCodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := invitecode.NewEnableLogic(r.Context(), svcCtx)
		resp, err := l.Enable(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
