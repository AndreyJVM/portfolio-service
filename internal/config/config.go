package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Default configuration constants acting as Single Source of Truth.
const (
	DefaultPort         = "8080"
	DefaultAppBaseURL   = "https://vorobevaqa.ru"
	DefaultGitHubUser   = "AndreyJVM"
	DefaultPagesFile    = "content/pages.json"
	DefaultTemplatesDir = "web/templates"
	DefaultStaticDir    = "web/static"
)

type Config struct {
	Port         string
	AppBaseURL   string
	GitHubUser   string
	PagesFile    string
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
		Port:         getEnvOrDefault("PORT", DefaultPort),
		AppBaseURL:   strings.TrimRight(getEnvOrDefault("APP_BASE_URL", DefaultAppBaseURL), "/"),
		GitHubUser:   getEnvOrDefault("GITHUB_USER", DefaultGitHubUser),
		PagesFile:    getEnvOrDefault("PAGES_FILE", DefaultPagesFile),
		TemplatesDir: getEnvOrDefault("TEMPLATES_DIR", DefaultTemplatesDir),
		StaticDir:    getEnvOrDefault("STATIC_DIR", DefaultStaticDir),
	}
}

// Validate verifies that critical configuration parameters are present and well-formed.
func (c *Config) Validate() error {
	// Validate Port
	portNum, err := strconv.Atoi(c.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("invalid port '%s': must be an integer between 1 and 65535", c.Port)
	}

	// Validate AppBaseURL
	if c.AppBaseURL == "" {
		return fmt.Errorf("app base URL cannot be empty")
	}
	parsedURL, err := url.ParseRequestURI(c.AppBaseURL)
	if err != nil {
		return fmt.Errorf("invalid app base URL '%s': %w", c.AppBaseURL, err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid scheme '%s' in app base URL '%s': only http and https are allowed", parsedURL.Scheme, c.AppBaseURL)
	}
	if parsedURL.Host == "" {
		return fmt.Errorf("missing host in app base URL '%s'", c.AppBaseURL)
	}

	// Validate required paths
	if c.PagesFile == "" {
		return fmt.Errorf("pages file cannot be empty")
	}
	if c.TemplatesDir == "" {
		return fmt.Errorf("templates dir cannot be empty")
	}
	if c.StaticDir == "" {
		return fmt.Errorf("static dir cannot be empty")
	}

	return nil
}
