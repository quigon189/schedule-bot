package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	SevrverPort string

	CodeCharset      string
	CodeLength       int
	MaxGenerateTries int
}

func Load() *Config {
	codeLengthStr := getEnv("CODE_LENGTH", "6")
	codeLength, err := strconv.Atoi(codeLengthStr)
	if err != nil {
		codeLength = 6
	}

	maxTriesStr := getEnv("MAX_GENERATE_TRIES", "10")
	maxTries, err := strconv.Atoi(maxTriesStr)
	if err != nil {
		maxTries = 10
	}
	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "auth_service"),
		SevrverPort:      getEnv("SERVER_PORT", "8080"),
		CodeCharset:      getEnv("CODE_CHARSET", "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"),
		CodeLength:       codeLength,
		MaxGenerateTries: maxTries,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
