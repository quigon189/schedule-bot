package main

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/config"
	"web-ui/internal/router"
	"web-ui/internal/session"

	"github.com/gorilla/sessions"
)

func main() {
	cfg := config.Load()

	cookieStore := sessions.NewCookieStore([]byte(cfg.CookieSecret))
	cookieStore.Options = &sessions.Options{
		Path: "/",
		MaxAge: cfg.SessionMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sessionManager := session.NewSessionManager(cookieStore, "user-session")
	coreClient := api.NewCoreClient("http://localhost:8088", time.Duration(cfg.CoreTimeout)*time.Second, sessionManager)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte("secret-key-from-config"),
	// 	csrf.Secure(false), // true для https
	// 	csrf.Path("/"),
	// 	csrf.RequestHeader("X-CSRF-Token"),
	// )

	r := router.NewRouter(coreClient, sessionManager)

	http.ListenAndServe(":8181", r.Handler())
}
