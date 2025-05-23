package store

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	DbName   string
	SSLMode  string
}

func InitDB(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.DbName, cfg.Port, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &ScanJob{}, &ScanResult{})
}
