package handler

import (
	"net/http"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// Index: Главная страница образовательной платформы и проектов инфраструктуры
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	courses := h.courseService.GetAllCourses()
	projectsData := h.integrationsService.GetProjectsData()

	data := struct {
		Meta    model.MetaData
		Courses []model.Course
		Data    model.ProjectsPageData
	}{
		Meta: model.MetaData{
			Title:         "Главная • Платформа по Контейнеризации",
			Description:   "Специализированная база знаний, курсы по Docker, Podman, LXC, ядру Linux и открытая инфраструктура от Андрея Воробьева.",
			Keywords:      "контейнеризация, Docker, Podman, LXC, Incus, Cgroups, Namespaces, уроки, курсы, инфраструктура, DevOps",
			OGTitle:       "Андрей Воробьев • Primary Cloud Lab | Обучение & Инфраструктура",
			OGDescription: "Глубокие практические руководства по Docker, Podman, LXC и защищенная инфраструктура на Linux.",
			OGImage:       h.baseURL + "/images/og-main.png",
			OGURL:         h.baseURL + "/",
			ActiveTab:     "home",
			CanonicalURL:  h.baseURL + "/",
		},
		Courses: courses,
		Data:    projectsData,
	}

	h.renderTemplate(w, "index.html", data)
}

// Projects: постоянный редирект (HTTP 301) на секцию проектов главной страницы
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#projects", http.StatusMovedPermanently)
}

// Education: Квалификация и сертификаты
func (h *Handler) Education(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Meta model.MetaData
	}{
		Meta: model.MetaData{
			Title:         "Квалификация и Сертификации",
			Description:   "Квалификация Андрея Воробьева: системное администрирование Linux (Слёрм), Docker, Podman, LXC, CI/CD, диплом гос. образца.",
			Keywords:      "DevOps, Linux, Слёрм, Docker, Podman, LXC, Ansible, сертификаты, квалификация",
			OGTitle:       "Квалификация и Сертификации | Андрей Воробьев",
			OGDescription: "Дипломы государственного образца, профессиональные сертификации Слёрм и подтвержденные инженерные компетенции.",
			OGImage:       h.baseURL + "/images/og-about.png",
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
			OGImage:       h.baseURL + "/images/avatar.jpg",
			OGURL:         h.baseURL + "/mustang",
			ActiveTab:     "mustang",
			CanonicalURL:  h.baseURL + "/mustang",
		},
	}

	h.renderTemplate(w, "mustang.html", data)
}
