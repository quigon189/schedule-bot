package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID           pgtype.UUID `json:"uuid"`
	RefreshToken string      `json:"-"`
	UserAgent    string      `json:"user_agent"`
	ClientIP     string      `json:"client_ip"`
	User         *User       `json:"user"`
	CreatedAt    time.Time   `json:"created_at"`
}
