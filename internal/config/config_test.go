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
