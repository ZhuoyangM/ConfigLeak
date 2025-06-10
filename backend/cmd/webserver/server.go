package main

import (
	"os"

	"github.com/ZhuoyangM/ConfigLeak/internal/controllers"
	store "github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/ZhuoyangM/ConfigLeak/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

type Config struct {
	DB store.DBConfig
}

func loadEnvConfig() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.DB = store.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DbName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	return cfg, nil
}

func main() {
	// load config
	cfg, err := loadEnvConfig()
	if err != nil {
		panic("failed to load env config")
	}

	// Initialize database
	db, err := store.InitDB(cfg.DB)
	if err != nil {
		panic("failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sqlDB")
	}
	defer sqlDB.Close()

	//Migrate the database
	err = store.Migrate(db)
	if err != nil {
		panic("failed to migrate database")
	}

	//setup services
	userService := store.NewUserService(db)
	scanService := store.NewScanService(db)

	redisConfig := asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	asynqClient := asynq.NewClient(redisConfig)

	//setup controllers
	userController := controllers.UserController{
		UserService: userService,
	}
	scanController := controllers.ScanController{
		ScanService: scanService,
		AsynqClient: asynqClient,
	}

	// setup gin router
	router := gin.Default()
	public := router.Group("/api")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World",
			})
		})

		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)

		auth := public.Group("/user", utils.JWTMiddleware())
		{
			auth.GET("/profile", userController.GetUserInfo)
		}

		scan := public.Group("/jobs", utils.JWTMiddleware())
		{
			scan.POST("/", scanController.StartScan)
			scan.GET("/", scanController.GetAllScanJobs)
			scan.GET("/:id", scanController.GetScanJob)
			scan.DELETE("/:id", scanController.DeleteScanJob)
			scan.GET("/:id/results", scanController.GetScanResults)
		}

	}

	router.Run(":8000")
}
