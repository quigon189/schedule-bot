package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/views/components"
	"web-ui/views/pages"
)

type AdminUsersHandler struct {
	coreClient     *api.CoreClient
}

func NewAdminUsersHandler(client *api.CoreClient) *AdminUsersHandler {
	return &AdminUsersHandler{coreClient: client}
}

// GET /admin/users — страница со списком пользователей
func (h *AdminUsersHandler) ListUsersPage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, err := h.coreClient.GetCurrentUser(r.Context(), session)
	if err != nil {
		w.Header().Set("HX-Redirect", "/logout")
		return
	}
	if !user.HasRole("admin") {
		http.NotFound(w, r)
		return
	}

	// Параметры фильтрации из query
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	perPage := 20
	if pp := r.URL.Query().Get("per_page"); pp != "" {
		if val, err := strconv.Atoi(pp); err == nil && val > 0 {
			perPage = val
		}
	}
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	fullName := r.URL.Query().Get("full_name")
	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")
	groupName := r.URL.Query().Get("group_name")
	roleName := r.URL.Query().Get("role_name")

	params := models.PaginatedUsersQuery{
		Page:      &page,
		PerPage:   &perPage,
		SortBy:    stringPtrOrNil(sortBy),
		SortOrder: stringPtrOrNil(sortOrder),
		FullName:  stringPtrOrNil(fullName),
		Username:  stringPtrOrNil(username),
		Email:     stringPtrOrNil(email),
		GroupName: stringPtrOrNil(groupName),
		RoleName:  stringPtrOrNil(roleName),
	}

	paginated, err := h.coreClient.GetPaginatedUsers(r.Context(), session, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим всю страницу (базовый шаблон + контент)
	pages.AdminUsersPage(user, paginated, params).Render(r.Context(), w)
}

// GET /admin/users/table?page=...&full_name=... — HTMX-фрагмент таблицы
func (h *AdminUsersHandler) TableFragment(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 10
	}
	params := models.PaginatedUsersQuery{
		Page: &page, PerPage: &perPage,
		SortBy:    stringPtrOrNil(r.URL.Query().Get("sort_by")),
		SortOrder: stringPtrOrNil(r.URL.Query().Get("sort_order")),
		FullName:  stringPtrOrNil(r.URL.Query().Get("full_name")),
		Username:  stringPtrOrNil(r.URL.Query().Get("username")),
		Email:     stringPtrOrNil(r.URL.Query().Get("email")),
		GroupName: stringPtrOrNil(r.URL.Query().Get("group_name")),
		RoleName:  stringPtrOrNil(r.URL.Query().Get("role_name")),
	}

	fullURL := "/admin/users?"
	if params.SortBy != nil {
		fullURL += fmt.Sprintf("sort_by=%s&", *params.SortBy)
	}
	if params.SortOrder != nil {
		fullURL += fmt.Sprintf("sort_order=%s&", *params.SortOrder)
	}
	if params.FullName != nil {
		fullURL += fmt.Sprintf("full_name=%s&", *params.FullName)
	}
	if params.Username != nil {
		fullURL += fmt.Sprintf("username=%s&", *params.Username)
	}
	if params.Email != nil {
		fullURL += fmt.Sprintf("email=%s&", *params.Email)
	}
	if params.GroupName != nil {
		fullURL += fmt.Sprintf("group_name=%s&", *params.GroupName)
	}
	if params.RoleName != nil {
		fullURL += fmt.Sprintf("role_name=%s&", *params.RoleName)
	}
	fullURL += fmt.Sprintf("page=%d&per_page=%d", *params.Page, *params.PerPage)

	w.Header().Set("HX-Push-Url", fullURL)

	paginated, err := h.coreClient.GetPaginatedUsers(r.Context(), session, params)
	if err != nil {
		log.Printf("failed to get paginated users")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	components.UsersTable(paginated, params).Render(r.Context(), w)
}


// GET /admin/users/new — форма создания пользователя
func (h *AdminUsersHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	groups, err := h.coreClient.GetGroups(r.Context(), session)
	if err != nil {
		groups = []models.Group{}
	}
	roles := []models.Role{
		models.Role{Name: "student"},
		models.Role{Name: "teacher"},
	}
    components.NewUserForm(roles, groups).Render(r.Context(), w)
}

// POST /admin/users — создание пользователя
func (h *AdminUsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("user-session").(*api.Session)
    req := models.CreateUserRequest{
        Username: r.FormValue("username"),
        Password: r.FormValue("password"),
        FullName: r.FormValue("full_name"),
        Email:    r.FormValue("email"),
		Role: r.FormValue("role"),
    }
    if gid := r.FormValue("group_id"); gid != "" {
        id, _ := strconv.Atoi(gid)
        req.GroupID = &id
    }

    err := h.coreClient.CreateUser(r.Context(), session, req)
    if err != nil {
        components.ErrorAlert(err.Error()).Render(r.Context(), w)
		w.WriteHeader(http.StatusBadRequest)
        return
    }
    // После успешного создания редиректим на список
    w.Header().Set("HX-Redirect", "/admin/users")
    w.WriteHeader(http.StatusOK)
}

// GET /admin/users/{id}/edit — форма редактирования
// func (h *AdminUsersHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
//     idStr := chi.URLParam(r, "id")
//     id, _ := strconv.Atoi(idStr)
//     user, err := h.coreClient.GetUser(w, r, id)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusNotFound)
//         return
//     }
//     components.UserForm(user, nil, nil).Render(r.Context(), w)
// }
//
// // PUT /admin/users/{id} — обновление
// func (h *AdminUsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
//     idStr := chi.URLParam(r, "id")
//     id, _ := strconv.Atoi(idStr)
//     req := models.UpdateUserRequest{
//         FullName: r.FormValue("full_name"),
//         Email:    r.FormValue("email"),
//     }
//     if pwd := r.FormValue("password"); pwd != "" {
//         req.Password = &pwd
//     }
//     // role_ids, group_id аналогично
//     _, err := h.coreClient.UpdateUser(w, r, id, req)
//     if err != nil {
//         components.ErrorAlert(err.Error()).Render(r.Context(), w)
//         return
//     }
//     w.Header().Set("HX-Redirect", "/admin/users")
//     w.WriteHeader(http.StatusOK)
// }
//
// // DELETE /admin/users/{id} — удаление
// func (h *AdminUsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
//     idStr := chi.URLParam(r, "id")
//     id, _ := strconv.Atoi(idStr)
//     err := h.coreClient.DeleteUser(w, r, id)
//     if err != nil {
//         w.WriteHeader(http.StatusInternalServerError)
//         w.Write([]byte(err.Error()))
//         return
//     }
//     w.WriteHeader(http.StatusOK)
// }

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
