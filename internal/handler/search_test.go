package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"ru.vorobevaqa/devops-platform/internal/model"
)

func TestSearchAPI(t *testing.T) {
	r := setupTestRouter()

	// 1. Fetch full search index
	req := httptest.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	var items []model.SearchItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("Ошибка декодирования search JSON: %v", err)
	}

	if len(items) == 0 {
		t.Errorf("Ожидался непустой индекс поиска")
	}

	// 2. Server-side filtering by query "docker"
	reqQuery := httptest.NewRequest("GET", "/api/search?q=docker", nil)
	wQuery := httptest.NewRecorder()
	r.ServeHTTP(wQuery, reqQuery)

	if wQuery.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200 для ?q=docker, получен %d", wQuery.Code)
	}

	var queryItems []model.SearchItem
	if err := json.NewDecoder(wQuery.Body).Decode(&queryItems); err != nil {
		t.Fatalf("Ошибка декодирования отфильтрованного search JSON: %v", err)
	}

	if len(queryItems) == 0 {
		t.Errorf("Ожидались результаты для запроса 'docker'")
	}

	// Verify that all results mention docker
	for _, it := range queryItems {
		found := strings.Contains(strings.ToLower(it.Title), "docker") ||
			strings.Contains(strings.ToLower(it.Description), "docker") ||
			strings.Contains(strings.ToLower(it.Snippet), "docker") ||
			strings.Contains(strings.ToLower(it.CourseTitle), "docker")
		if !found {
			t.Errorf("Элемент %s не соответствует запросу 'docker'", it.Title)
		}
	}
}

func TestSecurity_SearchEdgeCases(t *testing.T) {
	r := setupTestRouter()

	t.Run("XSS payload in search query", func(t *testing.T) {
		xssPayload := "<script>alert('xss')</script>"
		req := httptest.NewRequest("GET", "/api/search?q="+xssPayload, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}

		var items []model.SearchItem
		if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
			t.Fatalf("Ошибка десериализации JSON при XSS запросе: %v", err)
		}
	})

	t.Run("Huge query string (50000 chars) does not crash or hang", func(t *testing.T) {
		hugeQuery := strings.Repeat("A", 50000)
		req := httptest.NewRequest("GET", "/api/search?q="+hugeQuery, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}
	})

	t.Run("Special chars and regex symbols in search query", func(t *testing.T) {
		specialQueries := []string{"[", "]", "(", ")", "*", "+", "?", "\\", "$", "^", "\"", "'", "NULL", "%00"}
		for _, sq := range specialQueries {
			req := httptest.NewRequest("GET", "/api/search?q="+sq, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Сбой поиска при спецсимволе '%s': статус %d", sq, w.Code)
			}
		}
	})
}

func TestSearchAPI_MultibyteUTF8Truncation(t *testing.T) {
	r := setupTestRouter()

	// Long Cyrillic string (>128 runes)
	longQuery := strings.Repeat("контейнеры ", 20)
	req := httptest.NewRequest("GET", "/api/search?q="+url.QueryEscape(longQuery), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", w.Code)
	}

	body := w.Body.Bytes()
	if !utf8.Valid(body) {
		t.Errorf("Ответ API поиска содержит невалидный UTF-8 из-за обрезки рун")
	}

	var results []model.SearchItem
	if err := json.Unmarshal(body, &results); err != nil {
		t.Fatalf("Не удалось распарсить JSON ответ: %v", err)
	}
}

func TestSearchAPI_CacheRebuildAndTagMatch(t *testing.T) {
	r := setupTestRouter()

	t.Run("Cache rebuild when expired", func(t *testing.T) {
		searchIndexMu.Lock()
		lastIndexBuild = time.Now().Add(-2 * time.Hour)
		searchIndexMu.Unlock()

		req := httptest.NewRequest("GET", "/api/search", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200 при пересчете кэша, получен %d", w.Code)
		}
	})

	t.Run("Search by tag", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=cgroups", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200 при поиске по тегу, получен %d", w.Code)
		}
		var results []model.SearchItem
		if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
			t.Fatalf("Ошибка парсинга JSON: %v", err)
		}
		if len(results) == 0 {
			t.Errorf("Ожидались результаты поиска по тегу cgroups")
		}
	})
}
