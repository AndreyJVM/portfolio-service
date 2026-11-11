/**
 * Primary Cloud Lab — Интерактивная система обучения:
 * - Учет и сохранение прогресса изучения курсов в LocalStorage
 * - Динамические прогресс-бары в каталоге курсов и на страницах модулей
 * - Интерактивные мини-квизы с мгновенной валидацией и закреплением материала
 */
(function () {
  const STORAGE_KEY = 'pcl_learning_progress_v1';

  function getStorageData() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return { completedLessons: [] };
      const parsed = JSON.parse(raw);
      if (!Array.isArray(parsed.completedLessons)) {
        return { completedLessons: [] };
      }
      return parsed;
    } catch (e) {
      console.warn('Ошибка чтения прогресса обучения из localStorage:', e);
      return { completedLessons: [] };
    }
  }

  function saveStorageData(data) {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    } catch (e) {
      console.warn('Ошибка сохранения прогресса обучения в localStorage:', e);
    }
  }

  function makeKey(courseSlug, lessonSlug) {
    return `${courseSlug}/${lessonSlug}`;
  }

  window.LearningProgress = {
    isCompleted(courseSlug, lessonSlug) {
      const data = getStorageData();
      return data.completedLessons.includes(makeKey(courseSlug, lessonSlug));
    },

    setCompleted(courseSlug, lessonSlug, completed = true) {
      const data = getStorageData();
      const key = makeKey(courseSlug, lessonSlug);
      const idx = data.completedLessons.indexOf(key);

      if (completed && idx === -1) {
        data.completedLessons.push(key);
      } else if (!completed && idx !== -1) {
        data.completedLessons.splice(idx, 1);
      }

      saveStorageData(data);
      this.updateUI();
    },

    toggleCompleted(courseSlug, lessonSlug) {
      const current = this.isCompleted(courseSlug, lessonSlug);
      const next = !current;
      this.setCompleted(courseSlug, lessonSlug, next);
      return next;
    },

    getCompletedCount(courseSlug) {
      const data = getStorageData();
      const prefix = `${courseSlug}/`;
      return data.completedLessons.filter(k => k.startsWith(prefix)).length;
    },

    getCoursePercent(courseSlug, totalLessons) {
      if (!totalLessons || totalLessons <= 0) return 0;
      const count = this.getCompletedCount(courseSlug);
      return Math.min(100, Math.round((count / totalLessons) * 100));
    },

    updateUI() {
      // 1. Обновление прогресс-баров в каталоге курсов (/courses)
      document.querySelectorAll('[data-course-progress]').forEach(container => {
        const courseSlug = container.getAttribute('data-course-slug');
        const total = parseInt(container.getAttribute('data-total-lessons'), 10) || 0;
        const percent = this.getCoursePercent(courseSlug, total);
        const count = this.getCompletedCount(courseSlug);

        const bar = container.querySelector('.progress-bar');
        const label = container.querySelector('.course-progress-label');
        const countEl = container.querySelector('.course-progress-count');

        if (bar) bar.style.width = `${percent}%`;
        if (label) {
          label.textContent = percent === 100 ? 'Завершено! 🏆' : `${percent}%`;
          if (percent === 100) {
            label.className = 'small fw-bold text-success course-progress-label';
          }
        }
        if (countEl) {
          countEl.textContent = `${count} из ${total}`;
        }
      });

      // 2. Обновление страницы курса (/courses/{slug})
      const detailContainer = document.getElementById('courseDetailProgress');
      if (detailContainer) {
        const slug = detailContainer.getAttribute('data-course-slug');
        const total = parseInt(detailContainer.getAttribute('data-total-lessons'), 10) || 0;
        const percent = this.getCoursePercent(slug, total);
        const count = this.getCompletedCount(slug);

        const bar = document.getElementById('courseDetailProgressBar');
        const percentEl = document.getElementById('courseDetailProgressPercent');
        const countEl = document.getElementById('courseDetailProgressCount');

        if (bar) bar.style.width = `${percent}%`;
        if (percentEl) {
          percentEl.textContent = percent === 100 ? '100% Завершено! 🏆' : `${percent}%`;
        }
        if (countEl) {
          countEl.textContent = `Пройдено уроков: ${count} из ${total}`;
        }
      }

      // Обновление бейджей статуса уроков в списке
      document.querySelectorAll('.lesson-card-item').forEach(card => {
        const courseSlug = card.getAttribute('data-course-slug');
        const lessonSlug = card.getAttribute('data-lesson-slug');
        const badge = card.querySelector('.lesson-status-badge');

        if (badge && courseSlug && lessonSlug) {
          const completed = this.isCompleted(courseSlug, lessonSlug);
          if (completed) {
            badge.className = 'badge lesson-status-badge bg-success-subtle text-success border border-success-subtle';
            badge.innerHTML = '<i class="fas fa-check-circle me-1"></i>Пройдено';
          } else {
            badge.className = 'badge lesson-status-badge bg-body-secondary text-body-secondary border';
            badge.innerHTML = '<i class="far fa-circle me-1"></i>К изучению';
          }
        }
      });

      // 3. Обновление кнопки на странице конкретного урока (/courses/{course}/{lesson})
      const completeBtn = document.getElementById('completeLessonBtn');
      if (completeBtn) {
        const courseSlug = completeBtn.getAttribute('data-course-slug');
        const lessonSlug = completeBtn.getAttribute('data-lesson-slug');
        if (courseSlug && lessonSlug) {
          const completed = this.isCompleted(courseSlug, lessonSlug);
          if (completed) {
            completeBtn.className = 'btn btn-success rounded-pill px-4 shadow-sm';
            completeBtn.innerHTML = '<i class="fas fa-check-double me-1"></i> Урок изучен';
          } else {
            completeBtn.className = 'btn btn-outline-success rounded-pill px-4';
            completeBtn.innerHTML = '<i class="far fa-check-circle me-1"></i> Отметить как пройденный';
          }
        }
      }
    }
  };

  // Инициализация интерактивных элементов после загрузки DOM
  document.addEventListener('DOMContentLoaded', () => {
    window.LearningProgress.updateUI();

    // Слушатель кнопки "Урок изучен"
    const completeBtn = document.getElementById('completeLessonBtn');
    if (completeBtn) {
      completeBtn.addEventListener('click', () => {
        const courseSlug = completeBtn.getAttribute('data-course-slug');
        const lessonSlug = completeBtn.getAttribute('data-lesson-slug');
        if (!courseSlug || !lessonSlug) return;

        const isNowCompleted = window.LearningProgress.toggleCompleted(courseSlug, lessonSlug);
        if (window.showToast) {
          if (isNowCompleted) {
            window.showToast('Отлично! Урок отмечен как пройденный.', 'success');
          } else {
            window.showToast('Отметка о прохождении урока снята.', 'info');
          }
        }
      });
    }

    // Инициализация интерактивных проверочных квизов
    initQuizzes();
  });

  function initQuizzes() {
    document.querySelectorAll('.quiz-option-btn').forEach(btn => {
      btn.addEventListener('click', function () {
        const parentQuestion = this.closest('.quiz-question');
        if (!parentQuestion) return;

        // Если вопрос уже успешно отвечен, не реагируем
        if (parentQuestion.getAttribute('data-answered') === 'correct') return;

        const isCorrect = this.getAttribute('data-correct') === 'true';
        const feedbackEl = parentQuestion.querySelector('.quiz-feedback');

        // Сбрасываем стили вариантов
        parentQuestion.querySelectorAll('.quiz-option-btn').forEach(b => {
          b.classList.remove('btn-success', 'btn-danger');
          b.classList.add('btn-outline-secondary');
        });

        if (isCorrect) {
          this.classList.remove('btn-outline-secondary');
          this.classList.add('btn-success');
          parentQuestion.setAttribute('data-answered', 'correct');

          if (feedbackEl) {
            feedbackEl.className = 'quiz-feedback alert alert-success mt-2 py-2 px-3 small border-0 rounded-3 d-flex align-items-center';
            feedbackEl.innerHTML = '<i class="fas fa-check-circle me-2 fa-lg"></i><span><strong>Верно!</strong> ' + (this.getAttribute('data-explanation') || 'Отличный ответ.') + '</span>';
            feedbackEl.classList.remove('d-none');
          }

          // Проверяем, пройдены ли все вопросы квиза
          const quizCard = parentQuestion.closest('.lesson-quiz-card');
          if (quizCard) {
            const allQuestions = quizCard.querySelectorAll('.quiz-question');
            const correctQuestions = quizCard.querySelectorAll('.quiz-question[data-answered="correct"]');
            if (allQuestions.length > 0 && allQuestions.length === correctQuestions.length) {
              const celebration = quizCard.querySelector('.quiz-celebration');
              if (celebration) celebration.classList.remove('d-none');

              const completeBtn = document.getElementById('completeLessonBtn');
              if (completeBtn) {
                const cSlug = completeBtn.getAttribute('data-course-slug');
                const lSlug = completeBtn.getAttribute('data-lesson-slug');
                if (cSlug && lSlug && !window.LearningProgress.isCompleted(cSlug, lSlug)) {
                  window.LearningProgress.setCompleted(cSlug, lSlug, true);
                  if (window.showToast) {
                    window.showToast('🎉 Поздравляем! Все вопросы квиза решены, урок зачтен!', 'success');
                  }
                }
              }
            }
          }
        } else {
          this.classList.remove('btn-outline-secondary');
          this.classList.add('btn-danger');

          if (feedbackEl) {
            feedbackEl.className = 'quiz-feedback alert alert-danger mt-2 py-2 px-3 small border-0 rounded-3 d-flex align-items-center';
            feedbackEl.innerHTML = '<i class="fas fa-circle-xmark me-2 fa-lg"></i><span><strong>Не совсем так.</strong> Попробуйте другой вариант ответа.</span>';
            feedbackEl.classList.remove('d-none');
          }
        }
      });
    });
  }
})();
