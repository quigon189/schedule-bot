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

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	GroupID  *int   `json:"group_id"`
}

type Group struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Specialty     string `json:"specialty"`
	AdmissionYear int    `json:"admission_year"`
	Students      []User `json:"students"`
}

type GroupWithCurriculumResponse struct {
	Group    Group                   `json:"group"`
	Subjects []Subject               `json:"subjects"`
	Students []StudentCreationResult `json:"students"`
}

type StudentCreationResult struct {
	User     User   `json:"user"`
	Password string `json:"password"`
	GroupID  int    `json:"group_id"`
}

type TeacherCreationResult struct {
	User     User   `json:"user"`
	Password string `json:"password"`
}

func (u *User) HasRole(role string) bool {
	return slices.ContainsFunc(u.Roles, func(r Role) bool {
		return r.Name == role
	})
}
