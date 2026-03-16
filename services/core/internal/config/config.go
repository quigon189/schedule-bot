package config

import "os"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

func Load() *Config {
	db := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", ""),
	}

	server := ServerConfig{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnv("SERVER_PORT", "8080"),
	}
	return &Config{
		DB:     db,
		Server: server,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
