// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"net/http"

	"deutsch/internal/logic/progress"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdatePreferencesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePreferencesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := progress.NewUpdatePreferencesLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePreferences(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
