package config

import "os"

type Config struct {
	ScheduleServiceURL string
	Timeout            string
	ServerPort         string
}

func Load() *Config {
	return &Config{
		ScheduleServiceURL: getEnv("SCHEDULE_SERVICE_URL", "http://schedule-service:8080"),
		Timeout:            getEnv("TIMEOUT", "30"),
		ServerPort:         "8080",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
