package handlers

import (
	"log"
	"net/http"
	"web-ui/internal/api"
	"web-ui/internal/session"
	"web-ui/views/components"
	"web-ui/views/pages"
)

type AuthHandler struct {
	coreClient     *api.CoreClient
	sessionManager *session.SessionManager
}

func NewAuthHandler(client *api.CoreClient, sm *session.SessionManager) *AuthHandler {
	return &AuthHandler{coreClient: client, sessionManager: sm}
}

// GET /login - отображает страницу входа
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	pages.LoginPage("").Render(r.Context(), w)
}

// POST /login - обрабатывает форму
func (h *AuthHandler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		errorMsg := "Заполните обя поля"
		components.ErrorAlert(errorMsg).Render(r.Context(), w)
		return
	}

	token, err := h.coreClient.Login(r, username, password)
	if err != nil {
		log.Printf("error: %v", err)
		errorMsg := "Неверный логин или пароль"
		components.ErrorAlert(errorMsg).Render(r.Context(), w)
		return
	}

	if err := h.sessionManager.Set(w, r, "jwt", token); err != nil {
		log.Printf("failed to set jwt in session: %v", err)
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// POST /logout — выход (удаляем cookie)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.coreClient.Logout(w, r)
	// Для HTMX редиректим на логин
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}
