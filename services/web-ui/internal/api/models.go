package api

type JWT struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExiresAt     int64    `json:"expires_at"`
	SessionID    string   `json:"session_id"`
	Roles        []string `json:"roles"`
}
