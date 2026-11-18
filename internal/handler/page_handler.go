package handler

import (
	"net/http"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// Index: Main page of the learning platform — course catalog and learning roadmap
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
			OGImage:       h.baseURL + "/images/og-main.png",
			OGURL:         h.baseURL + "/",
			ActiveTab:     "courses",
			CanonicalURL:  h.baseURL + "/",
		},
		Courses: courses,
	}

	h.renderTemplate(w, "courses.html", data)
}

// Projects: permanent redirect (HTTP 301) to the projects and infrastructure section on the "About" page
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/education#projects", http.StatusMovedPermanently)
}

// Education: Page about the author, experience, projects, and qualifications
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

// RoadmapRedirect: HTTP 301 Permanent Redirect to the roadmap section
func (h *Handler) RoadmapRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#roadmap", http.StatusMovedPermanently)
}

// Mustang: Crowdfunding page for Shelby GT500
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
