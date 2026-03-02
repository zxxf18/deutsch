package httpparse

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

// PathVar 从 request 提取 path 中指定名称的变量
func PathVar(r *http.Request, name string) string {
	if vars := pathvar.Vars(r); vars != nil {
		if v := vars[name]; v != "" {
			return v
		}
	}
	return ""
}

// PathID 从 request 提取 path 中的 id，优先用 pathvar，否则从 URL 解析
func PathID(r *http.Request) string {
	if vars := pathvar.Vars(r); vars != nil {
		if id := vars["id"]; id != "" {
			return id
		}
	}
	// pathvar 未设置时从 URL 解析，如 /api/v1/user/xxx 或 /api/v1/user/xxx/enable
	path := r.URL.Path
	for _, prefix := range []string{"/api/v1/user/", "/api/v1/invitecode/", "/api/v1/invitecode/validate/"} {
		if after, ok := strings.CutPrefix(path, prefix); ok && after != "" {
			if idx := strings.Index(after, "/"); idx >= 0 {
				return after[:idx]
			}
			return after
		}
	}
	return ""
}
