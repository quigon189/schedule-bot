package dto

type CreateStudentRequest struct {
	CreateUserRequest
	GroupID int `json:"group_id" validate:"required"`
}

type CreateScheduleTemplateRequest struct {
	DayOfWeek        int `json:"day_of_week" validate:"required,min=1,max=7"`
	Number           int `json:"number" validate:"min=0,max=8"`
	WeekType         int `json:"week_type" validate:"min=0,max=2"`
	SubjectID        int `json:"subject_id" validate:"required"`
	TeacherID        int `json:"teacher_id" validate:"required"`
	AudienceID       int `json:"audience_id" validate:"required"`
	AcademicPeriodID int `json:"academic_period_id" validate:"required"`
}

type UpdateScheduleTemplateRequest struct {
	DayOfWeek        *int `json:"day_of_week"`
	Number           *int `json:"Number"`
	WeekType         *int `json:"week_type"`
	SubjectID        *int `json:"subject_id"`
	TeacherID        *int `json:"teacher_id"`
	AudienceID       *int `json:"audience_id"`
	AcademicPeriodID *int `json:"academic_period_id"`
}

type ScheduleFiltersRequest struct {
	GroupID          *int `json:"group_id"`
	SubjectID        *int `json:"subject_id"`
	TeacherID        *int `json:"teacher_id"`
	AudienceID       *int `json:"audience_id"`
	DayOfWeek        *int `json:"day_of_week"`
	WeekType         *int `json:"week_type"`
	AcademicPeriodID *int `json:"academic_period_id"`
}

type ScheduleEntryRequest struct {
	DayOfWeek  int `json:"day_of_week" validate:"min=1,max=7"`
	Number     int `json:"number" validate:"min=0,max=8"`
	WeekType   int `json:"week_type" validate:"min=0,max=2"`
	SubjectID  int `json:"subject_id" validate:"required"`
	TeacherID  int `json:"teacher_id" validate:"required"`
	AudienceID int `json:"audience_id" validate:"required"`
}

type CreateSemesterScheduleRequest struct {
	GroupID          int                    `json:"group_id" validate:"required"`
	AcademicPeriodID int                    `json:"academic_period_id" validate:"required"`
	ScheduleEntries  []ScheduleEntryRequest `json:"schedule_entries" validate:"required,min=1"`
}

type BulkCreateScheduleRequest struct {
	Templates []CreateScheduleTemplateRequest `json:"templates" validate:"required,min=1,dive"`
}

type BulkDeleteScheduleRequest struct {
	IDs []int `json:"ids" validate:"required,min=1"`
}
