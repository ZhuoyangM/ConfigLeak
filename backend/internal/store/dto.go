package store

import (
	"html"

	"golang.org/x/crypto/bcrypt"
)

// user related DTOs
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GetUserResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// scan job related DTOs
type CreateScanJobRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	TargetUrl string `json:"target_url" binding:"required,url"`
	Status    string `json:"status" binding:"required,oneof=running completed failed"`
}

type GetScanJobResponse struct {
	JobID     uint   `json:"job_id"`
	UserID    uint   `json:"user_id"`
	TargetUrl string `json:"target_url"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"` // need conversion
}

type GetScanResultResponse struct {
	ResultID      uint   `json:"result_id"`
	ScanJobID     uint   `json:"scan_job_id"`
	ScanUrl       string `json:"scan_url"`
	Status        string `json:"status"`
	Code          *int   `json:"code,omitempty"`           // HTTP status code, nullable
	ContentLength *int   `json:"content_length,omitempty"` // in bytes, nullable
}

// user-related conversion functions
func ToUser(req *RegisterRequest) (*User, error) {
	var user User
	user.Username = html.EscapeString(req.Username)
	user.Email = html.EscapeString(req.Email)

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPwd)
	return &user, nil
}

func ToGetUserResponse(user *User) *GetUserResponse {
	return &GetUserResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

// scan job-related conversion functions
func ToAllScanJobsResponse(jobs []ScanJob) []GetScanJobResponse {
	var responses []GetScanJobResponse
	for _, job := range jobs {
		responses = append(responses, GetScanJobResponse{
			JobID:     job.ID,
			UserID:    job.UserID,
			TargetUrl: job.TargetUrl,
			Status:    job.Status,
			CreatedAt: job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return responses
}

func ToScanJob(req *CreateScanJobRequest) *ScanJob {
	return &ScanJob{
		UserID:    req.UserID,
		TargetUrl: html.EscapeString(req.TargetUrl),
		Status:    req.Status,
	}
}
