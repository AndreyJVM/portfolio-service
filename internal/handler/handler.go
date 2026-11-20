package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"ru.vorobevaqa/devops-platform/internal/model"
	"ru.vorobevaqa/devops-platform/internal/service"
)

type Handler struct {
	integrationsService *service.IntegrationsService
	pageService         *service.PageService
	templatesDir        string
	staticDir           string
	baseURL             string

	templateMu    sync.RWMutex
	templateCache map[string]*template.Template
}

func NewHandler(
	is *service.IntegrationsService,
	ps *service.PageService,
	templatesDir string,
	staticDir string,
	baseURL string,
) *Handler {
	if ps == nil {
		ps = service.NewPageService("")
	}
	return &Handler{
		integrationsService: is,
		pageService:         ps,
		templatesDir:        templatesDir,
		staticDir:           staticDir,
		baseURL:             baseURL,
		templateCache:       make(map[string]*template.Template),
	}
}

// getMeta builds model.MetaData for a page using PageService and the base URL.
func (h *Handler) getMeta(pageKey, path string) model.MetaData {
	if h.pageService == nil {
		return model.MetaData{}
	}
	return h.pageService.GetMetaData(pageKey, path, h.baseURL)
}

// getTemplate returns a compiled template from a thread-safe in-memory cache,
// eliminating disk I/O on every HTTP request.
func (h *Handler) getTemplate(templateName string) (*template.Template, error) {
	h.templateMu.RLock()
	if tmpl, ok := h.templateCache[templateName]; ok {
		h.templateMu.RUnlock()
		return tmpl, nil
	}
	h.templateMu.RUnlock()

	h.templateMu.Lock()
	defer h.templateMu.Unlock()

	// Double-Checked Locking
	if tmpl, ok := h.templateCache[templateName]; ok {
		return tmpl, nil
	}

	layoutPath := filepath.Join(h.templatesDir, "layout.html")
	targetPath := filepath.Join(h.templatesDir, templateName)

	parsedTmpl, err := template.ParseFiles(layoutPath, targetPath)
	if err != nil {
		return nil, err
	}

	h.templateCache[templateName] = parsedTmpl
	return parsedTmpl, nil
}

// renderTemplate renders a cached template via a buffer,
// preventing partial delivery of malformed HTML on execution errors.
func (h *Handler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, err := h.getTemplate(templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parse error: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

// Health checks the service status.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP","service":"portfolio-service"}`))
}

// NotFound handles 404 errors, returning branded HTML for browsers and JSON for APIs.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found","message":"Requested API endpoint does not exist"}`))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	h.renderTemplate(w, "404.html", struct {
		Meta model.MetaData
	}{
		Meta: h.getMeta("404", r.URL.Path),
	})
}

// MethodNotAllowed handles 405 errors, returning JSON for API requests and standard text for web.
func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"method_not_allowed","message":"HTTP method is not supported for this endpoint"}`))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("405 Method Not Allowed\n"))
}
