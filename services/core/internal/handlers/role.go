package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"net/http"
)

type RoleHandler struct {
	scheduleService *services.ScheduleService
}

func NewRoleHandler(scheduleService *services.ScheduleService) *RoleHandler {
	return &RoleHandler{scheduleService: scheduleService}
}

func (h *RoleHandler) AssignStudent(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.scheduleService.AssignStudent(r.Context(), req.UserID, req.GroupID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(w, "student assigned", nil)
}

func (h *RoleHandler) AssignTeacher(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.scheduleService.AssignTeacher(r.Context(), req.UserID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(w, "teacher assigned", nil)
}

func (h *RoleHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.scheduleService.AssignRole(r.Context(), req.UserID, req.RoleName); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(w, "role assigned", nil)
}

func (h *RoleHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	var req dto.RemoveRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	switch req.RoleName {
	case "student":
		if err := h.scheduleService.RemoveStudent(r.Context(), req.UserID); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	case "teacher":
		if err := h.scheduleService.RemoveTeacher(r.Context(), req.UserID); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	default:
		if err := h.scheduleService.RemoveRole(r.Context(), req.UserID, req.RoleName); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	utils.SuccessResponse(w, "role removed", nil)
}
