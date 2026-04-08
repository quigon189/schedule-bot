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

type CompleteScheduleRequest struct {
	Date    string `json:"date" validate:"required,datetime=2006-01-02"`
	Comment string `json:"comment"`
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

// StatisticsRequest - фильтры для получения статистики занятий
type StatisticsRequest struct {
	GroupID          *int       `json:"group_id"`
	TeacherID        *int       `json:"teacher_id"`
	SubjectID        *int       `json:"subject_id"`
	AcademicPeriodID *int       `json:"academic_period_id"`
}

// SubjectStatistics - статистика по одному предмету
type SubjectStatistics struct {
	SubjectID          int     `json:"subject_id"`
	SubjectTitle       string  `json:"subject_title"`
	GroupID            int     `json:"group_id"`
	GroupName          string  `json:"group_name"`
	TotalHours         int     `json:"total_hours"`          // общее количество часов
	TotalLessons       int     `json:"total_lessons"`        // всего пар (1 пара = 2 часа)
	CompletedLessons   int     `json:"completed_lessons"`    // проведено пар
	RescheduledLessons int     `json:"rescheduled_lessons"`  // перенесено пар
	CancelledLessons   int     `json:"cancelled_lessons"`    // отменено пар
	PlannedLessons     int     `json:"planned_lessons"`      // запланировано на будущее
	RemainingLessons   int     `json:"remaining_lessons"`    // осталось провести (всего - проведено - отменено)
	CompletionPercent  float64 `json:"completion_percent"`   // процент выполнения (проведено / всего)
	IsOnSchedule       bool    `json:"is_on_schedule"`       // идём ли по графику
}

// StatisticsResponse - ответ со статистикой
type StatisticsResponse struct {
	TotalSubjects       int                 `json:"total_subjects"`
	OverallProgress     float64             `json:"overall_progress"`       // общий процент выполнения
	TotalLessonsAll     int                 `json:"total_lessons_all"`      // всего пар по всем предметам
	CompletedLessonsAll int                 `json:"completed_lessons_all"`  // проведено всего
	RemainingLessonsAll int                 `json:"remaining_lessons_all"`  // осталось провести всего
	Subjects            []SubjectStatistics `json:"subjects"`
}
