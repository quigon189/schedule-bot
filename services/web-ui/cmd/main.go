package main

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/handlers"
	"web-ui/internal/middlewares"
	"web-ui/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

func main() {
	coreClient := api.NewCoreClient("http://localhost:8088", 30*time.Second)
	cookieStore := sessions.NewCookieStore([]byte("secret-key-from-config"))
	cookieStore.Options = &sessions.Options{
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sessionManager := session.NewSessionManager(cookieStore, "user-session")

	authHandler := handlers.NewAuthHandler(coreClient, sessionManager)
	

	csrfMiddleware := csrf.Protect(
		[]byte("secret-key-from-config"),
		csrf.Secure(false), // true для https
		csrf.Path("/"),
		csrf.RequestHeader("X-CSRF-Token"),
	)

	r := chi.NewRouter()

	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.LoginSubmit)

	r.Group(func(r chi.Router) {
		r.Use(csrfMiddleware)
		r.Use(middlewares.AuthMiddleware(sessionManager))
		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello!"))
		})
	})

	http.ListenAndServe(":8181", r)
}
