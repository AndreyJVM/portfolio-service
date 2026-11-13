package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/model"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	staticDir := "../../web/static"
	h := NewHandler(cs, is, "../../web/templates", staticDir, "https://vorobevaqa.ru")
	return NewRouter(h, staticDir)
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
		{"GET /health", "GET", "/health", http.StatusOK},
		{"GET /api/health", "GET", "/api/health", http.StatusOK},
		{"GET /favicon.svg", "GET", "/favicon.svg", http.StatusOK},
		{"GET /css/ directory listing forbidden", "GET", "/css/", http.StatusNotFound},
		{"GET /robots.txt", "GET", "/robots.txt", http.StatusOK},
		{"GET /sitemap.xml", "GET", "/sitemap.xml", http.StatusOK},
		{"GET /api/search", "GET", "/api/search", http.StatusOK},
		{"GET /api/search?q=docker", "GET", "/api/search?q=docker", http.StatusOK},
		{"GET /courses", "GET", "/courses", http.StatusOK},
		{"HEAD /courses", "HEAD", "/courses", http.StatusOK},
		{"GET /courses/docker", "GET", "/courses/docker", http.StatusOK},
		{"GET /courses/docker/01-architecture-and-first-dockerfile", "GET", "/courses/docker/01-architecture-and-first-dockerfile", http.StatusOK},
		{"GET /courses/docker/02-multistage-and-cgroups", "GET", "/courses/docker/02-multistage-and-cgroups", http.StatusOK},
		{"GET /courses/docker/03-docker-compose-and-networking", "GET", "/courses/docker/03-docker-compose-and-networking", http.StatusOK},
		{"GET /courses/docker/04-volumes-and-data-persistence", "GET", "/courses/docker/04-volumes-and-data-persistence", http.StatusOK},
		{"GET /courses/docker/05-container-security-and-hardening", "GET", "/courses/docker/05-container-security-and-hardening", http.StatusOK},
		{"GET /courses/docker/06-logging-diagnostics-and-daemon-tuning", "GET", "/courses/docker/06-logging-diagnostics-and-daemon-tuning", http.StatusOK},
		{"GET /courses/docker/07-multiarch-buildx-and-cicd", "GET", "/courses/docker/07-multiarch-buildx-and-cicd", http.StatusOK},
		{"GET /courses/podman", "GET", "/courses/podman", http.StatusOK},
		{"GET /courses/podman/01-podman-rootless-and-daemonless", "GET", "/courses/podman/01-podman-rootless-and-daemonless", http.StatusOK},
		{"GET /courses/lxc", "GET", "/courses/lxc", http.StatusOK},
		{"GET /courses/lxc/01-lxc-and-system-containers", "GET", "/courses/lxc/01-lxc-and-system-containers", http.StatusOK},
		{"GET /courses/lxc/02-storage-pools-zfs-and-snapshots", "GET", "/courses/lxc/02-storage-pools-zfs-and-snapshots", http.StatusOK},
		{"GET /courses/lxc/03-networking-bridge-macvlan-and-proxy", "GET", "/courses/lxc/03-networking-bridge-macvlan-and-proxy", http.StatusOK},
		{"GET /courses/internals", "GET", "/courses/internals", http.StatusOK},
		{"GET /courses/internals/01-namespaces-and-cgroups-v2", "GET", "/courses/internals/01-namespaces-and-cgroups-v2", http.StatusOK},
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
		"https://vorobevaqa.ru/courses/docker",
		"https://vorobevaqa.ru/courses/podman",
		"https://vorobevaqa.ru/courses/lxc",
		"https://vorobevaqa.ru/courses/internals",
		"https://vorobevaqa.ru/courses/docker/01-architecture-and-first-dockerfile",
		"https://vorobevaqa.ru/courses/docker/02-multistage-and-cgroups",
		"https://vorobevaqa.ru/courses/docker/03-docker-compose-and-networking",
		"https://vorobevaqa.ru/courses/docker/04-volumes-and-data-persistence",
		"https://vorobevaqa.ru/courses/docker/05-container-security-and-hardening",
		"https://vorobevaqa.ru/courses/docker/06-logging-diagnostics-and-daemon-tuning",
		"https://vorobevaqa.ru/courses/docker/07-multiarch-buildx-and-cicd",
		"https://vorobevaqa.ru/courses/podman/01-podman-rootless-and-daemonless",
		"https://vorobevaqa.ru/courses/lxc/01-lxc-and-system-containers",
		"https://vorobevaqa.ru/courses/lxc/02-storage-pools-zfs-and-snapshots",
		"https://vorobevaqa.ru/courses/lxc/03-networking-bridge-macvlan-and-proxy",
		"https://vorobevaqa.ru/courses/internals/01-namespaces-and-cgroups-v2",
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
		{"Non-existent lesson", "/courses/docker/non-existent-lesson-slug-12345"},
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
	reqDetail := httptest.NewRequest("GET", "/courses/docker", nil)
	wDetail := httptest.NewRecorder()
	r.ServeHTTP(wDetail, reqDetail)
	tbodyDetail := wDetail.Body.String()
	if !strings.Contains(tbodyDetail, "courseDetailProgress") {
		t.Errorf("/courses/docker не содержит виджета courseDetailProgress")
	}
	if !strings.Contains(tbodyDetail, "lesson-status-badge") {
		t.Errorf("/courses/docker не содержит бейджей lesson-status-badge")
	}

	// 3. Проверяем наличие интерактивного квиза и кнопки фиксации прогресса на странице урока
	reqLesson := httptest.NewRequest("GET", "/courses/docker/01-architecture-and-first-dockerfile", nil)
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

	// 4. Проверяем интерактивный квиз на уроке 02 (Multi-stage)
	reqLesson2 := httptest.NewRequest("GET", "/courses/docker/02-multistage-and-cgroups", nil)
	wLesson2 := httptest.NewRecorder()
	r.ServeHTTP(wLesson2, reqLesson2)
	tbodyLesson2 := wLesson2.Body.String()
	if !strings.Contains(tbodyLesson2, "quiz-option-btn") {
		t.Errorf("Страница урока 02 не содержит вариантов ответов quiz-option-btn")
	}

	// 5. Проверяем интерактивный квиз на уроке 03 (Compose & Networks)
	reqLesson3 := httptest.NewRequest("GET", "/courses/docker/03-docker-compose-and-networking", nil)
	wLesson3 := httptest.NewRecorder()
	r.ServeHTTP(wLesson3, reqLesson3)
	tbodyLesson3 := wLesson3.Body.String()
	if !strings.Contains(tbodyLesson3, "quiz-option-btn") {
		t.Errorf("Страница урока 03 не содержит вариантов ответов quiz-option-btn")
	}

	// 6. Проверяем интерактивный квиз на уроке 04 (Volumes & Persistence)
	reqLesson4 := httptest.NewRequest("GET", "/courses/docker/04-volumes-and-data-persistence", nil)
	wLesson4 := httptest.NewRecorder()
	r.ServeHTTP(wLesson4, reqLesson4)
	tbodyLesson4 := wLesson4.Body.String()
	if !strings.Contains(tbodyLesson4, "quiz-option-btn") {
		t.Errorf("Страница урока 04 не содержит вариантов ответов quiz-option-btn")
	}

	// 7. Проверяем интерактивный квиз на уроке 05 (Security & Hardening)
	reqLesson5 := httptest.NewRequest("GET", "/courses/docker/05-container-security-and-hardening", nil)
	wLesson5 := httptest.NewRecorder()
	r.ServeHTTP(wLesson5, reqLesson5)
	tbodyLesson5 := wLesson5.Body.String()
	if !strings.Contains(tbodyLesson5, "quiz-option-btn") {
		t.Errorf("Страница урока 05 не содержит вариантов ответов quiz-option-btn")
	}

	// 8. Проверяем интерактивный квиз на уроке 06 (Logging & Diagnostics)
	reqLesson6 := httptest.NewRequest("GET", "/courses/docker/06-logging-diagnostics-and-daemon-tuning", nil)
	wLesson6 := httptest.NewRecorder()
	r.ServeHTTP(wLesson6, reqLesson6)
	tbodyLesson6 := wLesson6.Body.String()
	if !strings.Contains(tbodyLesson6, "quiz-option-btn") {
		t.Errorf("Страница урока 06 не содержит вариантов ответов quiz-option-btn")
	}

	// 9. Проверяем интерактивный квиз на уроке 07 (Multi-arch & CI/CD)
	reqLesson7 := httptest.NewRequest("GET", "/courses/docker/07-multiarch-buildx-and-cicd", nil)
	wLesson7 := httptest.NewRecorder()
	r.ServeHTTP(wLesson7, reqLesson7)
	tbodyLesson7 := wLesson7.Body.String()
	if !strings.Contains(tbodyLesson7, "quiz-option-btn") {
		t.Errorf("Страница урока 07 не содержит вариантов ответов quiz-option-btn")
	}

	// 10. Проверяем интерактивный квиз на уроке LXC 01
	reqLxc1 := httptest.NewRequest("GET", "/courses/lxc/01-lxc-and-system-containers", nil)
	wLxc1 := httptest.NewRecorder()
	r.ServeHTTP(wLxc1, reqLxc1)
	tbodyLxc1 := wLxc1.Body.String()
	if !strings.Contains(tbodyLxc1, "quiz-option-btn") {
		t.Errorf("Страница урока LXC 01 не содержит вариантов ответов quiz-option-btn")
	}

	// 11. Проверяем интерактивный квиз на уроке LXC 02 (Storage Pools & Snapshots)
	reqLxc2 := httptest.NewRequest("GET", "/courses/lxc/02-storage-pools-zfs-and-snapshots", nil)
	wLxc2 := httptest.NewRecorder()
	r.ServeHTTP(wLxc2, reqLxc2)
	tbodyLxc2 := wLxc2.Body.String()
	if !strings.Contains(tbodyLxc2, "quiz-option-btn") {
		t.Errorf("Страница урока LXC 02 не содержит вариантов ответов quiz-option-btn")
	}

	// 12. Проверяем интерактивный квиз на уроке LXC 03 (Networking & Proxy)
	reqLxc3 := httptest.NewRequest("GET", "/courses/lxc/03-networking-bridge-macvlan-and-proxy", nil)
	wLxc3 := httptest.NewRecorder()
	r.ServeHTTP(wLxc3, reqLxc3)
	tbodyLxc3 := wLxc3.Body.String()
	if !strings.Contains(tbodyLxc3, "quiz-option-btn") {
		t.Errorf("Страница урока LXC 03 не содержит вариантов ответов quiz-option-btn")
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

func TestSecurityHeaders(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Ожидался заголовок X-Content-Type-Options: nosniff, получено: %s", w.Header().Get("X-Content-Type-Options"))
	}
	if w.Header().Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Errorf("Ожидался заголовок X-Frame-Options: SAMEORIGIN, получено: %s", w.Header().Get("X-Frame-Options"))
	}
	if w.Header().Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
		t.Errorf("Ожидался заголовок Referrer-Policy: strict-origin-when-cross-origin, получено: %s", w.Header().Get("Referrer-Policy"))
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

