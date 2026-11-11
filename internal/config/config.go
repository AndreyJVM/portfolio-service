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

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "https://vorobevaqa.ru"
	}

	return &Config{
		Port:         port,
		AppBaseURL:   baseURL,
		GitHubUser:   "AndreyJVM",
		DockerUser:   "andreyvorobevaqa",
		ContentDir:   "content/courses",
		TemplatesDir: "web/templates",
		StaticDir:    "web/static",
	}
}
