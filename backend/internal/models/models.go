package models

import (
	"time"
)

// pointer indicates that the field is nullable
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type ScanJob struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TargetUrl string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Status    string    `gorm:"not null"` // "running", "completed", "failed"
}

type ScanResult struct {
	ID            uint    `gorm:"primaryKey"`
	ScanJobID     uint    `gorm:"not null"`
	ScanJob       ScanJob `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScanUrl       string  `gorm:"not null"`
	Status        string  `gorm:"not null"` // "unscanned", "scanned", "timeout"
	Code          *int    // HTTP status code
	ContentLength *int    // in bytes
}
