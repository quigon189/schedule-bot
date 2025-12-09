package dto

import (
	"time"
)

type CreateRegistrationCodeRequest struct {
	RoleName  string    `json:"role_name" enums:"student,teacher,manager" example:"student"`
	GroupName *string   `json:"group_name,omitempty" example:"СА-501"`
	MaxUses   int       `json:"max_uses" maximum:"50" default:"1" example:"10"`
	ExpiresAt time.Time `json:"expires_at" default:"current time + 6 hours" example:"2006-01-02T03:04:05Z"`
	CreatedBy int64     `json:"created_by" example:"12345678"`
}

type ValidateCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

type RegisterUserRequest struct {
	Code       string  `json:"code" example:"ABC123"`
	TelegramID int64   `json:"telegram_id" example:"12345678"`
	Username   string  `json:"username"`
	FullName   string  `json:"full_name"`
}

type ValidateCodeResponse struct {
	Valid         bool         `json:"valid"`
	Role          RoleResponse `json:"role"`
	GroupName     *string      `json:"group_name,omitempty"`
	RemainingUses int          `json:"remaining_uses,omitempty"`
}

type RegistrationCodeResponse struct {
	ID        int          `json:"id"`
	Code      string       `json:"code"`
	Role      RoleResponse `json:"role"`
	GroupName string       `json:"group_name,omitempty"`
	MaxUses   int          `json:"max_uses"`
	Creater   UserResponse `json:"created_by"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
}
