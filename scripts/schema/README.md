# 数据库脚本

## 初始化数据

首次部署或重建库后，执行 `init_data.sql` 插入初始数据：

```bash
mysql -u <user> -p <dbname> < init_data.sql
```

或使用 MySQL 客户端：

```sql
source /path/to/init_data.sql;
```

**前置条件**：表需已存在（应用首次启动时会执行 AutoMigrate 创建表）。

- **german_states**：联邦州数据（16 条），支持重复执行（按 slug 幂等）
- **users**、**invite_codes**：可选，取消注释后执行
