package store

import (
	"context"

	"gorm.io/gorm"
)

type ScanService struct {
	db  *gorm.DB
	ctx context.Context
}

// TODO: configure the context
func NewScanService(db *gorm.DB) *ScanService {
	return &ScanService{
		db:  db,
		ctx: context.Background(),
	}
}

func (service *ScanService) CreateScanJob(req *CreateScanJobRequest) error {
	scanJob := ToScanJob(req)
	return service.db.WithContext(service.ctx).Create(scanJob).Error
}

func (service *ScanService) GetScanJobByID(id uint) (*ScanJob, error) {
	var scanJob ScanJob
	err := service.db.WithContext(service.ctx).First(&scanJob, id).Error
	if err != nil {
		return nil, err
	}
	return &scanJob, nil
}

func (service *ScanService) GetScanJobsByUserID(userID uint) ([]ScanJob, error) {
	var scanJobs []ScanJob
	err := service.db.WithContext(service.ctx).Where("user_id = ?", userID).Find(&scanJobs).Error
	if err != nil {
		return nil, err
	}
	return scanJobs, nil
}

func (service *ScanService) CreateScanResult(scanResult *ScanResult) error {
	return service.db.WithContext(service.ctx).Create(scanResult).Error
}
func (service *ScanService) GetScanResultByID(id uint) (*ScanResult, error) {
	var scanResult ScanResult
	err := service.db.WithContext(service.ctx).First(&scanResult, id).Error
	if err != nil {
		return nil, err
	}
	return &scanResult, nil
}
func (service *ScanService) GetScanResultsByJobID(jobID uint) ([]ScanResult, error) {
	var scanResults []ScanResult
	err := service.db.WithContext(service.ctx).Where("scan_job_id = ?", jobID).Find(&scanResults).Error
	if err != nil {
		return nil, err
	}
	return scanResults, nil
}
