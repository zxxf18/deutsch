package invitecode

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/svc"
)

func GetInviteCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := invitecode.NewGetInviteCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetInviteCode()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
