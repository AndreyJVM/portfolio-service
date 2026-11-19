package service

import (
	"bytes"
	"encoding/json"
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
	coursesFile string
	mdParser    goldmark.Markdown
	lessonCache map[string]*model.Lesson
	lessonMu    sync.RWMutex
}

var mdTagRegex = regexp.MustCompile(`(?s)<[^>]+>|[#*` + "`" + `~_\[\]()-]`)

func NewCourseService(contentDir string, coursesFile ...string) *CourseService {
	filePath := ""
	if len(coursesFile) > 0 && coursesFile[0] != "" {
		filePath = coursesFile[0]
	} else {
		filePath = filepath.Join(filepath.Dir(contentDir), "courses.json")
	}

	return &CourseService{
		contentDir:  contentDir,
		coursesFile: filePath,
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

// GetAllCourses returns a list of all available courses from courses.json or defaults.
func (s *CourseService) GetAllCourses() []model.Course {
	var courses []model.Course
	if s.coursesFile != "" {
		if data, err := os.ReadFile(s.coursesFile); err == nil && len(data) > 0 {
			var loaded []model.Course
			if err := json.Unmarshal(data, &loaded); err == nil && len(loaded) > 0 {
				courses = loaded
			}
		}
	}

	if len(courses) == 0 {
		courses = s.defaultCourses()
	}

	for i := range courses {
		courses[i].LessonCount = s.countLessons(courses[i].Slug)
	}

	return courses
}

func (s *CourseService) defaultCourses() []model.Course {
	return []model.Course{
		{
			ID:          "docker",
			Slug:        "docker",
			Title:       "Модуль 1: Архитектура и фундаментальные основы",
			Icon:        "fab fa-docker",
			Color:       "#0284c7",
			Level:       "Уровень 1 • Базовый",
			Description: "Физика контейнеризации против VM, стек dockerd → containerd → runc, спецификация OCI, 7 пространств имен Linux и запуск первого контейнера.",
			Tags:        []string{"Docker", "Architecture", "containerd", "runc", "OCI", "Namespaces"},
		},
		{
			ID:          "docker-dockerfile",
			Slug:        "docker-dockerfile",
			Title:       "Модуль 2: Сборка образов и искусство Dockerfile",
			Icon:        "fas fa-layer-group",
			Color:       "#0ea5e9",
			Level:       "Уровень 1 • Базовый → Middle",
			Description: "Анатомия инструкций, оптимизация кэширования слоев, Multi-stage сборка, Google Distroless, Scratch, Alpine и движок BuildKit со скрытыми секретами.",
			Tags:        []string{"Dockerfile", "Multi-stage", "BuildKit", "Scratch", "Distroless", "Cache"},
		},
		{
			ID:          "docker-storage",
			Slug:        "docker-storage",
			Title:       "Модуль 3: Тома, хранилище и постоянство данных",
			Icon:        "fas fa-hard-drive",
			Color:       "#3b82f6",
			Level:       "Уровень 2 • Middle",
			Description: "Архитектура драйвера Overlay2, Copy-on-Write (lowerdir/upperdir/merged), Named Volumes, Bind Mounts, tmpfs и разбор извечной ошибки Permission Denied (UID/GID).",
			Tags:        []string{"Storage", "Volumes", "Bind Mounts", "Overlay2", "CoW", "UID/GID"},
		},
		{
			ID:          "docker-networking",
			Slug:        "docker-networking",
			Title:       "Модуль 4: Сетевой стек и DNS Discovery",
			Icon:        "fas fa-network-wired",
			Color:       "#6366f1",
			Level:       "Уровень 2 • Middle",
			Description: "Сетевые драйверы bridge, host, macvlan, ipvlan, встроенный DNS-сервер Docker, межконтейнерная маршрутизация и взаимодействие с iptables/nftables и UFW.",
			Tags:        []string{"Networking", "Bridge", "DNS", "Service Discovery", "iptables", "Macvlan"},
		},
		{
			ID:          "docker-compose",
			Slug:        "docker-compose",
			Title:       "Модуль 5: Мультиконтейнерные стеки (Compose v2)",
			Icon:        "fas fa-cubes",
			Color:       "#8b5cf6",
			Level:       "Уровень 2 • Middle",
			Description: "Стандарты Compose v2 спецификации, изоляция сетей микросервисов, интерполяция .env переменных, healthchecks и управление зависимостями depends_on.",
			Tags:        []string{"Compose v2", "Microservices", "Healthcheck", "Environment", "depends_on"},
		},
		{
			ID:          "docker-security",
			Slug:        "docker-security",
			Title:       "Модуль 6: Безопасность и Production-харденинг",
			Icon:        "fas fa-shield-halved",
			Color:       "#ef4444",
			Level:       "Уровень 3 • Middle → Senior",
			Description: "Запуск от непривилегированного пользователя (Non-Root), сброс привилегий ядра cap_drop: [ALL], Read-Only rootfs, профили Seccomp/AppArmor и сканирование уязвимостей Trivy.",
			Tags:        []string{"Security", "Hardening", "Non-Root", "Capabilities", "Seccomp", "Trivy"},
		},
		{
			ID:          "docker-resources",
			Slug:        "docker-resources",
			Title:       "Модуль 7: Изоляция ядра и ресурсы (Cgroups v2)",
			Icon:        "fas fa-microchip",
			Color:       "#f59e0b",
			Level:       "Уровень 3 • Middle → Senior",
			Description: "Контроллеры Cgroups v2: ограничение памяти, CPU квоты, CPU Pinning (cpuset), предотвращение OOM Killer (Exit Code 137) и троттлинг дискового ввода-вывода (Block I/O).",
			Tags:        []string{"Cgroups v2", "Memory Limits", "CPU Quotas", "OOM Killer", "Exit 137", "Kernel"},
		},
		{
			ID:          "docker-production",
			Slug:        "docker-production",
			Title:       "Модуль 8: Тюнинг демона, логи и траблшутинг",
			Icon:        "fas fa-wrench",
			Color:       "#10b981",
			Level:       "Уровень 3 • Senior",
			Description: "Конфигурация /etc/docker/daemon.json, live-restore: true, драйверы логирования (json-file, journald, Loki), мониторинг cAdvisor и диагностика аварий (Exit коды 125–143).",
			Tags:        []string{"Production", "daemon.json", "live-restore", "Logging", "cAdvisor", "Troubleshooting"},
		},
		{
			ID:          "docker-cicd",
			Slug:        "docker-cicd",
			Title:       "Модуль 9: Multi-arch сборки и CI/CD",
			Icon:        "fas fa-code-branch",
			Color:       "#06b6d4",
			Level:       "Уровень 3 • Senior",
			Description: "Кроссплатформенная сборка под amd64/arm64 через Docker Buildx и QEMU, кэширование слоев в GitHub Actions (type=gha, type=registry) и Zero-Downtime Blue-Green деплой.",
			Tags:        []string{"Buildx", "Multi-arch", "QEMU", "CI/CD", "GitHub Actions", "Blue-Green"},
		},
		{
			ID:          "docker-cookbook",
			Slug:        "docker-cookbook",
			Title:       "Модуль 10: Практический справочник и рецепты",
			Icon:        "fas fa-book-bookmark",
			Color:       "#ec4899",
			Level:       "Справочник • Все уровни",
			Description: "Шпаргалка консольных команд (Cheat Sheet), безопасная чистка диска docker system prune, готовые production-шаблоны Dockerfile для Go, Node.js, Python, Rust.",
			Tags:        []string{"Cheat Sheet", "Recipes", "Best Practices", "Prune", "Templates"},
		},
	}
}

// GetCourseBySlug finds a course by its slug and loads lessons from content/courses/<slug>.
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

	// Backward compatibility alias: map "docker" to first docker course if not found
	if foundCourse == nil && slug == "docker" {
		for i := range courses {
			if strings.HasPrefix(courses[i].Slug, "docker") {
				foundCourse = &courses[i]
				break
			}
		}
	}

	// Backward compatibility: if not in manifest, but directory exists on disk, serve it dynamically
	if foundCourse == nil {
		courseDir := filepath.Join(s.contentDir, slug)
		if info, err := os.Stat(courseDir); err == nil && info.IsDir() {
			dynamicCourse := model.Course{
				ID:          slug,
				Slug:        slug,
				Title:       titleCase(strings.ReplaceAll(slug, "-", " ")),
				Icon:        "fas fa-book",
				Color:       "#0284c7",
				Level:       "Базовый",
				Description: fmt.Sprintf("Материалы по теме %s", slug),
			}
			foundCourse = &dynamicCourse
		}
	}

	if foundCourse == nil {
		return nil, fmt.Errorf("курс '%s' не найден", slug)
	}

	lessons, err := s.loadLessonsForCourse(foundCourse.Slug)
	if err == nil {
		foundCourse.Lessons = lessons
		foundCourse.LessonCount = len(lessons)
	}

	return foundCourse, nil
}

// lessonMeta stores metadata extracted from frontmatter or filename.
type lessonMeta struct {
	Title       string
	Description string
	Duration    string
	Level       string
	Order       int
	Tags        []string
}

// parseFrontmatter extracts optional YAML-like frontmatter delimited by "---".
func parseFrontmatter(raw []byte) (lessonMeta, []byte) {
	meta := lessonMeta{}
	content := string(raw)

	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return meta, raw
	}

	var bodyStart int
	var frontmatterStr string
	if idx := strings.Index(content[4:], "\n---\n"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 5
	} else if idx := strings.Index(content[4:], "\r\n---\r\n"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 7
	} else if idx := strings.Index(content[4:], "\n---"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 4
		if bodyStart < len(content) && content[bodyStart] == '\n' {
			bodyStart++
		}
	} else {
		return meta, raw
	}

	lines := strings.Split(frontmatterStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		switch key {
		case "title":
			meta.Title = val
		case "description":
			meta.Description = val
		case "duration":
			meta.Duration = val
		case "level":
			meta.Level = val
		case "order":
			if o, err := strconv.Atoi(val); err == nil {
				meta.Order = o
			}
		case "tags":
			val = strings.Trim(val, "[]")
			for _, tag := range strings.Split(val, ",") {
				t := strings.Trim(strings.TrimSpace(tag), "\"'")
				if t != "" {
					meta.Tags = append(meta.Tags, t)
				}
			}
		}
	}

	return meta, []byte(content[bodyStart:])
}

