package middlewares

import (
	"context"
	"core/internal/models"
	"core/internal/services"
	"core/pkg/utils"
	"net/http"
	"slices"
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
			utils.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenStr := authHeader[7:]
		user, err := m.userService.ValidateToken(tokenStr)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "user", *user)
		ctx = context.WithValue(ctx, "access_token", tokenStr)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) AdminRequire(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(models.User)
		
		if !ok || !slices.ContainsFunc(user.Roles, func(r models.Role) bool {
			return r.Name == "admin"
		}) {
			utils.ErrorResponse(w, http.StatusForbidden, "access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}

