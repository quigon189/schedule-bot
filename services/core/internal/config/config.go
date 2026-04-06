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
	URL     string
	Model   string
	Timeout int
}

type GigaChatConfig struct {
	ClientID         string
	ClientSecret     string
	AuthorizationKey string
	BaseURL          string // https://gigachat.devices.sberbank.ru/api/v1
	Model            string // GigaChat, GigaChat-Pro, GigaChat-Lite
	Timeout          int
}

type LLMConfig struct {
	Provider string
	Ollama   OllamaConfig
	GigaChat GigaChatConfig
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	LLM    LLMConfig
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
		URL:     getEnv("OLLAMA_URL", "localhost:11434"),
		Model:   getEnv("OLLAMA_MODEL", "qwen3.5:9b"),
		Timeout: getIntEnv("OLLAMA_TIMEOUT", 600),
	}

	gigaChat := GigaChatConfig{
		ClientID:         getEnv("GIGACHAT_CLIENT_ID", ""),
		ClientSecret:     getEnv("GIGACHAT_CLIENT_SECRET", ""),
		AuthorizationKey: getEnv("GIGACHAT_AUTHORIZATION_KEY", ""),
		BaseURL:          getEnv("GIGACHAT_BASE_URL", "https://gigachat.devices.sberbank.ru/api/v1"),
		Model:            getEnv("GIGACHAT_MODEL", "GigaChat-2-Lite"),
	}

	llm := LLMConfig{
		Provider: getEnv("LLM_PROVIDER", "ollama"), //ollama или gigachat
		Ollama:   ollama,
		GigaChat: gigaChat,
	}

	if db.User == "" || db.Password == "" || db.Name == "" || jwt.Secret == "" {
		log.Fatal("DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET must be set")
	}

	return &Config{
		DB:     db,
		Server: server,
		JWT:    jwt,
		LLM:    llm,
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
