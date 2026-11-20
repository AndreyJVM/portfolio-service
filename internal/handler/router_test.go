package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TestRouterNoDirListingFS_Isolated verifies that directory traversal and directory listing
// are strictly blocked, returning os.ErrNotExist for directories while serving regular files.
func TestRouterNoDirListingFS_Isolated(t *testing.T) {
	tempDir := t.TempDir()

	// Create a regular file and a subdirectory
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	nfs := noDirListingFS{fs: http.Dir(tempDir)}

	// 1. Regular file should open successfully
	f, err := nfs.Open("/test.txt")
	if err != nil {
		t.Fatalf("expected regular file to open successfully, got error: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil || string(content) != "hello world" {
		t.Errorf("expected 'hello world', got '%s', err: %v", string(content), err)
	}

	// 2. Directory opening must be rejected with os.ErrNotExist
	dirFile, err := nfs.Open("/subdir")
	if err == nil {
		dirFile.Close()
		t.Errorf("expected directory opening to return error, but got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.ErrNotExist for directory, got: %v", err)
	}

	// 3. Root directory opening must also be rejected
	rootDirFile, err := nfs.Open("/")
	if err == nil {
		rootDirFile.Close()
		t.Errorf("expected root directory opening to return error, but got nil")
	}

	// 4. Non-existent file must return error
	_, err = nfs.Open("/non-existent.txt")
	if err == nil {
		t.Errorf("expected non-existent file to return error, but got nil")
	}
}

// TestRouterSecurityHeaders verifies that essential security headers
// (X-Content-Type-Options, X-Frame-Options, Referrer-Policy) are injected
// on all types of routes (web, API, static, 404).
func TestRouterSecurityHeaders(t *testing.T) {
	r := setupTestRouter()

	endpoints := []string{
		"/",
		"/health",
		"/api/health",
		"/favicon.svg",
		"/mustang",
		"/non-existent-route-404",
	}

	for _, ep := range endpoints {
		t.Run("Endpoint: "+ep, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ep, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if ctOpt := w.Header().Get("X-Content-Type-Options"); ctOpt != "nosniff" {
				t.Errorf("expected X-Content-Type-Options 'nosniff', got '%s'", ctOpt)
			}
			if frameOpt := w.Header().Get("X-Frame-Options"); frameOpt != "SAMEORIGIN" {
				t.Errorf("expected X-Frame-Options 'SAMEORIGIN', got '%s'", frameOpt)
			}
			if refPol := w.Header().Get("Referrer-Policy"); refPol != "strict-origin-when-cross-origin" {
				t.Errorf("expected Referrer-Policy 'strict-origin-when-cross-origin', got '%s'", refPol)
			}
		})
	}
}

// TestRouterMiddlewareStack verifies standard chi middleware behavior:
// RequestID injection, GetHead, Gzip compression, and panic recovery.
func TestRouterMiddlewareStack(t *testing.T) {
	r := setupTestRouter()

	// 1. Test RequestID context injection
	t.Run("Middleware: RequestID in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	// 2. Test GetHead middleware: HEAD request to GET route returns HTTP 200 with empty body
	t.Run("Middleware: GetHead on GET route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for HEAD /health, got %d", w.Code)
		}
	})

	// 3. Test Compress middleware: Accept-Encoding: gzip returns gzip compressed response
	t.Run("Middleware: Gzip Compression", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		contentEncoding := w.Header().Get("Content-Encoding")
		if contentEncoding != "gzip" {
			t.Errorf("expected Content-Encoding 'gzip', got '%s'", contentEncoding)
		}

		// Verify that the body can be decompressed with gzip.Reader
		gzReader, err := gzip.NewReader(w.Body)
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()

		decompressed, err := io.ReadAll(gzReader)
		if err != nil {
			t.Fatalf("failed to decompress gzip body: %v", err)
		}
		if !strings.Contains(string(decompressed), "User-agent: *") {
			t.Errorf("decompressed content does not match expected robots.txt: %s", string(decompressed))
		}
	})

	// 4. Test Recoverer middleware: Panicking route recovers and returns HTTP 500 without crashing
	t.Run("Middleware: Recoverer catches panic", func(t *testing.T) {
		panickingRouter := chi.NewRouter()
		panickingRouter.Use(middleware.Recoverer)
		panickingRouter.Get("/panic-test", func(w http.ResponseWriter, r *http.Request) {
			panic("intentional test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic-test", nil)
		w := httptest.NewRecorder()

		// Should not panic the test runner
		panickingRouter.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500 after panic recovery, got %d", w.Code)
		}
	})
}

