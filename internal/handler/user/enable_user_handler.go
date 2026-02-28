package user

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/user"
	"deutsch/internal/pkg/httpparse"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func EnableUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			IsEnabled bool `json:"is_enabled"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		req := &types.EnableUserRequest{ID: httpparse.PathID(r), IsEnabled: body.IsEnabled}
		l := user.NewEnableUserLogic(r.Context(), svcCtx)
		resp, err := l.EnableUser(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
