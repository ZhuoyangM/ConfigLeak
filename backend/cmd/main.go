package main

import (
	"fmt"
	"os"

	"github.com/ZhuoyangM/ConfigLeak/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	DbName   string
	SSLMode  string
}

type Config struct {
	db DbConfig
}

func loadEnvConfig() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.db = DbConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DbName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	return cfg, nil
}

func initDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.db.Host, cfg.db.Username, cfg.db.Password, cfg.db.DbName, cfg.db.Port, cfg.db.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	err = db.AutoMigrate(&models.User{}, &models.ScanJob{}, &models.ScanResult{})
	if err != nil {
		return db, err
	}
	return db, nil
}

func main() {
	// load config
	cfg, err := loadEnvConfig()
	if err != nil {
		panic("failed to load env config")
	}

	// Initialize database
	_, err = initDB(cfg)
	if err != nil {
		panic("failed to initialize database")
	}

	// setup gin router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080")
}
