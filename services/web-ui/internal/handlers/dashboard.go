package handlers

import (
	"net/http"
	"web-ui/internal/api"
	"web-ui/internal/session"
	"web-ui/views/components"
	"web-ui/views/pages"

	"github.com/gorilla/csrf"
)

type DashboardHandler struct {
	coreClient     *api.CoreClient
	sessionManager *session.SessionManager
}

func NewDashboardHandler(client *api.CoreClient, sm *session.SessionManager) *DashboardHandler {
	return &DashboardHandler{coreClient: client, sessionManager: sm}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user, err := h.coreClient.GetCurrentUser(w, r)
	if err != nil {
		h.sessionManager.Logout(w, r)
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}

	csrfToken := csrf.Token(r)

	if user.HasRole("admin") {
		pages.AdminDashboardPage(csrfToken, user).Render(r.Context(), w)
	} else if user.HasRole("student") {
		var props components.ScheduleViewProps
		student, err := h.coreClient.GetStudent(w, r, user)
		if err == nil {
			scheduleData, err := h.coreClient.GetGroupSchedule(w, r, student.Group.ID, nil)
			if err == nil {
				props.Data = *scheduleData
			}
		}
		props.ShowTeacher = true
		props.ShowAudience = true
		pages.StudentDashboardPage(csrfToken, user, props).Render(r.Context(), w)
	} else {
		h.sessionManager.Logout(w, r)
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	}
}
