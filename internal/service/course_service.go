package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"ru.vorobevaqa/devops-platform/internal/model"
)

type CourseService struct {
	contentDir string
	mdParser   goldmark.Markdown
}

func NewCourseService(contentDir string) *CourseService {
	return &CourseService{
		contentDir: contentDir,
		mdParser:   goldmark.New(),
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
			Level:       "Junior / Middle",
			LessonCount: 3,
			Description: "Namespaces и cgroups, многоэтапные Dockerfile, кэширование слоев, сети и оркестрация в Docker Compose с healthcheck.",
			Tags:        []string{"Docker", "Compose", "Multi-stage", "Networking"},
		},
		{
			ID:          "cicd",
			Slug:        "cicd",
			Title:       "CI/CD: GitHub Actions & Деплой",
			Icon:        "fas fa-infinity",
			Color:       "#10b981",
			Level:       "Middle",
			LessonCount: 2,
			Description: "Автоматизация тестирования и сборки, paths-filter, Zero-Downtime Blue-Green деплой на VDS, безопасность секретов.",
			Tags:        []string{"CI/CD", "GitHub Actions", "Blue-Green", "Zero-Downtime"},
		},
		{
			ID:          "ansible",
			Slug:        "ansible",
			Title:       "Ansible: Инфраструктура как Код (IaC)",
			Icon:        "fas fa-network-wired",
			Color:       "#ef4444",
			Level:       "Middle",
			LessonCount: 2,
			Description: "Декларативное управление серверами, написание ролей, инвентарь, хэндлеры, шифрование Vault и публикация в Ansible Galaxy.",
			Tags:        []string{"Ansible", "IaC", "Automation", "Vault"},
		},
	}
}

// GetCourseBySlug находит курс и загружает список его уроков.
func (s *CourseService) GetCourseBySlug(slug string) (*model.Course, error) {
	courses := s.GetAllCourses()
	for _, c := range courses {
		if c.Slug == slug {
			lessons, err := s.getLessonsForCourse(slug)
			if err == nil {
				c.Lessons = lessons
				c.LessonCount = len(lessons)
			}
			return &c, nil
		}
	}
	return nil, fmt.Errorf("курс '%s' не найден", slug)
}

// GetLesson находит и рендерит урок из Markdown файла.
func (s *CourseService) GetLesson(courseSlug, lessonSlug string) (*model.Course, *model.Lesson, error) {
	course, err := s.GetCourseBySlug(courseSlug)
	if err != nil {
		return nil, nil, err
	}

	for _, l := range course.Lessons {
		if l.Slug == lessonSlug {
			filePath := filepath.Join(s.contentDir, courseSlug, lessonSlug+".md")
			mdBytes, err := os.ReadFile(filePath)
			if err != nil {
				return course, &l, nil
			}

			var buf bytes.Buffer
			if err := s.mdParser.Convert(mdBytes, &buf); err != nil {
				l.ContentHTML = string(mdBytes)
			} else {
				l.ContentHTML = buf.String()
			}

			return course, &l, nil
		}
	}

	return nil, nil, fmt.Errorf("урок '%s' не найден в курсе '%s'", lessonSlug, courseSlug)
}

func (s *CourseService) getLessonsForCourse(courseSlug string) ([]model.Lesson, error) {
	dirPath := filepath.Join(s.contentDir, courseSlug)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var lessons []model.Lesson
	order := 1
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			slug := strings.TrimSuffix(entry.Name(), ".md")
			title := formatLessonTitle(slug)
			lessons = append(lessons, model.Lesson{
				ID:          slug,
				CourseSlug:  courseSlug,
				Slug:        slug,
				Title:       title,
				Order:       order,
				Duration:    "10-15 мин",
				Description: "Практический материал и примеры конфигураций.",
			})
			order++
		}
	}

	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Slug < lessons[j].Slug
	})

	return lessons, nil
}

func formatLessonTitle(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
