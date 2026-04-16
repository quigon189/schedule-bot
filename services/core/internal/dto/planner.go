package dto

import "core/internal/models"

type PlanScheduleRequest struct {
	AcademicPeriodID int              `json:"academic_period_id" validate:"required"`
	Seed             int              `json:"seed" validate:"min=0"`
	SubjectList      []SubjectRequest `json:"subject_list" validate:"required,min=1,dive"`
}

type SubjectRequest struct {
	SubjectID    int `json:"subject_id" validate:"required"`
	TeacherID    int `json:"teacher_id" validate:"required"`
	AudienceID   int `json:"audience_id" validate:"required"`
	Priority     int `json:"priority" validate:"min=0,max=100"`
	LessonsCount int `json:"lessons_count" validate:"required"`
}

type ScheduleCell struct {
	Subject  models.Subject  `json:"subject"`
	Teacher  models.Teacher  `json:"teacher"`
	Audience models.Audience `json:"audience"`
}

type WeeklySchedule struct {
	Seed int                            `json:"seed"`
	Grid map[int]map[int][]ScheduleCell `json:"grid"`
}
