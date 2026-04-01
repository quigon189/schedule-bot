package dto

import "time"

// CreateLessonLogRequest - создание записи в журнале
type CreateLessonLogRequest struct {
	Date             time.Time `json:"date" validate:"required"`
	Number           int       `json:"number" validate:"min=0,max=8"`
	Status           string    `json:"status" validate:"required,oneof=planned conducted cancelled replaced"`
	Comment          string    `json:"comment"`
	SubjectID        int       `json:"subject_id" validate:"required"`
	TeacherID        int       `json:"teacher_id" validate:"required"`
	AudienceID       int       `json:"audience_id" validate:"required"`
	AcademicPeriodID int       `json:"academic_period_id" validate:"required"`
	IsFromTemplate   bool      `json:"is_from_template"`
}

// UpdateLessonLogRequest - обновление записи в журнале
type UpdateLessonLogRequest struct {
	Date       *time.Time `json:"date"`
	Number     *int       `json:"number"`
	Status     *string    `json:"status"`
	Comment    *string    `json:"comment"`
	SubjectID  *int       `json:"subject_id"`
	TeacherID  *int       `json:"teacher_id"`
	AudienceID *int       `json:"audience_id"`
}

// CancelLessonRequest - отмена занятия
type CancelLessonRequest struct {
	Comment string `json:"comment" validate:"required"`
}

// ReplaceScheduleRequest - замена расписания
type ReplaceScheduleRequest struct {
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	Number     int    `json:"number" validate:"min=0,max=8"`
	SubjectID  int    `json:"subject_id" validate:"required"`
	TeacherID  int    `json:"teacher_id" validate:"required"`
	AudienceID int    `json:"audience_id" validate:"required"`
	Comment    string `json:"comment"`
}

// LessonLogFiltersRequest - фильтры для получения журнала
type LessonLogFiltersRequest struct {
	GroupID          *int       `json:"group_id"`
	SubjectID        *int       `json:"subject_id"`
	TeacherID        *int       `json:"teacher_id"`
	AudienceID       *int       `json:"audience_id"`
	AcademicPeriodID *int       `json:"academic_period_id"`
	DateFrom         *time.Time `json:"date_from"`
	DateTo           *time.Time `json:"date_to"`
	Status           *string    `json:"status"`
	Number           *int       `json:"number"`
}

// GenerateScheduleRequest - генерация расписания из шаблона
type GenerateScheduleRequest struct {
	AcademicPeriodID int `json:"academic_period_id" validate:"required"`
}
