package dto

type LoginRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	SessionID    string `json:"session_id" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
