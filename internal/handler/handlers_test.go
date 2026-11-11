package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	staticDir := "../../web/static"
	h := NewHandler(cs, is, "../../web/templates", staticDir, "https://vorobevaqa.ru")

	r := chi.NewRouter()

	// Статические файлы
	staticFS := http.FileServer(http.Dir(staticDir))
	for _, prefix := range []string{"/css", "/js", "/images", "/pdf"} {
		prefixPattern := prefix + "/*"
		r.Handle(prefixPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=86400")
			staticFS.ServeHTTP(w, r)
		}))
	}

	r.Get("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/favicon.svg")
	})

	r.Get("/", h.Index)
	r.Get("/courses", h.CoursesList)
	r.Get("/courses/{courseSlug}", h.CourseDetail)
	r.Get("/courses/{courseSlug}/{lessonSlug}", h.Lesson)
	r.Get("/projects", h.Projects)
	r.Get("/education", h.Education)
	r.Get("/about", h.RedirectAbout)
	r.Get("/qr", h.QRPage)
	r.Get("/mustang", h.Mustang)
	r.Get("/api/qr", h.GenerateQREndpoint)
	r.Post("/api/qr", h.GenerateQREndpoint)
	r.Get("/api/pdf/resume", h.ResumePDF)
	r.Get("/api/pdf/health", h.PdfHealth)
	return r
}

func TestHandlers(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"GET /", "GET", "/", http.StatusOK},
		{"GET /courses", "GET", "/courses", http.StatusOK},
		{"GET /courses/linux", "GET", "/courses/linux", http.StatusOK},
		{"GET /courses/linux/01-permissions-and-systemd", "GET", "/courses/linux/01-permissions-and-systemd", http.StatusOK},
		{"GET /courses/docker", "GET", "/courses/docker", http.StatusOK},
		{"GET /courses/docker/01-multistage-and-cgroups", "GET", "/courses/docker/01-multistage-and-cgroups", http.StatusOK},
		{"GET /courses/ansible", "GET", "/courses/ansible", http.StatusOK},
		{"GET /courses/ansible/01-ansible-basics-and-playbooks", "GET", "/courses/ansible/01-ansible-basics-and-playbooks", http.StatusOK},
		{"GET /courses/cicd", "GET", "/courses/cicd", http.StatusOK},
		{"GET /courses/cicd/01-github-actions-pipeline", "GET", "/courses/cicd/01-github-actions-pipeline", http.StatusOK},
		{"GET /projects", "GET", "/projects", http.StatusOK},
		{"GET /education", "GET", "/education", http.StatusOK},
		{"GET /about redirect", "GET", "/about", http.StatusMovedPermanently},
		{"GET /qr", "GET", "/qr", http.StatusOK},
		{"GET /mustang", "GET", "/mustang", http.StatusOK},
		{"GET /api/qr", "GET", "/api/qr?text=https://test.com", http.StatusOK},
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
		{"Theme toggle JS", "/js/theme-toggle.js", "javascript"},
		{"Main JS", "/js/main.js", "javascript"},
		{"About JS", "/js/pages/about.js", "javascript"},
		{"Mustang JS", "/js/pages/mustang.js", "javascript"},
		{"QR JS", "/js/pages/qr.js", "javascript"},
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
		"themeToggle",         // Кнопка переключения темы
		"/js/theme-toggle.js", // Скрипт переключения темы
		"/css/navbar.css",     // Стили боковой и верхней панели
		"/css/themes.css",     // Стили ночной и дневной темы
		"/css/main.css",       // Основные стили
		"liveToast",           // Контейнер уведомлений Toast
		"/images/avatar.jpg",  // Аватарка в Hero
		"navbar-brand",        // Бренд платформы
		"Primary Cloud Lab",   // Название бренда
		"app-sidebar",         // Боковая панель навигации слева
		"sidebar-nav",         // Список навигационных ссылок
		">Главная<",           // Ссылка Главная без иконок
		">Курсы<",             // Ссылка Курсы без иконок
		">Проекты<",           // Ссылка Проекты без иконок
		">Сертификаты<",       // Ссылка Сертификаты без иконок
		">QR-генератор<",      // Ссылка QR-генератор без иконок
		">Shelby<",            // Ссылка Shelby без иконок
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(body, snippet) {
			t.Errorf("Главная страница не содержит обязательный элемент: '%s'", snippet)
		}
	}

	// Проверяем, что в блоке sidebar-nav отсутствуют значки <i> (минималистичный дизайн ссылок)
	sidebarNavIdx := strings.Index(body, `sidebar-nav`)
	if sidebarNavIdx != -1 {
		sidebarNavEnd := strings.Index(body[sidebarNavIdx:], `</ul>`)
		if sidebarNavEnd != -1 {
			navHtml := body[sidebarNavIdx : sidebarNavIdx+sidebarNavEnd]
			if strings.Contains(navHtml, "<i ") {
				t.Errorf("Навигация в боковой панели не должна содержать тегов иконок <i>: %s", navHtml)
			}
		}
	}
}

func TestQRJsonEndpoint(t *testing.T) {
	r := setupTestRouter()

	body := `{"url":"https://vorobevaqa.ru"}`
	req := httptest.NewRequest("POST", "/api/qr", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	var resp struct {
		QRCode string `json:"qrCode"`
		URL    string `json:"url"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Ошибка парсинга JSON: %v", err)
	}
	if resp.QRCode == "" {
		t.Errorf("Ожидался base64 qrCode, получена пустая строка")
	}
	if resp.URL != "https://vorobevaqa.ru" {
		t.Errorf("Ожидался url 'https://vorobevaqa.ru', получен '%s'", resp.URL)
	}
}
