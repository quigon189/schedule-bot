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

type GroupHandler struct {
	scheduleService *services.ScheduleService
}

func NewGroupHandler(scheduleService *services.ScheduleService) *GroupHandler {
	return &GroupHandler{scheduleService: scheduleService}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	group, err := h.scheduleService.CreateGroup(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create group: %v", err))
		return
	}
	utils.SuccessResponse(w, "group created", group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	group, err := h.scheduleService.GetGroup(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get group: %v", err))
		return
	}
	if group == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "group not found")
		return
	}
	utils.SuccessResponse(w, "group retrieved", group)
}

func (h *GroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.scheduleService.GetAllGroups(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get groups: %v", err))
		return
	}
	utils.SuccessResponse(w, "groups retrieved", groups)
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	var req dto.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	group, err := h.scheduleService.UpdateGroup(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update group: %v", err))
		return
	}

	utils.SuccessResponse(w, "group updated", group)
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	if err := h.scheduleService.DeleteGroup(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete group: %v", err))
		return
	}
	utils.SuccessResponse(w, "group deleted", nil)
}

func (h *GroupHandler) CreateGroupWtihCurriculum(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGroupWithCurriculumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	resp, err := h.scheduleService.CreateGroupWithCurriculumAndStudents(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create group with curriculum: %v", err))
		return
	}

	utils.SuccessResponse(w, "group with curriculum and students created", resp)
}
