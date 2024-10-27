package storage

import (
	"fmt"
	"log"
	"sentinel/server/config"
	"sentinel/server/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	conf := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.PGConfig.Host, conf.PGConfig.Port, conf.PGConfig.User, conf.PGConfig.Password, conf.PGConfig.DB)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库表
	if err := DB.AutoMigrate(&model.GotifyMessage{}); err != nil {
		return err
	}

	log.Println("PostgreSQL Database initialized and migrated successfully.")
	return nil
}
