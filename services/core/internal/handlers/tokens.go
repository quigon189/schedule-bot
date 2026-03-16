func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization") // Ожидаем "Bearer <token>"
		if len(authHeader) < 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[7:]
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Теперь можно проверить роли, например:
		// if !contains(claims.Roles, "admin") { ... }

		next.ServeHTTP(w, r)
	})
}

