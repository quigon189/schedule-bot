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

type SubjectsHandler struct {
	coreClient *api.CoreClient
}

func NewSubjectsHandler(client *api.CoreClient) *SubjectsHandler {
	return &SubjectsHandler{coreClient: client}
}

// GET /admin/subjects — страница со списком дисциплин для выбранной группы
func (h *SubjectsHandler) SubjectsPage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, _ := r.Context().Value("user").(*models.User)

	// Получаем список всех групп для селекта
	groups, err := h.coreClient.GetGroups(r.Context(), session)
	if err != nil {
		RenderInternalError(w, r, err)
		return
	}

	// Получаем groupID из query (если есть)
	var groupID int
	var subjects []models.Subject
	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr != "" {
		groupID, err = strconv.Atoi(groupIDStr)
		if err == nil && groupID > 0 {
			subjects, err = h.coreClient.GetGroupSubjects(r.Context(), session, groupID)
			if err != nil {
				RenderInternalError(w, r, err)
				return
			}
		}
	}

	csrfToken := "" // позже добавим
	pages.AdminSubjectsPage(csrfToken, user, groups, groupID, subjects).Render(r.Context(), w)
}

// GET /admin/subjects/table — HTMX-фрагмент таблицы дисциплин для выбранной группы
func (h *SubjectsHandler) SubjectsTable(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("group_id required"))
		return
	}
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil || groupID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid group_id"))
		return
	}

	subjects, err := h.coreClient.GetGroupSubjects(r.Context(), session, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Push-Url", fmt.Sprintf("/admin/subjects?group_id=%d", groupID))
	components.SubjectsTable(subjects).Render(r.Context(), w)
}