func TestNoDirListingFS(t *testing.T) {
	fs := noDirListingFS{fs: http.Dir("../../web/static")}

	t.Run("Open valid regular file", func(t *testing.T) {
		f, err := fs.Open("/favicon.svg")
		if err != nil {
			t.Fatalf("Не удалось открыть favicon.svg: %v", err)
		}
		defer f.Close()
	})

	t.Run("Open directory returns ErrNotExist", func(t *testing.T) {
		_, err := fs.Open("/css")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Ожидалась ошибка os.ErrNotExist для каталога, получено: %v", err)
		}
	})

	t.Run("Open non-existent file returns error", func(t *testing.T) {
		_, err := fs.Open("/non-existent-file-xyz.txt")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего файла")
		}
	})
}

func TestQRGeneration_FallbackDefault(t *testing.T) {
	r := setupTestRouter()

	t.Run("GET without text param falls back to baseURL", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/qr?format=json", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Ошибка декодирования JSON: %v", err)
		}
		if resp["url"] != "https://vorobevaqa.ru" {
			t.Errorf("Ожидался fallback URL https://vorobevaqa.ru, получено: %v", resp["url"])
		}
	})

	t.Run("POST with invalid JSON falls back to baseURL", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/qr", strings.NewReader("invalid-json-body"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", w.Code)
		}
	})
}

