package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
	exportService   *services.ScheduleExportService
}

func NewScheduleHandler(scheduleService *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
		exportService:   services.NewScheduleExportService(scheduleService),
	}
}

// CreateScheduleTemplate создает новую запись в расписании
func (h *ScheduleHandler) CreateScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateScheduleTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validation error: %v", err))
		return
	}

	template, err := h.scheduleService.CreateScheduleTemplate(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create schedule template: %v", err))
		return
	}

	utils.SuccessResponse(w, "schedule template created", template)
}

// GetScheduleTemplate получает запись расписания по ID
func (h *ScheduleHandler) GetScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	template, err := h.scheduleService.GetScheduleTemplate(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get schedule template: %v", err))
		return
	}
	if template == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "schedule template not found")
		return
	}

	utils.SuccessResponse(w, "schedule template retrieved", template)
}

// GetAllScheduleTemplates получает список записей расписания с фильтрацией
func (h *ScheduleHandler) GetAllScheduleTemplates(w http.ResponseWriter, r *http.Request) {
	var filters dto.ScheduleFiltersRequest

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
	if dayOfWeek := r.URL.Query().Get("day_of_week"); dayOfWeek != "" {
		day, err := strconv.Atoi(dayOfWeek)
		if err == nil {
			filters.DayOfWeek = &day
		}
	}
	if weekType := r.URL.Query().Get("week_type"); weekType != "" {
		wt, err := strconv.Atoi(weekType)
		if err == nil {
			filters.WeekType = &wt
		}
	}
	if periodID := r.URL.Query().Get("academic_period_id"); periodID != "" {
		id, err := strconv.Atoi(periodID)
		if err == nil {
			filters.AcademicPeriodID = &id
		}
	}

	templates, err := h.scheduleService.GetAllScheduleTemplates(r.Context(), &filters)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get schedule templates: %v", err))
		return
	}

	utils.SuccessResponse(w, "schedule templates retrieved", templates)
}

// UpdateScheduleTemplate обновляет запись расписания
func (h *ScheduleHandler) UpdateScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	var req dto.UpdateScheduleTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	template, err := h.scheduleService.UpdateScheduleTemplate(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update schedule template: %v", err))
		return
	}

	utils.SuccessResponse(w, "schedule template updated", template)
}

// DeleteScheduleTemplate удаляет запись расписания
func (h *ScheduleHandler) DeleteScheduleTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	if err := h.scheduleService.DeleteScheduleTemplate(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete schedule template: %v", err))
		return
	}

	utils.SuccessResponse(w, "schedule template deleted", nil)
}

// GetGroupSchedule получает расписание группы за учебный период
func (h *ScheduleHandler) GetGroupSchedule(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid group id: %v", err))
		return
	}

	var periodID *int
	if periodIDstr := r.URL.Query().Get("period_id"); periodIDstr != "" {
		periodIDint, err := strconv.Atoi(periodIDstr)
		if err == nil {
			periodID = &periodIDint
		}
	}

	schedule, err := h.scheduleService.GetGroupSchedule(r.Context(), groupID, periodID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get group schedule: %v", err))
		return
	}

	utils.SuccessResponse(w, "group schedule retrieved", schedule)
}

// GetTeacherSchedule получает расписание преподавателя за учебный период
func (h *ScheduleHandler) GetTeacherSchedule(w http.ResponseWriter, r *http.Request) {
	teacherID, err := strconv.Atoi(chi.URLParam(r, "teacher_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid teacher id: %v", err))
		return
	}

	var periodID *int
	if periodIDstr := r.URL.Query().Get("period_id"); periodIDstr != "" {
		periodIDint, err := strconv.Atoi(periodIDstr)
		if err == nil {
			periodID = &periodIDint
		}
	}

	schedule, err := h.scheduleService.GetTeacherSchedule(r.Context(), teacherID, periodID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get teacher schedule: %v", err))
		return
	}

	utils.SuccessResponse(w, "teacher schedule retrieved", schedule)
}

