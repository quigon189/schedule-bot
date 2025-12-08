package models

import "time"

type RegistrationCode struct {
	ID          int
	Code        string
	RoleID      int
	GroupName   *string
	MaxUses     int
	CurrentUses int
	CreatedBy   *int
	ExpiresAt   time.Time
	CreatedAt   time.Time

	Role    Role
	Creater *User
}
