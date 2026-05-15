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

type SubjectHandler struct {
	scheduleService *services.ScheduleService
}

func NewSubjectHandler(scheduleService *services.ScheduleService) *SubjectHandler {
	return &SubjectHandler{scheduleService: scheduleService}
}

func (h *SubjectHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}
	
	subject, err := h.scheduleService.CreateSubject(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create subject: %v", err))
		return
	}
	utils.SuccessResponse(w, "subject created", subject)
}

func (h *SubjectHandler) GetSubject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	subject, err := h.scheduleService.GetSubjectByID(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get subject: %v", err))
		return
	}
	if subject == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "subject not found")
		return
	}
	utils.SuccessResponse(w, "subject retrieved", subject)
}

func (h *SubjectHandler) GetAllSubjects(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	req := dto.PaginatedSubjectsRequest{
		Page:      page,
		PerPage:   perPage,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	resp, err := h.scheduleService.GetAllSubjects(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get subjects: %v", err))
		return
	}
	utils.SuccessResponse(w, "subjects retrieved", resp)
}

func (h *SubjectHandler) GetSubjectsByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid group_id")
		return
	}
	
	resp, err := h.scheduleService.GetSubjectsByGroupID(r.Context(), groupID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get subjects: %v", err))
		return
	}
	utils.SuccessResponse(w, "subjects retrieved", resp)
}

func (h *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	subject, err := h.scheduleService.UpdateSubject(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update subject: %v", err))
		return
	}
	utils.SuccessResponse(w, "subject updated", subject)
}

func (h *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.scheduleService.DeleteSubject(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete subject: %v", err))
		return
	}
	utils.SuccessResponse(w, "subject deleted", nil)
}
