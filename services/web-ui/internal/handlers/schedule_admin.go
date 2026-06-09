package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/views/components"
	"web-ui/views/pages"
)

type ScheduleAdminHandler struct {
	coreClient *api.CoreClient
}

func NewScheduleAdminHandler(client *api.CoreClient) *ScheduleAdminHandler {
	return &ScheduleAdminHandler{coreClient: client}
}

// GET /admin/schedule — главная страница
func (h *ScheduleAdminHandler) SchedulePage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, _ := r.Context().Value("user").(*models.User)

	// Получаем списки для выбора
	groups, _ := h.coreClient.GetGroups(r.Context(), session)
	audiences, _ := h.coreClient.GetAudiences(r.Context(), session)
	periods, _ := h.coreClient.GetAcademicPeriods(r.Context(), session)

	pages.AdminSchedulePage(user, groups, audiences, periods).Render(r.Context(), w)
}

// GET /admin/schedule/view — фрагмент таблицы расписания (HTMX)
func (h *ScheduleAdminHandler) ViewSchedule(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	entityType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")
	periodIDStr := r.URL.Query().Get("period_id")

	if entityType == "" || idStr == "" {
		http.Error(w, "Не указан тип или ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var periodID *int
	if periodIDStr != "" {
		pid, err := strconv.Atoi(periodIDStr)
		if err == nil && pid > 0 {
			periodID = &pid
		}
	}

	var templates []models.ScheduleTemplate
	switch entityType {
	case "group":
		templates, err = h.coreClient.GetGroupSchedule(r.Context(), session, id, periodID)
	case "teacher":
		templates, err = h.coreClient.GetTeacherSchedule(r.Context(), session, id, periodID)
	case "audience":
		templates, err = h.coreClient.GetAudienceSchedule(r.Context(), session, id, periodID)
	default:
		http.Error(w, "Неизвестный тип", http.StatusBadRequest)
		return
	}
	if err != nil {
		RenderInternalError(w, r, err)
		return
	}

	// Преобразуем в ScheduleData для компонента
	data := convertTemplatesToScheduleData(templates, periodID)

	props := components.ScheduleViewProps{
		Data:         data,
		ShowGroup:    entityType == "teacher" || entityType == "audience",
		ShowTeacher:  entityType == "group" || entityType == "audience",
		ShowAudience: entityType == "group" || entityType == "teacher",
	}
	components.ScheduleView(props).Render(r.Context(), w)
}

// GET /admin/schedule/template — скачать шаблон для планировщика (Excel)
func (h *ScheduleAdminHandler) DownloadPlannerTemplate(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	periodIDStr := r.URL.Query().Get("period_id")
	if periodIDStr == "" {
		http.Error(w, "Не указан период", http.StatusBadRequest)
		return
	}
	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil {
		http.Error(w, "Некорректный период", http.StatusBadRequest)
		return
	}

	data, err := h.coreClient.DownloadPlannerTemplate(r.Context(), session, periodID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=schedule_planner_period_%d.xlsx", periodID))
	w.Write(data)
}

// POST /admin/schedule/upload — загрузка Excel для генерации расписания
func (h *ScheduleAdminHandler) UploadScheduleExcel(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	file, header, err := r.FormFile("file")
	if err != nil {
		alerts.AddError(w, session.SessionID, "Файл не загружен")
		w.Header().Set("HX-Redirect", "/admin/schedule")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData := make([]byte, header.Size)
	if _, err := file.Read(fileData); err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка чтения файла")
		w.Header().Set("HX-Redirect", "/admin/schedule")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Загружаем и получаем результат генерации
	result, err := h.coreClient.UploadPlannerExcel(r.Context(), session, fileData, header.Filename)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка генерации: "+err.Error())
		w.Header().Set("HX-Redirect", "/admin/schedule")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отображаем результат: кнопку скачивания
	components.ScheduleGenerationResult(result).Render(r.Context(), w)
	w.Header().Set("HX-Trigger", "closeModal")
}

// GET /admin/schedule/download-result — скачать сгенерированный файл по пути
func (h *ScheduleAdminHandler) DownloadGeneratedFile(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Путь не указан", http.StatusBadRequest)
		return
	}
	// Выполняем запрос к core на получение файла (doRawWithAuth)
	data, contentType, err := h.coreClient.DoRawGetFile(r.Context(), session, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=generated_schedule.xlsx")
	w.Write(data)
}

// GET /admin/schedule/export — экспорт существующего расписания
func (h *ScheduleAdminHandler) ExportSchedule(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	entityType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")
	periodIDStr := r.URL.Query().Get("period_id")
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "xlsx"
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}
	var periodID *int
	if periodIDStr != "" {
		pid, err := strconv.Atoi(periodIDStr)
		if err == nil && pid > 0 {
			periodID = &pid
		}
	}

	var data []byte
	var contentType string
	switch entityType {
	case "group":
		data, contentType, err = h.coreClient.ExportGroupSchedule(r.Context(), session, id, periodID, format)
	case "teacher":
		data, contentType, err = h.coreClient.ExportTeacherSchedule(r.Context(), session, id, periodID, format)
	case "audience":
		data, contentType, err = h.coreClient.ExportAudienceSchedule(r.Context(), session, id, periodID, format)
	default:
		http.Error(w, "Неизвестный тип", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("schedule_%s_%d.%s", entityType, id, format)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

// POST /admin/schedule/save — сохранение сгенерированного расписания (опционально)
func (h *ScheduleAdminHandler) SaveGeneratedSchedule(w http.ResponseWriter, r *http.Request) {
	// Этот метод можно реализовать позже, если core поддерживает сохранение.
	// Сейчас просто заглушка.
	http.Error(w, "Метод не реализован", http.StatusNotImplemented)
}

func (h *ScheduleAdminHandler) UploadForm(w http.ResponseWriter, r *http.Request) {
	components.ScheduleUploadForm().Render(r.Context(), w)
}

// Вспомогательная функция: преобразование []ScheduleTemplate в ScheduleData
func convertTemplatesToScheduleData(templates []models.ScheduleTemplate, periodID *int) models.ScheduleData {
	entries := make([]models.ScheduleEntry, 0, len(templates))
	var period models.AcademicPeriod
	var groupName, teacherName, audienceNumber string

	for _, tmpl := range templates {
		entry := models.ScheduleEntry{
			DayOfWeek: tmpl.DayOfWeek,
			Number:    tmpl.Number,
			WeekType:  tmpl.WeekType,
			Subject:   tmpl.Subject,
			Teacher:   tmpl.Teacher,
			Audience:  tmpl.Audience,
		}
		entries = append(entries, entry)
		if period.ID == 0 && tmpl.AcademicPeriod.ID != 0 {
			period = tmpl.AcademicPeriod
		}
		if groupName == "" && tmpl.Subject.Group.Name != "" {
			groupName = tmpl.Subject.Group.Name
		}
		if teacherName == "" && tmpl.Teacher.FullName != "" {
			teacherName = tmpl.Teacher.FullName
		}
		if audienceNumber == "" && tmpl.Audience.Number != "" {
			audienceNumber = tmpl.Audience.Number
		}
	}
	return models.ScheduleData{
		Period:         period,
		GroupName:      groupName,
		TeacherName:    teacherName,
		AudienceNumber: audienceNumber,
		Entries:        entries,
	}
}
