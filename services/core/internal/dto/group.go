package dto

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
