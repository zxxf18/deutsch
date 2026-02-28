package user

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/user"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func ListUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseFilterFromQuery(r)
		l := user.NewListUserLogic(r.Context(), svcCtx)
		resp, err := l.ListUser(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func parseFilterFromQuery(r *http.Request) *types.Filter {
	query := r.URL.Query()
	pageNo := 1
	if v := query.Get("pageNo"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pageNo = n
		}
	}
	pageSize := 10
	if v := query.Get("pageSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	return &types.Filter{PageNo: pageNo, PageSize: pageSize}
}
