package gormdb

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局DB实例（在svc中注入）
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 生产环境Warn
	})
	if err != nil {
		return err
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移
	if err := AutoMigrate(); err != nil {
		return err
	}

	return nil
}

// BuildDSN 生成MySQL DSN字符串
// 参数:
//   - username: 用户名
//   - password: 密码
//   - host: 主机地址
//   - port: 端口，默认3306
//   - dbname: 数据库名
//   - opts: 额外选项（如"charset=utf8mb4&parseTime=True"）
func BuildDSN(username, password, host string, port int, dbname string, opts ...string) string {
	if port == 0 {
		port = 3306
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname)
	if len(opts) > 0 {
		dsn += "?" + strings.Join(opts, "&")
	}
	return dsn
}

// AutoMigrate 执行模型迁移
func AutoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&InviteCode{},
	)
}
