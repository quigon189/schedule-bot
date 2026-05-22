package models

type JWT struct {
	AccessToken string   `json:"access_token"`
	Refresh     string   `json:"refresh_token"`
	ExpiresAt   int64    `json:"expires_at"`
	SessionID   string   `json:"session_id"`
	Roles       []string `json:"roles"`
}
