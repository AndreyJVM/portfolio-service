package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCourseService(t *testing.T) {
	// Создаем временную директорию с курсом для тестов
	tempDir := t.TempDir()
	linuxDir := filepath.Join(tempDir, "linux")
	if err := os.MkdirAll(linuxDir, 0755); err != nil {
		t.Fatalf("Не удалось создать тестовую папку: %v", err)
	}

	testMd := "# Тестовый урок Linux\nПример содержимого урока с кодом `ls -la` и описанием."
	if err := os.WriteFile(filepath.Join(linuxDir, "01-test-lesson.md"), []byte(testMd), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый md: %v", err)
	}

	cs := NewCourseService(tempDir)

	t.Run("GetAllCourses", func(t *testing.T) {
		courses := cs.GetAllCourses()
		if len(courses) == 0 {
			t.Errorf("Ожидался непустой список курсов")
		}
		// Проверяем, что для linux количество уроков подсчитано верно (1)
		var foundLinux bool
		for _, c := range courses {
			if c.Slug == "linux" {
				foundLinux = true
				if c.LessonCount != 1 {
					t.Errorf("Ожидался 1 урок для linux, получено: %d", c.LessonCount)
				}
			}
		}
		if !foundLinux {
			t.Errorf("Курс linux не найден в списке курсов")
		}
	})

	t.Run("GetCourseBySlug - Found", func(t *testing.T) {
		course, err := cs.GetCourseBySlug("linux")
		if err != nil {
			t.Fatalf("Ошибка получения курса: %v", err)
		}
		if course.Slug != "linux" {
			t.Errorf("Ожидался slug linux, получено: %s", course.Slug)
		}
		if len(course.Lessons) != 1 {
			t.Errorf("Ожидался 1 урок, получено: %d", len(course.Lessons))
		}
		if course.Lessons[0].Title != "Тестовый урок Linux" {
			t.Errorf("Ожидался заголовок из markdown 'Тестовый урок Linux', получено: '%s'", course.Lessons[0].Title)
		}
	})

	t.Run("GetCourseBySlug - Not Found", func(t *testing.T) {
		_, err := cs.GetCourseBySlug("non-existent-course")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего курса")
		}
	})

	t.Run("GetLesson - Success and Cache", func(t *testing.T) {
		course, lesson, err := cs.GetLesson("linux", "01-test-lesson")
		if err != nil {
			t.Fatalf("Ошибка получения урока: %v", err)
		}
		if course == nil || lesson == nil {
			t.Fatal("Курс или урок равны nil")
		}
		if lesson.Slug != "01-test-lesson" {
			t.Errorf("Ожидался slug 01-test-lesson, получено: %s", lesson.Slug)
		}
		if lesson.Title != "Тестовый урок Linux" {
			t.Errorf("Ожидался заголовок 'Тестовый урок Linux', получено: '%s'", lesson.Title)
		}
		if lesson.ContentHTML == "" {
			t.Errorf("Ожидался отрендеренный HTML контент урока")
		}

		// Повторный вызов должен возвращать закэшированный экземпляр
		course2, lesson2, err2 := cs.GetLesson("linux", "01-test-lesson")
		if err2 != nil || course2 == nil || lesson2 != lesson {
			t.Errorf("Ожидалось получение того же инстанса из кэша")
		}
	})

	t.Run("GetLesson - Missing Course", func(t *testing.T) {
		_, _, err := cs.GetLesson("invalid-course", "01-test-lesson")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего курса")
		}
	})

	t.Run("GetLesson - Missing Lesson", func(t *testing.T) {
		_, _, err := cs.GetLesson("linux", "non-existent-lesson")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего урока")
		}
	})

	t.Run("BuildSearchIndex", func(t *testing.T) {
		index := cs.BuildSearchIndex()
		if len(index) == 0 {
			t.Fatalf("Индекс поиска не должен быть пустым")
		}

		var hasCourse, hasLesson bool
		for _, item := range index {
			if item.Type == "course" && item.CourseSlug == "linux" {
				hasCourse = true
			}
			if item.Type == "lesson" && item.CourseSlug == "linux" && item.LessonSlug == "01-test-lesson" {
				hasLesson = true
				if item.Title != "Тестовый урок Linux" {
					t.Errorf("Неверный заголовок в индексе: %s", item.Title)
				}
				if !strings.Contains(item.Snippet, "Пример содержимого урока") {
					t.Errorf("Ожидался текст урока в сниппете, получено: %s", item.Snippet)
				}
			}
		}

		if !hasCourse {
			t.Errorf("Курс linux не попал в индекс поиска")
		}
		if !hasLesson {
			t.Errorf("Урок 01-test-lesson не попал в индекс поиска")
		}
	})
}

