package handler

import (
	"net/http"

	"ru.vorobevaqa/devops-platform/internal/model"
)

// Index: Main page of Primary Cloud Lab — About the author, experience, projects, and certifications.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	projectsData := h.integrationsService.GetProjectsData()

	data := struct {
		Meta model.MetaData
		Data model.ProjectsPageData
	}{
		Meta: h.getMeta("index", "/"),
		Data: projectsData,
	}

	h.renderTemplate(w, "about.html", data)
}

// RedirectHome: HTTP 301 Permanent Redirect to the home page (About the author)
func (h *Handler) RedirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Projects: HTTP 301 Permanent Redirect to the projects section on the home page
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#projects", http.StatusMovedPermanently)
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
