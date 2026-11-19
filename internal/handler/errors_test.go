package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"ru.vorobevaqa/devops-platform/internal/service"
)

func TestNotFoundHandlers(t *testing.T) {
	r := setupTestRouter()

	notFoundTests := []struct {
		name string
		url  string
	}{
		{"Non-existent course", "/courses/non-existent-course-slug-12345"},
		{"Non-existent lesson", "/courses/docker/non-existent-lesson-slug-12345"},
		{"Non-existent course in lesson path", "/courses/non-existent-course/01-permissions"},
	}

	for _, tt := range notFoundTests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("%s: ожидался статус 404, получен %d", tt.url, w.Code)
			}
		})
	}
}

func TestTemplateErrors(t *testing.T) {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	// Pass invalid templates directory to check template compilation error branch
	h := NewHandler(cs, is, nil, "/tmp/non-existent-templates-dir-12345", "../../web/static", "https://vorobevaqa.ru")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.Index(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500 при ошибке парсинга шаблона, получен %d", w.Code)
	}
}

func TestRenderTemplate_BufferSafety(t *testing.T) {
	tempDir := t.TempDir()

	// Create layout.html with non-existent field call to trigger runtime template execution error
	layoutContent := `{{define "layout"}}{{.InvalidField.SubField}}{{end}}`
	if err := os.WriteFile(filepath.Join(tempDir, "layout.html"), []byte(layoutContent), 0644); err != nil {
		t.Fatalf("Не удалось создать layout: %v", err)
	}
	indexContent := `{{define "content"}}Index{{end}}`
	if err := os.WriteFile(filepath.Join(tempDir, "broken.html"), []byte(indexContent), 0644); err != nil {
		t.Fatalf("Не удалось создать broken.html: %v", err)
	}

	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	h := NewHandler(cs, is, nil, tempDir, "../../web/static", "https://vorobevaqa.ru")

	w := httptest.NewRecorder()
	// Pass struct without InvalidField to trigger ExecuteTemplate execution error
	h.renderTemplate(w, "broken.html", struct{ Dummy string }{Dummy: "test"})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500 при ошибке исполнения шаблона, получен %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Template execution error") {
		t.Errorf("Ожидалось сообщение об ошибке исполнения шаблона, получено: %s", w.Body.String())
	}
}

func TestConcurrency_HighLoad(t *testing.T) {
	r := setupTestRouter()
	const workers = 50

	t.Run("Concurrent Search API requests", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func(idx int) {
				defer wg.Done()
				query := "docker"
				if idx%2 == 0 {
					query = "linux"
				}
				req := httptest.NewRequest("GET", "/api/search?q="+query, nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Ошибка конкурентного запроса к /api/search: статус %d", w.Code)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent lesson and page renders", func(t *testing.T) {
		urls := []string{
			"/",
			"/courses",
			"/courses/docker",
			"/courses/docker/01-architecture-and-first-dockerfile",
			"/courses/podman",
			"/projects",
			"/education",
			"/mustang",
		}

		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func(idx int) {
				defer wg.Done()
				url := urls[idx%len(urls)]
				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				if w.Code != http.StatusOK && w.Code != http.StatusMovedPermanently {
					t.Errorf("Ошибка конкурентного рендеринга %s: статус %d", url, w.Code)
				}
			}(i)
		}

		wg.Wait()
	})
}

func BenchmarkRenderIndexCached(b *testing.B) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func TestCustomErrorHandlers(t *testing.T) {
	r := setupTestRouter()

	t.Run("404 Web HTML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/non-existent-page-12345", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
			t.Errorf("expected Content-Type text/html, got %s", ct)
		}
		body := w.Body.String()
		if !strings.Contains(body, "404") || !strings.Contains(body, "Страница не найдена") {
			t.Errorf("expected 404 HTML body with 'Страница не найдена', got: %s", body)
		}
	})

	t.Run("404 API JSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/unknown-endpoint", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}
		if resp["error"] != "not_found" {
			t.Errorf("expected error 'not_found', got '%s'", resp["error"])
		}
	})

	t.Run("405 API Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/search", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}
		if resp["error"] != "method_not_allowed" {
			t.Errorf("expected error 'method_not_allowed', got '%s'", resp["error"])
		}
	})

	t.Run("405 Web Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
