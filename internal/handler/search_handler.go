package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"ru.vorobevaqa/devops-platform/internal/model"
)

var (
	searchIndexCache []model.SearchItem
	searchIndexMu    sync.RWMutex
	lastIndexBuild   time.Time
)

// SearchAPI returns a JSON index of lessons and courses for instant in-browser search.
// Supports both returning the full index (for fast client-side fuzzy search)
// and server-side filtering via the ?q= parameter.
func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {
	searchIndexMu.RLock()
	// Cache the index for 1 hour or serve from memory
	if len(searchIndexCache) > 0 && time.Since(lastIndexBuild) < 1*time.Hour {
		items := searchIndexCache
		searchIndexMu.RUnlock()
		h.respondSearch(w, r, items)
		return
	}
	searchIndexMu.RUnlock()

	searchIndexMu.Lock()
	// Double-Checked Locking: check if another concurrent request already built the index
	if len(searchIndexCache) > 0 && time.Since(lastIndexBuild) < 1*time.Hour {
		items := searchIndexCache
		searchIndexMu.Unlock()
		h.respondSearch(w, r, items)
		return
	}

	searchIndexCache = h.courseService.BuildSearchIndex()
	lastIndexBuild = time.Now()
	items := searchIndexCache
	searchIndexMu.Unlock()

	h.respondSearch(w, r, items)
}

func (h *Handler) respondSearch(w http.ResponseWriter, r *http.Request, items []model.SearchItem) {
	query := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	// Limit search query length by runes to prevent CPU exhaustion and maintain valid UTF-8
	runes := []rune(query)
	if len(runes) > 128 {
		query = string(runes[:128])
	}

	var result []model.SearchItem
	if query == "" {
		result = items
	} else {
		for _, item := range items {
			match := strings.Contains(strings.ToLower(item.Title), query) ||
				strings.Contains(strings.ToLower(item.CourseTitle), query) ||
				strings.Contains(strings.ToLower(item.Description), query) ||
				strings.Contains(strings.ToLower(item.Snippet), query)

			if !match {
				for _, tag := range item.Tags {
					if strings.Contains(strings.ToLower(tag), query) {
						match = true
						break
					}
				}
			}

			if match {
				result = append(result, item)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