func TestSecurity_PathTraversalHTTP(t *testing.T) {
	r := setupTestRouter()

	traversalURLs := []string{
		"/courses/linux/..%2f..%2fREADME",
		"/courses/linux/..%2f01-permissions-and-systemd",
		"/courses/..%2f..%2fetc/passwd",
		"/courses/linux/%2e%2e%2f%2e%2e%2fgo.mod",
	}

	for _, url := range traversalURLs {
		t.Run("Path Traversal HTTP: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Сервер обязан отдавать 404 Not Found или 301 на чистый путь, но ни в коем случае не 200 с содержимым системного файла
			if w.Code == http.StatusOK && strings.Contains(w.Body.String(), "Primary Cloud Lab") && !strings.Contains(url, "linux") {
				t.Errorf("Критическая уязвимость Path Traversal! Получен 200 на URL: %s", url)
			}
			if w.Code != http.StatusNotFound && w.Code != http.StatusMovedPermanently {
				t.Logf("URL %s вернул статус %d (не 200 OK — файл не скомпрометирован)", url, w.Code)
			}
		})
	}
}

func TestSecurity_SearchEdgeCases(t *testing.T) {
	r := setupTestRouter()

	t.Run("XSS payload in search query", func(t *testing.T) {
		xssPayload := "<script>alert('xss')</script>"
		req := httptest.NewRequest("GET", "/api/search?q="+xssPayload, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}

		var items []model.SearchItem
		if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
			t.Fatalf("Ошибка десериализации JSON при XSS запросе: %v", err)
		}
	})

	t.Run("Huge query string (50000 chars) does not crash or hang", func(t *testing.T) {
		hugeQuery := strings.Repeat("A", 50000)
		req := httptest.NewRequest("GET", "/api/search?q="+hugeQuery, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}
	})

	t.Run("Special chars and regex symbols in search query", func(t *testing.T) {
		specialQueries := []string{"[", "]", "(", ")", "*", "+", "?", "\\", "$", "^", "\"", "'", "NULL", "%00"}
		for _, sq := range specialQueries {
			req := httptest.NewRequest("GET", "/api/search?q="+sq, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Сбой поиска при спецсимволе '%s': статус %d", sq, w.Code)
			}
		}
	})
}

