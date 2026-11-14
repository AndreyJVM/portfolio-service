package handler

import (
	"net/http"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// Index: Главная страница с превью обучающих треков
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
			Title:         "Главная • Платформа по Контейнеризации",
			Description:   "Специализированная база знаний и практические курсы по Docker, Podman, LXC и изоляции в ядре Linux от Андрея Воробьева.",
			Keywords:      "контейнеризация, Docker, Podman, LXC, Incus, Cgroups, Namespaces, уроки, курсы",
			OGTitle:       "Андрей Воробьев • Containerization & Infrastructure Engineer",
			OGDescription: "Глубокие практические руководства по Docker, Podman, LXC и низкоуровневой изоляции в Linux.",
			OGImage:       "/images/og-main.png",
			OGURL:         h.baseURL + "/",
			ActiveTab:     "home",
			CanonicalURL:  h.baseURL + "/",
		},
		FeaturedCourses: featured,
	}

	h.renderTemplate(w, "index.html", data)
}

// Projects: Страница проектов с интеграцией GitHub, DockerHub, Ansible Galaxy
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

// Education: Квалификация и сертификаты
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

// RoadmapRedirect: HTTP 301 Permanent Redirect на раздел дорожной карты
func (h *Handler) RoadmapRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/courses#roadmap", http.StatusMovedPermanently)
}

// Mustang: Страница сбора на Shelby GT500
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
