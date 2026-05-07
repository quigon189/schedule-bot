package models

import "slices"

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`
	Group    *Group `json:"group"`
}

type Group struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Specialty     string `json:"specialty"`
	AdmissionYear int    `json:"admission_year"`
}

func (u *User) HasRole(role string) bool {
	return slices.ContainsFunc(u.Roles, func(r Role) bool {
		return r.Name == role
	})
}
