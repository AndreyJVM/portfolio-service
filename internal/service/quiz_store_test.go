package service

import (
	"testing"
)

func TestGetQuizForLesson_Found(t *testing.T) {
	testCases := []struct {
		slug              string
		expectedQuestions int
	}{
		{"01-architecture-and-first-dockerfile", 3},
		{"architecture-and-first-dockerfile", 3},
		{"02-multistage-and-cgroups", 3},
		{"01-lxc-and-system-containers", 3},
		{"01-podman-rootless-and-daemonless", 3},
		{"01-namespaces-and-cgroups-v2", 3},
		{"01-permissions-and-systemd", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.slug, func(t *testing.T) {
			quiz := GetQuizForLesson(tc.slug)
			if quiz == nil {
				t.Fatalf("Ожидался не nil квиз для урока %s", tc.slug)
			}
			if len(quiz.Questions) != tc.expectedQuestions {
				t.Errorf("Для урока %s ожидалось %d вопросов, получено %d", tc.slug, tc.expectedQuestions, len(quiz.Questions))
			}
			for _, q := range quiz.Questions {
				if q.ID == "" {
					t.Errorf("У вопроса отсутствует ID")
				}
				if q.Title == "" {
					t.Errorf("У вопроса отсутствует Title")
				}
				if len(q.Options) < 2 {
					t.Errorf("У вопроса %s менее 2 вариантов ответа", q.ID)
				}
				hasCorrect := false
				for _, opt := range q.Options {
					if opt.IsCorrect {
						hasCorrect = true
						break
					}
				}
				if !hasCorrect {
					t.Errorf("У вопроса %s нет ни одного правильного ответа", q.ID)
				}
			}
		})
	}
}

func TestGetQuizForLesson_Fallback(t *testing.T) {
	quiz := GetQuizForLesson("non-existent-lesson-slug")
	if quiz == nil {
		t.Fatalf("Ожидался fallback квиз по умолчанию")
	}
	if len(quiz.Questions) != 1 {
		t.Errorf("Ожидался 1 fallback вопрос, получено %d", len(quiz.Questions))
	}
	if quiz.Questions[0].ID != "q1" {
		t.Errorf("ID fallback вопроса должен быть q1")
	}
}
