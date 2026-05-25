package handlers

import (
	"net/http"
	"web-ui/internal/api"
	"web-ui/views/components"
	"web-ui/views/pages"

	"github.com/gorilla/csrf"
)

type DashboardHandler struct {
	coreClient *api.CoreClient
}

func NewDashboardHandler(client *api.CoreClient) *DashboardHandler {
	return &DashboardHandler{coreClient: client}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, err := h.coreClient.GetCurrentUser(r.Context(), session)
	if err != nil {
		w.Header().Set("HX-Redirect", "/logout")
		w.WriteHeader(http.StatusOK)
		return
	}

	csrfToken := csrf.Token(r)

	if user.HasRole("admin") {
		pages.AdminDashboardPage(csrfToken, user).Render(r.Context(), w)
	} else if user.HasRole("student") {
		var props components.ScheduleViewProps
		if user.Group != nil {
			scheduleData, err := h.coreClient.GetGroupSchedule(r.Context(), session, user.Group.ID, nil)
			if err == nil {
				props.Data = *scheduleData
			}
		}
		props.ShowTeacher = true
		props.ShowAudience = true
		pages.StudentDashboardPage(csrfToken, user, props).Render(r.Context(), w)
	} else {
		w.Header().Set("HX-Redirect", "/logout")
		w.WriteHeader(http.StatusOK)
	}
}
