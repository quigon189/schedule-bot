package handlers

import (
	"core/internal/models"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
}

func NewScheduleHandler(scheduleService *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

func (h *ScheduleHandler) CreateAudience(w http.ResponseWriter, r *http.Request) {
	var audience models.Audience
	if err := json.NewDecoder(r.Body).Decode(&audience); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	if err := h.scheduleService.CreateAudience(r.Context(), &audience); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience created", audience)
}

func (h *ScheduleHandler) GetAudience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	audience, err := h.scheduleService.GetAudience(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience getted", audience)
}

func (h *ScheduleHandler) GetAllAudience(w http.ResponseWriter, r *http.Request) {
	audiences, err := h.scheduleService.GetAllAudience(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get audiences: %v", err))
		return
	}

	utils.SuccessResponse(w, "audiences getted", audiences)
}

func (h *ScheduleHandler) UpdateAudience(w http.ResponseWriter, r *http.Request) {
	var audience models.Audience
	if err := json.NewDecoder(r.Body).Decode(&audience); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	if err := h.scheduleService.UpdateAudience(r.Context(), &audience); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience updated", audience)
}

func (h *ScheduleHandler) DeleteAudience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	if err := h.scheduleService.DeleteAudience(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete audience: %v", err))
		return
	}	

	utils.SuccessResponse(w, "audience deleted", nil)
}
