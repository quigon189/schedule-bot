package models

import (
	"slices"
	"time"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"username"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Roles []Role `json:"roles"`
}

type Group struct {
	ID            int
	Name          string
	Specialty     string
	AdmissionYear int
}

type Student struct {
	User  User
	Group *Group
}

type Teacher struct {
	User User
}

type Subject struct {
	ID        int
	Name      string
	HoursLoad int
	Semester  int
	Group     Group
	Teacher   *Teacher
}

func (u *User) RequireRole(role string) bool {
	return slices.ContainsFunc(u.Roles, func(r Role) bool {
		return r.Name == role
	})
}
