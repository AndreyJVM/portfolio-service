package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"ru.vorobevaqa/devops-platform/internal/model"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	staticDir := "../../web/static"
	h := NewHandler(cs, is, "../../web/templates", staticDir, "https://vorobevaqa.ru")

	r := chi.NewRouter()
	r.Use(middleware.GetHead)

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

	r.Get("/robots.txt", h.RobotsTxt)
	r.Get("/sitemap.xml", h.SitemapXML)

	r.Get("/", h.Index)
	r.Get("/courses", h.CoursesList)
	r.Get("/courses/{courseSlug}", h.CourseDetail)
	r.Get("/courses/{courseSlug}/{lessonSlug}", h.Lesson)
	r.Get("/projects", h.Projects)
	r.Get("/education", h.Education)
	r.Get("/about", h.RedirectAbout)
	r.Get("/qr", h.QRPage)
	r.Get("/mustang", h.Mustang)
	r.Get("/api/search", h.SearchAPI)
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
		{"HEAD /", "HEAD", "/", http.StatusOK},
		{"GET /robots.txt", "GET", "/robots.txt", http.StatusOK},
		{"GET /sitemap.xml", "GET", "/sitemap.xml", http.StatusOK},
		{"GET /api/search", "GET", "/api/search", http.StatusOK},
		{"GET /api/search?q=docker", "GET", "/api/search?q=docker", http.StatusOK},
		{"GET /courses", "GET", "/courses", http.StatusOK},
		{"HEAD /courses", "HEAD", "/courses", http.StatusOK},
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

func TestSearchAPI(t *testing.T) {
	r := setupTestRouter()

	// 1. Получение полного индекса поиска
	req := httptest.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	var items []model.SearchItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("Ошибка декодирования search JSON: %v", err)
	}

	if len(items) == 0 {
		t.Errorf("Ожидался непустой индекс поиска")
	}

	// 2. Серверная фильтрация по запросу "docker"
	reqQuery := httptest.NewRequest("GET", "/api/search?q=docker", nil)
	wQuery := httptest.NewRecorder()
	r.ServeHTTP(wQuery, reqQuery)

	if wQuery.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для ?q=docker, получен %d", wQuery.Code)
	}

	var queryItems []model.SearchItem
	if err := json.NewDecoder(wQuery.Body).Decode(&queryItems); err != nil {
		t.Fatalf("Ошибка декодирования отфильтрованного search JSON: %v", err)
	}

	if len(queryItems) == 0 {
		t.Errorf("Ожидались результаты для запроса 'docker'")
	}

	// Проверяем, что во всех результатах есть упоминание docker
	for _, it := range queryItems {
		found := strings.Contains(strings.ToLower(it.Title), "docker") ||
			strings.Contains(strings.ToLower(it.Description), "docker") ||
			strings.Contains(strings.ToLower(it.Snippet), "docker") ||
			strings.Contains(strings.ToLower(it.CourseTitle), "docker")
		if !found {
			t.Errorf("Элемент %s не соответствует запросу 'docker'", it.Title)
		}
	}
}

func TestSEOEndpoints(t *testing.T) {
	r := setupTestRouter()

	// 1. Тестируем robots.txt
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

	// 2. Тестируем динамический sitemap.xml
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

	// Проверяем наличие всех 4 курсов и основных страниц
	expectedLinks := []string{
		"https://vorobevaqa.ru/",
		"https://vorobevaqa.ru/courses",
		"https://vorobevaqa.ru/projects",
		"https://vorobevaqa.ru/education",
		"https://vorobevaqa.ru/qr",
		"https://vorobevaqa.ru/mustang",
		"https://vorobevaqa.ru/courses/linux",
		"https://vorobevaqa.ru/courses/docker",
		"https://vorobevaqa.ru/courses/ansible",
		"https://vorobevaqa.ru/courses/cicd",
		"https://vorobevaqa.ru/courses/linux/01-permissions-and-systemd",
		"https://vorobevaqa.ru/courses/docker/01-multistage-and-cgroups",
		"https://vorobevaqa.ru/courses/ansible/01-ansible-basics-and-playbooks",
		"https://vorobevaqa.ru/courses/cicd/01-github-actions-pipeline",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(sitemapBody, "<loc>"+link+"</loc>") {
			t.Errorf("sitemap.xml не содержит обязательную ссылку: %s", link)
		}
	}
}

