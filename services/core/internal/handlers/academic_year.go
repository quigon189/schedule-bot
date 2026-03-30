package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AcademicPeriodHandler struct {
	scheduleService *services.ScheduleService
}

func NewAcademicPeriodHandler(scheduleService *services.ScheduleService) *AcademicPeriodHandler {
	return &AcademicPeriodHandler{scheduleService: scheduleService}
}

func (h *AcademicPeriodHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAcademicPeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	academicPeriod, err := h.scheduleService.CreateAcademicPeriod(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("create academic period: %v", err))
		return
	}

	utils.SuccessResponse(w, "academic period created", academicPeriod)
}

func (h *AcademicPeriodHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	academicPeriod, err := h.scheduleService.GetAcademicPeriodByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("get academic period: %v", err))
		return
	}

	utils.SuccessResponse(w, "academic period", academicPeriod)
}

func (h *AcademicPeriodHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	academicPeriods, err := h.scheduleService.GetAllAcademicPeriods(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("get all academic period: %v", err))
		return
	}

	utils.SuccessResponse(w, "academic periods", academicPeriods)

}

func (h *AcademicPeriodHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	var req dto.UpdateAcademicPeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	academicPeriod, err := h.scheduleService.UpdateAcademicPeriod(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("update academic period: %v", err))
		return
	}

	utils.SuccessResponse(w, "academic period updated", academicPeriod)
}

func (h *AcademicPeriodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	err = h.scheduleService.DeleteAcademicPeriod(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("delete academic period: %v", err))
		return
	}

	utils.SuccessResponse(w, "academic period deleted", nil)

}
