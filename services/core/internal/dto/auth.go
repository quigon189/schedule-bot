package dto

type LoginRequest struct {
	Username  string
	Password  string
	UserAgent string
	ClientIP  string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}
