package middlewares

import (
	"net/http"
	"web-ui/internal/models"
	"web-ui/internal/session"
)

func AuthMiddleware(sm *session.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := sm.Get(r, "jwt").(models.JWT)
			if !ok {
				// Нет токена — редирект на логин
				// Для запросов HTMX отдаём редирект-заголовок
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				}
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
