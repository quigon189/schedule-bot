package main

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	coreClient := api.NewCoreClient("http://localhost:8088", 30*time.Second)

	authHandler := handlers.NewAuthHandler(coreClient)

	r := chi.NewRouter()

	csrfMiddleware := csrf.Protect(
		[]byte("secret-key-from-config"),
		csrf.Secure(false), // true для https
		csrf.Path("/"),
		csrf.RequestHeader("X-CSRF-Token"),
	)

	r.Use(csrfMiddleware)

	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.LoginSubmit)



	http.ListenAndServe(":8181", r)
}
