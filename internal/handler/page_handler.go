package handler

import (
	"net/http"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// Index: Главная страница образовательной платформы — каталог курсов и дорожная карта обучения
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	courses := h.courseService.GetAllCourses()

	data := struct {
		Meta    model.MetaData
		Courses []model.Course
	}{
		Meta: model.MetaData{
			Title:         "Курсы • Обучение Контейнеризации и Linux",
			Description:   "Специализированная программа обучения: Docker, Podman, LXC, Cgroups, Namespaces и Zero-Downtime деплой на платформе Primary Cloud Lab.",
			Keywords:      "контейнеризация, Docker, Podman, LXC, Incus, Cgroups, Namespaces, уроки, курсы, Linux, DevOps",
			OGTitle:       "Курсы по Контейнеризации и Linux • Primary Cloud Lab",
			OGDescription: "Глубокая практическая программа: от первого контейнера до изоляции на уровне ядра Linux.",
			OGImage:       h.baseURL + "/images/og-courses.png",
			OGURL:         h.baseURL + "/",
			ActiveTab:     "courses",
			CanonicalURL:  h.baseURL + "/",
		},
		Courses: courses,
	}

	h.renderTemplate(w, "courses.html", data)
}

// Projects: постоянный редирект (HTTP 301) на секцию проектов и инфраструктуры на странице "Об авторе"
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/education#projects", http.StatusMovedPermanently)
}

// Education: Страница об авторе, опыте, проектах и квалификации
func (h *Handler) Education(w http.ResponseWriter, r *http.Request) {
	projectsData := h.integrationsService.GetProjectsData()

	data := struct {
		Meta model.MetaData
		Data model.ProjectsPageData
	}{
		Meta: model.MetaData{
			Title:         "Об авторе • Андрей Воробьев | Primary Cloud Lab",
			Description:   "Об авторе платформы Primary Cloud Lab Андрее Воробьеве: инженерный опыт, системное администрирование Linux (Слёрм), Docker, Podman, LXC, CI/CD, репозитории и квалификация.",
			Keywords:      "Андрей Воробьев, об авторе, DevOps, Linux, Слёрм, Docker, Podman, LXC, Ansible, сертификаты, квалификация, инфраструктура",
			OGTitle:       "Об авторе • Андрей Воробьев | Primary Cloud Lab",
			OGDescription: "Создатель Primary Cloud Lab. Архитектура инфраструктуры, open-source репозитории, профессиональные сертификации Слёрм и дипломы.",
			OGImage:       h.baseURL + "/images/og-about.png",
			OGURL:         h.baseURL + "/education",
			ActiveTab:     "education",
			CanonicalURL:  h.baseURL + "/education",
		},
		Data: projectsData,
	}

	h.renderTemplate(w, "education.html", data)
}

// RedirectAbout: HTTP 301 Permanent Redirect
func (h *Handler) RedirectAbout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/education", http.StatusMovedPermanently)
}

// RoadmapRedirect: HTTP 301 Permanent Redirect на раздел дорожной карты
func (h *Handler) RoadmapRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#roadmap", http.StatusMovedPermanently)
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
			OGImage:       h.baseURL + "/images/avatar.jpg",
			OGURL:         h.baseURL + "/mustang",
			ActiveTab:     "mustang",
			CanonicalURL:  h.baseURL + "/mustang",
		},
	}

	h.renderTemplate(w, "mustang.html", data)
}
