package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"ru.vorobevaqa/devops-platform/internal/model"
)

type CourseService struct {
	contentDir string
	mdParser   goldmark.Markdown
}

func NewCourseService(contentDir string) *CourseService {
	return &CourseService{
		contentDir: contentDir,
		mdParser: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithUnsafe(),
			),
		),
	}
}

// GetAllCourses возвращает список всех доступных курсов.
func (s *CourseService) GetAllCourses() []model.Course {
	return []model.Course{
		{
			ID:          "linux",
			Slug:        "linux",
			Title:       "Linux & Bash с нуля до Pro",
			Icon:        "fab fa-linux",
			Color:       "#f59e0b",
			Level:       "Junior / Middle",
			LessonCount: 3,
			Description: "Устройство ОС Linux, работа с ядром, права доступа (chmod, chown), systemd сервисы, cron демоны и автоматизация на Bash.",
			Tags:        []string{"Linux", "Ubuntu", "Bash", "Systemd", "SSH"},
		},
		{
			ID:          "docker",
			Slug:        "docker",
			Title:       "Docker & Контейнеризация",
			Icon:        "fab fa-docker",
			Color:       "#0284c7",
			Level:       "Middle",
			LessonCount: 2,
			Description: "Многоэтапные сборки, оптимизация кэша слоев, изоляция cgroups, сети bridge/overlay, Docker Compose v2 и Swarm.",
			Tags:        []string{"Docker", "Compose", "cgroups", "Multi-stage"},
		},
		{
			ID:          "ansible",
			Slug:        "ansible",
			Title:       "Ansible & Инфраструктура как код (IaC)",
			Icon:        "fas fa-network-wired",
			Color:       "#dc2626",
			Level:       "Middle / Senior",
			LessonCount: 2,
			Description: "Идемпотентные плейбуки, создание модульных ролей для Galaxy, шифрование с Ansible Vault, управление inventory.",
			Tags:        []string{"Ansible", "IaC", "Galaxy", "YAML"},
		},
		{
			ID:          "cicd",
			Slug:        "cicd",
			Title:       "CI/CD: GitHub Actions & GitLab",
			Icon:        "fas fa-infinity",
			Color:       "#8b5cf6",
			Level:       "Middle / Senior",
			LessonCount: 2,
			Description: "Проектирование надежных пайплайнов, Blue-Green деплой без простоя, авто-откат (rollback) и проверка здоровья.",
			Tags:        []string{"CI/CD", "GitHub Actions", "GitLab CI", "Blue-Green"},
		},
	}
}

// GetCourseBySlug находит курс по слагу и загружает уроки из каталога content/courses/<slug>.
func (s *CourseService) GetCourseBySlug(slug string) (*model.Course, error) {
	courses := s.GetAllCourses()
	var foundCourse *model.Course
	for i := range courses {
		if courses[i].Slug == slug {
			foundCourse = &courses[i]
			break
		}
	}
	if foundCourse == nil {
		return nil, fmt.Errorf("курс '%s' не найден", slug)
	}

	lessons, err := s.loadLessonsForCourse(slug)
	if err == nil {
		foundCourse.Lessons = lessons
		foundCourse.LessonCount = len(lessons)
	}

	return foundCourse, nil
}

// GetLesson загружает конкретный урок курса и парсит Markdown в HTML.
func (s *CourseService) GetLesson(courseSlug, lessonSlug string) (*model.Course, *model.Lesson, error) {
	course, err := s.GetCourseBySlug(courseSlug)
	if err != nil {
		return nil, nil, err
	}

	lessonPath := filepath.Join(s.contentDir, courseSlug, lessonSlug+".md")
	contentBytes, err := os.ReadFile(lessonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("урок '%s' не найден: %w", lessonSlug, err)
	}

	var buf bytes.Buffer
	if err := s.mdParser.Convert(contentBytes, &buf); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга markdown: %w", err)
	}

	title := strings.ReplaceAll(lessonSlug, "-", " ")
	title = strings.Title(title)

	lesson := &model.Lesson{
		CourseSlug:  courseSlug,
		Slug:        lessonSlug,
		Title:       title,
		Duration:    "15 мин",
		Description: "Практическое руководство и лабораторная работа",
		ContentHTML: buf.String(),
	}

	return course, lesson, nil
}

func (s *CourseService) loadLessonsForCourse(courseSlug string) ([]model.Lesson, error) {
	courseDir := filepath.Join(s.contentDir, courseSlug)
	entries, err := os.ReadDir(courseDir)
	if err != nil {
		return nil, err
	}

	var lessons []model.Lesson
	for idx, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			slug := strings.TrimSuffix(entry.Name(), ".md")
			title := strings.ReplaceAll(slug, "-", " ")
			// Убираем числовой префикс если есть (напр. "01 ")
			if len(title) > 3 && title[2] == ' ' {
				title = title[3:]
			}
			title = strings.Title(title)

			lessons = append(lessons, model.Lesson{
				CourseSlug:  courseSlug,
				Slug:        slug,
				Title:       title,
				Order:       idx + 1,
				Duration:    "15-20 мин",
				Description: fmt.Sprintf("Лабораторная работа и разбор темы %s", title),
			})
		}
	}

	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Order < lessons[j].Order
	})

	return lessons, nil
}
