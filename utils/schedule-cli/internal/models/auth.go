package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	SessionID    string   `json:"session_id"`
	ExiresAt     int64    `json:"exires_at"`
	Roles        []string `json:"roles"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}
