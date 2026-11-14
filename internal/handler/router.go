package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// noDirListingFS оборачивает http.FileSystem, блокируя листинг содержимого директорий.
type noDirListingFS struct {
	fs http.FileSystem
}

func (nfs noDirListingFS) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if s.IsDir() {
		_ = f.Close()
		return nil, os.ErrNotExist
	}
	return f, nil
}

// NewRouter инициализирует маршрутизатор chi со всеми middleware, защитой статики и эндпоинтами.
func NewRouter(h *Handler, staticDir string) *chi.Mux {
	r := chi.NewRouter()

	// Middleware стек
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.GetHead)
	r.Use(middleware.Timeout(15 * time.Second))

	// Заголовки безопасности
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	})

	// Безопасная раздача статических файлов (CSS, JS, Favicon, Сертификаты, Резюме)
	fileServer := http.FileServer(noDirListingFS{fs: http.Dir(staticDir)})
	for _, prefix := range []string{"/css", "/js", "/images", "/pdf"} {
		prefixPattern := prefix + "/*"
		r.Handle(prefixPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Дополнительная проверка на завершающий слэш для запрета обращения к директориям
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=86400")
			fileServer.ServeHTTP(w, r)
		}))
	}

	r.Get("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, filepath.Join(staticDir, "favicon.svg"))
	})

	// Healthcheck для Docker, docker-compose и Kubernetes/Nginx
	r.Get("/health", h.Health)
	r.Get("/api/health", h.Health)

	// Динамические SEO-эндпоинты
	r.Get("/robots.txt", h.RobotsTxt)
	r.Get("/sitemap.xml", h.SitemapXML)

	// Страницы платформы
	r.Get("/", h.Index)
	r.Get("/courses", h.CoursesList)
	r.Get("/courses/{courseSlug}", h.CourseDetail)
	r.Get("/courses/{courseSlug}/{lessonSlug}", h.Lesson)

	r.Get("/projects", h.Projects)
	r.Get("/education", h.Education)
	r.Get("/about", h.RedirectAbout)
	r.Get("/roadmap", h.RoadmapRedirect)
	r.Get("/qr", h.QRPage)
	r.Get("/mustang", h.Mustang)

	// Полнотекстовый поиск по материалам и урокам
	r.Get("/api/search", h.SearchAPI)

	// API генерации QR-кодов (Go native engine)
	r.Get("/api/qr", h.GenerateQREndpoint)
	r.Post("/api/qr", h.GenerateQREndpoint)

	// API PDF резюме
	r.Get("/api/pdf/resume", h.ResumePDF)
	r.Get("/api/pdf/health", h.PdfHealth)

	return r
}
