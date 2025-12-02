package dto

import (
	"schedule-service/internal/models"
	"time"
)

type AddChangeRequest struct {
	Date        string `json:"date" example:"2025-09-01" validate:"required,date"`
	ImgURL      string `json:"image_url" example:"http://example.com/change.jpg" validate:"required,url"`
	Description string `json:"description"`
}

type ChangeResponse struct {
	ID          int       `json:"id" example:"1"`
	Date        string    `json:"date" example:"2025-09-01"`
	ImgURL      string    `json:"image_url" example:"http://example.com/change.jpg"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToChangeResponse(sc *models.ScheduleChange) *ChangeResponse {
	return &ChangeResponse{
		ID:          sc.ID,
		Date:        sc.Date.Format("2006-01-02"),
		ImgURL:      sc.ImgURL,
		Description: sc.Description,
		CreatedAt:   sc.CreatedAt,
	}
}
