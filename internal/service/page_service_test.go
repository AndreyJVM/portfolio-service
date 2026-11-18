package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPageService_DefaultsWhenFileNotFound(t *testing.T) {
	ps := NewPageService("/non-existent/path/pages.json")

	meta := ps.GetMetaData("index", "/", "https://test.local")
	if meta.Title != defaultPages["index"].Title {
		t.Errorf("Ожидался дефолтный Title '%s', получен '%s'", defaultPages["index"].Title, meta.Title)
	}
	if meta.OGURL != "https://test.local/" {
		t.Errorf("Ожидался OGURL 'https://test.local/', получен '%s'", meta.OGURL)
	}
	if meta.CanonicalURL != "https://test.local/" {
		t.Errorf("Ожидался CanonicalURL 'https://test.local/', получен '%s'", meta.CanonicalURL)
	}
	if meta.OGImage != "https://test.local/images/og-main.png" {
		t.Errorf("Ожидался OGImage с префиксом baseURL, получен '%s'", meta.OGImage)
	}
}

func TestPageService_LoadFromJSONFile(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "pages.json")
	content := `{
		"custom": {
			"title": "Custom Title",
			"description": "Custom Description",
			"keywords": "k1, k2",
			"ogImage": "/custom.png",
			"activeTab": "custom"
		}
	}`
	if err := os.WriteFile(jsonPath, []byte(content), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый JSON: %v", err)
	}

	ps := NewPageService(jsonPath)
	meta := ps.GetMetaData("custom", "/custom", "https://example.com")

	if meta.Title != "Custom Title" {
		t.Errorf("Ожидался Title 'Custom Title', получен '%s'", meta.Title)
	}
	if meta.OGTitle != "Custom Title" {
		t.Errorf("Ожидался OGTitle равный Title при отсутствии явного, получен '%s'", meta.OGTitle)
	}
	if meta.OGDescription != "Custom Description" {
		t.Errorf("Ожидался OGDescription равный Description, получен '%s'", meta.OGDescription)
	}
	if meta.OGImage != "https://example.com/custom.png" {
		t.Errorf("Ожидался OGImage 'https://example.com/custom.png', получен '%s'", meta.OGImage)
	}
	if meta.ActiveTab != "custom" {
		t.Errorf("Ожидался ActiveTab 'custom', получен '%s'", meta.ActiveTab)
	}
}
