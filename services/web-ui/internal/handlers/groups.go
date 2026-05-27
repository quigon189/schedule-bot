package handlers

import (
	"net/http"
	"sort"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/views/pages"
)

type GroupsHandler struct {
	coreClient *api.CoreClient
}

func NewGroupsHandler(client *api.CoreClient) *GroupsHandler {
	return &GroupsHandler{coreClient: client}
}

// GET /admin/groups – страница со списком групп (карточки)
func (h *GroupsHandler) ListGroupsPage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, _ := r.Context().Value("user").(*models.User)

	groups, err := h.coreClient.GetGroups(r.Context(), session)
	if err != nil {
		RenderInternalError(w, r, err)
		return
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	// CSRF-токен пока передаём пустым, как в admin_users
	pages.AdminGroupsPage("", user, groups).Render(r.Context(), w)
}
