package config

import (
	"os"
	"testing"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Очищаем переменные окружения для проверки дефолтов
	keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "TEMPLATES_DIR", "STATIC_DIR"}
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
	_ = os.Setenv("TEMPLATES_DIR", "/custom/templates")
	_ = os.Setenv("STATIC_DIR", "/custom/static")

	defer func() {
		keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "TEMPLATES_DIR", "STATIC_DIR"}
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
	// 1. Описываем колонки и строки нашей таблицы проверок
	tests := []struct {
		name        string // Описание тест-кейса
		envKey      string // Ключ переменной
		envVal      string // Значение (если пустое — не устанавливаем)
		fallbackVal string // Дефолтное значение
		expected    string // Что должна вернуть функция
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

	// 2. Итерируемся по таблице и запускаем каждый кейс как подтест
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
