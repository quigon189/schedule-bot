package dto

import (
	"time"
)

type AddGroupScheduleRequest struct {
	AcademicYear   string `json:"academic_year" example:"2025/2026" validate:"required,academic_year"`
	HalfYear       int    `json:"half_year" example:"1" validate:"required,min=1,max=2"`
	GroupName      string `json:"group_name" example:"СА-501" validate:"required,group_name"`
	Semester       int    `json:"semester" example:"5" validate:"required,min=1,max=10"`
	ScheduleImgURL string `json:"schedule_img_url" example:"https://example.com/schedule.jpg" validate:"required,url"`
}

type UpdateGroupScheduleRequest struct {
	ScheduleImgURL string `json:"schedule_img_url" example:"https://example.com/schedule.jpg" validate:"required,url"`
}

type GroupScheduleQueryParams struct {
	AcademicYear string `schema:"academic_year" example:"2025/2026" validate:"omitempty,academic_year"`
	HalfYear     int    `schema:"half_year" example:"1" validate:"omitempty,min=1,max=2"`
	GroupName    string `schema:"group_name" example:"CA-501" validate:"omitempty,group_name"`
}

type GroupScheduleResponse struct {
	GroupName      string    `json:"group_name" example:"СА-501"`
	AcademicYear   string    `json:"academic_year" example:"2025/2026"`
	HalfYear       int       `json:"half_year" example:"1"`
	Semester       int       `json:"semester" example:"5"`
	ScheduleImgURL string    `json:"schedule_img_url" example:"https://example.com/schedule.jpg"`
	CreatedAt      time.Time `json:"created_at"`
}
