package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/config"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	is := service.NewIntegrationsService(config.DefaultGitHubUser, config.DefaultDockerUser)
	ps := service.NewPageService("../../content/pages.json")
	staticDir := "../../web/static"
	h := NewHandler(is, ps, "../../web/templates", staticDir, config.DefaultAppBaseURL)
	return NewRouter(h)
}

func TestHandlers(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"GET / (About the author)", "GET", "/", http.StatusOK},
		{"HEAD /", "HEAD", "/", http.StatusOK},
		{"GET /mustang", "GET", "/mustang", http.StatusOK},
		{"HEAD /mustang", "HEAD", "/mustang", http.StatusOK},
		{"GET /health", "GET", "/health", http.StatusOK},
		{"GET /api/health", "GET", "/api/health", http.StatusOK},
		{"GET /favicon.svg", "GET", "/favicon.svg", http.StatusOK},
		{"GET /css/ directory listing forbidden", "GET", "/css/", http.StatusNotFound},
		{"GET /robots.txt", "GET", "/robots.txt", http.StatusOK},
		{"GET /sitemap.xml", "GET", "/sitemap.xml", http.StatusOK},
		{"GET /about redirect to home", "GET", "/about", http.StatusMovedPermanently},
		{"GET /education redirect to home", "GET", "/education", http.StatusMovedPermanently},
		{"GET /projects redirect to anchor", "GET", "/projects", http.StatusMovedPermanently},
		{"GET /courses redirect to home", "GET", "/courses", http.StatusMovedPermanently},
		{"GET /courses/docker redirect to home", "GET", "/courses/docker", http.StatusMovedPermanently},
		{"GET /roadmap redirect to home", "GET", "/roadmap", http.StatusMovedPermanently},
		{"GET /api/pdf/resume", "GET", "/api/pdf/resume", http.StatusOK},
		{"GET /api/pdf/health", "GET", "/api/pdf/health", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: ожидался статус %d, получен %d", tt.url, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRedirects(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		url              string
		expectedLocation string
	}{
		{"/about", "/"},
		{"/education", "/"},
		{"/projects", "/#projects"},
		{"/courses", "/"},
		{"/courses/docker", "/"},
		{"/roadmap", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusMovedPermanently {
				t.Fatalf("Ожидался статус 301 для %s, получен %d", tt.url, w.Code)
			}
			if loc := w.Header().Get("Location"); loc != tt.expectedLocation {
				t.Errorf("Для %s ожидался Location '%s', получен '%s'", tt.url, tt.expectedLocation, loc)
			}
		})
	}
}

func TestSEOEndpoints(t *testing.T) {
	r := setupTestRouter()

	// 1. Test robots.txt
	reqRobots := httptest.NewRequest("GET", "/robots.txt", nil)
	wRobots := httptest.NewRecorder()
	r.ServeHTTP(wRobots, reqRobots)

	if wRobots.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для /robots.txt, получен %d", wRobots.Code)
	}
	if wRobots.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("Неверный Content-Type для robots.txt: %s", wRobots.Header().Get("Content-Type"))
	}
	robotsBody := wRobots.Body.String()
	if !strings.Contains(robotsBody, "Sitemap: https://vorobevaqa.ru/sitemap.xml") {
		t.Errorf("robots.txt не содержит правильную ссылку на sitemap.xml: %s", robotsBody)
	}

	// 2. Test dynamic sitemap.xml
	reqSitemap := httptest.NewRequest("GET", "/sitemap.xml", nil)
	wSitemap := httptest.NewRecorder()
	r.ServeHTTP(wSitemap, reqSitemap)

	if wSitemap.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для /sitemap.xml, получен %d", wSitemap.Code)
	}
	if !strings.Contains(wSitemap.Header().Get("Content-Type"), "application/xml") {
		t.Errorf("Неверный Content-Type для sitemap.xml: %s", wSitemap.Header().Get("Content-Type"))
	}

	sitemapBody := wSitemap.Body.String()
	expectedLinks := []string{
		"https://vorobevaqa.ru/",
		"https://vorobevaqa.ru/mustang",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(sitemapBody, "<loc>"+link+"</loc>") {
			t.Errorf("sitemap.xml не содержит обязательную ссылку: %s", link)
		}
	}
}

