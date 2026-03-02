// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package progress

import (
	"net/http"
	"strconv"

	"deutsch/internal/logic/progress"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListExamRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseFilterFromQuery(r)
		l := progress.NewListExamRecordsLogic(r.Context(), svcCtx)
		resp, err := l.ListExamRecords(req)
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
