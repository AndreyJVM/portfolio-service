package config

import (
	"os"
	"testing"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Clear environment variables to check defaults
	keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "PAGES_FILE", "TEMPLATES_DIR", "STATIC_DIR"}
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Ожидался Port '8080', получен '%s'", cfg.Port)
	}
	if cfg.AppBaseURL != "https://vorobevaqa.ru" {
		t.Errorf("Ожидался AppBaseURL 'https://vorobevaqa.ru', получен '%s'", cfg.AppBaseURL)
	}
	if cfg.GitHubUser != "AndreyJVM" {
		t.Errorf("Ожидался GitHubUser 'AndreyJVM', получен '%s'", cfg.GitHubUser)
	}
	if cfg.DockerUser != "andreyvorobevaqa" {
		t.Errorf("Ожидался DockerUser 'andreyvorobevaqa', получен '%s'", cfg.DockerUser)
	}
	if cfg.ContentDir != "content/courses" {
		t.Errorf("Ожидался ContentDir 'content/courses', получен '%s'", cfg.ContentDir)
	}
	if cfg.PagesFile != "content/pages.json" {
		t.Errorf("Ожидался PagesFile 'content/pages.json', получен '%s'", cfg.PagesFile)
	}
	if cfg.TemplatesDir != "web/templates" {
		t.Errorf("Ожидался TemplatesDir 'web/templates', получен '%s'", cfg.TemplatesDir)
	}
	if cfg.StaticDir != "web/static" {
		t.Errorf("Ожидался StaticDir 'web/static', получен '%s'", cfg.StaticDir)
	}
}

func TestConfig_CustomEnvValues(t *testing.T) {
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("APP_BASE_URL", "http://localhost:9090")
	_ = os.Setenv("GITHUB_USER", "CustomUser")
	_ = os.Setenv("DOCKER_USER", "customdocker")
	_ = os.Setenv("CONTENT_DIR", "/custom/content")
	_ = os.Setenv("PAGES_FILE", "/custom/pages.json")
	_ = os.Setenv("TEMPLATES_DIR", "/custom/templates")
	_ = os.Setenv("STATIC_DIR", "/custom/static")

	defer func() {
		keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "PAGES_FILE", "TEMPLATES_DIR", "STATIC_DIR"}
		for _, k := range keys {
			_ = os.Unsetenv(k)
		}
	}()

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Ожидался Port '9090', получен '%s'", cfg.Port)
	}
	if cfg.AppBaseURL != "http://localhost:9090" {
		t.Errorf("Ожидался AppBaseURL 'http://localhost:9090', получен '%s'", cfg.AppBaseURL)
	}
	if cfg.GitHubUser != "CustomUser" {
		t.Errorf("Ожидался GitHubUser 'CustomUser', получен '%s'", cfg.GitHubUser)
	}
	if cfg.DockerUser != "customdocker" {
		t.Errorf("Ожидался DockerUser 'customdocker', получен '%s'", cfg.DockerUser)
	}
	if cfg.ContentDir != "/custom/content" {
		t.Errorf("Ожидался ContentDir '/custom/content', получен '%s'", cfg.ContentDir)
	}
	if cfg.PagesFile != "/custom/pages.json" {
		t.Errorf("Ожидался PagesFile '/custom/pages.json', получен '%s'", cfg.PagesFile)
	}
	if cfg.TemplatesDir != "/custom/templates" {
		t.Errorf("Ожидался TemplatesDir '/custom/templates', получен '%s'", cfg.TemplatesDir)
	}
	if cfg.StaticDir != "/custom/static" {
		t.Errorf("Ожидался StaticDir '/custom/static', получен '%s'", cfg.StaticDir)
	}
}

func TestGetEnvOrDefaults(t *testing.T) {
	t.Run("fallback when unset", func(t *testing.T) {
		result := getEnvOrDefault("NON_EXISTING_KEY_12345", "my-default-value")
		if result != "my-default-value" {
			t.Errorf("expected 'my-default-value', but got '%s'", result)
		}
	})

	t.Run("value when set", func(t *testing.T) {
		t.Setenv("TEST_EXISTING_VAR", "actual_value")
		result := getEnvOrDefault("TEST_EXISTING_VAR", "default_val")
		if result != "actual_value" {
			t.Errorf("expected 'actual_value', got '%s'", result)
		}
	})
}

func TestGetEnvOrDefault_TableDriven(t *testing.T) {
	// 1. Define columns and rows for our test table
	tests := []struct {
		name        string // Test case description
		envKey      string // Environment variable key
		envVal      string // Value (do not set if empty)
		fallbackVal string // Default value
		expected    string // Expected return value of the function
	}{
		{
			name:        "returns fallback when env var does not exist",
			envKey:      "TEST_NOT_SET_KEY",
			envVal:      "",
			fallbackVal: "default_1",
			expected:    "default_1",
		},
		{
			name:        "returns env value when set",
			envKey:      "TEST_SET_KEY",
			envVal:      "custom_value_42",
			fallbackVal: "default_1",
			expected:    "custom_value_42",
		},
		{
			name:        "returns fallback when env var is explicitly empty string",
			envKey:      "TEST_EMPTY_KEY",
			envVal:      "",
			fallbackVal: "fallback_for_empty",
			expected:    "fallback_for_empty",
		},
	}

	// 2. Iterate through the table and run each case as a subtest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.envKey, tt.envVal)
			}

			actual := getEnvOrDefault(tt.envKey, tt.fallbackVal)

			if actual != tt.expected {
				t.Errorf("case '%s': expected '%s', got '%s'", tt.name, tt.expected, actual)
			}
		})
	}
}
