package controllers

import (
	"github.com/ZhuoyangM/ConfigLeak/internal/models"
	"github.com/ZhuoyangM/ConfigLeak/internal/store"
)

type ScanController struct {
	ScanService store.ScanService
}

func NewScanController(scanService store.ScanService) *ScanController {
	return &ScanController{
		ScanService: scanService,
	}
}

// GET /scan
func (c *ScanController) GetAllScanJobs(scanJobID int) (*models.ScanJob, error) {
	return nil, nil
}

// POST /scan
func (c *ScanController) StartScan(scanJobID int) error {
	return nil
}

// GET /scan/:id
func (c *ScanController) GetScanJob(scanJobID int) (*models.ScanJob, error) {
	return nil, nil
}

// DELETE /scan/:id
func (c *ScanController) DeleteScanJob(scanJobID int) error {
	return nil
}

// GET /scan/:id/result
func (c *ScanController) GetScanResults(scanJobID int) ([]models.ScanResult, error) {
	return nil, nil
}
