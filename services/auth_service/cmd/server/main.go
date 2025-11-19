package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"auth_service/internal/config"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	"auth_service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := repository.NewDB(connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := Migrations(db.DB, "./migrations"); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	log.Printf("Server starting on port %s", cfg.SevrverPort)
	log.Fatal(http.ListenAndServe(":"+cfg.SevrverPort, r))
}

func Migrations(db *sql.DB, path string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return nil
	}

	return goose.Up(db, path)
}
