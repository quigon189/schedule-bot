package models

import "time"

type Group struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type StudyPeriod struct {
	ID           int
	HalfYear     int
	AcademicYear string
	StrartDate   time.Time
	EndDate      time.Time
	CreatedAt    time.Time
}

type GroupSchedule struct {
	StudyPeriod    StudyPeriod
	Group          Group
	Semester       int
	ScheduleImgURL string
	CreatedAt      time.Time
}

type GroupScheduleFilter struct {
	GroupName    string
	AcademicYear string
	HalfYear     int
	Semester     int
}

type Change struct {
	ID          int
	Date        string
	Description string
	ImgURL      string
	CreatedAt   time.Time
}
