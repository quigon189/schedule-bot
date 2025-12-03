package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"schedule-service/internal/dto"
	"schedule-service/internal/repository"
	"schedule-service/internal/service"
	"schedule-service/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ChangeHandler struct {
	service  *service.ChangeService
	validate *validator.Validate
}

func NewChangeHandler(db *repository.DB) *ChangeHandler {
	repo := repository.NewChangeRepo(db)
	return &ChangeHandler{
		service:  service.NewChangeService(repo),
		validate: dto.NewChangeValidator(),
	}
}

// AddChange добавляет изменение в расписании на заданную дату
//
//	@Summary		Добавить изменения в расписании
//	@Description	Добавляет изменения в расписании
//	@Tags			changes
//	@Accept			json
//	@Produce		json
//	@Param			change	body		dto.AddChangeRequest	true	"Запрос на добавление изменения расписания"
//	@Success		200		{object}	utils.Response{data=dto.ChangeResponse}
//	@Failure		400		{object}	utils.Response
//	@Failure		500		{object}	utils.Response
//	@Router			/changes [post]
func (h *ChangeHandler) AddChange(w http.ResponseWriter, r *http.Request) {
	req := dto.AddChangeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		log.Printf("Invalid request body: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	sc, err := h.service.AddChange(req)
	if err != nil {
		log.Printf("Internal error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "Changes added", dto.ToChangeResponse(sc))
}

// GetChange возвращает изменения в расписании на заданную дату
//
//	@Summary		Получить изменения в расписании
//	@Description	Возвращает изменения в расписании
//	@Description	Если параметры не заданы, показывает изменения на текущую дату
//	@Tags			changes
//	@Accept			json
//	@Produce		json
//	@Param			date		query		string	false	"Дата изменений ISO	Example: ГГГГ-ММ-ДД"
//	@Param			image_url	query		string	false	"url изображения изменений Example: http://example.com/img.jpg"
//	@Param			description	query		string	false	"Дополнительное описание"
//	@Success		200			{object}	utils.Response{data=[]dto.ChangeResponse}
//	@Failure		400			{object}	utils.Response
//	@Failure		500			{object}	utils.Response
//	@Router			/changes [get]
func (h *ChangeHandler) GetChange(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")

	if err := h.validate.Var(dateStr, "omitempty,date"); err != nil {
		log.Printf("Invalid date format: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	chs, err := h.service.GetChanges(dateStr)
	if err != nil {
		log.Printf("Invalid date format: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	resp := []dto.ChangeResponse{}
	for _, sc := range chs {
		resp = append(resp, *dto.ToChangeResponse(&sc))
	}

	utils.SuccessResponse(w, "Changes sent", resp)
}

// RemoveChange  удаляет изменение в расписании по id
//
//	@Summary		Удалить изменения в расписании
//	@Description	Удаляет изменения в расписании по id
//	@Tags			changes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id удаляемого изменения"	example=1
//	@Success		200	{object}	utils.Response
//	@Failure		400	{object}	utils.Response
//	@Failure		500	{object}	utils.Response
//	@Router			/changes/{id} [delete]
func (h *ChangeHandler) RemoveChange(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Failed to parse id: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid id")
		return
	}

	if err := h.service.RemoveChange(id); err != nil {
		log.Printf("Failed to remove change: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "Change removed", nil)
}
