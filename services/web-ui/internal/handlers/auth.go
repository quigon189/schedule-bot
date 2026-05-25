package handlers

import (
	"log"
	"net/http"
	"web-ui/internal/api"
	"web-ui/views/components"
	"web-ui/views/pages"

	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	coreClient *api.CoreClient
	store      sessions.Store
}

func NewAuthHandler(client *api.CoreClient, store sessions.Store) *AuthHandler {
	return &AuthHandler{coreClient: client, store: store}
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

	session, err := h.coreClient.Login(r, username, password)
	if err != nil {
		log.Printf("error: %v", err)
		errorMsg := "Неверный логин или пароль"
		components.ErrorAlert(errorMsg).Render(r.Context(), w)
		return
	}

	userSession, _ := h.store.Get(r, "user-session")
	userSession.Values["session"] = *session
	err = userSession.Save(r, w)
	if err != nil {
		log.Printf("error save session: %v", err)
		errorMsg := "Ошибка авторизации, попробуйте позже"
		components.ErrorAlert(errorMsg).Render(r.Context(), w)
		return

	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// POST /logout — выход (удаляем cookie)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*api.Session)
	if ok {
		h.coreClient.Logout(r.Context(), session)
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}
