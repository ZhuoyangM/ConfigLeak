package models

import (
	"time"

	"github.com/google/uuid"
)

type ScanJob struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TargetUrl string     `json:"target_url"`
	CreatedAt *time.Time `json:"created_at"`
	Status    string     `json:"status"`
}
