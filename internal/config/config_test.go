package config

import (
	"os"
	"testing"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Clear environment variables to check defaults
	keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "COURSES_FILE", "PAGES_FILE", "TEMPLATES_DIR", "STATIC_DIR"}
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}

	cfg := Load()

	if cfg.Port != DefaultPort {
		t.Errorf("Ожидался Port '%s', получен '%s'", DefaultPort, cfg.Port)
	}
	if cfg.AppBaseURL != DefaultAppBaseURL {
		t.Errorf("Ожидался AppBaseURL '%s', получен '%s'", DefaultAppBaseURL, cfg.AppBaseURL)
	}
	if cfg.GitHubUser != DefaultGitHubUser {
		t.Errorf("Ожидался GitHubUser '%s', получен '%s'", DefaultGitHubUser, cfg.GitHubUser)
	}
	if cfg.DockerUser != DefaultDockerUser {
		t.Errorf("Ожидался DockerUser '%s', получен '%s'", DefaultDockerUser, cfg.DockerUser)
	}
	if cfg.ContentDir != DefaultContentDir {
		t.Errorf("Ожидался ContentDir '%s', получен '%s'", DefaultContentDir, cfg.ContentDir)
	}
	if cfg.CoursesFile != DefaultCoursesFile {
		t.Errorf("Ожидался CoursesFile '%s', получен '%s'", DefaultCoursesFile, cfg.CoursesFile)
	}
	if cfg.PagesFile != DefaultPagesFile {
		t.Errorf("Ожидался PagesFile '%s', получен '%s'", DefaultPagesFile, cfg.PagesFile)
	}
	if cfg.TemplatesDir != DefaultTemplatesDir {
		t.Errorf("Ожидался TemplatesDir '%s', получен '%s'", DefaultTemplatesDir, cfg.TemplatesDir)
	}
	if cfg.StaticDir != DefaultStaticDir {
		t.Errorf("Ожидался StaticDir '%s', получен '%s'", DefaultStaticDir, cfg.StaticDir)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Ожидалась валидная дефолтная конфигурация, получена ошибка: %v", err)
	}
}

func TestConfig_CustomEnvValues(t *testing.T) {
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("APP_BASE_URL", "http://localhost:9090")
	_ = os.Setenv("GITHUB_USER", "CustomUser")
	_ = os.Setenv("DOCKER_USER", "customdocker")
	_ = os.Setenv("CONTENT_DIR", "/custom/content")
	_ = os.Setenv("COURSES_FILE", "/custom/courses.json")
	_ = os.Setenv("PAGES_FILE", "/custom/pages.json")
	_ = os.Setenv("TEMPLATES_DIR", "/custom/templates")
	_ = os.Setenv("STATIC_DIR", "/custom/static")

	defer func() {
		keys := []string{"PORT", "APP_BASE_URL", "GITHUB_USER", "DOCKER_USER", "CONTENT_DIR", "COURSES_FILE", "PAGES_FILE", "TEMPLATES_DIR", "STATIC_DIR"}
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
	if cfg.CoursesFile != "/custom/courses.json" {
		t.Errorf("Ожидался CoursesFile '/custom/courses.json', получен '%s'", cfg.CoursesFile)
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

func TestConfig_Validate(t *testing.T) {
	validCfg := func() *Config {
		return &Config{
			Port:         DefaultPort,
			AppBaseURL:   DefaultAppBaseURL,
			GitHubUser:   DefaultGitHubUser,
			DockerUser:   DefaultDockerUser,
			ContentDir:   DefaultContentDir,
			PagesFile:    DefaultPagesFile,
			TemplatesDir: DefaultTemplatesDir,
			StaticDir:    DefaultStaticDir,
		}
	}

	tests := []struct {
		name    string
		modify  func(c *Config)
		wantErr bool
	}{
		{
			name:    "valid configuration",
			modify:  func(c *Config) {},
			wantErr: false,
		},
		{
			name:    "invalid port non-numeric",
			modify:  func(c *Config) { c.Port = "abc" },
			wantErr: true,
		},
		{
			name:    "invalid port zero",
			modify:  func(c *Config) { c.Port = "0" },
			wantErr: true,
		},
		{
			name:    "invalid port negative",
			modify:  func(c *Config) { c.Port = "-1" },
			wantErr: true,
		},
		{
			name:    "invalid port too high",
			modify:  func(c *Config) { c.Port = "70000" },
			wantErr: true,
		},
		{
			name:    "empty app base URL",
			modify:  func(c *Config) { c.AppBaseURL = "" },
			wantErr: true,
		},
		{
			name:    "invalid app base URL scheme (ftp)",
			modify:  func(c *Config) { c.AppBaseURL = "ftp://vorobevaqa.ru" },
			wantErr: true,
		},
		{
			name:    "invalid app base URL format",
			modify:  func(c *Config) { c.AppBaseURL = "not-a-valid-url" },
			wantErr: true,
		},
		{
			name:    "missing host in app base URL",
			modify:  func(c *Config) { c.AppBaseURL = "https://" },
			wantErr: true,
		},
		{
			name:    "empty content dir",
			modify:  func(c *Config) { c.ContentDir = "" },
			wantErr: true,
		},
		{
			name:    "empty pages file",
			modify:  func(c *Config) { c.PagesFile = "" },
			wantErr: true,
		},
		{
			name:    "empty templates dir",
			modify:  func(c *Config) { c.TemplatesDir = "" },
			wantErr: true,
		},
		{
			name:    "empty static dir",
			modify:  func(c *Config) { c.StaticDir = "" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validCfg()
			tt.modify(cfg)
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
