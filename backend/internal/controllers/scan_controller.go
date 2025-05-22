package controllers

import (
	"github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/gin-gonic/gin"
)

type ScanController struct {
	ScanService store.ScanService
}

func NewScanController(scanService store.ScanService) *ScanController {
	return &ScanController{
		ScanService: scanService,
	}
}

// GET /api/jobs
func (sc *ScanController) GetAllScanJobs(c *gin.Context) {

}

// POST /api/jobs
func (sc *ScanController) StartScan(c *gin.Context) {

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
