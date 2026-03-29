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

type TeacherHandler struct {
	scheduleService *services.ScheduleService
}

func NewTeacherHandler(scheduleService *services.ScheduleService) *TeacherHandler {
	return &TeacherHandler{scheduleService: scheduleService}
}

func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
	}

	teacher, err := h.scheduleService.CreateTeacher(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create teacher: %v", err))
		return
	}
	utils.SuccessResponse(w, "teacher created", teacher)
}

func (h *TeacherHandler) GetTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	teacher, err := h.scheduleService.GetTeacher(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get teacher: %v", err))
		return
	}
	if teacher == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "teacher not found")
		return
	}
	utils.SuccessResponse(w, "teacher retrieved", teacher)
}

func (h *TeacherHandler) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.scheduleService.GetAllTeachers(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get teachers: %v", err))
		return
	}
	utils.SuccessResponse(w, "teachers retrieved", teachers)
}

func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	if err := h.scheduleService.DeleteTeacher(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete teacher: %v", err))
		return
	}
	utils.SuccessResponse(w, "teacher deleted", nil)
}
