package invitecode

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/pkg/httpparse"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func GetInviteCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &types.GetInviteCodeRequest{ID: httpparse.PathID(r)}
		l := invitecode.NewGetInviteCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetInviteCode(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
