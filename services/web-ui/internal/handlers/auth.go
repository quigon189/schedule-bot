package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"web-ui/internal/api"
	"web-ui/views/components"
	"web-ui/views/pages"
)

type AuthHandler struct {
	coreClient *api.CoreClient
}

func NewAuthHandler(client *api.CoreClient) *AuthHandler {
	return &AuthHandler{coreClient: client}
}

//GET /login - отображает страницу входа
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

	token, err := h.coreClient.Login(r.Context(), username, password)
	if err != nil {
		log.Printf("error: %v", err)
		errorMsg := "Неверный логин или пароль"
		components.ErrorAlert(errorMsg).Render(r.Context(), w)
		return
	}

	jsonJWT, _ := json.Marshal(token)

	http.SetCookie(w, &http.Cookie{
        Name:     "jwt",
        Value:    string(jsonJWT),
        HttpOnly: true,
        Path:     "/",
        MaxAge:   86400,
        SameSite: http.SameSiteLaxMode,
        Secure:   false,
    })

	w.Header().Set("HX-Redirect", "/")
    w.WriteHeader(http.StatusOK)
}

// POST /logout — выход (удаляем cookie)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:     "jwt",
        Value:    "",
        HttpOnly: true,
        Path:     "/",
        MaxAge:   -1,
    })
    // Для HTMX редиректим на логин
    w.Header().Set("HX-Redirect", "/login")
    w.WriteHeader(http.StatusOK)
}
