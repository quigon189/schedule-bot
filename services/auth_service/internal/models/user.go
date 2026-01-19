package models

import "time"

const (
	AdminRole   = "admin"
	ManagerRole = "manager"
	StudentRole = "student"
	TeacherRole = "teacher"
	UserRole    = "user"
)

type User struct {
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsActive   bool      `json:"is_active"`

	Roles   []Role   `json:"roles,omitempty"`
	Student *Student `json:"student,omitempty"`
}

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Student struct {
	UserID    int64  `json:"user_id"`
	GroupName string `json:"group_name"`
}
