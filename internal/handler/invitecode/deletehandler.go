// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package invitecode

import (
	"net/http"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := invitecode.NewDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
