package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"ru.vorobevaqa/devops-platform/internal/service"
)

type Handler struct {
	courseService       *service.CourseService
	integrationsService *service.IntegrationsService
	templatesDir        string
	staticDir           string
	baseURL             string

	templateMu    sync.RWMutex
	templateCache map[string]*template.Template
}

func NewHandler(
	cs *service.CourseService,
	is *service.IntegrationsService,
	templatesDir string,
	staticDir string,
	baseURL string,
) *Handler {
	return &Handler{
		courseService:       cs,
		integrationsService: is,
		templatesDir:        templatesDir,
		staticDir:           staticDir,
		baseURL:             baseURL,
		templateCache:       make(map[string]*template.Template),
	}
}

// getTemplate возвращает скомпилированный шаблон из потокобезопасного кэша памяти,
// исключая дисковый ввод-вывод (I/O) на каждом HTTP-запросе.
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

// renderTemplate выполняет рендеринг кэшированного шаблона.
func (h *Handler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, err := h.getTemplate(templateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parse error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
	}
}

// Health возвращает статус здоровья сервиса для healthcheck (Docker, Nginx, K8s).
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP","service":"portfolio-service"}`))
}
