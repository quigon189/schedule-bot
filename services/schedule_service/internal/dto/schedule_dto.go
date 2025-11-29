package dto

type AddGroupScheduleRequest struct {
	AcademicYear string `json:"academic_year" example:"2025/2026"`
	HalfYear int `json:"half_year" example:"1"`
	GroupName string `json:"group_name" example:"СА-501"`
	Semester int `json:"semester" example:"5"`
	ScheduleImgURL int `json:"schedule_ing_url" example:"https://example.com/schedule.jpg"`
}

type GroupScheduleResponse struct {

}
