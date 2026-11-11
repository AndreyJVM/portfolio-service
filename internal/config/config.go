package config

import (
	"os"
)

type Config struct {
	Port         string
	AppBaseURL   string
	GitHubUser   string
	DockerUser   string
	ContentDir   string
	TemplatesDir string
	StaticDir    string
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func Load() *Config {
	return &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		AppBaseURL:   getEnvOrDefault("APP_BASE_URL", "https://vorobevaqa.ru"),
		GitHubUser:   getEnvOrDefault("GITHUB_USER", "AndreyJVM"),
		DockerUser:   getEnvOrDefault("DOCKER_USER", "andreyvorobevaqa"),
		ContentDir:   getEnvOrDefault("CONTENT_DIR", "content/courses"),
		TemplatesDir: getEnvOrDefault("TEMPLATES_DIR", "web/templates"),
		StaticDir:    getEnvOrDefault("STATIC_DIR", "web/static"),
	}
}
