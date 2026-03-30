package dto

type CreateAudienceRequest struct {
	Name   string `json:"name" validate:"required"`
	Number string `json:"number" validate:"required"`
}

type UpdateAudienceRequest struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}
