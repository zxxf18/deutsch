package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic"
	"deutsch/internal/svc"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPingLogic(r.Context(), svcCtx)
		err := l.Ping()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "pong"})
		}
	}
}