func TestStaticAssets(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name                string
		url                 string
		expectedContentType string
	}{
		{"Avatar image", "/images/avatar.jpg", "image/jpeg"},
		{"Slurm certificate", "/images/certificates/slurm-linux.jpg", "image/jpeg"},
		{"Docker certificate", "/images/certificates/stepik-docker-ansible.jpg", "image/jpeg"},
		{"Favicon SVG", "/favicon.svg", "image/svg+xml"},
		{"Main CSS", "/css/main.css", "text/css"},
		{"Navbar CSS", "/css/navbar.css", "text/css"},
		{"Themes CSS", "/css/themes.css", "text/css"},
		{"About CSS", "/css/pages/about.css", "text/css"},
		{"Mustang CSS", "/css/pages/mustang.css", "text/css"},
		{"Main JS", "/js/main.js", "javascript"},
		{"About JS", "/js/pages/about.js", "javascript"},
		{"Mustang JS", "/js/pages/mustang.js", "javascript"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("%s: ожидался статус 200, получен %d", tt.url, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if !strings.Contains(contentType, tt.expectedContentType) {
				t.Errorf("%s: Content-Type '%s' не содержит '%s'", tt.url, contentType, tt.expectedContentType)
			}
		})
	}
}

func TestLayoutIncludes(t *testing.T) {
	r := setupTestRouter()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	body := w.Body.String()

	requiredSnippets := []string{
		"data-bs-theme=\"dark\"", // Permanent dark theme attribute
		"data-theme=\"dark\"",    // Custom dark theme attribute
		"/css/navbar.css",        // Sidebar and top navbar styles
		"/css/themes.css",        // Permanent focused dark theme styles
		"/css/main.css",          // Main styles
		"liveToast",              // Toast notifications container
		"navbar-brand",           // Platform brand
		"Primary Cloud Lab",      // Brand name
		"app-sidebar",            // Left navigation sidebar
		"sidebar-nav",            // Navigation links list
		"Об авторе",              // About author link
		"На мечту",               // "На мечту" link
		"/images/avatar.jpg",     // Author avatar
		"id=\"projects\"",        // Projects anchor
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(body, snippet) {
			t.Errorf("Главная страница не содержит обязательный элемент: '%s'", snippet)
		}
	}

	// Browser title is strictly Primary Cloud Lab
	if !strings.Contains(body, "<title>Primary Cloud Lab</title>") {
		t.Errorf("Заголовок страницы должен быть <title>Primary Cloud Lab</title>")
	}

	// OpenGraph and Twitter Card meta tags
	if !strings.Contains(body, `property="og:site_name" content="Primary Cloud Lab"`) {
		t.Errorf("Страница не содержит тег og:site_name")
	}
	if !strings.Contains(body, `name="twitter:card" content="summary_large_image"`) {
		t.Errorf("Страница не содержит тег twitter:card")
	}
}

func TestResumePDF_MissingFile(t *testing.T) {
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	// Pass non-existent static directory to trigger os.Stat error branch
	h := NewHandler(is, nil, "../../web/templates", "/tmp/non-existent-static-dir-12345", "https://vorobevaqa.ru")

	req := httptest.NewRequest("GET", "/api/pdf/resume", nil)
	w := httptest.NewRecorder()
	h.ResumePDF(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404 при отсутствии файла резюме, получен %d", w.Code)
	}
}

func TestResumePDF_DownloadAttachment(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/api/pdf/resume?download=true", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	disposition := w.Header().Get("Content-Disposition")
	if !strings.HasPrefix(disposition, "attachment;") {
		t.Errorf("Ожидался заголовок Content-Disposition c attachment, получено: %s", disposition)
	}
}
