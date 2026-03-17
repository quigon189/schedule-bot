package dto

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}
