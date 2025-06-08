package task

import (
	"context"
	"encoding/json"
	"log"
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

type ScanHandler struct {
	Client *http.Client
	Paths  []string
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

func (h *ScanHandler) HandleScanTask(ctx context.Context, t *asynq.Task) error {
	var payload ScanTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("Error unmarshalling task payload: %v\n", err)
		return err
	}

	urls, err := utils.BuildScanUrls(payload.TargetUrl, h.Paths)
	if err != nil {
		log.Printf("Error building scan URLs: %v\n", err)
		return err
	}

	log.Printf("Starting scan for user %d on target URL %s with %d paths\n", payload.UserID, payload.TargetUrl, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		go scanUrl(h.Client, &wg, url)
	}
	wg.Wait()

	log.Printf("Scan completed for user %d on target URL %s\n", payload.UserID, payload.TargetUrl)
	return nil
}

func NewScanHandler(client *http.Client, paths []string) *ScanHandler {
	return &ScanHandler{
		Client: client,
		Paths:  paths,
	}
}

func scanUrl(client *http.Client, wg *sync.WaitGroup, url string) {
	defer wg.Done()
	response, err := client.Get(url)
	if err != nil {
		log.Printf("Error scanning URL %s: %v\n", url, err)
		return
	}
	defer response.Body.Close()
	log.Printf("Scanned URL %s: Status Code %d\n", url, response.StatusCode)
}
