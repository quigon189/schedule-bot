package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to decode body: %v", err))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	utils.SuccessResponse(w, "user created", user)
}
