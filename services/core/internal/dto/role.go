package dto

type AssignStudentRequest struct {
	UserID  int `json:"user_id" validate:"required"`
	GroupID int `json:"group_id" validate:"required"`
}

type AssignTeacherRequest struct {
	UserID int `json:"user_id" validate:"required"`
}

type AssignRoleRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required,oneof=admin manager"`
}

type RemoveRoleRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	RoleName string `json:"role_name" validate:"required,oneof=admin manager student teacher"`
}
