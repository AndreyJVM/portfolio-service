package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	qrcode "github.com/skip2/go-qrcode"
	"ru.vorobevaqa/devops-platform/internal/model"
	"ru.vorobevaqa/devops-platform/internal/service"
)

type Handler struct {
	courseService       *service.CourseService
	integrationsService *service.IntegrationsService
	templatesDir        string
	staticDir           string
	baseURL             string
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
	}
}

// renderTemplate загружает базовый layout и целевой шаблон.
func (h *Handler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	layoutPath := filepath.Join(h.templatesDir, "layout.html")
	targetPath := filepath.Join(h.templatesDir, templateName)

	tmpl, err := template.ParseFiles(layoutPath, targetPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parse error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
	}
}

// IndexHandler: Главная страница с превью обучающих треков
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	courses := h.courseService.GetAllCourses()
	var featured []model.Course
	if len(courses) > 2 {
		featured = courses[:2]
	} else {
		featured = courses
	}

	data := struct {
		Meta            model.MetaData
		FeaturedCourses []model.Course
	}{
		Meta: model.MetaData{
			Title:         "Главная • Обучающая платформа DevOps",
			Description:   "Бесплатная платформа по системному администрированию Linux, Docker, Ansible и CI/CD от Андрея Воробьева.",
			Keywords:      "DevOps, Linux, Docker, Ansible, CI/CD, Go, Golang, уроки, курсы",
			OGTitle:       "Андрей Воробьев • DevOps Platform & Engineer",
			OGDescription: "Практические материалы и курсы по Linux, Docker, Ansible и автоматизации релизов.",
			OGImage:       "/images/og-main.png",
			OGURL:         h.baseURL + "/",
			ActiveTab:     "home",
			CanonicalURL:  h.baseURL + "/",
		},
		FeaturedCourses: featured,
	}

	h.renderTemplate(w, "index.html", data)
}

// CoursesListHandler: Каталог обучающих треков
func (h *Handler) CoursesList(w http.ResponseWriter, r *http.Request) {
	courses := h.courseService.GetAllCourses()

	data := struct {
		Meta    model.MetaData
		Courses []model.Course
	}{
		Meta: model.MetaData{
			Title:         "Курсы и Треки обучения DevOps",
			Description:   "Бесплатные структурированные уроки по Linux, Docker, Kubernetes, Ansible и CI/CD.",
			Keywords:      "курсы devops, linux обучение, docker с нуля, ansible уроки, ci cd github actions",
			OGTitle:       "Курсы DevOps & Linux Infrastructure",
			OGDescription: "Изучайте администрирование и автоматизацию инфраструктуры бесплатно на практике.",
			OGImage:       "/images/og-main.png",
			OGURL:         h.baseURL + "/courses",
			ActiveTab:     "courses",
			CanonicalURL:  h.baseURL + "/courses",
		},
		Courses: courses,
	}

	h.renderTemplate(w, "courses.html", data)
}

// CourseDetailHandler: Страница конкретного курса со списком уроков
func (h *Handler) CourseDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "courseSlug")
	course, err := h.courseService.GetCourseBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Meta   model.MetaData
		Course *model.Course
	}{
		Meta: model.MetaData{
			Title:         course.Title + " • Курс DevOps",
			Description:   course.Description,
			Keywords:      fmt.Sprintf("%s, devops, уроки, обучение", course.Title),
			OGTitle:       course.Title,
			OGDescription: course.Description,
			OGImage:       "/images/og-main.png",
			OGURL:         fmt.Sprintf("%s/courses/%s", h.baseURL, slug),
			ActiveTab:     "courses",
			CanonicalURL:  fmt.Sprintf("%s/courses/%s", h.baseURL, slug),
		},
		Course: course,
	}

	h.renderTemplate(w, "course_detail.html", data)
}

// LessonHandler: Страница просмотра и изучения отдельного урока
func (h *Handler) Lesson(w http.ResponseWriter, r *http.Request) {
	courseSlug := chi.URLParam(r, "courseSlug")
	lessonSlug := chi.URLParam(r, "lessonSlug")

	course, lesson, err := h.courseService.GetLesson(courseSlug, lessonSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Meta              model.MetaData
		Course            *model.Course
		Lesson            *model.Lesson
		LessonContentSafe template.HTML
	}{
		Meta: model.MetaData{
			Title:         lesson.Title + " • " + course.Title,
			Description:   lesson.Description,
			Keywords:      fmt.Sprintf("%s, %s, devops", lesson.Title, course.Title),
			OGTitle:       lesson.Title,
			OGDescription: lesson.Description,
			OGImage:       "/images/og-main.png",
			OGURL:         fmt.Sprintf("%s/courses/%s/%s", h.baseURL, courseSlug, lessonSlug),
			ActiveTab:     "courses",
			CanonicalURL:  fmt.Sprintf("%s/courses/%s/%s", h.baseURL, courseSlug, lessonSlug),
		},
		Course:            course,
		Lesson:            lesson,
		LessonContentSafe: template.HTML(lesson.ContentHTML),
	}

	h.renderTemplate(w, "lesson.html", data)
}

