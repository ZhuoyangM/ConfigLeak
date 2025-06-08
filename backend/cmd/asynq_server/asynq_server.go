package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "embed"

	"github.com/ZhuoyangM/ConfigLeak/internal/task"
	"github.com/ZhuoyangM/ConfigLeak/internal/utils"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
		return
	}

	redisConfig := asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	asynqConfig := asynq.Config{
		Concurrency: 20,
	}
	server := asynq.NewServer(
		redisConfig,
		asynqConfig,
	)

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	paths, err := utils.LoadPathsFromFile(os.Getenv("YAML_FILE_PATH"))
	if err != nil {
		log.Printf("Error loading paths from file: %v\n", err)
		return
	}
	taskHandler := task.NewScanHandler(httpClient, paths)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeScanJob, taskHandler.HandleScanTask)

	if err := server.Run(mux); err != nil {
		log.Printf("Error running asynq server: %v\n", err)
	}

}
