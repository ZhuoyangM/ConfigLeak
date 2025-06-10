package controllers

import (
	"strconv"

	"github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type ScanController struct {
	ScanService *store.ScanService
	AsynqClient *asynq.Client
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
		UserID:    2,
		TargetUrl: "https://example2.com",
		Status:    "running",
	}
	err := sc.ScanService.CreateScanJob(&req)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create scan job"})
		return
	}
	c.JSON(201, gin.H{"message": "scan job created successfully"})
}

// GET /api/jobs/:id
func (sc *ScanController) GetScanJob(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(400, gin.H{"error": "user not logged in"})
		return
	}
	jobId := c.Param("id")
	intJobId, err := strconv.ParseUint(jobId, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid job ID"})
		return
	}
	job, err := sc.ScanService.GetScanJobByID(uint(intJobId))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve scan job"})
		return
	}
	if job.UserID != userId.(uint) {
		c.JSON(403, gin.H{"error": "forbidden: you do not have access to this job"})
		return
	}
	response := store.ToGetScanJobResponse(job)
	c.JSON(200, gin.H{"job": response})

}

// DELETE /api/jobs/:id
func (sc *ScanController) DeleteScanJob(c *gin.Context) {

}

// GET /api/jobs/:id/results
func (sc *ScanController) GetScanResults(c *gin.Context) {

}
