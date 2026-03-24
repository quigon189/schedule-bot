package dto

type CreateStudentRequest struct {
	CreateUserRequest
	GroupID  int    `json:"group_id" validate:"required"`
}
