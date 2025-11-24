package config

import (
	"fmt"
	"myapp/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtSecret = []byte("your_super_secret_key") // 全局 JWT 密钥

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ Failed to connect to database")
	}

	// 自动迁移
	DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Announcement{}, &models.Mistake{})
	fmt.Println("✅ Database migrated successfully")
}