.PHONY: api build

# 重新生成 API 代码（生成后自动将 interface{} 替换为 any）
api:
	goctl api go -api api/deutsch.api -dir . --style go_zero
	sed -i '' 's/interface{}/any/g' internal/types/types.go

build:
	go build ./...
