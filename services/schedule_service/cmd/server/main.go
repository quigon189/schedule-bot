package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"schedule-service/internal/config"
	"schedule-service/internal/handlers"
	"schedule-service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := repository.NewPostgresDB(connString)
	if err != nil {
		log.Fatalf("Failed to connetc to DB: %v", err)
	}
	defer db.Close()

	if err := Migrations(db.DB, "./migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	schedulesHandler := handlers.NewScheduleHandler(db)
	changesHandler := handlers.NewChangeHandler(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/group_schedules", func(r chi.Router) {
			r.Get("/", schedulesHandler.GetGroupSchedule)
			r.Post("/", schedulesHandler.AddGroupSchedule)
			r.Delete("/", schedulesHandler.RemoveGroupSchedule)
		})
		r.Route("/changes", func(r chi.Router) {
			r.Get("/", changesHandler.GetChange)
			r.Post("/", changesHandler.AddChange)
			r.Delete("/{id}", changesHandler.RemoveChange)
		})
	})

	r.Get("/health", handlers.HealthCheck)

	server := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	log.Printf("Start schedules server on port %s", cfg.ServerPort)
	go func() {
		if err := server.ListenAndServe(); err != nil || err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	log.Println("Server stoped")
}

func Migrations(db *sql.DB, path string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, path)
}
