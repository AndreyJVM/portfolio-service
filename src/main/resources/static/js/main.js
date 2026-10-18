/**
 * Основной скрипт глобальной инициализации приложения:
 * - Инициализация компонентов Bootstrap (тултипы, поповеры)
 * - Глобальная система Toast-уведомлений (window.showToast)
 * - Микро-анимации появления элементов при скролле (Scroll Reveal / IntersectionObserver)
 * - Индикатор прогресса чтения страницы и плавающая кнопка "Наверх"
 */
document.addEventListener('DOMContentLoaded', () => {
  // 1. Инициализация тултипов Bootstrap
  const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
  tooltipTriggerList.forEach((el) => new bootstrap.Tooltip(el));

  // 2. Инициализация поповеров Bootstrap
  const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
  popoverTriggerList.forEach((el) => new bootstrap.Popover(el, { html: true }));

  // 3. Плавное появление контента при прокрутке
  initScrollReveal();

  // 4. Индикатор прогресса чтения и кнопка "Наверх"
  initScrollControls();
});

/**
 * Глобальная функция показа стильных всплывающих Toast-уведомлений Bootstrap.
 * @param {string} message - Текст уведомления
 * @param {'success'|'danger'|'error'|'warning'|'info'} [type='success'] - Тип уведомления
 */
window.showToast = function (message, type = 'success') {
  const toastEl = document.getElementById('liveToast');
  const toastMsg = document.getElementById('toastMessage');
  const toastIcon = document.getElementById('toastIcon');

  if (!toastEl || !toastMsg || !toastIcon) {
    return;
  }

  toastMsg.textContent = message;

  if (type === 'danger' || type === 'error') {
    toastIcon.className = 'fas fa-circle-exclamation text-danger me-2 fa-lg';
  } else if (type === 'warning') {
    toastIcon.className = 'fas fa-triangle-exclamation text-warning me-2 fa-lg';
  } else if (type === 'info') {
    toastIcon.className = 'fas fa-circle-info text-primary me-2 fa-lg';
  } else {
    toastIcon.className = 'fas fa-circle-check text-success me-2 fa-lg';
  }

  const toast = bootstrap.Toast.getOrCreateInstance(toastEl, { delay: 3500 });
  toast.show();
};

/**
 * Инициализирует IntersectionObserver для элементов с классом .reveal.
 * Элементы плавно всплывают при попадании в зону видимости окна.
 */
function initScrollReveal() {
  const revealElements = document.querySelectorAll('.reveal');
  if (!revealElements.length) return;

  if (!('IntersectionObserver' in window)) {
    revealElements.forEach((el) => el.classList.add('revealed'));
    return;
  }

  const observer = new IntersectionObserver((entries, obs) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('revealed');
        obs.unobserve(entry.target);
      }
    });
  }, {
    threshold: 0.08,
    rootMargin: '0px 0px -25px 0px'
  });

  revealElements.forEach((el) => observer.observe(el));
}

/**
 * Инициализирует полосу прогресса прокрутки страницы и плавающую кнопку "Наверх".
 */
function initScrollControls() {
  const progressBar = document.getElementById('scrollProgressBar');
  const backToTopBtn = document.getElementById('backToTopBtn');

  if (!progressBar && !backToTopBtn) return;

  let ticking = false;

  function onScroll() {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const docHeight = document.documentElement.scrollHeight - document.documentElement.clientHeight;

    if (progressBar && docHeight > 0) {
      const scrollPercent = Math.min(100, Math.max(0, (scrollTop / docHeight) * 100));
      progressBar.style.width = `${scrollPercent}%`;
    }

    if (backToTopBtn) {
      if (scrollTop > 300) {
        backToTopBtn.classList.add('show');
      } else {
        backToTopBtn.classList.remove('show');
      }
    }

    ticking = false;
  }

  window.addEventListener('scroll', () => {
    if (!ticking) {
      window.requestAnimationFrame(onScroll);
      ticking = true;
    }
  }, { passive: true });

  if (backToTopBtn) {
    backToTopBtn.addEventListener('click', () => {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    });
  }

  onScroll();
}
