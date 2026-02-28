package invitecode

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/logic/invitecode"
	"deutsch/internal/svc"
	"deutsch/internal/types"
)

func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseListReqFromQuery(r)
		l := invitecode.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func parseListReqFromQuery(r *http.Request) *types.ListInviteCodeRequest {
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
	availableOnly := false
	if v := query.Get("availableOnly"); v == "1" || v == "true" || v == "True" {
		availableOnly = true
	}
	return &types.ListInviteCodeRequest{
		Filter:        types.Filter{PageNo: pageNo, PageSize: pageSize},
		AvailableOnly: availableOnly,
	}
}
