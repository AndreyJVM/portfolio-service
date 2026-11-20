package handler

import (
	"encoding/xml"
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

// SitemapXML generates a dynamic sitemap.xml containing the home page (About the author) and Mustang page.
func (h *Handler) SitemapXML(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")

	urlSet := sitemapURLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: []sitemapURL{
			{Loc: h.baseURL + "/", LastMod: today, ChangeFreq: "weekly", Priority: "1.0"},
			{Loc: h.baseURL + "/mustang", LastMod: today, ChangeFreq: "monthly", Priority: "0.8"},
		},
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(urlSet)
}

// RobotsTxt serves dynamic robots.txt pointing to the canonical sitemap.
func (h *Handler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	content := "User-agent: *\nAllow: /\n\nSitemap: " + h.baseURL + "/sitemap.xml\n"
	_, _ = w.Write([]byte(content))
}
