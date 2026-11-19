package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/config"
	"ru.vorobevaqa/devops-platform/internal/service"
)

func setupTestRouter() *chi.Mux {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService(config.DefaultGitHubUser, config.DefaultDockerUser)
	ps := service.NewPageService("../../content/pages.json")
	staticDir := "../../web/static"
	h := NewHandler(cs, is, ps, "../../web/templates", staticDir, config.DefaultAppBaseURL)
	return NewRouter(h)
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
		{"GET /courses/linux", "GET", "/courses/linux", http.StatusOK},
		{"GET /courses/linux/01-permissions-and-systemd", "GET", "/courses/linux/01-permissions-and-systemd", http.StatusOK},
		{"GET /courses/linux/02-networking-and-firewalls", "GET", "/courses/linux/02-networking-and-firewalls", http.StatusOK},
		{"GET /courses/linux/03-processes-and-system-limits", "GET", "/courses/linux/03-processes-and-system-limits", http.StatusOK},
		{"GET /courses/linux/04-bash-automation-and-scripting", "GET", "/courses/linux/04-bash-automation-and-scripting", http.StatusOK},
		{"GET /courses/linux/05-disks-filesystems-and-lvm", "GET", "/courses/linux/05-disks-filesystems-and-lvm", http.StatusOK},
		{"GET /courses/docker", "GET", "/courses/docker", http.StatusOK},
		{"GET /courses/docker/01-architecture-and-first-dockerfile", "GET", "/courses/docker/01-architecture-and-first-dockerfile", http.StatusOK},
		{"GET /courses/docker/02-multistage-and-cgroups", "GET", "/courses/docker/02-multistage-and-cgroups", http.StatusOK},
		{"GET /courses/docker/03-docker-compose-and-networking", "GET", "/courses/docker/03-docker-compose-and-networking", http.StatusOK},
		{"GET /courses/docker/04-volumes-and-data-persistence", "GET", "/courses/docker/04-volumes-and-data-persistence", http.StatusOK},
		{"GET /courses/docker/05-container-security-and-hardening", "GET", "/courses/docker/05-container-security-and-hardening", http.StatusOK},
		{"GET /courses/docker/06-logging-diagnostics-and-daemon-tuning", "GET", "/courses/docker/06-logging-diagnostics-and-daemon-tuning", http.StatusOK},
		{"GET /courses/docker/07-multiarch-buildx-and-cicd", "GET", "/courses/docker/07-multiarch-buildx-and-cicd", http.StatusOK},
		{"GET /courses/docker/08-storage-drivers-and-cow", "GET", "/courses/docker/08-storage-drivers-and-cow", http.StatusOK},
		{"GET /courses/podman", "GET", "/courses/podman", http.StatusOK},
		{"GET /courses/podman/01-podman-rootless-and-daemonless", "GET", "/courses/podman/01-podman-rootless-and-daemonless", http.StatusOK},
		{"GET /courses/podman/02-podman-quadlet-and-systemd", "GET", "/courses/podman/02-podman-quadlet-and-systemd", http.StatusOK},
		{"GET /courses/podman/03-podman-pods-and-kube-generate", "GET", "/courses/podman/03-podman-pods-and-kube-generate", http.StatusOK},
		{"GET /courses/podman/04-podman-storage-volumes-and-userns", "GET", "/courses/podman/04-podman-storage-volumes-and-userns", http.StatusOK},
		{"GET /courses/podman/05-podman-compose-and-production-security", "GET", "/courses/podman/05-podman-compose-and-production-security", http.StatusOK},
		{"GET /courses/lxc", "GET", "/courses/lxc", http.StatusOK},
		{"GET /courses/lxc/01-lxc-and-system-containers", "GET", "/courses/lxc/01-lxc-and-system-containers", http.StatusOK},
		{"GET /courses/lxc/02-storage-pools-zfs-and-snapshots", "GET", "/courses/lxc/02-storage-pools-zfs-and-snapshots", http.StatusOK},
		{"GET /courses/lxc/03-networking-bridge-macvlan-and-proxy", "GET", "/courses/lxc/03-networking-bridge-macvlan-and-proxy", http.StatusOK},
		{"GET /courses/lxc/04-profiles-cgroups-and-autostart", "GET", "/courses/lxc/04-profiles-cgroups-and-autostart", http.StatusOK},
		{"GET /courses/lxc/05-security-unprivileged-and-docker-nesting", "GET", "/courses/lxc/05-security-unprivileged-and-docker-nesting", http.StatusOK},
		{"GET /courses/internals", "GET", "/courses/internals", http.StatusOK},
		{"GET /courses/internals/01-namespaces-and-cgroups-v2", "GET", "/courses/internals/01-namespaces-and-cgroups-v2", http.StatusOK},
		{"GET /courses/internals/02-syscalls-clone-unshare-and-pivot-root", "GET", "/courses/internals/02-syscalls-clone-unshare-and-pivot-root", http.StatusOK},
		{"GET /courses/internals/03-building-mini-container-runtime-in-go", "GET", "/courses/internals/03-building-mini-container-runtime-in-go", http.StatusOK},
		{"GET /courses/internals/04-oci-spec-runc-and-crun", "GET", "/courses/internals/04-oci-spec-runc-and-crun", http.StatusOK},
		{"GET /courses/internals/05-seccomp-capabilities-and-ebpf-tracing", "GET", "/courses/internals/05-seccomp-capabilities-and-ebpf-tracing", http.StatusOK},
		{"GET /courses/ansible", "GET", "/courses/ansible", http.StatusOK},
		{"GET /courses/ansible/01-ansible-basics-and-playbooks", "GET", "/courses/ansible/01-ansible-basics-and-playbooks", http.StatusOK},
		{"GET /courses/ansible/02-roles-and-galaxy", "GET", "/courses/ansible/02-roles-and-galaxy", http.StatusOK},
		{"GET /courses/cicd", "GET", "/courses/cicd", http.StatusOK},
		{"GET /courses/cicd/01-github-actions-pipeline", "GET", "/courses/cicd/01-github-actions-pipeline", http.StatusOK},
		{"GET /courses/cicd/02-zero-downtime-and-blue-green", "GET", "/courses/cicd/02-zero-downtime-and-blue-green", http.StatusOK},
		{"GET /projects redirect", "GET", "/projects", http.StatusMovedPermanently},
		{"GET /education", "GET", "/education", http.StatusOK},
		{"GET /about redirect", "GET", "/about", http.StatusMovedPermanently},
		{"GET /roadmap redirect", "GET", "/roadmap", http.StatusMovedPermanently},
		{"GET /mustang", "GET", "/mustang", http.StatusOK},
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

func TestSEOEndpoints(t *testing.T) {
	r := setupTestRouter()

	// 1. Test robots.txt
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

	// 2. Test dynamic sitemap.xml
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

	// Check presence of Docker Library modules and main pages
	expectedLinks := []string{
		"https://vorobevaqa.ru/",
		"https://vorobevaqa.ru/courses",
		"https://vorobevaqa.ru/education",
		"https://vorobevaqa.ru/mustang",
		"https://vorobevaqa.ru/courses/docker",
		"https://vorobevaqa.ru/courses/docker-dockerfile",
		"https://vorobevaqa.ru/courses/docker-storage",
		"https://vorobevaqa.ru/courses/docker-networking",
		"https://vorobevaqa.ru/courses/docker-compose",
		"https://vorobevaqa.ru/courses/docker-security",
		"https://vorobevaqa.ru/courses/docker-resources",
		"https://vorobevaqa.ru/courses/docker-production",
		"https://vorobevaqa.ru/courses/docker-cicd",
		"https://vorobevaqa.ru/courses/docker-cookbook",
		"https://vorobevaqa.ru/courses/docker/01-architecture-and-first-dockerfile",
		"https://vorobevaqa.ru/courses/docker/02-multistage-and-cgroups",
		"https://vorobevaqa.ru/courses/docker-dockerfile/01-dockerfile-instructions-and-layers",
		"https://vorobevaqa.ru/courses/docker-storage/01-volumes-and-data-persistence",
		"https://vorobevaqa.ru/courses/docker-networking/01-network-drivers-and-dns",
		"https://vorobevaqa.ru/courses/docker-compose/01-compose-v2-specification",
		"https://vorobevaqa.ru/courses/docker-security/01-container-security-and-hardening",
		"https://vorobevaqa.ru/courses/docker-resources/01-cgroups-v2-memory-and-cpu-limits",
		"https://vorobevaqa.ru/courses/docker-production/01-logging-diagnostics-and-daemon-tuning",
		"https://vorobevaqa.ru/courses/docker-cicd/01-multiarch-buildx-and-cicd",
		"https://vorobevaqa.ru/courses/docker-cookbook/01-docker-cli-cheatsheet",
	}

	for _, link := range expectedLinks {
		if !strings.Contains(sitemapBody, "<loc>"+link+"</loc>") {
			t.Errorf("sitemap.xml не содержит обязательную ссылку: %s", link)
		}
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
		{"Lesson CSS", "/css/pages/lesson.css", "text/css"},
		{"Theme toggle JS", "/js/theme-toggle.js", "javascript"},
		{"Learning progress JS", "/js/learning-progress.js", "javascript"},
		{"Search JS", "/js/search.js", "javascript"},
		{"Main JS", "/js/main.js", "javascript"},
		{"About JS", "/js/pages/about.js", "javascript"},
		{"Mustang JS", "/js/pages/mustang.js", "javascript"},
		{"Lesson JS", "/js/pages/lesson.js", "javascript"},
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
		"data-bs-theme=\"dark\"",     // Permanent dark theme attribute
		"data-theme=\"dark\"",        // Custom dark theme attribute
		"/js/learning-progress.js",   // Learning progress tracking script
		"/js/search.js",              // Search script
		"searchModal",                // Search modal
		"searchTriggerSidebar",       // Sidebar search button
		"/css/navbar.css",            // Sidebar and top navbar styles
		"/css/themes.css",            // Permanent focused dark theme styles
		"/css/main.css",              // Main styles
		"liveToast",                  // Toast notifications container
		"navbar-brand",               // Platform brand
		"Primary Cloud Lab",          // Brand name
		"app-sidebar",                // Left navigation sidebar
		"sidebar-nav",                // Navigation links list
		">Docker: Библиотека<",       // Docker Library link without icons (main page)
		"sidebar-subnav",             // Docker Library modules submenu
		"/courses/docker",            // Quick link: Architecture & Fundamentals
		"/courses/docker-dockerfile", // Quick link: Dockerfile & BuildKit
		"/courses/docker-storage",    // Quick link: Storage & Volumes
		"/courses/docker-networking", // Quick link: Networking & DNS
		"/courses/docker-compose",    // Quick link: Compose v2
		"/courses/docker-security",   // Quick link: Security
		"/courses/docker-resources",  // Quick link: Cgroups & Resources
		"/courses/docker-production", // Quick link: Production & Troubleshooting
		"/courses/docker-cicd",       // Quick link: CI/CD & Buildx
		"/courses/docker-cookbook",   // Quick link: Cookbook
		">Об авторе<",                // About link without icons
		">На мечту<",                 // "На мечту" link without icons
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(body, snippet) {
			t.Errorf("Главная страница не содержит обязательный элемент: '%s'", snippet)
		}
	}

	// Check active course highlighting in submenu when navigating to course page
	reqCourse := httptest.NewRequest("GET", "/courses/docker", nil)
	wCourse := httptest.NewRecorder()
	r.ServeHTTP(wCourse, reqCourse)
	if !strings.Contains(wCourse.Body.String(), `href="/courses/docker"`) || !strings.Contains(wCourse.Body.String(), `active`) {
		t.Errorf("Страница /courses/docker должна содержать активную ссылку на курс")
	}

	// About page contains avatar in Hero section
	reqEdu := httptest.NewRequest("GET", "/education", nil)
	wEdu := httptest.NewRecorder()
	r.ServeHTTP(wEdu, reqEdu)
	if !strings.Contains(wEdu.Body.String(), "/images/avatar.jpg") {
		t.Errorf("Страница 'Об авторе' должна содержать '/images/avatar.jpg'")
	}

	// Browser title is strictly Primary Cloud Lab
	if !strings.Contains(body, "<title>Primary Cloud Lab</title>") {
		t.Errorf("Заголовок страницы должен быть <title>Primary Cloud Lab</title>")
	}

	// OpenGraph and Twitter Card meta tags
	if !strings.Contains(body, `property="og:site_name" content="Primary Cloud Lab"`) {
		t.Errorf("Страница не содержит тег og:site_name")
	}
	if !strings.Contains(body, `name="twitter:card" content="summary_large_image"`) {
		t.Errorf("Страница не содержит тег twitter:card")
	}

	// DEVOPS PLATFORM is removed from below the logo
	if strings.Contains(body, "DEVOPS PLATFORM") {
		t.Errorf("Подпись 'DEVOPS PLATFORM' под логотипом должна быть удалена")
	}

	// The word Navigation is removed from the sidebar menu
	if strings.Contains(body, "Навигация") {
		t.Errorf("Слово 'Навигация' в боковом меню должно быть удалено")
	}

	// Go version badge is removed
	if strings.Contains(body, "Go 1.24") || strings.Contains(body, "Go 1.") {
		t.Errorf("Бейдж версии Go должен быть удален из бокового меню")
	}

	// 'Popular topics' section is removed from the main page
	if strings.Contains(body, "Популярные направления") {
		t.Errorf("Раздел 'Популярные направления' должен быть удален с главной страницы")
	}

	// Verify that sidebar-nav block has no <i> icon tags (minimalist link design)
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

	// 1. Verify that all pages have strictly <title>Primary Cloud Lab</title> in browser tab
	for _, path := range []string{"/", "/courses", "/education", "/mustang"} {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		body := w.Body.String()

		if !strings.Contains(body, "<title>Primary Cloud Lab</title>") {
			t.Errorf("%s: заголовок должен быть строго <title>Primary Cloud Lab</title>", path)
		}
	}

	// 2. Projects page redirects to the projects section of About page (301)
	reqProj := httptest.NewRequest("GET", "/projects", nil)
	wProj := httptest.NewRecorder()
	r.ServeHTTP(wProj, reqProj)
	if wProj.Code != http.StatusMovedPermanently {
		t.Errorf("Ожидался статус 301 для /projects, получен %d", wProj.Code)
	}
	if loc := wProj.Header().Get("Location"); loc != "/education#projects" {
		t.Errorf("Ожидался редирект на /education#projects, получен %s", loc)
	}
}

func TestLearningProgressUI(t *testing.T) {
	r := setupTestRouter()

	// 1. Check for presence of progress container and roadmap in course catalog
	reqCourses := httptest.NewRequest("GET", "/courses", nil)
	wCourses := httptest.NewRecorder()
	r.ServeHTTP(wCourses, reqCourses)
	coursesBody := wCourses.Body.String()
	if !strings.Contains(coursesBody, "data-course-progress") {
		t.Errorf("/courses не содержит контейнеров прогресса data-course-progress")
	}
	if !strings.Contains(coursesBody, "roadmap-timeline") {
		t.Errorf("/courses не содержит контейнера roadmap-timeline")
	}
	if !strings.Contains(coursesBody, "Дорожная карта") {
		t.Errorf("/courses не содержит раздела Дорожная карта")
	}

	// 2. Check for presence of progress widget and status badges on course detail page
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

	// 3. Check for presence of content and progress toggle button on lesson page
	lessonPaths := []string{
		"/courses/docker/01-architecture-and-first-dockerfile",
		"/courses/docker/02-multistage-and-cgroups",
		"/courses/docker/03-docker-compose-and-networking",
		"/courses/docker/04-volumes-and-data-persistence",
		"/courses/docker/05-container-security-and-hardening",
		"/courses/docker/06-logging-diagnostics-and-daemon-tuning",
		"/courses/docker/07-multiarch-buildx-and-cicd",
		"/courses/lxc/01-lxc-and-system-containers",
		"/courses/lxc/02-storage-pools-zfs-and-snapshots",
		"/courses/lxc/03-networking-bridge-macvlan-and-proxy",
		"/courses/lxc/04-profiles-cgroups-and-autostart",
		"/courses/lxc/05-security-unprivileged-and-docker-nesting",
		"/courses/podman/01-podman-rootless-and-daemonless",
		"/courses/podman/02-podman-quadlet-and-systemd",
		"/courses/podman/03-podman-pods-and-kube-generate",
		"/courses/podman/04-podman-storage-volumes-and-userns",
		"/courses/podman/05-podman-compose-and-production-security",
		"/courses/internals/01-namespaces-and-cgroups-v2",
		"/courses/internals/02-syscalls-clone-unshare-and-pivot-root",
		"/courses/internals/03-building-mini-container-runtime-in-go",
		"/courses/internals/04-oci-spec-runc-and-crun",
		"/courses/internals/05-seccomp-capabilities-and-ebpf-tracing",
		"/courses/linux/01-permissions-and-systemd",
		"/courses/ansible/01-ansible-basics-and-playbooks",
		"/courses/cicd/01-github-actions-pipeline",
	}

	for _, path := range lessonPaths {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Урок %s вернул статус %d, ожидался 200", path, w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "completeLessonBtn") {
			t.Errorf("Страница урока %s не содержит кнопки completeLessonBtn", path)
		}
		if !strings.Contains(body, "markdown-content") {
			t.Errorf("Страница урока %s не содержит контейнера markdown-content", path)
		}
	}
}

