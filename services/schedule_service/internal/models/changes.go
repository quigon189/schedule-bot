package models

import "time"

type ScheduleChange struct {
	ID          int
	Date        time.Time
	ImgURL      string
	Description string
	CreatedAt   time.Time
}
