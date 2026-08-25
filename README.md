# Deutsch

德国入籍考试题库与学习平台后端服务。

## 功能概述

- **题库管理**：300 道通用题 + 16 州各 10 道州题，支持按州/通用查询
- **模拟考试**：随机抽题（30 通用 + 3 州题），提交阅卷，返回正确答案与解析
- **学习进度**：练习记录、考试记录、错题本
- **用户系统**：注册（需邀请码）、登录、JWT 鉴权
- **管理后台**：用户管理、邀请码生成与管理（需 admin 角色）

## 技术栈

- **框架**：go-zero
- **数据库**：MySQL
- **缓存**：Redis
- **ORM**：GORM

## 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis

## 快速开始

### 1. 配置

复制并编辑 `etc/deutsch.yaml`：

```yaml
Name: deutsch
Host: 0.0.0.0
Port: 8888

JWTAuth:
  AccessSecret: "your-secret"
  AccessExpire: 7200

MySQL:
  DataSource: "user:pass@tcp(host:3306)/deutsch?charset=utf8mb4&parseTime=True"

PasswordEncryption:
  # 使用高熵密钥材料；程序通过 SHA-256 派生 AES-256 密钥
  # 真实值只放在未跟踪的 etc/deutsch.yaml 中
  Key: "replace-with-random-key-material"

Redis:
  Host: "localhost:6379"
  Type: "node"
  Pass: ""
```

### 2. 初始化数据库

```bash
mysql --default-character-set=utf8mb4 -u root -p < scripts/schema/init_data.sql

# 创建默认管理员：admin / admin@example.com / admin120420
# 固定默认密码仅用于初始化，部署后应立即替换
go run ./cmd/dbinit -f etc/deutsch.yaml
```

### 3. 运行

```bash
make build
./deutsch -f etc/deutsch.yaml
```

或开发模式：

```bash
go run deutsch.go -f etc/deutsch.yaml
```

## 命令

| 命令 | 说明 |
|------|------|
| `make api` | 根据 api 定义重新生成代码 |
| `make build` | 编译项目 |
| `make build-linux-amd64` | 交叉编译 Linux amd64 服务和数据库初始化工具 |
| `make init-db` | 创建默认管理员（可通过 ADMIN_PASSWORD 等变量覆盖） |

## 接口文档

详见 [doc/api.md](doc/api.md)。

## 项目结构

```
deutsch/
├── api/              # API 定义（goctl）
├── etc/               # 配置文件
├── internal/          # 业务逻辑
│   ├── handler/       # HTTP 处理器
│   ├── logic/        # 业务逻辑
│   └── middleware/   # 中间件
├── model/             # 数据模型与仓库
├── scripts/schema/    # 数据库初始化脚本
└── doc/               # 接口文档
```
