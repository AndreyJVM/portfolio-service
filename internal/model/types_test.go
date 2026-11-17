package model

import (
	"encoding/json"
	"html/template"
	"testing"
)

func TestQuizQuestion_TitleHTML(t *testing.T) {
	q := QuizQuestion{
		ID:    "q1",
		Title: "Вопрос с <code>CMD [\"./app\"]</code> и <b>жирным текстом</b>",
	}

	expected := template.HTML("Вопрос с <code>CMD [\"./app\"]</code> и <b>жирным текстом</b>")
	if q.TitleHTML() != expected {
		t.Errorf("Ожидалось %v, получено %v", expected, q.TitleHTML())
	}
}

func TestQuizOption_TextHTML(t *testing.T) {
	opt := QuizOption{
		Text:        "Вариант ответа с <code>type=gha</code>",
		IsCorrect:   true,
		Explanation: "Пояснение",
	}

	expected := template.HTML("Вариант ответа с <code>type=gha</code>")
	if opt.TextHTML() != expected {
		t.Errorf("Ожидалось %v, получено %v", expected, opt.TextHTML())
	}
}

func TestModelJSONSerialization(t *testing.T) {
	quiz := Quiz{
		Questions: []QuizQuestion{
			{
				ID:    "q1",
				Title: "Test question",
				Options: []QuizOption{
					{
						Text:        "Option 1",
						IsCorrect:   true,
						Explanation: "Explanation 1",
					},
					{
						Text:      "Option 2",
						IsCorrect: false,
					},
				},
			},
		},
	}

	data, err := json.Marshal(quiz)
	if err != nil {
		t.Fatalf("Ошибка сериализации Quiz в JSON: %v", err)
	}

	var deserialized Quiz
	if err := json.Unmarshal(data, &deserialized); err != nil {
		t.Fatalf("Ошибка десериализации Quiz из JSON: %v", err)
	}

	if len(deserialized.Questions) != 1 {
		t.Fatalf("Ожидался 1 вопрос, получено %d", len(deserialized.Questions))
	}
	if deserialized.Questions[0].ID != "q1" {
		t.Errorf("Ожидался ID 'q1', получено '%s'", deserialized.Questions[0].ID)
	}
	if len(deserialized.Questions[0].Options) != 2 {
		t.Errorf("Ожидалось 2 варианта ответа, получено %d", len(deserialized.Questions[0].Options))
	}
	if !deserialized.Questions[0].Options[0].IsCorrect {
		t.Errorf("Первый вариант должен быть правильным")
	}
	if deserialized.Questions[0].Options[1].IsCorrect {
		t.Errorf("Второй вариант должен быть неправильным")
	}
}

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
