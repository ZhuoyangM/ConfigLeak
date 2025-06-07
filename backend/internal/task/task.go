package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ZhuoyangM/ConfigLeak/internal/utils"
	"github.com/hibiken/asynq"
)

const (
	TypeScanJob = "scanjob"
)

type ScanTaskPayload struct {
	UserID    uint
	TargetUrl string
}

func NewScanTask(userID uint, targetUrl string) (*asynq.Task, error) {
	task := ScanTaskPayload{
		UserID:    userID,
		TargetUrl: targetUrl,
	}
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScanJob, payload), nil
}

func HandleScanTask(ctx *context.Context, t *asynq.Task) error {
	var payload ScanTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	paths, err := utils.LoadPathsFromFile("backend/config/paths.yaml")
	if err != nil {
		return err
	}

	urls, err := utils.BuildScanUrls(payload.TargetUrl, paths)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		go scanUrl(&wg, url)
	}
	wg.Wait()

	return nil
}

func scanUrl(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error scanning URL %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()
	fmt.Printf("Scanned URL %s: Status Code %d\n", url, response.StatusCode)
}
