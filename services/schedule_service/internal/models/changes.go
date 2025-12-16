package models

import "time"

type ScheduleChange struct {
	ID          int
	Date        time.Time
	ImgURLs      []string
	Description string
	CreatedAt   time.Time
}
