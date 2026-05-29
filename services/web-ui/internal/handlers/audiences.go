package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/views/components"
	"web-ui/views/pages"

	"github.com/go-chi/chi/v5"
)

type AudiencesHandler struct {
	coreClient *api.CoreClient
}

func NewAudiencesHandler(client *api.CoreClient) *AudiencesHandler {
	return &AudiencesHandler{coreClient: client}
}

// GET /admin/audiences – страница со списком аудиторий
func (h *AudiencesHandler) ListAudiencesPage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, _ := r.Context().Value("user").(*models.User)

	// Получаем фильтры из query
	filterNumber := r.URL.Query().Get("number")
	filterName := r.URL.Query().Get("name")

	audiences, err := h.coreClient.GetAudiences(r.Context(), session)
	if err != nil {
		RenderInternalError(w, r, err)
		return
	}

	// Фильтрация на стороне сервера (API не поддерживает фильтры)
	filtered := filterAudiences(audiences, filterNumber, filterName)
	sortAudiences(filtered)

	pages.AdminAudiencesPage(user, filtered, filterNumber, filterName).Render(r.Context(), w)
}

// GET /admin/audiences/table – HTMX-фрагмент таблицы
func (h *AudiencesHandler) TableFragment(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	filterNumber := r.URL.Query().Get("number")
	filterName := r.URL.Query().Get("name")

	audiences, err := h.coreClient.GetAudiences(r.Context(), session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filtered := filterAudiences(audiences, filterNumber, filterName)
	sortAudiences(filtered)

	// Обновляем URL в адресной строке (без перезагрузки страницы)
	w.Header().Set("HX-Push-Url", "/admin/audiences?"+r.URL.RawQuery)
	components.AudiencesTable(filtered, filterNumber, filterName).Render(r.Context(), w)
}

// GET /admin/audiences/new – форма создания аудитории (модальное окно)
func (h *AudiencesHandler) NewAudienceForm(w http.ResponseWriter, r *http.Request) {
	components.AudienceForm(nil).Render(r.Context(), w)
}

// GET /admin/audiences/{id}/edit – форма редактирования
func (h *AudiencesHandler) EditAudienceForm(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	audience, err := h.coreClient.GetAudience(r.Context(), session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	components.AudienceForm(audience).Render(r.Context(), w)
}

// POST /admin/audiences – создание
func (h *AudiencesHandler) CreateAudience(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	name := r.FormValue("name")
	number := r.FormValue("number")

	req := api.CreateAudienceRequest{Name: name, Number: number}
	_, err := h.coreClient.CreateAudience(r.Context(), session, req)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка создания: "+err.Error())
		w.Header().Set("HX-Redirect", "/admin/audiences")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alerts.AddSuccess(w, session.SessionID, "Аудитория успешно создана")
	w.Header().Set("HX-Redirect", "/admin/audiences")
	w.WriteHeader(http.StatusOK)
}

// PUT /admin/audiences/{id} – обновление
func (h *AudiencesHandler) UpdateAudience(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	number := r.FormValue("number")
	req := api.UpdateAudienceRequest{Name: &name, Number: &number}
	_, err = h.coreClient.UpdateAudience(r.Context(), session, id, req)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка обновления: "+err.Error())
		w.Header().Set("HX-Redirect", "/admin/audiences")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alerts.AddSuccess(w, session.SessionID, "Аудитория обновлена")
	w.Header().Set("HX-Redirect", "/admin/audiences")
	w.WriteHeader(http.StatusOK)
}

// DELETE /admin/audiences/{id} – удаление
func (h *AudiencesHandler) DeleteAudience(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.coreClient.DeleteAudience(r.Context(), session, id)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка удаления: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	alerts.AddSuccess(w, session.SessionID, "Аудитория удалена")
	w.WriteHeader(http.StatusOK)
}

// GET /admin/audiences/upload-form
func (h *AudiencesHandler) UploadForm(w http.ResponseWriter, r *http.Request) {
	components.AudienceUploadForm().Render(r.Context(), w)
}

// GET /admin/audiences/template – скачать шаблон Excel
func (h *AudiencesHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	data, err := h.coreClient.DownloadAudienceTemplate(r.Context(), session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=audiences_template.xlsx")
	w.Write(data)
}

// POST /admin/audiences/upload – загрузка Excel
func (h *AudiencesHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	file, header, err := r.FormFile("file")
	if err != nil {
		alerts.AddError(w, session.SessionID, "Файл не загружен")
		w.Header().Set("HX-Redirect", "/admin/audiences")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем файл в байты
	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка чтения файла")
		w.Header().Set("HX-Redirect", "/admin/audiences")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.coreClient.UploadAudiencesExcel(r.Context(), session, fileData, header.Filename)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка загрузки: "+err.Error())
	} else {
		alerts.AddSuccess(w, session.SessionID, "Аудитории успешно загружены")
	}
	w.Header().Set("HX-Redirect", "/admin/audiences")
	w.WriteHeader(http.StatusOK)
}

// Вспомогательная фильтрация (локальная)
func filterAudiences(audiences []models.Audience, number, name string) []models.Audience {
	if number == "" && name == "" {
		return audiences
	}
	filtered := make([]models.Audience, 0, len(audiences))
	for _, a := range audiences {
		matchNumber := number == "" || strings.Contains(strings.ToLower(a.Number), strings.ToLower(number))
		matchName := name == "" || strings.Contains(strings.ToLower(a.Name), strings.ToLower(name))
		if matchNumber && matchName {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func sortAudiences(audiences []models.Audience) {
	sort.Slice(audiences, func(i, j int) bool {
		return audiences[i].Number < audiences[j].Number
	})
}