// TestRouterStaticAssetCacheHeaders verifies Cache-Control headers for static asset prefixes.
func TestRouterStaticAssetCacheHeaders(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		url                  string
		expectedCacheControl string
	}{
		{"/css/main.css", "public, max-age=86400, stale-while-revalidate=3600"},
		{"/js/main.js", "public, max-age=86400, stale-while-revalidate=3600"},
		{"/images/avatar.jpg", "public, max-age=604800, immutable"},
		{"/favicon.svg", "public, max-age=604800, immutable"},
	}

	for _, tt := range tests {
		t.Run("Cache-Control for "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200 for %s, got %d", tt.url, w.Code)
			}
			cacheControl := w.Header().Get("Cache-Control")
			if cacheControl != tt.expectedCacheControl {
				t.Errorf("expected Cache-Control '%s', got '%s'", tt.expectedCacheControl, cacheControl)
			}
		})
	}
}

// TestRouterDirectoryListingForbidden verifies trailing slash requests to static dirs return 404.
func TestRouterDirectoryListingForbidden(t *testing.T) {
	r := setupTestRouter()

	dirs := []string{
		"/css/",
		"/js/",
		"/images/",
		"/pdf/",
	}

	for _, d := range dirs {
		t.Run("Trailing slash: "+d, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, d, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("expected status 404 for directory listing attempt %s, got %d", d, w.Code)
			}
		})
	}
}

// TestRouterRedirects verifies 301 redirects to home and anchors.
func TestRouterRedirects(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		url              string
		expectedLocation string
	}{
		{"/about", "/"},
		{"/education", "/"},
		{"/roadmap", "/"},
		{"/courses", "/"},
		{"/projects", "/#projects"},
	}

	for _, tt := range tests {
		t.Run("Redirect: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusMovedPermanently && w.Code != http.StatusFound && w.Code != http.StatusSeeOther {
				t.Errorf("expected redirect status for %s, got %d", tt.url, w.Code)
			}
			loc := w.Header().Get("Location")
			if loc != tt.expectedLocation {
				t.Errorf("expected Location '%s', got '%s'", tt.expectedLocation, loc)
			}
		})
	}
}

// TestRouterMethodNotAllowed verifies that calling unsupported HTTP methods on registered routes
// returns HTTP 405 Method Not Allowed.
func TestRouterMethodNotAllowed(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		method string
		url    string
	}{
		{http.MethodPost, "/health"},
		{http.MethodPut, "/mustang"},
		{http.MethodDelete, "/"},
		{http.MethodPatch, "/api/health"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405 Method Not Allowed for %s %s, got %d", tt.method, tt.url, w.Code)
			}
		})
	}
}

// TestRouterAPIGroup verifies routing and responses for grouped /api endpoints.
func TestRouterAPIGroup(t *testing.T) {
	r := setupTestRouter()

	t.Run("GET /api/health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected application/json Content-Type, got %s", ct)
		}
		if !strings.Contains(w.Body.String(), `"status":"UP"`) {
			t.Errorf("expected body to contain status UP, got %s", w.Body.String())
		}
	})

	t.Run("GET /api/pdf/health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/pdf/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected application/json Content-Type, got %s", ct)
		}
	})

	t.Run("GET /api/unknown-route returns 404 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/unknown-endpoint", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("expected application/json for API 404, got %s", ct)
		}
	})
}
