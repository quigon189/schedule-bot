package dto

import "core/internal/models"

type CreateGroupRequest struct {
	Name          string `json:"name" validate:"required"`
	Specialty     string `json:"specialty" validate:"required"`
	AdmissionYear int    `json:"admission_year" validate:"required"`
}

type UpdateGroupRequest struct {
	Name          string `json:"name"`
	Specialty     string `json:"specialty"`
	AdmissionYear int    `json:"admission_year"`
}

type CreateGroupSubjectRequest struct {
	Title     string `json:"title" validate:"required"`
	Semester  int    `json:"semester" validate:"required,min=1,max=10"`
	HoursLoad int    `json:"hours_load" validate:"required,min=1"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

// CreateGroupWithCurriculumRequest - запрос на создание группы с дисциплинами и студентами
type CreateGroupWithCurriculumRequest struct {
	Group    CreateGroupRequest                    `json:"group" validate:"required"`
	Subjects []CreateGroupSubjectRequest           `json:"subjects" validate:"required,min=1,dive"`
	Students []CreateStudentWithCredentialsRequest `json:"students" validate:"dive"`
}

// CreateStudentWithCredentialsRequest - данные студента без указания пароля (генерируется)
type CreateStudentWithCredentialsRequest struct {
	FullName string `json:"full_name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
}

// GroupWithCurriculumResponse - ответ после создания
type GroupWithCurriculumResponse struct {
	Group    models.Group            `json:"group"`
	Subjects []models.Subject        `json:"subjects"`
	Students []StudentCreationResult `json:"students"`
}

type StudentCreationResult struct {
	User     models.User `json:"user"`
	Password string      `json:"password"`
	GroupID  int         `json:"group_id"`
}
