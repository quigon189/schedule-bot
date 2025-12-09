package dto

import (
	"time"
)

type CreateRegistrationCodeRequest struct {
	RoleName  string    `json:"role_name" validate:"required"`
	GroupName *string   `json:"group_name,omitempty"`
	MaxUses   int       `json:"max_uses" validate:"required,min=1"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
	CreatedBy int64     `json:"created_by" validate:"required"`
}

type ValidateCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

type RegisterUserRequest struct {
	Code       string  `json:"code" validate:"required"`
	TelegramID int64   `json:"telegram_id" validate:"required"`
	Username   string  `json:"username"`
	FullName   string  `json:"full_name" validate:"required"`
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
