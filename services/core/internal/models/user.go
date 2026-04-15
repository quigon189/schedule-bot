package models

import (
	"fmt"
	"slices"
	"strings"
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
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Specialty     string `json:"specialty"`
	AdmissionYear int    `json:"admission_year"`
	Students      []User `json:"students,omitempty"`
}

type Student struct {
	User  User   `json:"user"`
	Group *Group `json:"group"`
}

type Teacher struct {
	User User `json:"user"`
}

func (u *User) GetShortName() string {
	fields := strings.Fields(u.FullName)
	if len(fields) == 1 {
		return fields[0]
	}

	if len(fields) == 2 {
		return fmt.Sprintf("%s %s.", fields[0], string([]rune(fields[1])[0]))
	}

	if len(fields) == 3 {
		return fmt.Sprintf("%s %s.%s.", fields[0], string([]rune(fields[1])[0]), string([]rune(fields[2])[0]))
	}

	if len(fields) > 3 {
		return fmt.Sprintf("%s %s.", fields[0], string([]rune(fields[1])[0]))
	}

	return ""
}

func (u *User) RequireRole(role string) bool {
	return slices.ContainsFunc(u.Roles, func(r Role) bool {
		return r.Name == role
	})
}
