package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"ru.vorobevaqa/devops-platform/internal/model"
)

// CoursesList: Каталог обучающих треков (алиас главной страницы)
func (h *Handler) CoursesList(w http.ResponseWriter, r *http.Request) {
	h.Index(w, r)
}

// CourseDetail: Страница конкретного курса со списком уроков
func (h *Handler) CourseDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "courseSlug")
	course, err := h.courseService.GetCourseBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Meta   model.MetaData
		Course *model.Course
	}{
		Meta: model.MetaData{
			Title:         course.Title + " • Курс DevOps",
			Description:   course.Description,
			Keywords:      fmt.Sprintf("%s, devops, уроки, обучение", course.Title),
			OGTitle:       course.Title,
			OGDescription: course.Description,
			OGImage:       h.baseURL + "/images/og-main.png",
			OGURL:         fmt.Sprintf("%s/courses/%s", h.baseURL, slug),
			ActiveTab:     "courses",
			CanonicalURL:  fmt.Sprintf("%s/courses/%s", h.baseURL, slug),
		},
		Course: course,
	}

	h.renderTemplate(w, "course_detail.html", data)
}

// Lesson: Страница просмотра и изучения отдельного урока
func (h *Handler) Lesson(w http.ResponseWriter, r *http.Request) {
	courseSlug := chi.URLParam(r, "courseSlug")
	lessonSlug := chi.URLParam(r, "lessonSlug")

	course, lesson, err := h.courseService.GetLesson(courseSlug, lessonSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var prevLesson, nextLesson *model.Lesson
	for i, l := range course.Lessons {
		if l.Slug == lessonSlug || l.Slug == lesson.Slug {
			if i > 0 {
				prevLesson = &course.Lessons[i-1]
			}
			if i+1 < len(course.Lessons) {
				nextLesson = &course.Lessons[i+1]
			}
			break
		}
	}

	data := struct {
		Meta              model.MetaData
		Course            *model.Course
		Lesson            *model.Lesson
		LessonContentSafe template.HTML
		PrevLesson        *model.Lesson
		NextLesson        *model.Lesson
	}{
		Meta: model.MetaData{
			Title:         lesson.Title + " • " + course.Title,
			Description:   lesson.Description,
			Keywords:      fmt.Sprintf("%s, %s, devops", lesson.Title, course.Title),
			OGTitle:       lesson.Title,
			OGDescription: lesson.Description,
			OGImage:       h.baseURL + "/images/og-main.png",
			OGURL:         fmt.Sprintf("%s/courses/%s/%s", h.baseURL, courseSlug, lessonSlug),
			ActiveTab:     "courses",
			CanonicalURL:  fmt.Sprintf("%s/courses/%s/%s", h.baseURL, courseSlug, lessonSlug),
		},
		Course:            course,
		Lesson:            lesson,
		LessonContentSafe: template.HTML(lesson.ContentHTML),
		PrevLesson:        prevLesson,
		NextLesson:        nextLesson,
	}

	h.renderTemplate(w, "lesson.html", data)
}
