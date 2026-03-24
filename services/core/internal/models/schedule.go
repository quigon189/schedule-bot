package models

import "time"

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
	ID         int `json:"id"`
	DayOfWeek  int `json:"day_of_week"`
	Number     int `json:"number"`
	WeekType   int `json:"week_type"`
	SubjectID  int `json:"subject_id"`
	TeacherID  int `json:"teacher_id"`
	AudienceID int `json:"audience_id"`

	Subject  Subject  `json:"subject"`
	Teacher  Teacher  `json:"teacher"`
	Audience Audience `json:"audience"`
}

type LessonLog struct {
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Comment    string    `json:"comment"`
	SubjectID  int       `json:"subject_id"`
	TeacherID  int       `json:"teacher_id"`
	AudienceID int       `json:"audience_id"`

	Subject  Subject  `json:"subject"`
	Teacher  Teacher  `json:"teacher"`
	Audience Audience `json:"audience"`
}
