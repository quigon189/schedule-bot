package main

import (
	"database/sql"
	"fmt"
	"log"
	"schedule-service/internal/config"
	"schedule-service/internal/repository"

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
}

func Migrations(db *sql.DB, path string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, path)
}
