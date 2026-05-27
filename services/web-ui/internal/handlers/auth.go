package handlers

import (
	"log"
	"net/http"
	"time"
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
	maxAge := int(time.Until(session.ExiresAt.Add(30 * 24 * time.Hour)).Seconds())
	if maxAge < 0 {
		maxAge = -1
	}
	userSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
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

	if w.Header().Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
