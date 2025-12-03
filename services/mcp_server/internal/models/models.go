package models

import "time"

type ChangesRequest struct {
	Date time.Time `json:"date" jsonschema:"description=Дата изменений в формате ISO 8601, если значение пустое, то возвращает изменения на текущую дату"`
}

type ChangesResponse struct {
	ID          int       `json:"id"`
	Date        string    `json:"date"`
	ImgURL      string    `json:"image_url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    []ChangesResponse `json:"data"`
	Error   string          `json:"error,omitempty"`
}
