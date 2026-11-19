package service

import (
	"strings"
	"testing"
)

func TestExtractTitleAndTitleCase(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		fallbackSlug string
		expected     string
	}{
		{
			name:         "Markdown with H1",
			content:      "# Архитектура Linux ядра\nТекст урока",
			fallbackSlug: "01-kernel-architecture",
			expected:     "Архитектура Linux ядра",
		},
		{
			name:         "Markdown with empty lines before H1",
			content:      "\n\n# Права доступа chmod\nПараграф",
			fallbackSlug: "01-permissions",
			expected:     "Права доступа chmod",
		},
		{
			name:         "No H1 header - fallback with numeric prefix",
			content:      "Просто текст без решетки заголовка",
			fallbackSlug: "02-systemd-units-guide",
			expected:     "Systemd Units Guide",
		},
		{
			name:         "No H1 header - fallback without prefix",
			content:      "",
			fallbackSlug: "docker-compose-intro",
			expected:     "Docker Compose Intro",
		},
		{
			name:         "No H1 header - single digit prefix",
			content:      "",
			fallbackSlug: "1-bash-scripting",
			expected:     "Bash Scripting",
		},
		{
			name:         "No H1 header - three digits prefix",
			content:      "",
			fallbackSlug: "007-advanced-kubernetes",
			expected:     "Advanced Kubernetes",
		},
		{
			name:         "No H1 header - unicode slug without digits",
			content:      "",
			fallbackSlug: "настройка-сети",
			expected:     "Настройка Сети",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle([]byte(tt.content), tt.fallbackSlug)
			if result != tt.expected {
				t.Errorf("extractTitle() = '%s', ожидалось '%s'", result, tt.expected)
			}
		})
	}

	t.Run("titleCase function", func(t *testing.T) {
		cases := []struct {
			input    string
			expected string
		}{
			{"hello world", "Hello World"},
			{"devops engineer", "Devops Engineer"},
			{"", ""},
		}

		for _, c := range cases {
			actual := titleCase(c.input)
			if actual != c.expected {
				t.Errorf("titleCase('%s') = '%s', ожидалось '%s'", c.input, actual, c.expected)
			}
		}
	})
}

func TestParseFrontmatter(t *testing.T) {
	t.Run("Full frontmatter with tags", func(t *testing.T) {
		raw := `---
title: "Сетевые драйверы Docker"
description: "Разбор veth-пар и DNS"
duration: "18 мин"
level: "Middle"
order: 3
tags: ["Docker", "Networking", "Bridge"]
---
# Заголовок урока
Тело статьи без frontmatter.`

		meta, body := parseFrontmatter([]byte(raw))
		if meta.Title != "Сетевые драйверы Docker" {
			t.Errorf("Ожидался title 'Сетевые драйверы Docker', получено: '%s'", meta.Title)
		}
		if meta.Description != "Разбор veth-пар и DNS" {
			t.Errorf("Ожидался description, получено: '%s'", meta.Description)
		}
		if meta.Duration != "18 мин" {
			t.Errorf("Ожидалась duration '18 мин', получено: '%s'", meta.Duration)
		}
		if meta.Level != "Middle" {
			t.Errorf("Ожидался level 'Middle', получено: '%s'", meta.Level)
		}
		if meta.Order != 3 {
			t.Errorf("Ожидался order 3, получено: %d", meta.Order)
		}
		if len(meta.Tags) != 3 || meta.Tags[0] != "Docker" || meta.Tags[1] != "Networking" || meta.Tags[2] != "Bridge" {
			t.Errorf("Неверные теги: %v", meta.Tags)
		}
		if strings.Contains(string(body), "---") {
			t.Errorf("Тело статьи не должно содержать разделители frontmatter: %s", string(body))
		}
		if !strings.HasPrefix(strings.TrimSpace(string(body)), "# Заголовок урока") {
			t.Errorf("Тело статьи должно начинаться с заголовка: %s", string(body))
		}
	})

	t.Run("Without frontmatter", func(t *testing.T) {
		raw := "# Обычный заголовок\nБез метаданных."
		meta, body := parseFrontmatter([]byte(raw))
		if meta.Title != "" {
			t.Errorf("Ожидался пустой title, получено: '%s'", meta.Title)
		}
		if string(body) != raw {
			t.Errorf("Тело должно совпадать с исходным текстом")
		}
	})

	t.Run("With Windows CRLF line endings", func(t *testing.T) {
		raw := "---\r\ntitle: \"CRLF Заголовок\"\r\norder: 5\r\n---\r\n# Текст с CRLF"
		meta, body := parseFrontmatter([]byte(raw))
		if meta.Title != "CRLF Заголовок" {
			t.Errorf("Ожидался title 'CRLF Заголовок', получено: '%s'", meta.Title)
		}
		if meta.Order != 5 {
			t.Errorf("Ожидался order 5, получено: %d", meta.Order)
		}
		if !strings.Contains(string(body), "# Текст с CRLF") {
			t.Errorf("Тело должно содержать заголовок")
		}
	})
}