func TestLesson_PrevNextNavigation(t *testing.T) {
	r := setupTestRouter()

	// First lesson: must only contain "Next lesson"
	reqFirst := httptest.NewRequest("GET", "/courses/docker/01-architecture-and-first-dockerfile", nil)
	wFirst := httptest.NewRecorder()
	r.ServeHTTP(wFirst, reqFirst)
	bodyFirst := wFirst.Body.String()

	if !strings.Contains(bodyFirst, "Следующий урок") {
		t.Errorf("Первый урок курса должен содержать ссылку на следующий урок")
	}
	if strings.Contains(bodyFirst, "Предыдущий урок") {
		t.Errorf("Первый урок курса не должен содержать ссылку на предыдущий урок")
	}

	// Second lesson: must contain both "Previous lesson" and "Next lesson"
	reqMid := httptest.NewRequest("GET", "/courses/docker/02-multistage-and-cgroups", nil)
	wMid := httptest.NewRecorder()
	r.ServeHTTP(wMid, reqMid)
	bodyMid := wMid.Body.String()

	if !strings.Contains(bodyMid, "Следующий урок") {
		t.Errorf("Промежуточный урок должен содержать ссылку на следующий урок")
	}
	if !strings.Contains(bodyMid, "Предыдущий урок") {
		t.Errorf("Промежуточный урок должен содержать ссылку на предыдущий урок")
	}
}

