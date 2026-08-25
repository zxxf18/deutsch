# 数据库脚本

## 初始化数据

首次部署或重建库后，执行 `init_data.sql` 插入初始数据：

```bash
mysql --default-character-set=utf8mb4 -u <user> -p <dbname> < init_data.sql
```

或使用 MySQL 客户端：

```sql
source /path/to/init_data.sql;
```

**前置条件**：表需已存在（应用首次启动时会执行 AutoMigrate 创建表）。若 `question_options` 报错 "Column count doesn't match"，请先启动应用做一次迁移，或在 MySQL 8.0.12+ 下执行脚本（含 `ADD COLUMN IF NOT EXISTS` 自动补列）。

- **german_states**：联邦州数据（16 条），支持重复执行（按 slug 幂等）
- **questions**、**question_options**：题目数据（300 通用 + 160 州题）；带图题目（has_image=1）的选项含 image_path（如 wappen/berlin.svg），图片存于 `deutsch/assets/wappen/`，通过 `GET /api/v1/assets/wappen/:file` 获取
- **州徽图片**：运行 `cd .. && python3 fetch_and_update.py` 从 Wikipedia Commons 下载并保存到 `assets/wappen/`，同时更新 init_data.sql 中的 image_path
- **users**、**invite_codes**：可选，取消注释后执行

## 创建或重置管理员

管理员密码由应用使用 AES-256-GCM 加密，不能直接写入 SQL。配置好未跟踪的
`etc/deutsch.yaml` 后执行：

```bash
go run ./cmd/dbinit -f etc/deutsch.yaml
```

默认管理员为 `admin / admin@example.com / 123456678`。固定默认密码仅适用于初始化，
部署后应立即通过 `DEUTSCH_ADMIN_PASSWORD` 或 `-admin-password` 覆盖。

如需清空现有用户、学习进度、考试记录、错题和邀请码后重建管理员，追加
`-reset-users`。此参数会删除用户相关数据，仅用于明确的重置操作。

**题目数据**：从 deutsch-fe 的 `general-questions.md`、`state-questions.md` 提取。
