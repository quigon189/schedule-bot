package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID           pgtype.UUID
	RefreshToken string
	UserAgent    string
	ClientIP     string
	User         *User
	CreatedAt    time.Time
}
