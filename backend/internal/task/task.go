package task

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

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

	before := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, 30)
	wg.Add(len(urls))
	for _, url := range urls {
		sem <- struct{}{}
		go func(url string) {
			defer func() { <-sem }()
			defer wg.Done()
			scanUrl(h.Client, url)
		}(url)
	}
	wg.Wait()
	elapsed := time.Since(before)
	log.Printf("Scan completed for user %d on target URL %s in %s\n", payload.UserID, payload.TargetUrl, elapsed)
	return nil
}

func NewScanHandler(client *http.Client, paths []string) *ScanHandler {
	return &ScanHandler{
		Client: client,
		Paths:  paths,
	}
}

func scanUrl(client *http.Client, url string) {
	response, err := client.Get(url)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("URL %s does not exist", url)
		} else if os.IsTimeout(err) {
			log.Printf("Timeout while scanning URL %s", url)
		} else {
			log.Printf("Error scanning URL %s: %v\n", url, err)
		}
		return
	}
	defer response.Body.Close()
}