// estimateDuration calculates reading time based on word count.
func estimateDuration(body []byte) string {
	words := len(strings.Fields(string(body)))
	mins := words / 160
	if mins < 5 {
		mins = 10
	} else if mins > 45 {
		mins = 45
	}
	return fmt.Sprintf("%d мин", mins)
}

// extractOrderFromSlug parses numeric prefix (e.g., "01-basics" -> 1).
func extractOrderFromSlug(slug string, fallback int) int {
	parts := strings.SplitN(slug, "-", 2)
	if len(parts) > 0 {
		if num, err := strconv.Atoi(parts[0]); err == nil && num > 0 {
			return num
		}
	}
	return fallback
}

// GetLesson loads a specific course lesson and parses Markdown to HTML with in-memory caching.
func (s *CourseService) GetLesson(courseSlug, lessonSlug string) (*model.Course, *model.Lesson, error) {
	if strings.Contains(lessonSlug, "..") || strings.ContainsAny(lessonSlug, `/\`) {
		return nil, nil, fmt.Errorf("недопустимый идентификатор урока: %s", lessonSlug)
	}

	course, err := s.GetCourseBySlug(courseSlug)
	if err != nil {
		return nil, nil, err
	}

	courseDir := filepath.Clean(filepath.Join(s.contentDir, course.Slug))
	lessonPath := filepath.Clean(filepath.Join(courseDir, lessonSlug+".md"))
	if !strings.HasPrefix(lessonPath, courseDir+string(filepath.Separator)) {
		return nil, nil, fmt.Errorf("попытка выхода за пределы каталога курса")
	}

	// Fallback: if not found in primary course directory, look in the legacy "docker" directory
	if _, err := os.Stat(lessonPath); os.IsNotExist(err) {
		fallbackPath := filepath.Clean(filepath.Join(s.contentDir, "docker", lessonSlug+".md"))
		if _, err2 := os.Stat(fallbackPath); err2 == nil {
			lessonPath = fallbackPath
		}
	}

	cacheKey := course.Slug + "/" + lessonSlug
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

	meta, bodyBytes := parseFrontmatter(contentBytes)

	var buf bytes.Buffer
	if err := s.mdParser.Convert(bodyBytes, &buf); err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга markdown: %w", err)
	}

	title := meta.Title
	if title == "" {
		title = extractTitle(bodyBytes, lessonSlug)
	}

	duration := meta.Duration
	if duration == "" {
		duration = estimateDuration(bodyBytes)
	}

	description := meta.Description
	if description == "" {
		description = fmt.Sprintf("Практическое руководство и разбор темы: %s", title)
	}

	lesson := &model.Lesson{
		CourseSlug:  course.Slug,
		Slug:        lessonSlug,
		Title:       title,
		Order:       meta.Order,
		Duration:    duration,
		Description: description,
		Level:       meta.Level,
		Tags:        meta.Tags,
		ContentHTML: buf.String(),
	}

	s.lessonMu.Lock()
	s.lessonCache[cacheKey] = lesson
	s.lessonMu.Unlock()

	return course, lesson, nil
}

// BuildSearchIndex creates a flat search index across all knowledge base courses and lessons.
func (s *CourseService) BuildSearchIndex() []model.SearchItem {
	var index []model.SearchItem

	courses := s.GetAllCourses()
	for _, c := range courses {
		// Add course card itself
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

		// Load all course lessons
		lessons, err := s.loadLessonsForCourse(c.Slug)
		if err != nil {
			continue
		}

		for _, l := range lessons {
			lessonPath := filepath.Join(s.contentDir, c.Slug, l.Slug+".md")
			rawBytes, err := os.ReadFile(lessonPath)
			snippet := ""
			if err == nil {
				_, bodyBytes := parseFrontmatter(rawBytes)
				cleaned := mdTagRegex.ReplaceAllString(string(bodyBytes), " ")
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
			meta, body := parseFrontmatter(rawBytes)

			title := meta.Title
			if title == "" {
				title = extractTitle(body, slug)
			}

			order := meta.Order
			if order == 0 {
				order = extractOrderFromSlug(slug, idx+1)
			}

			duration := meta.Duration
			if duration == "" {
				duration = estimateDuration(body)
			}

			desc := meta.Description
			if desc == "" {
				desc = fmt.Sprintf("Практическое руководство и разбор темы %s", title)
			}

			lessons = append(lessons, model.Lesson{
				CourseSlug:  courseSlug,
				Slug:        slug,
				Title:       title,
				Order:       order,
				Duration:    duration,
				Description: desc,
				Level:       meta.Level,
				Tags:        meta.Tags,
			})
		}
	}

	sort.Slice(lessons, func(i, j int) bool {
		if lessons[i].Order != lessons[j].Order {
			return lessons[i].Order < lessons[j].Order
		}
		return lessons[i].Slug < lessons[j].Slug
	})

	return lessons, nil
}

// extractTitle extracts the h1 title from the first markdown line or formats the fallback slug.
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

// titleCase capitalizes the first letter of each word (replacement for deprecated strings.Title).
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
