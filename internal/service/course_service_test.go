package service

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"
)

func TestCourseService(t *testing.T) {
	// Create temporary directory with a course for tests
	tempDir := t.TempDir()
	dockerDir := filepath.Join(tempDir, "docker")
	if err := os.MkdirAll(dockerDir, 0755); err != nil {
		t.Fatalf("Не удалось создать тестовую папку: %v", err)
	}

	testMd := "# Тестовый урок Docker\nПример содержимого урока с кодом `docker run` и описанием."
	if err := os.WriteFile(filepath.Join(dockerDir, "01-test-lesson.md"), []byte(testMd), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый md: %v", err)
	}

	cs := NewCourseService(tempDir)

	t.Run("GetAllCourses", func(t *testing.T) {
		courses := cs.GetAllCourses()
		if len(courses) == 0 {
			t.Errorf("Ожидался непустой список курсов")
		}
		// Verify that lesson count is calculated correctly for docker (1)
		var foundDocker bool
		for _, c := range courses {
			if c.Slug == "docker" {
				foundDocker = true
				if c.LessonCount != 1 {
					t.Errorf("Ожидался 1 урок для docker, получено: %d", c.LessonCount)
				}
			}
		}
		if !foundDocker {
			t.Errorf("Курс docker не найден в списке курсов")
		}
	})

	t.Run("GetCourseBySlug - Found", func(t *testing.T) {
		course, err := cs.GetCourseBySlug("docker")
		if err != nil {
			t.Fatalf("Ошибка получения курса: %v", err)
		}
		if course.Slug != "docker" {
			t.Errorf("Ожидался slug docker, получено: %s", course.Slug)
		}
		if len(course.Lessons) != 1 {
			t.Errorf("Ожидался 1 урок, получено: %d", len(course.Lessons))
		}
		if course.Lessons[0].Title != "Тестовый урок Docker" {
			t.Errorf("Ожидался заголовок из markdown 'Тестовый урок Docker', получено: '%s'", course.Lessons[0].Title)
		}
	})

	t.Run("GetCourseBySlug - Not Found", func(t *testing.T) {
		_, err := cs.GetCourseBySlug("non-existent-course")
		if err == nil {
			t.Errorf("Ожидалась ошибка для несуществующего курса")
		}
	})

	t.Run("GetLesson - Success and Cache", func(t *testing.T) {
		course, lesson, err := cs.GetLesson("docker", "01-test-lesson")
		if err != nil {
			t.Fatalf("Ошибка получения урока: %v", err)
		}
		if course == nil || lesson == nil {
			t.Fatal("Курс или урок равны nil")
		}
		if lesson.Slug != "01-test-lesson" {
			t.Errorf("Ожидался slug 01-test-lesson, получено: %s", lesson.Slug)
		}
		if lesson.Title != "Тестовый урок Docker" {
			t.Errorf("Ожидался заголовок 'Тестовый урок Docker', получено: '%s'", lesson.Title)
		}
		if lesson.ContentHTML == "" {
			t.Errorf("Ожидался отрендеренный HTML контент урока")
		}

		// Repeated call must return the cached instance
		course2, lesson2, err2 := cs.GetLesson("docker", "01-test-lesson")
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
		_, _, err := cs.GetLesson("docker", "non-existent-lesson")
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
			if item.Type == "course" && item.CourseSlug == "docker" {
				hasCourse = true
			}
			if item.Type == "lesson" && item.CourseSlug == "docker" && item.LessonSlug == "01-test-lesson" {
				hasLesson = true
				if item.Title != "Тестовый урок Docker" {
					t.Errorf("Неверный заголовок в индексе: %s", item.Title)
				}
				if !strings.Contains(item.Snippet, "Пример содержимого урока") {
					t.Errorf("Ожидался текст урока в сниппете, получено: %s", item.Snippet)
				}
			}
		}

		if !hasCourse {
			t.Errorf("Курс docker не попал в индекс поиска")
		}
		if !hasLesson {
			t.Errorf("Урок 01-test-lesson не попал в индекс поиска")
		}
	})
}

