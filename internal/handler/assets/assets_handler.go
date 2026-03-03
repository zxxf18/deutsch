package assets

import (
	"net/http"
	"path/filepath"
	"strings"

	"deutsch/internal/svc"
)

// AssetsHandler 返回静态文件 handler，路径 /api/v1/assets/wappen/:file
func AssetsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := strings.TrimPrefix(r.URL.Path, "/api/v1/assets/wappen/")
		file = strings.TrimPrefix(file, "/")
		if file == "" || strings.Contains(file, "..") || strings.Contains(file, "/") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		fullPath := filepath.Join(svcCtx.AssetsDir, "wappen", file)
		rel, err := filepath.Rel(svcCtx.AssetsDir, fullPath)
		if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if strings.HasSuffix(strings.ToLower(file), ".svg") {
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		http.ServeFile(w, r, fullPath)
	}
}
