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

// noDirListingFS wraps http.FileSystem, blocking directory listing.
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

// NewRouter initializes the chi router with all middleware, static file protection, and endpoints.
func NewRouter(h *Handler, staticDir string) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.GetHead)
	r.Use(middleware.Timeout(15 * time.Second))

	// Security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	})

	// Secure static file serving (CSS, JS, Favicon, Certificates, Resume)
	fileServer := http.FileServer(noDirListingFS{fs: http.Dir(staticDir)})
	for _, prefix := range []string{"/css", "/js", "/images", "/pdf"} {
		prefixPattern := prefix + "/*"
		r.Handle(prefixPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Additional trailing slash check to prevent access to directories
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

	// Health check for Docker, docker-compose, and Kubernetes/Nginx
	r.Get("/health", h.Health)
	r.Get("/api/health", h.Health)

	// Dynamic SEO endpoints
	r.Get("/robots.txt", h.RobotsTxt)
	r.Get("/sitemap.xml", h.SitemapXML)

	// Platform pages
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

	// Full-text search across materials and lessons
	r.Get("/api/search", h.SearchAPI)

	// QR code generation API (Go native engine)
	r.Get("/api/qr", h.GenerateQREndpoint)
	r.Post("/api/qr", h.GenerateQREndpoint)

	// PDF resume API
	r.Get("/api/pdf/resume", h.ResumePDF)
	r.Get("/api/pdf/health", h.PdfHealth)

	return r
}
