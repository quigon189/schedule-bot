package handlers

import (
	"core/internal/dto"
	"core/internal/models"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	utils.SuccessResponse(w, "user created", user)
}

func (h *UserHandler) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	var req dto.PagiantedUserRequest
	req.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	req.PerPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))
	req.SortBy = r.URL.Query().Get("sort_by")
	req.SortOrder = r.URL.Query().Get("sort_order")

	users, err := h.userService.GetPaginatedUsers(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get users: %v", err))
		return
	}

	utils.SuccessResponse(w, "paginated users", users)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req dto.UserFilter
	if fullName := r.URL.Query().Get("full_name"); fullName != "" {
		req.FullName = &fullName
	}
	if username := r.URL.Query().Get("username"); username != "" {
		req.Username = &username
	}

	users, err := h.userService.GetAllUsers(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get users: %v", err))
		return
	}

	utils.SuccessResponse(w, "ok", users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid path value: %v", err))
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get user wtih id %d: %v", id, err))
		return
	}

	utils.SuccessResponse(w, "get user success", user)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to get current user")
		return
	}

	user, err := h.userService.GetUser(r.Context(), currentUser.ID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get user: %v", err))
		return
	}

	utils.SuccessResponse(w, "get user success", user)
}

func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	currentUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to get current user")
		return
	}

	isAdmin := currentUser.RequireRole("admin")

	if !isAdmin && currentUser.ID != req.UserID {
		utils.ErrorResponse(w, http.StatusForbidden, "access denied")
		return
	}

	userID := currentUser.ID
	if isAdmin {
		if req.UserID < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "invalid user_id")
			return
		}

		if err := h.userService.UpdatePasswordAdmin(r.Context(), req.UserID, req.NewPassword); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get user: %v", err))
			return
		}
	} else {
		user, err := h.userService.GetUser(r.Context(), userID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get user: %v", err))
			return
		}

		if err := h.userService.UpdatePassword(r.Context(), user, req.NewPassword, req.OldPassword); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update password: %v", err))
			return
		}
	}

	utils.SuccessResponse(w, "password changed", nil)
}
