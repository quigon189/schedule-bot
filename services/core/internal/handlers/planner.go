package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/internal/storage"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type PlannerHandler struct {
	plannderService *services.PlannerService
	storage         *storage.MemoryStorage
}

func NewPlannerHandler(svc *services.ScheduleService, stg *storage.MemoryStorage) *PlannerHandler {
	return &PlannerHandler{
		plannderService: services.NewPlannerService(svc),
		storage:         stg,
	}
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

func (h *PlannerHandler) UploadPlannerExcel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	schedule, err := h.plannderService.GenerateScheduleFromExcel(r.Context(), file)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate schedule from excel failed: %v", err))
		return
	}

	exlService := services.NewExcelService()

	resFile, err := exlService.GenerateGroupSchedule(schedule)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate schedule file failed: %v", err))
		return
	}

	filename := fmt.Sprintf("schedule_%d.xlsx", schedule.Seed)
	mime := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	id, exires := h.storage.PutFile(filename, mime, resFile)

	result := dto.GenerationResult{
		Seed: schedule.Seed,
		FilePath: fmt.Sprintf("/files/%s", id),
		ExiresAt: exires,
	}

	utils.SuccessResponse(w, "schedule generated from Excel", result)
}
