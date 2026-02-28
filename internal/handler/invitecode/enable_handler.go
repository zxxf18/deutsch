package invitecode

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/pkg/httpparse"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func EnableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			IsEnabled bool `json:"is_enabled"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		req := &types.EnableInviteCodeRequest{ID: httpparse.PathID(r), IsEnabled: body.IsEnabled}
		l := invitecode.NewEnableLogic(r.Context(), svcCtx)
		resp, err := l.Enable(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
