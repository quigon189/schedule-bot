package models

import "time"

type Role struct {
	ID          int
	Name        string
	Description string
}

type User struct {
	ID           int
	Email        string
	Name         string
	FullName     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool

	Roles []Role
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
