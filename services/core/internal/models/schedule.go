package models

import "time"

type AcademicPeriod struct {
	ID        int       `json:"id"`
	Year      string    `json:"yaer"`
	Semester  int       `json:"semester"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Audience struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type Subject struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Semester  int       `json:"semester"`
	HoursLoad int       `json:"hours_load"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GroupID   int       `json:"group_id"`

	Group Group `json:"group"`
}

type ScheduleTemplate struct {
	ID               int `json:"id"`
	DayOfWeek        int `json:"day_of_week"`
	Number           int `json:"number"`
	WeekType         int `json:"week_type"`
	SubjectID        int `json:"subject_id"`
	TeacherID        int `json:"teacher_id"`
	AudienceID       int `json:"audience_id"`
	AcademicPeriodID int `json:"academic_period_id"`

	AcademicPeriod AcademicPeriod `json:"academic_period"`
	Subject        Subject        `json:"subject"`
	Teacher        Teacher        `json:"teacher"`
	Audience       Audience       `json:"audience"`
}

type LessonLog struct {
	ID               int       `json:"id"`
	Date             time.Time `json:"date"`
	Number           int       `json:"number"`
	Status           string    `json:"status"`
	Comment          string    `json:"comment"`
	SubjectID        int       `json:"subject_id"`
	TeacherID        int       `json:"teacher_id"`
	AudienceID       int       `json:"audience_id"`
	AcademicPeriodID int       `json:"academic_period_id"`
	IsFromTemplate   bool      `json:"id_form_template"`

	AcademicPeriod AcademicPeriod `json:"academic_period"`
	Subject        Subject        `json:"subject"`
	Teacher        Teacher        `json:"teacher"`
	Audience       Audience       `json:"audience"`
}

const (
	LessonStatusPlanned     = "planned"
	LessonStatusCompleted   = "completed"
	LessonStatusCanceled   = "canceled"
	LessonStatusRescheduled = "rescheduled"
)

type SubjectProgress struct {
	SubjectID          int     `json:"subject_id"`
	Title              string  `json:"title"`
	TotalHours         int     `json:"total_hours"`
	TotalLessons       int     `json:"total_lessons"`        // всего пар (1 пара = 2 часа)
	CompletedLessons   int     `json:"completed_lessons"`    // проведено пар
	RemainingLessons   int     `json:"remaining_lessons"`    // осталось пар
	PlannedLessons     int     `json:"planned_lessons"`      // запланировано
	CancelledLessons   int     `json:"cancelled_lessons"`    // отменено
	CompletionPercent  float64 `json:"completion_percent"`   // процент выполнения
	IsOnSchedule       bool    `json:"is_on_schedule"`       // идём ли по графику
	NeedMakeupLessons  int     `json:"need_makeup_lessons"`  // сколько нужно добавить
	Recommendations    string  `json:"recommendations"`      // рекомендации
}
