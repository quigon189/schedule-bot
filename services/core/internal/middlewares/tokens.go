package middlewares

import (
	"context"
	"core/internal/services"
	"net/http"
)

type AuthMiddleware struct {
	userService *services.UserService
}

func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService}
}

func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[7:]
		user, err := m.userService.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "access_token", tokenStr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

