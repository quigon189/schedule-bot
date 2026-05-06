package main

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/handlers"
	"web-ui/internal/middlewares"
	"web-ui/internal/session"

	"github.com/go-chi/chi/v5"
	//"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

func main() {
	cookieStore := sessions.NewCookieStore([]byte("secret-key-from-config"))
	cookieStore.Options = &sessions.Options{
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	sessionManager := session.NewSessionManager(cookieStore, "user-session")
	coreClient := api.NewCoreClient("http://localhost:8088", 30*time.Second, sessionManager)

	authHandler := handlers.NewAuthHandler(coreClient, sessionManager)
	dashboardHandler := handlers.NewDashboardHandler(coreClient, sessionManager)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte("secret-key-from-config"),
	// 	csrf.Secure(false), // true для https
	// 	csrf.Path("/"),
	// 	csrf.RequestHeader("X-CSRF-Token"),
	// )

	r := chi.NewRouter()

	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.LoginSubmit)

	r.Group(func(r chi.Router) {
		//r.Use(csrfMiddleware)
		r.Use(middlewares.AuthMiddleware(sessionManager))
		r.Get("/", dashboardHandler.Dashboard)
		r.Post("/logout", authHandler.Logout)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8181", r)
}
