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

	_ "auth_service/docs"
	"auth_service/internal/config"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	"auth_service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title			Schedule-bot auth service
// @version		0.1
// @description	Микросервис для управления учетными данными пользователей бота расписания
// @host			auth-service:8080
// @BasePath		/api/v1
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
	codeRepo := repository.NewCodeRepository(db.DB, cfg)

	userService := service.NewUserService(userRepo)
	codeService := service.NewRegistrationCodeService(codeRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	codeHahdler := handlers.NewRegistrationCodeHandler(codeService)

	initAdmins(cfg, userRepo)	

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", codeHahdler.RegisterWithCode)
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
		r.Route("/code", func(r chi.Router) {
			r.Post("/create", codeHahdler.CreateCode)
		})
	})

	r.Mount("/swagger", httpSwagger.WrapHandler)

	server := http.Server{
		Addr:    ":" + cfg.SevrverPort,
		Handler: r,
	}

	log.Printf("Auth server starting on port %s", cfg.SevrverPort)
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
	err := goose.SetDialect("postgres")
	if err != nil {
		return nil
	}

	return goose.Up(db, path)
}