func TestExtractTitleAndTitleCase(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		fallbackSlug string
		expected     string
	}{
		{
			name:         "Markdown with H1",
			content:      "# Архитектура Linux ядра\nТекст урока",
			fallbackSlug: "01-kernel-architecture",
			expected:     "Архитектура Linux ядра",
		},
		{
			name:         "Markdown with empty lines before H1",
			content:      "\n\n# Права доступа chmod\nПараграф",
			fallbackSlug: "01-permissions",
			expected:     "Права доступа chmod",
		},
		{
			name:         "No H1 header - fallback with numeric prefix",
			content:      "Просто текст без решетки заголовка",
			fallbackSlug: "02-systemd-units-guide",
			expected:     "Systemd Units Guide",
		},
		{
			name:         "No H1 header - fallback without prefix",
			content:      "",
			fallbackSlug: "docker-compose-intro",
			expected:     "Docker Compose Intro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle([]byte(tt.content), tt.fallbackSlug)
			if result != tt.expected {
				t.Errorf("extractTitle() = '%s', ожидалось '%s'", result, tt.expected)
			}
		})
	}

	t.Run("titleCase function", func(t *testing.T) {
		cases := []struct {
			input    string
			expected string
		}{
			{"hello world", "Hello World"},
			{"devops engineer", "Devops Engineer"},
			{"", ""},
		}

		for _, c := range cases {
			actual := titleCase(c.input)
			if actual != c.expected {
				t.Errorf("titleCase('%s') = '%s', ожидалось '%s'", c.input, actual, c.expected)
			}
		}
	})
}

func TestPathTraversal_Attacks(t *testing.T) {
	tempDir := t.TempDir()
	linuxDir := filepath.Join(tempDir, "linux")
	if err := os.MkdirAll(linuxDir, 0755); err != nil {
		t.Fatalf("Не удалось создать тестовую папку: %v", err)
	}

	// Создаем секретный файл вне каталога курса
	secretFile := filepath.Join(tempDir, "secret.md")
	if err := os.WriteFile(secretFile, []byte("# SECRET TOKEN"), 0644); err != nil {
		t.Fatalf("Не удалось записать secret.md: %v", err)
	}

	cs := NewCourseService(tempDir)

	maliciousLessonSlugs := []string{
		"../../secret",
		"../secret",
		"..\\secret",
		"/etc/passwd",
		"sub/lesson",
		"nested\\lesson",
	}

	for _, slug := range maliciousLessonSlugs {
		t.Run("Lesson slug traversal: "+slug, func(t *testing.T) {
			_, _, err := cs.GetLesson("linux", slug)
			if err == nil {
				t.Errorf("Уязвимость Path Traversal: вызов с '%s' не вернул ошибку!", slug)
			}
		})
	}

	maliciousCourseSlugs := []string{
		"../linux",
		"..\\linux",
		"/linux",
		"linux/../../",
	}

	for _, slug := range maliciousCourseSlugs {
		t.Run("Course slug traversal: "+slug, func(t *testing.T) {
			_, err := cs.GetCourseBySlug(slug)
			if err == nil {
				t.Errorf("Уязвимость Path Traversal: вызов с '%s' не вернул ошибку!", slug)
			}
		})
	}
}
