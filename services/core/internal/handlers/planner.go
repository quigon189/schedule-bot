package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type PlannerHandler struct {
	plannderService *services.PlannerService
}

func NewPlannerHandler(svc *services.ScheduleService) *PlannerHandler {
	return &PlannerHandler{plannderService: services.NewPlannerService(svc)}
}

func (h *PlannerHandler) GenerateWeeklySchedule(w http.ResponseWriter, r *http.Request) {
	var req dto.PlanScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalide body: %v", err))
		return
	}

	resp, err := h.plannderService.PlanWeeklySchedule(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get weekly schedule: %v", err))
		return
	}

	utils.SuccessResponse(w, "weekly schedule generated", resp)
}
