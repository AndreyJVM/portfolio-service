/**
 * Primary Cloud Lab — Интерактивная фильтрация каталога модулей Docker
 */
(function() {
  'use strict';

  document.addEventListener('DOMContentLoaded', () => {
    initCourseFilters();
  });

  function initCourseFilters() {
    const filterButtons = document.querySelectorAll('#courseLevelFilters .filter-btn');
    const searchInput = document.getElementById('courseQuickFilter');
    const clearBtn = document.getElementById('clearCourseFilter');
    const resetBtn = document.getElementById('resetFiltersBtn');
    const counterBadge = document.getElementById('catalogCounterBadge');
    const noCoursesMsg = document.getElementById('noCoursesMessage');
    const courseCards = document.querySelectorAll('.course-card-wrapper');

    if (!courseCards.length) return;

    let activeLevel = 'all';
    let searchQuery = '';

    const totalCount = courseCards.length;

    // 1. Клик по кнопкам уровней
    filterButtons.forEach(btn => {
      btn.addEventListener('click', () => {
        filterButtons.forEach(b => {
          b.classList.remove('active', 'btn-primary');
          b.classList.add('btn-outline-secondary');
        });
        btn.classList.add('active', 'btn-primary');
        btn.classList.remove('btn-outline-secondary');

        activeLevel = btn.getAttribute('data-filter') || 'all';
        applyFilters();
      });
    });

    // 2. Живой поиск в каталоге
    if (searchInput) {
      searchInput.addEventListener('input', (e) => {
        searchQuery = (e.target.value || '').trim().toLowerCase();
        if (clearBtn) {
          if (searchQuery.length > 0) {
            clearBtn.classList.remove('d-none');
          } else {
            clearBtn.classList.add('d-none');
          }
        }
        applyFilters();
      });
    }

    // 3. Кнопка очистки текстового поиска
    if (clearBtn && searchInput) {
      clearBtn.addEventListener('click', () => {
        searchInput.value = '';
        searchQuery = '';
        clearBtn.classList.add('d-none');
        searchInput.focus();
        applyFilters();
      });
    }

    // 4. Кнопка полного сброса
    if (resetBtn) {
      resetBtn.addEventListener('click', () => {
        resetFilters();
      });
    }

    function resetFilters() {
      activeLevel = 'all';
      searchQuery = '';
      if (searchInput) {
        searchInput.value = '';
      }
      if (clearBtn) {
        clearBtn.classList.add('d-none');
      }

      filterButtons.forEach(b => {
        if (b.getAttribute('data-filter') === 'all') {
          b.classList.add('active', 'btn-primary');
          b.classList.remove('btn-outline-secondary');
        } else {
          b.classList.remove('active', 'btn-primary');
          b.classList.add('btn-outline-secondary');
        }
      });

      applyFilters();
    }

    function applyFilters() {
      let visibleCount = 0;

      courseCards.forEach(card => {
        const level = (card.getAttribute('data-level') || '').toLowerCase();
        const slug = (card.getAttribute('data-slug') || '').toLowerCase();
        const title = (card.getAttribute('data-title') || '').toLowerCase();
        const desc = (card.getAttribute('data-description') || '').toLowerCase();
        const tags = (card.getAttribute('data-tags') || '').toLowerCase();

        // Проверка фильтра по уровню
        let matchesLevel = false;
        if (activeLevel === 'all') {
          matchesLevel = true;
        } else if (activeLevel === 'basic') {
          matchesLevel = level.includes('базовый') || level.includes('уровень 1');
        } else if (activeLevel === 'middle') {
          matchesLevel = level.includes('middle') || level.includes('уровень 2');
        } else if (activeLevel === 'senior') {
          matchesLevel = level.includes('senior') || level.includes('уровень 3');
        } else if (activeLevel === 'cookbook') {
          matchesLevel = slug.includes('cookbook') || level.includes('справочник');
        }

        // Проверка текстового запроса
        let matchesSearch = true;
        if (searchQuery) {
          matchesSearch = title.includes(searchQuery) ||
                          desc.includes(searchQuery) ||
                          tags.includes(searchQuery) ||
                          slug.includes(searchQuery);
        }

        if (matchesLevel && matchesSearch) {
          card.classList.remove('d-none');
          card.classList.add('course-fade-in');
          visibleCount++;
        } else {
          card.classList.add('d-none');
          card.classList.remove('course-fade-in');
        }
      });

      // Обновление счетчика
      if (counterBadge) {
        if (visibleCount === totalCount && !searchQuery) {
          counterBadge.textContent = '10 модулей • 35 практических уроков';
          counterBadge.className = 'badge bg-body-tertiary text-body-secondary border px-3 py-2 rounded-pill';
        } else {
          counterBadge.textContent = `Показано: ${visibleCount} из ${totalCount} модулей`;
          counterBadge.className = 'badge bg-primary-subtle text-primary border border-primary-subtle px-3 py-2 rounded-pill';
        }
      }

      // Сообщение при пустом результате
      if (noCoursesMsg) {
        if (visibleCount === 0) {
          noCoursesMsg.classList.remove('d-none');
        } else {
          noCoursesMsg.classList.add('d-none');
        }
      }
    }
  }
})();
