package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type sitemapURL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// SitemapXML генерирует динамическую карту сайта sitemap.xml со всеми статическими
// разделами, курсами и уроками базы знаний.
func (h *Handler) SitemapXML(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")

	urlSet := sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []sitemapURL{
			{Loc: h.baseURL + "/", LastMod: today, ChangeFreq: "weekly", Priority: "1.0"},
			{Loc: h.baseURL + "/courses", LastMod: today, ChangeFreq: "weekly", Priority: "0.9"},
			{Loc: h.baseURL + "/projects", LastMod: today, ChangeFreq: "weekly", Priority: "0.8"},
			{Loc: h.baseURL + "/education", LastMod: today, ChangeFreq: "monthly", Priority: "0.8"},
			{Loc: h.baseURL + "/qr", LastMod: today, ChangeFreq: "monthly", Priority: "0.7"},
			{Loc: h.baseURL + "/mustang", LastMod: today, ChangeFreq: "monthly", Priority: "0.7"},
		},
	}

	// Динамически обходим курсы и их уроки
	courses := h.courseService.GetAllCourses()
	for _, c := range courses {
		// Добавляем страницу курса
		urlSet.URLs = append(urlSet.URLs, sitemapURL{
			Loc:        fmt.Sprintf("%s/courses/%s", h.baseURL, c.Slug),
			LastMod:    today,
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})

		// Загружаем уроки курса
		courseWithLessons, err := h.courseService.GetCourseBySlug(c.Slug)
		if err == nil && courseWithLessons != nil {
			for _, lesson := range courseWithLessons.Lessons {
				urlSet.URLs = append(urlSet.URLs, sitemapURL{
					Loc:        fmt.Sprintf("%s/courses/%s/%s", h.baseURL, c.Slug, lesson.Slug),
					LastMod:    today,
					ChangeFreq: "monthly",
					Priority:   "0.7",
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(urlSet)
}

// RobotsTxt отдает динамический файл robots.txt с абсолютной ссылкой на sitemap.xml.
func (h *Handler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)

	content := fmt.Sprintf("User-agent: *\nAllow: /\nDisallow: /api/\n\nSitemap: %s/sitemap.xml\n", h.baseURL)
	_, _ = w.Write([]byte(content))
}
