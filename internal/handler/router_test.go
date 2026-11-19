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
		"/courses",
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

// TestRouterMiddlewareStack tests RequestID, GetHead, and Gzip compression middleware.
func TestRouterMiddlewareStack(t *testing.T) {
	r := setupTestRouter()

	t.Run("RequestID middleware", func(t *testing.T) {
		var capturedReqID string
		testMux := chi.NewRouter()
		testMux.Use(middleware.RequestID)
		testMux.Get("/check-id", func(w http.ResponseWriter, r *http.Request) {
			capturedReqID = middleware.GetReqID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/check-id", nil)
		w := httptest.NewRecorder()
		testMux.ServeHTTP(w, req)

		if capturedReqID == "" {
			t.Errorf("expected request ID to be injected into context, got empty")
		}
	})

	t.Run("HEAD request handling", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for HEAD /health, got %d", w.Code)
		}
	})

	t.Run("Gzip compression middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		if ce := w.Header().Get("Content-Encoding"); ce == "gzip" {
			gzReader, err := gzip.NewReader(w.Body)
			if err != nil {
				t.Fatalf("failed to create gzip reader: %v", err)
			}
			defer gzReader.Close()
			decompressed, err := io.ReadAll(gzReader)
			if err != nil {
				t.Fatalf("failed to read decompressed body: %v", err)
			}
			if !strings.Contains(string(decompressed), "Primary Cloud Lab") {
				t.Errorf("decompressed response does not contain 'Primary Cloud Lab'")
			}
		}
	})

	t.Run("Panic recovery middleware", func(t *testing.T) {
		// Create a sub-router inheriting Chi recoverer middleware to test recoverer specifically
		testMux := chi.NewRouter()
		testMux.Use(chi.Chain(setupTestRouter().Middlewares()...).Handler)
		testMux.Get("/panic-test", func(w http.ResponseWriter, r *http.Request) {
			panic("intentional test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic-test", nil)
		w := httptest.NewRecorder()

		defer func() {
			if rcv := recover(); rcv != nil {
				t.Errorf("panic was not caught by Recoverer middleware: %v", rcv)
			}
		}()

		testMux.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500 on panic recovery, got %d", w.Code)
		}
	})
}

// TestRouterStaticAssetCacheHeaders tests that static assets have the expected
// Cache-Control headers and caching policies.
func TestRouterStaticAssetCacheHeaders(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		url                  string
		expectedCacheControl string
		expectedStatus       int
	}{
		{"/css/main.css", "public, max-age=86400, stale-while-revalidate=3600", http.StatusOK},
		{"/js/main.js", "public, max-age=86400, stale-while-revalidate=3600", http.StatusOK},
		{"/images/avatar.jpg", "public, max-age=604800, immutable", http.StatusOK},
		{"/pdf/resume.pdf", "public, max-age=86400", http.StatusOK},
		{"/favicon.svg", "public, max-age=604800, immutable", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run("Asset: "+tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d for %s, got %d", tt.expectedStatus, tt.url, w.Code)
			}
			cc := w.Header().Get("Cache-Control")
			if cc != tt.expectedCacheControl {
				t.Errorf("expected Cache-Control '%s', got '%s'", tt.expectedCacheControl, cc)
			}
		})
	}
}

// TestRouterDirectoryListingForbidden verifies that trailing slash requests to asset
// directories return 404 rather than listing contents.
func TestRouterDirectoryListingForbidden(t *testing.T) {
	r := setupTestRouter()

	dirEndpoints := []string{
		"/css/",
		"/js/",
		"/images/",
		"/pdf/",
	}

	for _, ep := range dirEndpoints {
		t.Run("Trailing slash: "+ep, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, ep, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("expected status 404 for directory listing '%s', got %d", ep, w.Code)
			}
		})
	}
}

// TestRouterRedirects verifies that legacy and shorthand routes redirect to their canonical targets.
func TestRouterRedirects(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		url              string
		expectedLocation string
	}{
		{"/about", "/education"},
		{"/roadmap", "/#roadmap"},
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
		{http.MethodPut, "/courses"},
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

	t.Run("GET /api/search with query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/search?q=docker", nil)
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