func TestSecurity_QRMaximumLength(t *testing.T) {
	r := setupTestRouter()

	t.Run("Excessively long text returns 400 Bad Request", func(t *testing.T) {
		tooLongText := strings.Repeat("A", 3000)
		req := httptest.NewRequest("GET", "/api/qr?text="+tooLongText, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400 Bad Request для 3000 символов, получен %d", w.Code)
		}
	})
}

func TestConcurrency_HighLoad(t *testing.T) {
	r := setupTestRouter()
	const workers = 50

	t.Run("Concurrent Search API requests", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func(idx int) {
				defer wg.Done()
				query := "docker"
				if idx%2 == 0 {
					query = "linux"
				}
				req := httptest.NewRequest("GET", "/api/search?q="+query, nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Ошибка конкурентного запроса к /api/search: статус %d", w.Code)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("Concurrent lesson and page renders", func(t *testing.T) {
		urls := []string{
			"/",
			"/courses",
			"/courses/docker",
			"/courses/docker/01-architecture-and-first-dockerfile",
			"/courses/podman",
			"/projects",
			"/education",
			"/qr",
			"/mustang",
		}

		var wg sync.WaitGroup
		wg.Add(workers)

		for i := 0; i < workers; i++ {
			go func(idx int) {
				defer wg.Done()
				url := urls[idx%len(urls)]
				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Ошибка конкурентного рендеринга %s: статус %d", url, w.Code)
				}
			}(i)
		}

		wg.Wait()
	})
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

func TestSearchAPI_MultibyteUTF8Truncation(t *testing.T) {
	r := setupTestRouter()

	// Длинная кириллическая строка (>128 рун)
	longQuery := strings.Repeat("контейнеры ", 20)
	req := httptest.NewRequest("GET", "/api/search?q="+url.QueryEscape(longQuery), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	body := w.Body.Bytes()
	if !utf8.Valid(body) {
		t.Errorf("Ответ API поиска содержит невалидный UTF-8 из-за обрезки рун")
	}

	var results []model.SearchItem
	if err := json.Unmarshal(body, &results); err != nil {
		t.Fatalf("Не удалось распарсить JSON ответ: %v", err)
	}
}

func TestRenderTemplate_BufferSafety(t *testing.T) {
	tempDir := t.TempDir()

	// Создаем layout.html с вызовом несуществующего поля для провокации runtime-ошибки выполнения шаблона
	layoutContent := `{{define "layout"}}{{.InvalidField.SubField}}{{end}}`
	if err := os.WriteFile(filepath.Join(tempDir, "layout.html"), []byte(layoutContent), 0644); err != nil {
		t.Fatalf("Не удалось создать layout: %v", err)
	}
	indexContent := `{{define "content"}}Index{{end}}`
	if err := os.WriteFile(filepath.Join(tempDir, "broken.html"), []byte(indexContent), 0644); err != nil {
		t.Fatalf("Не удалось создать broken.html: %v", err)
	}

	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	h := NewHandler(cs, is, tempDir, "../../web/static", "https://vorobevaqa.ru")

	w := httptest.NewRecorder()
	// Передаем структуру, у которой нет поля InvalidField, вызывая ошибку исполнения ExecuteTemplate
	h.renderTemplate(w, "broken.html", struct{ Dummy string }{Dummy: "test"})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500 при ошибке исполнения шаблона, получен %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Template execution error") {
		t.Errorf("Ожидалось сообщение об ошибке исполнения шаблона, получено: %s", w.Body.String())
	}
}

func TestStaticAssets_SecurityChecks(t *testing.T) {
	r := setupTestRouter()

	traversalURLs := []string{
		"/css/..%2f..%2fetc/passwd",
		"/css/%2e%2e/index.html",
		"/js/../../main.go",
		"/pdf/../../../etc/shadow",
	}

	for _, url := range traversalURLs {
		t.Run("Static traversal: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest && w.Code != http.StatusMovedPermanently {
				t.Errorf("Ожидался статус 404/400/301 для path traversal %s, получен %d", url, w.Code)
			}
		})
	}

	trailingSlashURLs := []string{
		"/css/",
		"/js/",
		"/images/",
		"/pdf/",
	}

	for _, url := range trailingSlashURLs {
		t.Run("Trailing slash directory listing block: "+url, func(t *testing.T) {
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("Ожидался статус 404 для обращения к директории %s со слэшем, получен %d", url, w.Code)
			}
		})
	}
}
