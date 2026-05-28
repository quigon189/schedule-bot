package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
	"web-ui/views/pages"
)

// RecoveryWithHTML перехватывает panic и возвращает HTML-страницу 500.
func RecoveryWithHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				log.Printf("%s", debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
				// Игнорируем возможную ошибку рендеринга, так как запрос уже в панике
				pages.InternalServerErrorPage().Render(r.Context(), w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
