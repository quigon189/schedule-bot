package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	CoreURL       string
	CoreTimeout   int
	CookieSecret  string
	SessionMaxAge int
	Port          string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		CoreURL:       getEnv("CORE_URL", "http://localhost:8080"),
		CoreTimeout:   getIntEnv("CORE_TIMEOUT", 30),
		CookieSecret:  getEnv("COOKIE_SECRET", "very-secret-key"),
		SessionMaxAge: getIntEnv("SESSION_MAX_AGE", 3600),
		Port:          getEnv("PORT", "8181"),
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
