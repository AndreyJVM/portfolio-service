package model

import (
	"encoding/json"
	"testing"
)

func TestCourseAndLessonJSONSerialization(t *testing.T) {
	lesson := Lesson{
		ID:          "lesson-1",
		CourseSlug:  "docker",
		Slug:        "01-intro",
		Title:       "Введение",
		Order:       1,
		Duration:    "10 мин",
		Description: "Описание урока",
		ContentHTML: "<p>Контент</p>",
	}

	course := Course{
		ID:          "docker",
		Title:       "Docker",
		Slug:        "docker",
		Icon:        "fab fa-docker",
		Color:       "#0284c7",
		Description: "Курс по Docker",
		Level:       "Базовый",
		LessonCount: 1,
		Tags:        []string{"Docker", "OCI"},
		Lessons:     []Lesson{lesson},
	}

	data, err := json.Marshal(course)
	if err != nil {
		t.Fatalf("Ошибка сериализации Course в JSON: %v", err)
	}

	var parsed Course
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Ошибка десериализации Course из JSON: %v", err)
	}

	if parsed.Slug != "docker" || len(parsed.Lessons) != 1 {
		t.Errorf("Некорректная десериализация курса: %+v", parsed)
	}
	if parsed.Lessons[0].Title != "Введение" {
		t.Errorf("Некорректный заголовок урока: %s", parsed.Lessons[0].Title)
	}
}

func TestSearchItemAndMetaData(t *testing.T) {
	item := SearchItem{
		Type:        "lesson",
		Title:       "Namespaces",
		CourseTitle: "Internals",
		CourseSlug:  "internals",
		LessonSlug:  "01-namespaces",
		URL:         "/courses/internals/01-namespaces",
		Description: "Анатомия пространств имен",
		Tags:        []string{"Linux", "Kernel"},
		Snippet:     "Фрагмент текста",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Ошибка сериализации SearchItem: %v", err)
	}

	var parsed SearchItem
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Ошибка десериализации SearchItem: %v", err)
	}

	if parsed.Title != "Namespaces" || parsed.CourseSlug != "internals" {
		t.Errorf("Некорректная десериализация SearchItem: %+v", parsed)
	}

	meta := MetaData{
		Title:        "Title",
		ActiveTab:    "courses",
		ActiveCourse: "docker",
		CanonicalURL: "https://vorobevaqa.ru/",
	}
	if meta.ActiveCourse != "docker" || meta.ActiveTab != "courses" {
		t.Errorf("Некорректная MetaData: %+v", meta)
	}
}
