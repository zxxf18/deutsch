package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/user"
	"deutsch/internal/svc"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewDeleteUserLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
