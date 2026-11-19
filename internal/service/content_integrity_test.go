package service

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestContentIntegrity validates that all course markdown files across all modules
// have valid content, non-empty titles, valid frontmatter, and can be converted to HTML.
func TestContentIntegrity(t *testing.T) {
	contentDir := "../../content/courses"
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		t.Skip("content/courses not found at relative path, skipping")
	}

	cs := NewCourseService(contentDir)
	courses := cs.GetAllCourses()
	if len(courses) == 0 {
		t.Fatalf("Список курсов пуст")
	}

	totalLessons := 0
	for _, course := range courses {
		courseDir := filepath.Join(contentDir, course.Slug)
		entries, err := os.ReadDir(courseDir)
		if err != nil {
			t.Errorf("Каталог модуля '%s' недоступен: %v", course.Slug, err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			totalLessons++
			filePath := filepath.Join(courseDir, entry.Name())
			t.Run(course.Slug+"/"+entry.Name(), func(t *testing.T) {
				rawBytes, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Не удалось прочитать файл %s: %v", filePath, err)
				}
				if len(rawBytes) == 0 {
					t.Fatalf("Файл %s пуст", filePath)
				}

				meta, bodyBytes := parseFrontmatter(rawBytes)
				title := meta.Title
				if title == "" {
					slug := strings.TrimSuffix(entry.Name(), ".md")
					title = extractTitle(bodyBytes, slug)
				}

				if strings.TrimSpace(title) == "" {
					t.Errorf("Файл %s не имеет заголовка (в frontmatter или первом '# ')", filePath)
				}

				// Verify goldmark conversion succeeds without panic/error
				var buf bytes.Buffer
				if err := cs.mdParser.Convert(bodyBytes, &buf); err != nil {
					t.Errorf("Ошибка парсинга Markdown в %s: %v", filePath, err)
				}

				if buf.Len() == 0 {
					t.Errorf("Результат HTML парсинга пуст для %s", filePath)
				}
			})
		}
	}

	if totalLessons == 0 {
		t.Errorf("Не найдено ни одного урока в базе")
	}
}
