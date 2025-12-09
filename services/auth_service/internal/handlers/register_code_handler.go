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

// CreateCode создает код доступа для новых пользователей
//
//	@Summary		Создать код доступа
//	@Description	Создает новый код доступа для пользователя
//	@Description	Создать код могут только администраторы или менеджеры
//	@Description	Обязательные поля: created_by, role_name	
//	@Description	Поле  role_name имеет значения: student (group_name становится обязательным), teacher, manager
//	@Tags			code
//	@Accept			json
//	@Produce		json
//	@Param			code	body		dto.CreateRegistrationCodeRequest	true	"Параметры кода доступа"
//	@Success		200		{object}	utils.Response{data=dto.RegistrationCodeResponse}
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/code/create [post]
func (h *RegistrationCodeHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRegistrationCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Failed to validate request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
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

// RegisterWithCode Зарегистрировать нового или существующего пользователя по коду доступа
//
//	@Summary		Зарегистрировать пользователя
//	@Description	Регистрирует нового или уже существующего пользователя используя код доступа
//	@Description	Добавляет пользователю роль связанную с кодом
//	@Description	Если при создании кода указана роль student, то также добавляет пользователя в группу
//	@Tags			code, users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterUserRequest	true	"Параметр пользователя и код"
//	@Success		200		{object}	utils.Response{data=dto.UserResponse}
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/users/register [post]
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
