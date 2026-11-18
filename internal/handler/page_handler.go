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
		Meta:    h.getMeta("index", "/"),
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
		Meta: h.getMeta("education", "/education"),
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
		Meta: h.getMeta("mustang", "/mustang"),
	}

	h.renderTemplate(w, "mustang.html", data)
}
