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
