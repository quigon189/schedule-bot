package main

import (
	"context"
	"core/internal/config"
	"core/internal/dto"
	"core/internal/router"
	"core/pkg/postgres"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	dto.SetupValidator()

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

	r, err := router.New(cfg, pgPool)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	server := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r.Router(),
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
		return fmt.Errorf("set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, path); err != nil {
		return fmt.Errorf("goose Up migrations: %w", err)
	}

	return db.Close()
}