// ProjectsHandler: Страница проектов с интеграцией GitHub, DockerHub, Ansible Galaxy
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	projectsData := h.integrationsService.GetProjectsData()

	data := struct {
		Meta model.MetaData
		Data model.ProjectsPageData
	}{
		Meta: model.MetaData{
			Title:         "Проекты и Инфраструктура",
			Description:   "Архитектура продакшн-окружения: Linux VDS, Docker Compose, Nginx Reverse Proxy, Ansible и CI/CD.",
			Keywords:      "проекты, архитектура, Linux, Docker, Ansible, Nginx, CI/CD",
			OGTitle:       "Проекты и Инфраструктура | Андрей Воробьев",
			OGDescription: "Стек развертывания, Docker-контейнеры, Ansible-роли и открытый исходный код.",
			OGImage:       "/images/og-projects.png",
			OGURL:         h.baseURL + "/projects",
			ActiveTab:     "projects",
			CanonicalURL:  h.baseURL + "/projects",
		},
		Data: projectsData,
	}

	h.renderTemplate(w, "projects.html", data)
}

// EducationHandler: Квалификация и сертификаты
func (h *Handler) Education(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Meta model.MetaData
	}{
		Meta: model.MetaData{
			Title:         "Квалификация и Сертификации",
			Description:   "Квалификация Андрея Воробьева: системное администрирование Linux (Слёрм), Docker, Ansible, CI/CD.",
			Keywords:      "DevOps, Linux, Слёрм, Docker, Ansible, сертификаты",
			OGTitle:       "Квалификация и Сертификации | Андрей Воробьев",
			OGDescription: "Дипломы государственного образца и профессиональные сертификации.",
			OGImage:       "/images/og-about.png",
			OGURL:         h.baseURL + "/education",
			ActiveTab:     "education",
			CanonicalURL:  h.baseURL + "/education",
		},
	}

	h.renderTemplate(w, "education.html", data)
}

// RedirectAbout: HTTP 301 Permanent Redirect
func (h *Handler) RedirectAbout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/education", http.StatusMovedPermanently)
}

// QRPageHandler: QR Playground страница
func (h *Handler) QRPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Meta model.MetaData
	}{
		Meta: model.MetaData{
			Title:         "QR Code Playground",
			Description:   "Генератор QR-кодов на легковесном нативном Go движке.",
			Keywords:      "QR code generator, golang, qrcode",
			OGTitle:       "QR Code Playground",
			OGDescription: "Интерактивный генератор QR-кодов на Go.",
			OGImage:       "/images/og-qr.png",
			OGURL:         h.baseURL + "/qr",
			ActiveTab:     "qr",
			CanonicalURL:  h.baseURL + "/qr",
		},
	}

	h.renderTemplate(w, "qr.html", data)
}

// MustangHandler: Страница сбора на Shelby GT500
func (h *Handler) Mustang(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Meta model.MetaData
	}{
		Meta: model.MetaData{
			Title:         "Ford Mustang Shelby GT500 • Мечта DevOps-инженера",
			Description:   "Фонд спасения от рутины: краудфандинг на Ford Mustang Shelby GT500.",
			Keywords:      "Ford Mustang, Shelby GT500, краудфандинг, V8",
			OGTitle:       "Ford Mustang Shelby GT500 • Фонд мечты инженера",
			OGDescription: "Осталось закрыть единственный критический баг — отсутствие 760 л.с. в гараже.",
			OGImage:       "/favicon.svg",
			OGURL:         h.baseURL + "/mustang",
			ActiveTab:     "mustang",
			CanonicalURL:  h.baseURL + "/mustang",
		},
	}

	h.renderTemplate(w, "mustang.html", data)
}

// GenerateQREndpoint: Native Go QR Code generation PNG or JSON
func (h *Handler) GenerateQREndpoint(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		text = r.URL.Query().Get("url")
	}

	if r.Method == http.MethodPost {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.URL != "" {
			text = body.URL
		}
	}

	if text == "" {
		text = h.baseURL
	}

	pngBytes, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Не удалось сгенерировать QR-код", http.StatusInternalServerError)
		return
	}

	acceptsJSON := strings.Contains(r.Header.Get("Accept"), "application/json") ||
		strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
		r.Method == http.MethodPost ||
		r.URL.Query().Get("format") == "json"

	if acceptsJSON {
		base64Str := base64.StdEncoding.EncodeToString(pngBytes)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"qrCode": base64Str,
			"url":    text,
			"size":   256,
			"format": "PNG",
		})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngBytes)
}

// ResumePDF отдает PDF-файл резюме
func (h *Handler) ResumePDF(w http.ResponseWriter, r *http.Request) {
	pdfPath := filepath.Join(h.staticDir, "pdf", "resume.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		http.Error(w, "Файл резюме не найден", http.StatusNotFound)
		return
	}

	download := r.URL.Query().Get("download") == "true"
	disposition := "inline"
	if download {
		disposition = "attachment"
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\"", disposition))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, pdfPath)
}

// PdfHealth проверяет статус PDF сервиса
func (h *Handler) PdfHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP","service":"pdf-service"}`))
}
