package handlers

import (
	"auth_service/internal/dto"
	"auth_service/internal/service"
	"auth_service/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RegistrationCodeHandler struct {
	codeService *service.RegistrationCodeService
}

func NewRegistrationCodeHandler(codeService *service.RegistrationCodeService) *RegistrationCodeHandler {
	return &RegistrationCodeHandler{codeService: codeService}
}

func (h *RegistrationCodeHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRegistrationCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	code, err := h.codeService.CreateCode(&req)
	if err != nil {
		log.Printf("Failed to create code: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create code: %v", err))
		return
	}

	resp := dto.RegistrationCodeResponse{
		ID: code.ID,
		Code: code.Code,
		Role: dto.RoleResponse{
			ID: code.Role.ID,
			Name: code.Role.Name,
			Description: code.Role.Description,
		},
		GroupName: *code.GroupName,
		MaxUses: code.MaxUses,
		Creater: dto.UserResponse{
			TelegramID: code.Creater.TelegramID,
			Username: code.Creater.Username,
			FullName: code.Creater.FullName,
		},
		ExpiresAt: code.ExpiresAt,
		CreatedAt: code.CreatedAt,
	}
	utils.SuccessResponse(w, "registration code created", resp)
}

func (h *RegistrationCodeHandler) RegisterWithCode(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.codeService.RegisterWithCode(&req)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to register user: %v", err))
		return
	}

	utils.SuccessResponse(w, "user registred", dto.ToUserResponse(user))
}
