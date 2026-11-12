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

// SearchAPI возвращает JSON-индекс уроков и курсов для мгновенного поиска в браузере.
// Поддерживает как отдачу полного индекса (для быстрого client-side fuzzy поиска),
// так и server-side фильтрацию через параметр ?q=.
func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {
	searchIndexMu.RLock()
	// Кэшируем индекс на 1 час или отдаем из памяти
	if len(searchIndexCache) > 0 && time.Since(lastIndexBuild) < 1*time.Hour {
		items := searchIndexCache
		searchIndexMu.RUnlock()
		h.respondSearch(w, r, items)
		return
	}
	searchIndexMu.RUnlock()

	searchIndexMu.Lock()
	// Double-Checked Locking: проверяем, не построил ли индекс другой параллельный запрос
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
	// Ограничиваем длину поискового запроса по рунам для защиты от CPU exhaustion и сохранения валидности UTF-8
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
