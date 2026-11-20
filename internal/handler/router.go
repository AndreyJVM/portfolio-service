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
func NewRouter(h *Handler) *chi.Mux {
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

	// Custom 404 & 405 error handlers (JSON for API, branded HTML for web requests)
	r.NotFound(h.NotFound)
	r.MethodNotAllowed(h.MethodNotAllowed)

	// Secure static file serving (CSS, JS, Favicon, Certificates, Resume)
	fileServer := http.FileServer(noDirListingFS{fs: http.Dir(h.staticDir)})

	// Static assets with caching policies
	for _, prefix := range []string{"/css", "/js"} {
		prefixPattern := prefix + "/*"
		r.Handle(prefixPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				h.NotFound(w, r)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=3600")
			fileServer.ServeHTTP(w, r)
		}))
	}

	r.Handle("/images/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			h.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		fileServer.ServeHTTP(w, r)
	}))

	r.Handle("/pdf/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			h.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fileServer.ServeHTTP(w, r)
	}))

	r.Get("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		http.ServeFile(w, r, filepath.Join(h.staticDir, "favicon.svg"))
	})

	// Top-level health check (for Docker, K8s, Nginx)
	r.Get("/health", h.Health)

	// Dynamic SEO endpoints
	r.Get("/robots.txt", h.RobotsTxt)
	r.Get("/sitemap.xml", h.SitemapXML)

	// Platform web pages
	r.Get("/", h.Index)
	r.Get("/about", h.RedirectHome)
	r.Get("/education", h.RedirectHome)
	r.Get("/projects", h.Projects)
	r.Get("/courses", h.RedirectHome)
	r.Get("/courses/*", h.RedirectHome)
	r.Get("/roadmap", h.RedirectHome)
	r.Get("/mustang", h.Mustang)

	// Grouped API routes
	r.Route("/api", func(api chi.Router) {
		api.Get("/health", h.Health)

		// PDF resume API
		api.Route("/pdf", func(pdf chi.Router) {
			pdf.Get("/resume", h.ResumePDF)
			pdf.Get("/health", h.PdfHealth)
		})
	})

	return r
}
