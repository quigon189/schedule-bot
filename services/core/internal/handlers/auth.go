package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to decode body: %v", err))
		return
	}

	resp, err := h.userService.Login(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to login: %v", err))
		return
	}

	utils.SuccessResponse(w, "success login", resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	accessToken, ok := r.Context().Value("access_token").(string)
	if !ok {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get access token"))
		return
	}
	err := h.userService.Logout(r.Context(), accessToken)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to logout: %v", err))
		return
	}

	utils.SuccessResponse(w, "success logout", nil)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to decode body: %v", err))
		return
	}

	resp, err := h.userService.RefreshToken(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to refresh token: %v", err))
		return
	}

	utils.SuccessResponse(w, "success refresh token", resp)
}
