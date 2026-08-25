package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"deutsch/internal/config"
	"deutsch/internal/pkg/passwordcrypto"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/gorm"
)

var (
	configFile    = flag.String("f", "etc/deutsch.yaml", "runtime config file")
	resetUsers    = flag.Bool("reset-users", false, "delete all users and user-owned data before creating the admin")
	adminUser     = flag.String("admin-username", "admin", "administrator username")
	adminEmail    = flag.String("admin-email", "admin@example.com", "administrator login email")
	adminPassword = flag.String("admin-password", "123456678", "administrator password (prefer DEUTSCH_ADMIN_PASSWORD to avoid process-list exposure)")
)

func main() {
	flag.Parse()
	password := os.Getenv("DEUTSCH_ADMIN_PASSWORD")
	if password == "" {
		password = *adminPassword
	}
	var c config.Config
	conf.MustLoad(*configFile, &c)
	if c.MySQL.DataSource == "" {
		log.Fatal("MySQL DataSource is required")
	}
	passwordCipher, err := passwordcrypto.New(c.PasswordEncryption.Key)
	if err != nil {
		log.Fatalf("invalid password encryption config: %v", err)
	}
	if err := gormdb.InitDB(c.MySQL.DataSource); err != nil {
		log.Fatalf("initialize database: %v", err)
	}

	encrypted, err := passwordCipher.Encrypt(password)
	if err != nil {
		log.Fatalf("encrypt administrator password: %v", err)
	}
	err = gormdb.DB.Transaction(func(tx *gorm.DB) error {
		if *resetUsers {
			for _, table := range []string{
				"user_wrong_questions",
				"user_question_progress",
				"exam_records",
				"user_preferences",
				"invite_codes",
				"users",
			} {
				if err := tx.Exec("DELETE FROM " + table).Error; err != nil {
					return fmt.Errorf("clear %s: %w", table, err)
				}
			}
		}
		admin := &gormdb.User{
			Username:          *adminUser,
			Email:             *adminEmail,
			PasswordEncrypted: encrypted,
			Role:              string(gormdb.UserRoleAdmin),
			Nickname:          "Administrator",
			IsEnabled:         true,
		}
		return tx.Where("email = ?", *adminEmail).Assign(admin).FirstOrCreate(admin).Error
	})
	if err != nil {
		log.Fatalf("initialize administrator: %v", err)
	}
	fmt.Printf("database initialized; administrator=%s email=%s reset_users=%t\n", *adminUser, *adminEmail, *resetUsers)
}
