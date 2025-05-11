package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server struct {
		Port string
		Env  string
	}
	OpenRouter struct {
		APIKey string
		Model  string
	}
	Image struct {
		MaxSize      int64
		AllowedTypes []string
	}
}

func Load() (*Config, error) {
	config := &Config{}

	// Server config
	config.Server.Port = getEnv("PORT", "8080")
	config.Server.Env = getEnv("ENV", "development")

	// OpenRouter config
	config.OpenRouter.APIKey = getEnv("OPENROUTER_API_KEY", "")
	// 使用更快的模型，優先使用付費模型，如果沒有則使用免費模型
	config.OpenRouter.Model = getEnv("OPENROUTER_MODEL", "anthropic/claude-3-sonnet:free")

	// Image config
	maxSize, _ := strconv.ParseInt(getEnv("MAX_IMAGE_SIZE", "5242880"), 10, 64)
	config.Image.MaxSize = maxSize
	config.Image.AllowedTypes = strings.Split(getEnv("ALLOWED_IMAGE_TYPES", "image/jpeg,image/png"), ",")

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
