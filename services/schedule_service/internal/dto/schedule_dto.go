package dto

type AddGroupScheduleRequest struct {
	AcademicYear   string `json:"academic_year" example:"2025/2026"`
	HalfYear       int    `json:"half_year" example:"1"`
	GroupName      string `json:"group_name" example:"СА-501"`
	Semester       int    `json:"semester" example:"5"`
	ScheduleImgURL string `json:"schedule_img_url" example:"https://example.com/schedule.jpg"`
}

type UpdateGroupScheduleRequest struct {
	ScheduleImgURL string `json:"schedule_img_url" example:"https://example.com/schedule.jpg"`
}

type GroupScheduleQueryParams struct {
	AcademicYear string `json:"academic_year" example:"2025/2026"`
	HalfYear     int    `json:"half_year" example:"1"`
	GroupName    string `json:"group_name" example:"CA-501"`
}

type GroupScheduleResponse struct {
	GroupName      string `json:"group_name" example:"СА-501"`
	Semester       int    `json:"semester" example:"5"`
	ScheduleImgURL string `json:"schedule_img_url" example:"https://example.com/schedule.jpg"`
}
