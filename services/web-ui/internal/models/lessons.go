package models

import "time"

type LessonLog struct {
	ID               int            `json:"id"`
	Date             time.Time      `json:"date"`
	Number           int            `json:"number"`
	Status           string         `json:"status"`
	Comment          string         `json:"comment"`
	Subject          Subject        `json:"subject"`
	Teacher          Teacher        `json:"teacher"`
	Audience         Audience       `json:"audience"`
	AcademicPeriod   AcademicPeriod `json:"academic_period"`
}

type ScheduleTemplate struct {
	ID               int            `json:"id"`
	DayOfWeek        int            `json:"day_of_week"`
	Number           int            `json:"number"`
	WeekType         int            `json:"week_type"`
	Subject          Subject        `json:"subject"`
	Teacher          Teacher        `json:"teacher"`
	Audience         Audience       `json:"audience"`
	AcademicPeriod   AcademicPeriod `json:"academic_period"`
}
