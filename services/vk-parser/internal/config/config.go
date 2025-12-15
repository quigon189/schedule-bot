package config

import (
	"strconv"
	"os"
	"time"
)

type Config struct {
	VKToken            string
	VKGroupID          int
	CheckInterval      time.Duration
	Pattern            string
	ScheduleServiceURL string
	MaxPosts           int
	TesseractPath      string
}

func Load() *Config {
	interval, _ := strconv.Atoi(getEnv("CHECK_INTERVAL_MINUTES", "5"))
	maxPosts, _ := strconv.Atoi(getEnv("MAX_POSTS", "20"))
	groupID, _ := strconv.Atoi(getEnv(":VK_GROUP_ID", "1"))

	return &Config{
		VKToken: getEnv("VK_ACCESS_TOKEN", ""),
		VKGroupID: groupID,
		CheckInterval: time.Duration(interval) * time.Minute,
		Pattern: getEnv("PATTERN", "Изменение в расписании на"),
		ScheduleServiceURL: getEnv("SCHEDULE_SERVICE_URL", "http://schedule-service:8080"),
		MaxPosts: maxPosts,
		TesseractPath: getEnv("TESSERACT_PATH", "/usr/bin/tesseract"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