func TestIndexPage_ContainsCoursesAndProjects(t *testing.T) {
	r := setupTestRouter()

	// 1. Main page (/) is "Courses"
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для /, получен %d", w.Code)
	}

	body := w.Body.String()
	expectedCoursesSubstrings := []string{
		"id=\"catalog\"",
		"id=\"guide\"",
		"id=\"roadmap\"",
		"Каталог модулей Docker",
		"Дорожная карта",
		"data-course-progress",
		"Архитектура: Ядро Linux, пространства имен и стек OCI",
		"Dockerfile: Оптимизация слоев, Multi-stage и BuildKit",
	}

	for _, sub := range expectedCoursesSubstrings {
		if !strings.Contains(body, sub) {
			t.Errorf("Главная страница (Курсы) должна содержать '%s'", sub)
		}
	}

	// 2. About page (/education) — contains profile, architecture, and open-source projects
	reqEdu := httptest.NewRequest("GET", "/education", nil)
	wEdu := httptest.NewRecorder()
	r.ServeHTTP(wEdu, reqEdu)

	if wEdu.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для /education, получен %d", wEdu.Code)
	}

	bodyEdu := wEdu.Body.String()
	expectedAuthorSubstrings := []string{
		"id=\"projects\"",
		"id=\"docker-activity\"",
		"id=\"ansible-activity\"",
		"id=\"github-activity\"",
		"Docker Hub",
		"GitHub",
		"Ansible",
		"/images/avatar.jpg",
	}

	for _, sub := range expectedAuthorSubstrings {
		if !strings.Contains(bodyEdu, sub) {
			t.Errorf("Страница 'Об авторе' должна содержать '%s'", sub)
		}
	}

	// Infrastructure architecture block is removed from "About" tab
	if strings.Contains(bodyEdu, "Архитектура Инфраструктуры Платформы") {
		t.Errorf("Блок 'Архитектура Инфраструктуры Платформы' должен быть удален с вкладки 'Об авторе'")
	}
}

func TestResumePDF_MissingFile(t *testing.T) {
	cs := service.NewCourseService("../../content/courses")
	is := service.NewIntegrationsService("AndreyJVM", "andreyvorobevaqa")
	// Pass non-existent static directory to trigger os.Stat error branch
	h := NewHandler(cs, is, nil, "../../web/templates", "/tmp/non-existent-static-dir-12345", "https://vorobevaqa.ru")

	req := httptest.NewRequest("GET", "/api/pdf/resume", nil)
	w := httptest.NewRecorder()
	h.ResumePDF(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404 при отсутствии файла резюме, получен %d", w.Code)
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