// GetAudienceSchedule получает расписание аудитории за учебный период
func (h *ScheduleHandler) GetAudienceSchedule(w http.ResponseWriter, r *http.Request) {
	audienceID, err := strconv.Atoi(chi.URLParam(r, "audience_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid audience id: %v", err))
		return
	}

	var periodID *int
	if periodIDstr := r.URL.Query().Get("period_id"); periodIDstr != "" {
		periodIDint, err := strconv.Atoi(periodIDstr)
		if err == nil {
			periodID = &periodIDint
		}
	}

	schedule, err := h.scheduleService.GetAudienceSchedule(r.Context(), audienceID, periodID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get audience schedule: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience schedule retrieved", schedule)
}

func (h *ScheduleHandler) CreateSemesterSchedule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSemesterScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	err := h.scheduleService.CreateSemesterSchedule(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	utils.SuccessResponse(w, "semester schedule created successfully", nil)
}

func (h *ScheduleHandler) ExportGroupSchedule(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid group id")
	}

	periodIDStr := r.URL.Query().Get("period_id")
	var periodID *int
	if periodIDStr != "" {
		pid, err := strconv.Atoi(periodIDStr)
		if err == nil {
			periodID = &pid
		}
	}

	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		formatStr = "html"
	}

	format := services.ExportFormat(formatStr)
	if format != services.FormatHTML && format != services.FormatPDF && format != services.FormatPNG && format != services.FormatDOCX {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid format, must be html, docx, pdf or png")
		return
	}

	content, mime, err := h.exportService.ExportGroupSchedule(r.Context(), groupID, periodID, format)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export schedule: %v", err))
		return
	}

	filename := fmt.Sprintf("schedule_group_%d.%s", groupID, format)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(content)
}

func (h *ScheduleHandler) ExportTeacherSchedule(w http.ResponseWriter, r *http.Request) {
	teacherID, err := strconv.Atoi(chi.URLParam(r, "teacher_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid teacher id")
	}

	periodIDStr := r.URL.Query().Get("period_id")
	var periodID *int
	if periodIDStr != "" {
		pid, err := strconv.Atoi(periodIDStr)
		if err == nil {
			periodID = &pid
		}
	}

	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		formatStr = "html"
	}

	format := services.ExportFormat(formatStr)
	if format != services.FormatHTML && format != services.FormatPDF && format != services.FormatPNG && format != services.FormatDOCX {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid format, must be html, docx, pdf or png")
		return
	}

	content, mime, err := h.exportService.ExportTeacherSchedule(r.Context(), teacherID, periodID, format)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export schedule: %v", err))
		return
	}

	filename := fmt.Sprintf("schedule_teahcer_%d.%s", teacherID, format)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(content)
}

func (h *ScheduleHandler) ExportAudienceSchedule(w http.ResponseWriter, r *http.Request) {
	audienceID, err := strconv.Atoi(chi.URLParam(r, "audience_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid audience id")
	}

	periodIDStr := r.URL.Query().Get("period_id")
	var periodID *int
	if periodIDStr != "" {
		pid, err := strconv.Atoi(periodIDStr)
		if err == nil {
			periodID = &pid
		}
	}

	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		formatStr = "html"
	}

	format := services.ExportFormat(formatStr)
	if format != services.FormatHTML && format != services.FormatPDF && format != services.FormatPNG && format != services.FormatDOCX {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid format, must be html, docx, pdf or png")
		return
	}

	content, mime, err := h.exportService.ExportAudienceSchedule(r.Context(), audienceID, periodID, format)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export schedule: %v", err))
		return
	}

	filename := fmt.Sprintf("schedule_audience_%d.%s", audienceID, format)
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(content)
}

func (h *ScheduleHandler) DownloadPlannerTemplate(w http.ResponseWriter, r *http.Request) {
	periodID, err := strconv.Atoi(chi.URLParam(r, "period_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid period_id")
		return
	}
	period, err := h.scheduleService.GetAcademicPeriodByID(r.Context(), periodID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to ge academic period: %v", err))
		return
	}
    excelSvc := services.NewExcelService()
    data, err := excelSvc.GeneratePlannerTemplate(r.Context(), h.scheduleService, period)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate template: %v", err))
        return
    }
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=planner_template.xlsx")
    w.Write(data)
}
