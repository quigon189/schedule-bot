package dto

import (
	"core/internal/models"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,username"`
	FullName string `json:"full_name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

type PagiantedUserRequest struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type UpdateUserPasswordRequest struct {
	UserID      int    `json:"user_id"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
	OldPassword string `json:"old_password"`
}

type UserFilter struct {
	FullName *string
	Username *string
}
