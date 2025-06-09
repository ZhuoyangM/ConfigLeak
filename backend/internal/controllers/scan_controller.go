package controllers

import (
	"github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type ScanController struct {
	ScanService *store.ScanService
	AsynqClient *asynq.Client
}

func NewScanController(scanService *store.ScanService, asynqClient *asynq.Client) *ScanController {
	return &ScanController{
		ScanService: scanService,
		AsynqClient: asynqClient,
	}
}

// GET /api/jobs
func (sc *ScanController) GetAllScanJobs(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(400, gin.H{"error": "user not logged in"})
		return
	}
	jobs, err := sc.ScanService.GetScanJobsByUserID(userId.(uint))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve scan jobs"})
		return
	}
	response := store.ToAllScanJobsResponse(jobs)

	c.JSON(200, gin.H{"jobs": response})
}

// POST /api/jobs
func (sc *ScanController) StartScan(c *gin.Context) {
	// hardcode scan job for now
	req := store.CreateScanJobRequest{
		UserID:    1,
		TargetUrl: "https://example.com",
		Status:    "running",
	}
	err := sc.ScanService.CreateScanJob(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create scan job"})
		return
	}
}

// GET /api/jobs/:id
func (sc *ScanController) GetScanJob(c *gin.Context) {

}

// DELETE /api/jobs/:id
func (sc *ScanController) DeleteScanJob(c *gin.Context) {

}

// GET /api/jobs/:id/results
func (sc *ScanController) GetScanResults(c *gin.Context) {

}
