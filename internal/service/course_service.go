package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"ru.vorobevaqa/devops-platform/internal/model"
)

type CourseService struct {
	contentDir  string
	mdParser    goldmark.Markdown
	lessonCache map[string]*model.Lesson
	lessonMu    sync.RWMutex
}

var mdTagRegex = regexp.MustCompile(`(?s)<[^>]+>|[#*` + "`" + `~_\[\]()-]`)

func NewCourseService(contentDir string) *CourseService {
	return &CourseService{
		contentDir:  contentDir,
		lessonCache: make(map[string]*model.Lesson),
		mdParser: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
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
			LessonCount: 3,
			Description: "Многоэтапная сборка (Multi-stage), оптимизация слоев, безопасность, cgroups, неймспейсы ядра и docker-compose оркестрация.",
			Tags:        []string{"Docker", "Compose", "Containers", "DevOps", "Isolation"},
		},
		{
			ID:          "ansible",
			Slug:        "ansible",
			Title:       "Ansible & Управление конфигурацией",
			Icon:        "fas fa-network-wired",
			Color:       "#dc2626",
			Level:       "Middle",
			LessonCount: 3,
			Description: "Инфраструктура как код (IaC), декларативные плейбуки, роли, шаблонизация Jinja2, шифрование Ansible Vault и интеграция с CI/CD.",
			Tags:        []string{"Ansible", "IaC", "Automation", "Galaxy", "YAML"},
		},
		{
			ID:          "cicd",
			Slug:        "cicd",
			Title:       "CI/CD: GitHub Actions & Деплой",
			Icon:        "fas fa-infinity",
			Color:       "#10b981",
			Level:       "Middle / Senior",
			LessonCount: 3,
			Description: "Построение надежных конвейеров сборки, автотесты, публикация в Docker Hub, Zero-Downtime Blue-Green деплой на VDS хостинг.",
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

// GetLesson загружает конкретный урок курса и парсит Markdown в HTML с кэшированием в памяти.
func (s *CourseService) GetLesson(courseSlug, lessonSlug string) (*model.Course, *model.Lesson, error) {
	course, err := s.GetCourseBySlug(courseSlug)
	if err != nil {
		return nil, nil, err
	}

	cacheKey := courseSlug + "/" + lessonSlug
	s.lessonMu.RLock()
	if cached, ok := s.lessonCache[cacheKey]; ok {
		s.lessonMu.RUnlock()
		return course, cached, nil
	}
	s.lessonMu.RUnlock()

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
	title = titleCase(title)

	lesson := &model.Lesson{
		CourseSlug:  courseSlug,
		Slug:        lessonSlug,
		Title:       title,
		Duration:    "15 мин",
		Description: "Практическое руководство и лабораторная работа",
		ContentHTML: buf.String(),
	}

	s.lessonMu.Lock()
	s.lessonCache[cacheKey] = lesson
	s.lessonMu.Unlock()

	return course, lesson, nil
}

// BuildSearchIndex создает плоский поисковый индекс по всем курсам и урокам базы знаний.
func (s *CourseService) BuildSearchIndex() []model.SearchItem {
	var index []model.SearchItem

	courses := s.GetAllCourses()
	for _, c := range courses {
		// Добавляем карточку самого курса
		index = append(index, model.SearchItem{
			Type:        "course",
			Title:       c.Title,
			CourseTitle: c.Title,
			CourseSlug:  c.Slug,
			URL:         fmt.Sprintf("/courses/%s", c.Slug),
			Description: c.Description,
			Tags:        c.Tags,
			Snippet:     c.Description,
		})

		// Загружаем все уроки курса
		lessons, err := s.loadLessonsForCourse(c.Slug)
		if err != nil {
			continue
		}

		for _, l := range lessons {
			lessonPath := filepath.Join(s.contentDir, c.Slug, l.Slug+".md")
			rawBytes, err := os.ReadFile(lessonPath)
			snippet := ""
			if err == nil {
				cleaned := mdTagRegex.ReplaceAllString(string(rawBytes), " ")
				cleaned = strings.Join(strings.Fields(cleaned), " ")
				if len(cleaned) > 200 {
					cleaned = cleaned[:200] + "..."
				}
				snippet = cleaned
			} else {
				snippet = l.Description
			}

			index = append(index, model.SearchItem{
				Type:        "lesson",
				Title:       l.Title,
				CourseTitle: c.Title,
				CourseSlug:  c.Slug,
				LessonSlug:  l.Slug,
				URL:         fmt.Sprintf("/courses/%s/%s", c.Slug, l.Slug),
				Description: l.Description,
				Tags:        c.Tags,
				Snippet:     snippet,
			})
		}
	}

	return index
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
			title = titleCase(title)

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

// titleCase переводит первую букву каждого слова в верхний регистр (замена устаревшего strings.Title).
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			r := []rune(w)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
