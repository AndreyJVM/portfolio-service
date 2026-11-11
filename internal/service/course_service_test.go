package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCourseService(t *testing.T) {
	// Создаем временную директорию с курсом для тестов
	tempDir := t.TempDir()
	linuxDir := filepath.Join(tempDir, "linux")
	if err := os.MkdirAll(linuxDir, 0755); err != nil {
		t.Fatalf("Не удалось создать тестовую папку: %v", err)
	}

	testMd := "# Тестовый урок\nПример содержимого урока."
	if err := os.WriteFile(filepath.Join(linuxDir, "test-lesson.md"), []byte(testMd), 0644); err != nil {
		t.Fatalf("Не удалось записать тестовый md: %v", err)
	}

	cs := NewCourseService(tempDir)

	t.Run("GetAllCourses", func(t *testing.T) {
		courses := cs.GetAllCourses()
		if len(courses) == 0 {
			t.Errorf("Ожидался непустой список курсов")
		}
	})

	t.Run("GetCourseBySlug", func(t *testing.T) {
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
	})

	t.Run("GetLesson", func(t *testing.T) {
		course, lesson, err := cs.GetLesson("linux", "test-lesson")
		if err != nil {
			t.Fatalf("Ошибка получения урока: %v", err)
		}
		if course == nil || lesson == nil {
			t.Fatal("Курс или урок равны nil")
		}
		if lesson.Slug != "test-lesson" {
			t.Errorf("Ожидался slug test-lesson, получено: %s", lesson.Slug)
		}
		if lesson.ContentHTML == "" {
			t.Errorf("Ожидался отрендеренный HTML контент урока")
		}
	})
}
