package dto

import "core/internal/models"

type CreateSubjectRequest struct {
	Title string `json:"title" validate:"required"`
	Semester int `json:"semester" validate:"required"`
	HoursLoad int `json:"hours_load" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate string `json:"end_date" validate:"required"`
	GroupID int `json:"group_id" validate:"required"`
}

type PaginatedSubjectsRequest struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type PaginatedSubjectsResponse struct {
	Subjects   []models.Subject `json:"subjects"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}
