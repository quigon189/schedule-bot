package api

import (
	"time"
)

type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_toke"`
	SessionID    string    `json:"session_id"`
	Roles        []string  `json:"roles"`
	ExiresAt     time.Time `json:"expires_at"`
	Updated      bool      `json:"-"`
}

func (s *Session) Set(accessToken, refreshToken, sessionID string, exiresAt int64) {
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.SessionID = sessionID
	s.ExiresAt = time.Unix(exiresAt, 0)
}

func (s *Session) Exire(d time.Duration) bool {
	return time.Now().Add(d).After(s.ExiresAt)
}

func (s *Session) Update(accessToken, refreshToken string, exiresAt int64) {
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.ExiresAt = time.Unix(exiresAt, 0)
	s.Updated = true
}
