package main

import (
	"fmt"
	"os"

	"github.com/ZhuoyangM/ConfigLeak/internal/task"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func singleTask(client *asynq.Client) {
	scanTask, err := task.NewScanTask(1208, "http://example.com")
	if err != nil {
		fmt.Println("Error creating scan task:", err)
		return
	}

	info, err := client.Enqueue(scanTask)
	if err != nil {
		fmt.Println("Error enqueueing task:", err)
		return
	}

	fmt.Printf("Enqueued task: %s with ID: %s\n", scanTask.Type(), info.ID)

}

func manyTasks(client *asynq.Client) {
	for i := 0; i < 30; i++ {
		scanTask, err := task.NewScanTask(uint(i+1), fmt.Sprintf("http://example%d.com/", i+1))
		if err != nil {
			fmt.Println("Error creating scan task:", err)
			continue
		}

		info, err := client.Enqueue(scanTask)
		if err != nil {
			fmt.Println("Error enqueueing task:", err)
			continue
		}

		fmt.Printf("Enqueued task: %s with ID: %s\n", scanTask.Type(), info.ID)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	redisConfig := asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		DB:       0,
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	client := asynq.NewClient(redisConfig)

	manyTasks(client)

}
