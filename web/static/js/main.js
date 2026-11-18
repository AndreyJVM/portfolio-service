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

  // 5. Кнопки копирования кода и терминальная стилизация
  initCodeCopyButtons();

  // 6. Улучшения мобильной навигации
  initMobileNavEnhancements();
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
 * Элементы, уже находящиеся в поле зрения, отображаются мгновенно без задержки.
 */
function initScrollReveal() {
  const revealElements = document.querySelectorAll('.reveal');
  if (!revealElements.length) return;

  if (!('IntersectionObserver' in window)) {
    revealElements.forEach((el) => el.classList.add('revealed'));
    return;
  }

  // Мгновенно отображаем элементы, которые уже находятся на экране при загрузке
  const windowHeight = window.innerHeight || document.documentElement.clientHeight;
  revealElements.forEach((el) => {
    const rect = el.getBoundingClientRect();
    if (rect.top < windowHeight * 0.95) {
      el.classList.add('revealed');
    }
  });

  const observer = new IntersectionObserver((entries, obs) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('revealed');
        obs.unobserve(entry.target);
      }
    });
  }, {
    threshold: 0.02,
    rootMargin: '0px 0px 50px 0px'
  });

  revealElements.forEach((el) => {
    if (!el.classList.contains('revealed')) {
      observer.observe(el);
    }
  });
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

/**
 * Оборачивает блоки <pre> в карточки с терминальной шапкой и кнопкой быстрого копирования кода.
 */
function initCodeCopyButtons() {
  const preElements = document.querySelectorAll('.markdown-content pre');
  preElements.forEach((pre) => {
    if (pre.closest('.code-block-wrapper')) return;
    if (pre.querySelector('code.language-mermaid')) return;

    // Создаем контейнер-обертку
    const wrapper = document.createElement('div');
    wrapper.className = 'code-block-wrapper';

    // Создаем шапку
    const header = document.createElement('div');
    header.className = 'code-block-header';

    const dots = document.createElement('div');
    dots.className = 'code-block-dots';
    dots.innerHTML = '<span class="code-block-dot dot-red"></span><span class="code-block-dot dot-yellow"></span><span class="code-block-dot dot-green"></span>';

    // Определяем язык/подпись из class="language-bash" или подобного
    const codeEl = pre.querySelector('code');
    let lang = 'TERMINAL';
    let lineCount = 0;

    const rawText = (codeEl || pre).innerText.trim();
    if (rawText) {
      lineCount = rawText.split('\n').length;
    }

    if (codeEl) {
      const match = codeEl.className.match(/language-([a-zA-Z0-9_-]+)/);
      if (match) {
        const rawLang = match[1].toLowerCase();
        const langMap = {
          bash: 'BASH / SHELL',
          sh: 'SHELL',
          dockerfile: 'DOCKERFILE',
          docker: 'DOCKERFILE',
          yaml: 'YAML',
          yml: 'YAML',
          go: 'GO',
          golang: 'GO',
          json: 'JSON',
          ini: 'SYSTEMD / INI',
          systemd: 'SYSTEMD',
          c: 'C',
          mermaid: 'DIAGRAM',
          text: 'TEXT'
        };
        lang = langMap[rawLang] || rawLang.toUpperCase();
      }
    }

    const titleSpan = document.createElement('span');
    titleSpan.className = 'small fw-medium text-muted font-monospace code-lang-tag';
    titleSpan.textContent = lang;

    const linesSpan = document.createElement('span');
    linesSpan.className = 'small text-muted font-monospace d-none d-sm-inline code-line-count';
    linesSpan.style.fontSize = '0.7rem';
    linesSpan.style.opacity = '0.7';
    linesSpan.textContent = lineCount > 1 ? `${lineCount} строк` : '';

    const copyBtn = document.createElement('button');
    copyBtn.type = 'button';
    copyBtn.className = 'code-copy-btn';
    copyBtn.setAttribute('aria-label', 'Скопировать код');
    copyBtn.innerHTML = '<i class="far fa-copy"></i><span>Копировать</span>';

    copyBtn.addEventListener('click', async () => {
      const textToCopy = (codeEl || pre).innerText.trim();
      let success = false;

      if (navigator.clipboard && window.isSecureContext) {
        try {
          await navigator.clipboard.writeText(textToCopy);
          success = true;
        } catch (e) {
          success = false;
        }
      }

      if (!success) {
        // Fallback через скрытый textarea
        const textArea = document.createElement('textarea');
        textArea.value = textToCopy;
        textArea.style.position = 'fixed';
        textArea.style.left = '-9999px';
        document.body.appendChild(textArea);
        textArea.select();
        try {
          success = document.execCommand('copy');
        } catch (e) {
          success = false;
        }
        document.body.removeChild(textArea);
      }

      if (success) {
        copyBtn.classList.add('copied');
        copyBtn.innerHTML = '<i class="fas fa-check"></i><span>Скопировано!</span>';
        if (typeof window.showToast === 'function') {
          window.showToast('Код скопирован в буфер обмена', 'success');
        }
        setTimeout(() => {
          copyBtn.classList.remove('copied');
          copyBtn.innerHTML = '<i class="far fa-copy"></i><span>Копировать</span>';
        }, 2000);
      }
    });

    const headerRight = document.createElement('div');
    headerRight.className = 'd-flex align-items-center gap-2';
    if (linesSpan.textContent) {
      headerRight.appendChild(linesSpan);
    }
    headerRight.appendChild(titleSpan);
    headerRight.appendChild(copyBtn);

    header.appendChild(dots);
    header.appendChild(headerRight);

    pre.parentNode.insertBefore(wrapper, pre);
    wrapper.appendChild(header);
    wrapper.appendChild(pre);
  });
}

/**
 * Автоматическое закрытие мобильного offcanvas-сайдбара при переходе по ссылке.
 */
function initMobileNavEnhancements() {
  const sidebarNav = document.getElementById('sidebarNav');
  if (!sidebarNav) return;

  sidebarNav.querySelectorAll('.nav-link').forEach((link) => {
    link.addEventListener('click', () => {
      const bsOffcanvas = bootstrap.Offcanvas.getInstance(sidebarNav);
      if (bsOffcanvas) {
        bsOffcanvas.hide();
      }
    });
  });
}

