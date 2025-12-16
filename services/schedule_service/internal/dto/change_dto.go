package dto

import (
	"schedule-service/internal/models"
	"time"
)

type AddChangeRequest struct {
	Date        string   `json:"date" example:"2025-09-01" validate:"required,date"`
	ImgURLs     []string `json:"image_urls" example:"http://example.com/change.jpg" validate:"required"`
	Description string   `json:"description"`
}

type ChangeResponse struct {
	ID          int       `json:"id" example:"1"`
	Date        string    `json:"date" example:"2025-09-01"`
	ImgURLs     []string  `json:"image_urls" example:"[http://example.com/change.jpg]"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToChangeResponse(sc *models.ScheduleChange) *ChangeResponse {
	return &ChangeResponse{
		ID:          sc.ID,
		Date:        sc.Date.Format("2006-01-02"),
		ImgURLs:     sc.ImgURLs,
		Description: sc.Description,
		CreatedAt:   sc.CreatedAt,
	}
}
