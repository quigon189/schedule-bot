package config

import (
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	MigrationsPath string
}

type ServerConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret  string
	Expires int //seconds
}

type OllamaConfig struct {
	URL   string
	Model string
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Ollama OllamaConfig
}

func Load() *Config {
	db := DBConfig{
		Host:           getEnv("DB_HOST", "localhost"),
		Port:           getEnv("DB_PORT", "5432"),
		User:           getEnv("DB_USER", ""),
		Password:       getEnv("DB_PASSWORD", ""),
		Name:           getEnv("DB_NAME", ""),
		MigrationsPath: getEnv("DB_MIGRATIONS", "migrations"),
	}

	server := ServerConfig{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnv("SERVER_PORT", "8080"),
	}

	jwt := JWTConfig{
		Secret:  getEnv("JWT_SECRET", ""),
		Expires: getIntEnv("JWT_EXPIRES_SECONDS", 15*60),
	}

	ollama := OllamaConfig{
		URL: getEnv("OLLAMA_URL", "localhost:11434"),
		Model: getEnv("OLLAMA_MODEL", "qwen3.5:9b"),
	}

	if db.User == "" || db.Password == "" || db.Name == "" || jwt.Secret == "" {
		log.Fatal("DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET must be set")
	}

	return &Config{
		DB:     db,
		Server: server,
		JWT:    jwt,
		Ollama: ollama,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