func TestNotFoundHandlers(t *testing.T) {
	r := setupTestRouter()

	notFoundTests := []struct {
		name string
		url  string
	}{
		{"Non-existent course", "/courses/non-existent-course-slug-12345"},
		{"Non-existent lesson", "/courses/linux/non-existent-lesson-slug-12345"},
		{"Non-existent course in lesson path", "/courses/non-existent-course/01-permissions"},
	}

	for _, tt := range notFoundTests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("%s: ожидался статус 404, получен %d", tt.url, w.Code)
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
		{"Learning progress JS", "/js/learning-progress.js", "javascript"},
		{"Search JS", "/js/search.js", "javascript"},
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
		"themeToggle",              // Кнопка переключения темы
		"/js/theme-toggle.js",      // Скрипт переключения темы
		"/js/learning-progress.js", // Скрипт учета прогресса обучения
		"/js/search.js",            // Скрипт поиска
		"searchModal",              // Модальное окно поиска
		"searchTriggerSidebar",     // Кнопка поиска в сайдбаре
		"/css/navbar.css",          // Стили боковой и верхней панели
		"/css/themes.css",          // Стили ночной и дневной темы
		"/css/main.css",            // Основные стили
		"liveToast",                // Контейнер уведомлений Toast
		"/images/avatar.jpg",       // Аватарка в Hero
		"navbar-brand",             // Бренд платформы
		"Primary Cloud Lab",        // Название бренда
		"app-sidebar",              // Боковая панель навигации слева
		"sidebar-nav",              // Список навигационных ссылок
		">Главная<",                // Ссылка Главная без иконок
		">Курсы<",                  // Ссылка Курсы без иконок
		">Проекты<",                // Ссылка Проекты без иконок
		">Сертификаты<",            // Ссылка Сертификаты без иконок
		">QR-генератор<",           // Ссылка QR-генератор без иконок
		">Shelby<",                 // Ссылка Shelby без иконок
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(body, snippet) {
			t.Errorf("Главная страница не содержит обязательный элемент: '%s'", snippet)
		}
	}

	// Заголовок браузера строго Primary Cloud Lab
	if !strings.Contains(body, "<title>Primary Cloud Lab</title>") {
		t.Errorf("Заголовок страницы должен быть <title>Primary Cloud Lab</title>")
	}

	// Под логотипом отсутствует DEVOPS PLATFORM
	if strings.Contains(body, "DEVOPS PLATFORM") {
		t.Errorf("Подпись 'DEVOPS PLATFORM' под логотипом должна быть удалена")
	}

	// В боковом меню убрано слово Навигация
	if strings.Contains(body, "Навигация") {
		t.Errorf("Слово 'Навигация' в боковом меню должно быть удалено")
	}

	// Бейдж версии Go удален
	if strings.Contains(body, "Go 1.24") || strings.Contains(body, "Go 1.") {
		t.Errorf("Бейдж версии Go должен быть удален из бокового меню")
	}

	// Раздел 'Популярные направления' удален с главной страницы
	if strings.Contains(body, "Популярные направления") {
		t.Errorf("Раздел 'Популярные направления' должен быть удален с главной страницы")
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

func TestCleanPagesLayout(t *testing.T) {
	r := setupTestRouter()

	// 1. Проверяем, что во вкладке браузера всех страниц строго <title>Primary Cloud Lab</title>
	for _, path := range []string{"/", "/courses", "/projects", "/education", "/qr", "/mustang"} {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		body := w.Body.String()

		if !strings.Contains(body, "<title>Primary Cloud Lab</title>") {
			t.Errorf("%s: заголовок должен быть строго <title>Primary Cloud Lab</title>", path)
		}
	}

	// 2. Страница проектов не должна содержать лишнюю плашку курсов
	reqProj := httptest.NewRequest("GET", "/projects", nil)
	wProj := httptest.NewRecorder()
	r.ServeHTTP(wProj, reqProj)
	if strings.Contains(wProj.Body.String(), "Курсы DevOps") {
		t.Errorf("/projects не должна содержать плашку 'Курсы DevOps'")
	}

	// 3. Страница QR не должна содержать лишнюю плашку курсов
	reqQR := httptest.NewRequest("GET", "/qr", nil)
	wQR := httptest.NewRecorder()
	r.ServeHTTP(wQR, reqQR)
	if strings.Contains(wQR.Body.String(), "Курсы DevOps") {
		t.Errorf("/qr не должна содержать плашку 'Курсы DevOps'")
	}
}

func TestLearningProgressUI(t *testing.T) {
	r := setupTestRouter()

	// 1. Проверяем наличие контейнера прогресса в каталоге курсов
	reqCourses := httptest.NewRequest("GET", "/courses", nil)
	wCourses := httptest.NewRecorder()
	r.ServeHTTP(wCourses, reqCourses)
	if !strings.Contains(wCourses.Body.String(), "data-course-progress") {
		t.Errorf("/courses не содержит контейнеров прогресса data-course-progress")
	}

	// 2. Проверяем наличие виджета прогресса и бейджей статусов на детальной странице курса
	reqDetail := httptest.NewRequest("GET", "/courses/linux", nil)
	wDetail := httptest.NewRecorder()
	r.ServeHTTP(wDetail, reqDetail)
	tbodyDetail := wDetail.Body.String()
	if !strings.Contains(tbodyDetail, "courseDetailProgress") {
		t.Errorf("/courses/linux не содержит виджета courseDetailProgress")
	}
	if !strings.Contains(tbodyDetail, "lesson-status-badge") {
		t.Errorf("/courses/linux не содержит бейджей lesson-status-badge")
	}

	// 3. Проверяем наличие интерактивного квиза и кнопки фиксации прогресса на странице урока
	reqLesson := httptest.NewRequest("GET", "/courses/linux/01-permissions-and-systemd", nil)
	wLesson := httptest.NewRecorder()
	r.ServeHTTP(wLesson, reqLesson)
	tbodyLesson := wLesson.Body.String()
	if !strings.Contains(tbodyLesson, "lesson-quiz-card") {
		t.Errorf("Страница урока не содержит блока проверочного квиза lesson-quiz-card")
	}
	if !strings.Contains(tbodyLesson, "completeLessonBtn") {
		t.Errorf("Страница урока не содержит кнопки completeLessonBtn")
	}
	if !strings.Contains(tbodyLesson, "quiz-option-btn") {
		t.Errorf("Страница урока не содержит вариантов ответов quiz-option-btn")
	}
}

func TestQRJsonEndpoint(t *testing.T) {
	r := setupTestRouter()

	// 1. Успешный POST-запрос JSON
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

	// 2. Запрос без параметра text/url (дефолтное значение fallback на baseURL)
	emptyReq := httptest.NewRequest("GET", "/api/qr?format=json", nil)
	emptyW := httptest.NewRecorder()
	r.ServeHTTP(emptyW, emptyReq)
	if emptyW.Code != http.StatusOK {
		t.Errorf("Ожидался 200 с дефолтным url, получен %d", emptyW.Code)
	}

	// 3. Отдача сырого PNG при запросе без заголовков JSON
	pngReq := httptest.NewRequest("GET", "/api/qr?text=https://test.com", nil)
	pngW := httptest.NewRecorder()
	r.ServeHTTP(pngW, pngReq)
	if pngW.Code != http.StatusOK {
		t.Errorf("Ожидался 200 для PNG, получен %d", pngW.Code)
	}
	if pngW.Header().Get("Content-Type") != "image/png" {
		t.Errorf("Ожидался Content-Type image/png, получен %s", pngW.Header().Get("Content-Type"))
	}
}

func TestResumePDF_MissingFile(t *testing.T) {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	// Передаем несуществующую директорию статики, чтобы вызвать ветку ошибки os.Stat
	h := NewHandler(cs, is, "../../web/templates", "/tmp/non-existent-static-dir-12345", "https://vorobevaqa.ru")

	req := httptest.NewRequest("GET", "/api/pdf/resume", nil)
	w := httptest.NewRecorder()
	h.ResumePDF(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404 при отсутствии файла резюме, получен %d", w.Code)
	}
}

func TestTemplateErrors(t *testing.T) {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	// Передаем некорректную директорию шаблонов для проверки ветки ошибки компиляции
	h := NewHandler(cs, is, "/tmp/non-existent-templates-dir-12345", "../../web/static", "https://vorobevaqa.ru")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.Index(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500 при ошибке парсинга шаблона, получен %d", w.Code)
	}
}

func BenchmarkRenderIndexCached(b *testing.B) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
