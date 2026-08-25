.PHONY: api build build-linux-amd64 init-db

# 重新生成 API 代码（生成后自动将 interface{} 替换为 any）
api:
	goctl api go -api api/deutsch.api -dir . --style go_zero
	sed -i '' 's/interface{}/any/g' internal/types/types.go

build:
	go build ./...

build-linux-amd64:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/deutsch-linux-amd64 ./deutsch.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/deutsch-dbinit-linux-amd64 ./cmd/dbinit

init-db:
	@test -n "$(ADMIN_PASSWORD)" || (echo "ADMIN_PASSWORD is required" && exit 1)
	DEUTSCH_ADMIN_PASSWORD="$(ADMIN_PASSWORD)" go run ./cmd/dbinit -f "$${CONFIG_FILE:-etc/deutsch.yaml}" -admin-username "$${ADMIN_USERNAME:-admin}" -admin-email "$${ADMIN_EMAIL:-admin@deutsch.local}"
