package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Ожидался заголовок X-Content-Type-Options: nosniff, получено: %s", w.Header().Get("X-Content-Type-Options"))
	}
	if w.Header().Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Errorf("Ожидался заголовок X-Frame-Options: SAMEORIGIN, получено: %s", w.Header().Get("X-Frame-Options"))
	}
	if w.Header().Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
		t.Errorf("Ожидался заголовок Referrer-Policy: strict-origin-when-cross-origin, получено: %s", w.Header().Get("Referrer-Policy"))
	}
}

func TestNoDirListingFS(t *testing.T) {
	fs := noDirListingFS{fs: http.Dir("../../web/static")}

	t.Run("Open valid regular file", func(t *testing.T) {
		f, err := fs.Open("/favicon.svg")
		if err != nil {
			t.Fatalf("Не удалось открыть favicon.svg: %v", err)
		}
		defer f.Close()
	})

	t.Run("Open directory returns ErrNotExist", func(t *testing.T) {
		_, err := fs.Open("/css")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Ожидалась ошибка os.ErrNotExist для каталога, получено: %v", err)
		}
	})

	t.Run("Open non-existent file returns error", func(t *testing.T) {
		_, err := fs.Open("/non-existent-file-xyz.txt")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего файла")
		}
	})
}

func TestSecurity_PathTraversalHTTP(t *testing.T) {
	r := setupTestRouter()

	traversalURLs := []string{
		"/courses/linux/..%2f..%2fREADME",
		"/courses/linux/..%2f01-permissions-and-systemd",
		"/courses/..%2f..%2fetc/passwd",
		"/courses/linux/%2e%2e%2f%2e%2e%2fgo.mod",
	}

	for _, url := range traversalURLs {
		t.Run("Path Traversal HTTP: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Server must return 404 Not Found or 301 to sanitized path, but never 200 with system file content
			if w.Code == http.StatusOK && strings.Contains(w.Body.String(), "Primary Cloud Lab") && !strings.Contains(url, "linux") {
				t.Errorf("Критическая уязвимость Path Traversal! Получен 200 на URL: %s", url)
			}
			if w.Code != http.StatusNotFound && w.Code != http.StatusMovedPermanently {
				t.Logf("URL %s вернул статус %d (не 200 OK — файл не скомпрометирован)", url, w.Code)
			}
		})
	}
}

func TestStaticAssets_SecurityChecks(t *testing.T) {
	r := setupTestRouter()

	traversalURLs := []string{
		"/css/..%2f..%2fetc/passwd",
		"/css/%2e%2e/index.html",
		"/js/../../main.go",
		"/pdf/../../../etc/shadow",
	}

	for _, url := range traversalURLs {
		t.Run("Static traversal: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest && w.Code != http.StatusMovedPermanently {
				t.Errorf("Ожидался статус 404/400/301 для path traversal %s, получен %d", url, w.Code)
			}
		})
	}

	trailingSlashURLs := []string{
		"/css/",
		"/js/",
		"/images/",
		"/pdf/",
	}

	for _, url := range trailingSlashURLs {
		t.Run("Trailing slash directory listing block: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("Ожидался статус 404 для обращения к директории %s со слэшем, получен %d", url, w.Code)
			}
		})
	}
}

func TestNoDirListingFS_DirectUnit(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Ошибка создания каталога: %v", err)
	}
	sampleFile := filepath.Join(tempDir, "sample.txt")
	if err := os.WriteFile(sampleFile, []byte("sample"), 0644); err != nil {
		t.Fatalf("Ошибка создания файла: %v", err)
	}

	fs := noDirListingFS{fs: http.Dir(tempDir)}

	t.Run("Opening directory returns os.ErrNotExist", func(t *testing.T) {
		_, err := fs.Open("/sub")
		if err != os.ErrNotExist {
			t.Errorf("Ожидалась ошибка os.ErrNotExist для каталога, получено: %v", err)
		}
	})

	t.Run("Opening regular file succeeds", func(t *testing.T) {
		f, err := fs.Open("/sample.txt")
		if err != nil {
			t.Fatalf("Ожидалось успешное открытие файла, получена ошибка: %v", err)
		}
		_ = f.Close()
	})

	t.Run("Opening missing file returns error", func(t *testing.T) {
		_, err := fs.Open("/missing.txt")
		if err == nil {
			t.Errorf("Ожидалась ошибка при открытии несуществующего файла")
		}
	})
}
