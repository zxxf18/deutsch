package invitecode

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/pkg/httpparse"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &types.DeleteInviteCodeRequest{ID: httpparse.PathID(r)}
		l := invitecode.NewDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
