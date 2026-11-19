package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

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
