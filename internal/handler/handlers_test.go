package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	h := NewHandler(cs, is, "../../web/templates", "../../web/static", "https://vorobevaqa.ru")

	r := chi.NewRouter()
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
