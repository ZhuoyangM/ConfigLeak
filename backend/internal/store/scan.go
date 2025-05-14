package store

import (
	"context"

	"github.com/ZhuoyangM/ConfigLeak/internal/models"
	"gorm.io/gorm"
)

type ScanService struct {
	db  *gorm.DB
	ctx context.Context
}

func (service *ScanService) CreateScanJob(scanJob *models.ScanJob) error {
	return service.db.WithContext(service.ctx).Create(scanJob).Error
}

func (service *ScanService) GetScanJobByID(id uint) (*models.ScanJob, error) {
	var scanJob models.ScanJob
	err := service.db.WithContext(service.ctx).First(&scanJob, id).Error
	if err != nil {
		return nil, err
	}
	return &scanJob, nil
}

func (service *ScanService) GetScanJobsByUserID(userID uint) ([]models.ScanJob, error) {
	var scanJobs []models.ScanJob
	err := service.db.WithContext(service.ctx).Where("user_id = ?", userID).Find(&scanJobs).Error
	if err != nil {
		return nil, err
	}
	return scanJobs, nil
}

func (service *ScanService) CreateScanResult(scanResult *models.ScanResult) error {
	return service.db.WithContext(service.ctx).Create(scanResult).Error
}
func (service *ScanService) GetScanResultByID(id uint) (*models.ScanResult, error) {
	var scanResult models.ScanResult
	err := service.db.WithContext(service.ctx).First(&scanResult, id).Error
	if err != nil {
		return nil, err
	}
	return &scanResult, nil
}
func (service *ScanService) GetScanResultsByJobID(jobID uint) ([]models.ScanResult, error) {
	var scanResults []models.ScanResult
	err := service.db.WithContext(service.ctx).Where("scan_job_id = ?", jobID).Find(&scanResults).Error
	if err != nil {
		return nil, err
	}
	return scanResults, nil
}
