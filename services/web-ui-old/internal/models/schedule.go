package models

import "time"

type Student struct {
	User  User  `json:"user"`
	Group Group `json:"group"`
}

type AcademicPeriod struct {
	ID        int       `json:"id"`
	Year      string    `json:"year"`
	Semester  int       `json:"semester"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Audience struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type Subject struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"` // добавим для удобства
}

type Teacher struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

type ScheduleEntry struct {
	DayOfWeek int      `json:"day_of_week"`
	Number    int      `json:"number"`
	WeekType  int      `json:"week_type"` // 0 - обе, 1 - числитель, 2 - знаменатель
	Subject   Subject  `json:"subject"`
	Teacher   Teacher  `json:"teacher"`
	Audience  Audience `json:"audience"`
}

type ScheduleData struct {
	Period         AcademicPeriod  `json:"period"`
	GroupName      string          `json:"group_name,omitempty"`
	TeacherName    string          `json:"teacher_name,omitempty"`
	AudienceNumber string          `json:"audience_number,omitempty"`
	Entries        []ScheduleEntry `json:"entries"`
}
