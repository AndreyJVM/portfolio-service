package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	courses := []model.Course{
		{
			ID:          "docker",
			Slug:        "docker",
			Title:       "Docker: От основ до Production",
			Icon:        "fab fa-docker",
			Color:       "#0284c7",
			Level:       "Базовый → Pro",
			Description: "Устройство OverlayFS, кэширование слоев, Multi-stage сборка, сети, тома, Compose и безопасность контейнеров.",
			Tags:        []string{"Docker", "Compose", "BuildKit", "Production", "Multi-stage"},
		},
		{
			ID:          "podman",
			Slug:        "podman",
			Title:       "Podman & Quadlet: Rootless",
			Icon:        "fas fa-shield-halved",
			Color:       "#8b5cf6",
			Level:       "Middle → Senior",
			Description: "Daemonless контейнеризация без root-прав, изоляция через user namespaces и нативный запуск контейнеров в systemd через Quadlet.",
			Tags:        []string{"Podman", "Rootless", "Quadlet", "Systemd", "Security"},
		},
		{
			ID:          "lxc",
			Slug:        "lxc",
			Title:       "LXC / Incus: Системные контейнеры",
			Icon:        "fas fa-boxes-stacked",
			Color:       "#10b981",
			Level:       "Middle → Senior",
			Description: "Легковесная замена виртуальным машинам (KVM): системные контейнеры Linux, Proxmox, ZFS/Btrfs бекенды и прямое выделение ресурсов.",
			Tags:        []string{"LXC", "Incus", "LXD", "Proxmox", "System Containers"},
		},
		{
			ID:          "internals",
			Slug:        "internals",
			Title:       "Container Internals: Ядро и изоляция",
			Icon:        "fas fa-microchip",
			Color:       "#f59e0b",
			Level:       "Senior",
			Description: "Низкоуровневая анатомия контейнеризации: Namespaces (PID, NET, MNT, USER), Cgroups v2, runc, crun, seccomp и AppArmor.",
			Tags:        []string{"Namespaces", "Cgroups v2", "runc", "crun", "Kernel"},
		},
	}

	for i := range courses {
		courses[i].LessonCount = s.countLessons(courses[i].Slug)
	}

	return courses
}

// GetCourseBySlug находит курс по слагу и загружает уроки из каталога content/courses/<slug>.
func (s *CourseService) GetCourseBySlug(slug string) (*model.Course, error) {
	if strings.Contains(slug, "..") || strings.ContainsAny(slug, `/\`) {
		return nil, fmt.Errorf("недопустимый идентификатор курса: %s", slug)
	}

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
	if strings.Contains(lessonSlug, "..") || strings.ContainsAny(lessonSlug, `/\`) {
		return nil, nil, fmt.Errorf("недопустимый идентификатор урока: %s", lessonSlug)
	}

	course, err := s.GetCourseBySlug(courseSlug)
	if err != nil {
		return nil, nil, err
	}

	courseDir := filepath.Clean(filepath.Join(s.contentDir, courseSlug))
	lessonPath := filepath.Clean(filepath.Join(courseDir, lessonSlug+".md"))
	if !strings.HasPrefix(lessonPath, courseDir+string(filepath.Separator)) {
		return nil, nil, fmt.Errorf("попытка выхода за пределы каталога курса")
	}

	cacheKey := courseSlug + "/" + lessonSlug
	s.lessonMu.RLock()
	if cached, ok := s.lessonCache[cacheKey]; ok {
		s.lessonMu.RUnlock()
		return course, cached, nil
	}
	s.lessonMu.RUnlock()

	contentBytes, err := os.ReadFile(lessonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("урок '%s' не найден: %w", lessonSlug, err)
	}

	var buf bytes.Buffer
	if err := s.mdParser.Convert(contentBytes, &buf); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга markdown: %w", err)
	}

	title := extractTitle(contentBytes, lessonSlug)

	lesson := &model.Lesson{
		CourseSlug:  courseSlug,
		Slug:        lessonSlug,
		Title:       title,
		Duration:    "15-20 мин",
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
				runes := []rune(cleaned)
				if len(runes) > 200 {
					cleaned = string(runes[:200]) + "..."
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

func (s *CourseService) countLessons(slug string) int {
	entries, err := os.ReadDir(filepath.Join(s.contentDir, slug))
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			count++
		}
	}
	return count
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
			rawBytes, _ := os.ReadFile(filepath.Join(courseDir, entry.Name()))
			title := extractTitle(rawBytes, slug)

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

// extractTitle извлекает заголовок h1 из первой строки markdown или форматирует slug.
func extractTitle(content []byte, fallbackSlug string) string {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	parts := strings.SplitN(fallbackSlug, "-", 2)
	cleanedSlug := fallbackSlug
	if len(parts) == 2 {
		if _, err := strconv.Atoi(parts[0]); err == nil {
			cleanedSlug = parts[1]
		}
	}
	title := strings.ReplaceAll(cleanedSlug, "-", " ")
	return titleCase(title)
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
