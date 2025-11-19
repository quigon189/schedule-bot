package dto

import (
	"auth_service/internal/models"
	"time"
)

type CreateUserRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
}

type UpdateUserRequest struct {
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
}

type UpdateUserRolesRequest struct {
	Roles []string `json:"roles"`
}

type RoleResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StudentResponse struct {
	GroupName string `json:"group_name"`
}

type UserResponse struct {
	TelegramID int64            `json:"telegram_id"`
	Username   string           `json:"username"`
	FullName   string           `json:"full_name"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	IsActive   bool             `json:"id_active"`
	Roles      []RoleResponse   `json:"roles"`
	Student    *StudentResponse `json:"student,omitempty"`
}

func RequestToUser(id int64, req *UpdateUserRequest) *models.User {
	return &models.User{
		TelegramID: id,
		Username:   req.Username,
		FullName:   req.FullName,
	}
}

func ToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	var roles []RoleResponse
	for _, role := range user.Roles {
		roles = append(roles, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	var student *StudentResponse
	if user.Student != nil {
		student = &StudentResponse{
			GroupName: user.Student.GroupName,
		}
	}

	return &UserResponse{
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FullName:   user.FullName,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		IsActive:   user.IsActive,
		Roles:      roles,
		Student:    student,
	}
}
