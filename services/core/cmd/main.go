package main

import (
	"context"
	"core/internal/config"
	"core/internal/handlers"
	"core/internal/middlewares"
	"core/internal/repository"
	"core/internal/services"
	"core/pkg/postgres"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	pgPool, err := postgres.NewPostgresPool(dsn)
	if err != nil {
		log.Fatalf("Failed to get pgPool: %v", err)
	}

	if err := Migrations(pgPool, cfg.DB.MigrationsPath); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	userRepo := repository.NewUserRepo(pgPool)
	sessionRepo := repository.NewSessionRepo(pgPool)

	tokenService := services.NewJWTService([]byte("123"), 24*time.Hour)
	userService := services.NewUserService(userRepo, sessionRepo, tokenService)

	authMiddleware := middlewares.NewAuthMiddleware(userService)

	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", authHandler.Login)
	r.Post("/refresh", authHandler.RefreshToken)
	r.With(authMiddleware.ValidateToken).Group(func(r chi.Router) {
		r.Get("/logout", authHandler.Logout)

		r.With(authMiddleware.AdminRequire).Group(func(r chi.Router) {
			r.Post("/create_user", userHandler.CreateUser)
			r.Get("/get_users", userHandler.GetAllUsers)
		})
	})

	server := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	log.Printf("Auth server starting on port %s", cfg.Server.Port)
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

func Migrations(pool *pgxpool.Pool, path string) error {
	err := goose.SetDialect(string(goose.DialectPostgres))
	if err != nil {
		return nil
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, path); err != nil {
		return err
	}

	return db.Close()
}
