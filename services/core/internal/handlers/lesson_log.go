package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type LessonLogHandler struct {
	scheduleService *services.ScheduleService
}

func NewLessonLogHandler(scheduleService *services.ScheduleService) *LessonLogHandler {
	return &LessonLogHandler{scheduleService: scheduleService}
}

// GetLessonLogs получает записи журнала с фильтрацией
func (h *LessonLogHandler) GetLessonLogs(w http.ResponseWriter, r *http.Request) {
	var filters dto.LessonLogFiltersRequest
	
	// Парсим query параметры
	if groupID := r.URL.Query().Get("group_id"); groupID != "" {
		id, err := strconv.Atoi(groupID)
		if err == nil {
			filters.GroupID = &id
		}
	}
	if subjectID := r.URL.Query().Get("subject_id"); subjectID != "" {
		id, err := strconv.Atoi(subjectID)
		if err == nil {
			filters.SubjectID = &id
		}
	}
	if teacherID := r.URL.Query().Get("teacher_id"); teacherID != "" {
		id, err := strconv.Atoi(teacherID)
		if err == nil {
			filters.TeacherID = &id
		}
	}
	if audienceID := r.URL.Query().Get("audience_id"); audienceID != "" {
		id, err := strconv.Atoi(audienceID)
		if err == nil {
			filters.AudienceID = &id
		}
	}
	if periodID := r.URL.Query().Get("academic_period_id"); periodID != "" {
		id, err := strconv.Atoi(periodID)
		if err == nil {
			filters.AcademicPeriodID = &id
		}
	}
	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &t
		}
	}
	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &t
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters.Status = &status
	}
	if number := r.URL.Query().Get("number"); number != "" {
		n, err := strconv.Atoi(number)
		if err == nil {
			filters.Number = &n
		}
	}

	logs, err := h.scheduleService.GetLessonLogs(r.Context(), &filters)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get lesson logs: %v", err))
		return
	}

	utils.SuccessResponse(w, "lesson logs retrieved", logs)
}

func (h *LessonLogHandler) CancelLesson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	var req dto.CancelLessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := h.scheduleService.CancelLesson(r.Context(), id, &req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to cancel lesson: %v", err))
		return 
	}

	utils.SuccessResponse(w, fmt.Sprintf("lesson %d canceled", id), nil)
}

func (h *LessonLogHandler) RescheduleLesson(w http.ResponseWriter, r *http.Request) {
	var req dto.ReplaceScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := h.scheduleService.RescheduleLesson(r.Context(), &req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to reschedule lesson: %v", err))
		return 
	}

	utils.SuccessResponse(w, "lesson rescheduled", nil)

}

func (h *LessonLogHandler) CompleteLessonFromDate(w http.ResponseWriter, r *http.Request) {
	var req dto.CompleteScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad date format: %v", err))
		return
	}

	if err := h.scheduleService.CompleteLessonFromDate(r.Context(), date, req.Comment); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to reschedule lesson: %v", err))
		return 
	}

	utils.SuccessResponse(w, "lesson status changed to completed", nil)

}