func TestPathTraversal_Attacks(t *testing.T) {
	tempDir := t.TempDir()
	linuxDir := filepath.Join(tempDir, "linux")
	if err := os.MkdirAll(linuxDir, 0755); err != nil {
		t.Fatalf("Не удалось создать тестовую папку: %v", err)
	}

	// Create secret file outside course directory
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
			_, _, err := cs.GetLesson("docker", slug)
			if err == nil {
				t.Errorf("Уязвимость Path Traversal: вызов с '%s' не вернул ошибку!", slug)
			}
		})
	}

	maliciousCourseSlugs := []string{
		"../docker",
		"..\\docker",
		"/docker",
		"docker/../../",
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

func TestBuildSearchIndex_MultibyteUTF8(t *testing.T) {
	tempDir := t.TempDir()
	dockerDir := filepath.Join(tempDir, "docker")
	if err := os.MkdirAll(dockerDir, 0755); err != nil {
		t.Fatalf("Не удалось создать папку курса: %v", err)
	}

	// Create long Russian text (>200 runes) where 200th byte boundary falls inside a two-byte Cyrillic rune
	longRussianText := "# Управление памятью и процессами в контейнерах Docker\n" +
		strings.Repeat("Операционная система распределяет оперативную память виртуальными страницами через cgroups. ", 10)

	lessonFile := filepath.Join(dockerDir, "01-memory-cgroups.md")
	if err := os.WriteFile(lessonFile, []byte(longRussianText), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый файл: %v", err)
	}

	cs := NewCourseService(tempDir)
	index := cs.BuildSearchIndex()

	var snippet string
	for _, item := range index {
		if item.Type == "lesson" && item.LessonSlug == "01-memory-cgroups" {
			snippet = item.Snippet
			break
		}
	}

	if snippet == "" {
		t.Fatal("Сниппет для 01-memory-cgroups не найден в индексе")
	}

	if !utf8.ValidString(snippet) {
		t.Errorf("Сниппет содержит невалидный UTF-8 из-за разрыва многобайтовой руны: %q", snippet)
	}

	if !strings.HasSuffix(snippet, "...") {
		t.Errorf("Сниппет должен оканчиваться троеточием '...', получено: %s", snippet)
	}

	runes := []rune(snippet)
	// Must be 200 runes + 3 dot characters
	if len(runes) != 203 {
		t.Errorf("Ожидалось 203 руны (200 + ...), получено: %d", len(runes))
	}
}

func TestCourseService_NegativeAndConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	dockerDir := filepath.Join(tempDir, "docker")
	_ = os.MkdirAll(dockerDir, 0755)

	content := "# Урок по контейнеризации\nСодержимое урока для тестирования параллельного доступа."
	_ = os.WriteFile(filepath.Join(dockerDir, "01-containers.md"), []byte(content), 0644)

	cs := NewCourseService(tempDir)

	t.Run("Empty non-existent directory", func(t *testing.T) {
		emptyCS := NewCourseService(filepath.Join(tempDir, "non-existent-dir"))
		courses := emptyCS.GetAllCourses()
		if len(courses) == 0 {
			t.Errorf("GetAllCourses должен возвращать каталог курсов даже при пустой contentDir")
		}
		// All courses should have 0 lessons
		for _, c := range courses {
			if c.LessonCount != 0 {
				t.Errorf("Курс %s должен иметь 0 уроков, получено %d", c.Slug, c.LessonCount)
			}
		}
	})

	t.Run("Concurrent access to GetLesson and BuildSearchIndex", func(t *testing.T) {
		var wg sync.WaitGroup
		workers := 20
		wg.Add(workers * 2)

		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				_, _, _ = cs.GetLesson("docker", "01-containers")
			}()
			go func() {
				defer wg.Done()
				_ = cs.BuildSearchIndex()
			}()
		}

		wg.Wait()
	})
}

func TestCourseService_LoadCoursesFromJSON(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "custom-courses.json")
	coursesJSON := `[
		{
			"id": "docker-test",
			"slug": "docker-test",
			"title": "Модуль тестирования Docker",
			"icon": "fab fa-docker",
			"color": "#0284c7",
			"level": "Middle",
			"description": "Тестовый модуль из JSON",
			"tags": ["Docker", "Testing"]
		}
	]`
	if err := os.WriteFile(jsonPath, []byte(coursesJSON), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый JSON: %v", err)
	}

	cs := NewCourseService(tempDir, jsonPath)
	courses := cs.GetAllCourses()
	if len(courses) != 1 {
		t.Fatalf("Ожидался 1 курс из JSON, получено: %d", len(courses))
	}
	if courses[0].Slug != "docker-test" {
		t.Errorf("Ожидался slug 'docker-test', получено: '%s'", courses[0].Slug)
	}
	if courses[0].Title != "Модуль тестирования Docker" {
		t.Errorf("Ожидался title 'Модуль тестирования Docker', получено: '%s'", courses[0].Title)
	}
}
