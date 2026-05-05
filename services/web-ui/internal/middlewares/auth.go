package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"web-ui/internal/models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err == nil {
			var jwt models.JWT
			if err := json.Unmarshal([]byte(cookie.Value), &jwt); err == nil {
				ctx := context.WithValue(r.Context(), "jwt", jwt)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Нет токена — редирект на логин
		// Для запросов HTMX отдаём редирект-заголовок
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/login")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
}

func GetJWTFromContext(r *http.Request) *models.JWT {
	if jwt, ok := r.Context().Value("jwt").(models.JWT); ok {
		return &jwt
	}
	return nil
}
