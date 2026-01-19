package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"auth_service/internal/dto"
	"auth_service/internal/service"
	"auth_service/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser создает нового пользователя
//
//	@Summary		Создать пользователя
//	@Description	Создает нового пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateUserRequest	true	"Данные пользователя"
//	@Success		200		{object}	utils.Response{data=dto.UserResponse}
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("failed to validate request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "user created", dto.ToUserResponse(user))
}

// GetUser получить информацию о пользователе
//
//	@Summary		Получить пользователя
//	@Description	Возвращает информацию о пользователе по telegram id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			telegram_id	path		int64	true	"TelegramID пользователя"
//	@Success		200			{object}	utils.Response{data=dto.UserResponse}
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/users/{telegram_id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "user not found")
		return
	}

	utils.SuccessResponse(w, "user sent", dto.ToUserResponse(user))
}

// UpdateUser обновить данные пользователя
// @Summary Обновить пользователя
// @Description Обновить данные пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param UpdateUserRequest body dto.UpdateUserRequest true "Данные пользователя"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/users/{telegram_id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := dto.RequestToUser(id, &req)
	if err = h.userService.UpdateUser(user); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err = h.userService.GetUser(user.TelegramID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "user updated", dto.ToUserResponse(user))
}

func (h *UserHandler) UpdateUserRoles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateUserRolesRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err = h.userService.UpdateUserRoles(id, req.Roles); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "user roles updated", dto.ToUserResponse(user))
}

// DeleteUser удаленип пользователя
// @Summary Удалить пользователя
// @Description Удалить пользователя из базы данных
// @Tags users
// @Accept json
// @Produce json
// @Param telegram_id path int64 true "TelegramID пользователя"
// @Success 200 {object} utils.Response
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/users/{telegram_id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err = h.userService.DeleteUser(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "user deleted", nil)
}
