package dto

import "time"

type AddChangeRequest struct {
	Date        string `json:"date" example:"2025-09-01"`
	ImgURL      string `json:"image_url" example:"http://example.com/change.jpg"`
	Description string `json:"description"`
}

type ChangeResponse struct {
	ID          int       `json:"id" example:"1"`
	Date        time.Time `json:"date" example:"2025-09-01"`
	ImgURL      string    `json:"image_url" example:"http://example.com/change.jpg"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
