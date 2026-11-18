/**
 * Primary Cloud Lab — Интерактивный функционал страницы урока:
 * - Автоматическая генерация интерактивного оглавления (Table of Contents)
 * - ScrollSpy (подсветка активного раздела при скролле)
 * - Панель чтения: управление размером шрифта (localStorage), расчет времени чтения
 * - Быстрое копирование ссылки на урок / секцию
 * - Клавиатурная навигация между уроками (Alt + ← / Alt + →)
 * - Интеграция с Prism.js для подсветки синтаксиса
 */
(function () {
  'use strict';

  document.addEventListener('DOMContentLoaded', () => {
    initTableOfContents();
    initReadingToolbar();
    initKeyboardNavigation();
    initCodeHighlighting();
    initMermaidDiagrams();
  });

  /**
   * 1. Генерация интерактивного оглавления (Desktop & Mobile) + ScrollSpy
   */
  function initTableOfContents() {
    const content = document.querySelector('.markdown-content');
    if (!content) return;

    const headings = content.querySelectorAll('h2, h3');
    const desktopContainer = document.getElementById('desktopTocContainer');
    const desktopList = document.getElementById('desktopTocList');
    const mobileList = document.getElementById('mobileTocList');
    const mobileTocWrapper = document.getElementById('mobileTocWrapper');
    const countBadge = document.getElementById('tocCountBadge');
    const mobileCountBadge = document.getElementById('mobileTocCountBadge');

    if (!headings || headings.length < 2) {
      if (desktopContainer) desktopContainer.style.display = 'none';
      if (mobileTocWrapper) mobileTocWrapper.style.display = 'none';
      return;
    }

    if (countBadge) countBadge.textContent = `${headings.length} разделов`;
    if (mobileCountBadge) mobileCountBadge.textContent = `${headings.length} разделов`;

    const tocLinks = [];

    headings.forEach((heading, index) => {
      // Гарантируем наличие уникального ID для якорной навигации
      if (!heading.id) {
        heading.id = 'section-' + (index + 1);
      }

      const titleText = heading.textContent.replace(/^[#\s]+/, '').trim();
      const isSub = heading.tagName.toLowerCase() === 'h3';

      // Desktop link
      if (desktopList) {
        const li = document.createElement('li');
        li.className = 'toc-item';
        const a = document.createElement('a');
        a.href = '#' + heading.id;
        a.className = 'toc-link' + (isSub ? ' toc-sublink' : '');
        a.title = titleText;
        a.textContent = titleText;
        a.addEventListener('click', (e) => {
          e.preventDefault();
          scrollToHeading(heading);
          history.pushState(null, '', '#' + heading.id);
        });
        li.appendChild(a);
        desktopList.appendChild(li);
        tocLinks.push({ heading, link: a });
      }

      // Mobile link
      if (mobileList) {
        const li = document.createElement('li');
        li.className = 'toc-item mb-1';
        const a = document.createElement('a');
        a.href = '#' + heading.id;
        a.className = 'toc-link' + (isSub ? ' toc-sublink' : '');
        a.textContent = titleText;
        a.addEventListener('click', (e) => {
          e.preventDefault();
          scrollToHeading(heading);
          history.pushState(null, '', '#' + heading.id);
          // Сворачиваем мобильный аккордеон оглавления
          const collapseEl = document.getElementById('mobileTocCollapse');
          if (collapseEl && window.bootstrap) {
            const bsCollapse = bootstrap.Collapse.getInstance(collapseEl);
            if (bsCollapse) bsCollapse.hide();
          }
        });
        li.appendChild(a);
        mobileList.appendChild(li);
      }
    });

    // ScrollSpy: подсветка активного раздела при прокрутке
    initScrollSpy(tocLinks);
  }

  function scrollToHeading(heading) {
    const yOffset = -90; // Отступ под фиксированную шапку
    const y = heading.getBoundingClientRect().top + window.pageYOffset + yOffset;
    window.scrollTo({ top: y, behavior: 'smooth' });
  }

  function initScrollSpy(tocLinks) {
    if (!tocLinks.length) return;

    let ticking = false;

    function updateActiveLink() {
      const scrollY = window.pageYOffset;
      const offsetMargin = 120;
      let activeIndex = -1;

      for (let i = 0; i < tocLinks.length; i++) {
        const top = tocLinks[i].heading.getBoundingClientRect().top + scrollY;
        if (scrollY >= top - offsetMargin) {
          activeIndex = i;
        } else {
          break;
        }
      }

      tocLinks.forEach((item, idx) => {
        if (idx === activeIndex) {
          item.link.classList.add('active');
        } else {
          item.link.classList.remove('active');
        }
      });

      ticking = false;
    }

    window.addEventListener('scroll', () => {
      if (!ticking) {
        window.requestAnimationFrame(updateActiveLink);
        ticking = true;
      }
    }, { passive: true });

    updateActiveLink();
  }

  /**
   * 2. Панель управления чтением: размер шрифта, расчет времени, копирование ссылки
   */
  function initReadingToolbar() {
    const content = document.querySelector('.markdown-content');
    if (!content) return;

    // 2.1 Настройка размера шрифта
    const FONT_STORAGE_KEY = 'pcl_lesson_font_size';
    const savedSize = localStorage.getItem(FONT_STORAGE_KEY) || 'md';

    function setFontSize(size) {
      content.classList.remove('font-size-sm', 'font-size-md', 'font-size-lg');
      content.classList.add('font-size-' + size);
      document.querySelectorAll('[data-font-size]').forEach(btn => {
        if (btn.getAttribute('data-font-size') === size) {
          btn.classList.add('active');
        } else {
          btn.classList.remove('active');
        }
      });
      try {
        localStorage.setItem(FONT_STORAGE_KEY, size);
      } catch (e) {}
    }

    setFontSize(savedSize);

    document.querySelectorAll('[data-font-size]').forEach(btn => {
      btn.addEventListener('click', () => {
        const size = btn.getAttribute('data-font-size');
        setFontSize(size);
      });
    });

    // 2.2 Копирование ссылки на урок
    const copyLinkBtn = document.getElementById('copyLessonLinkBtn');
    if (copyLinkBtn) {
      copyLinkBtn.addEventListener('click', async () => {
        const cleanUrl = window.location.origin + window.location.pathname;
        let copied = false;
        if (navigator.clipboard && window.isSecureContext) {
          try {
            await navigator.clipboard.writeText(cleanUrl);
            copied = true;
          } catch (e) {
            copied = false;
          }
        }
        if (!copied) {
          const ta = document.createElement('textarea');
          ta.value = cleanUrl;
          ta.style.position = 'fixed';
          ta.style.left = '-9999px';
          document.body.appendChild(ta);
          ta.select();
          try {
            copied = document.execCommand('copy');
          } catch (e) {}
          document.body.removeChild(ta);
        }

        if (window.showToast) {
          window.showToast('Ссылка на урок скопирована в буфер обмена!', 'success');
        }
      });
    }

    // 2.3 Динамический расчет времени чтения (180 слов/мин)
    const durationBadge = document.getElementById('lessonDurationBadge');
    if (durationBadge) {
      const text = content.innerText || '';
      const words = text.trim().split(/\s+/).length;
      if (words > 150) {
        const minutes = Math.max(2, Math.round(words / 180));
        durationBadge.innerHTML = `<i class="far fa-clock me-1"></i>~${minutes} мин чтения`;
      }
    }
  }

  /**
   * 3. Клавиатурная навигация (Alt + Left / Alt + Right)
   */
  function initKeyboardNavigation() {
    window.addEventListener('keydown', (e) => {
      // Игнорируем нажатия, если пользователь вводит текст в input/textarea
      if (['INPUT', 'TEXTAREA', 'SELECT'].includes(document.activeElement.tagName)) return;

      if (e.altKey && e.key === 'ArrowLeft') {
        const prevBtn = document.querySelector('.prev-lesson-btn');
        if (prevBtn && prevBtn.href) {
          e.preventDefault();
          window.location.href = prevBtn.href;
        }
      } else if (e.altKey && e.key === 'ArrowRight') {
        const nextBtn = document.querySelector('.next-lesson-btn');
        if (nextBtn && nextBtn.href) {
          e.preventDefault();
          window.location.href = nextBtn.href;
        }
      }
    });
  }

  /**
   * 4. Подсветка синтаксиса кода с помощью Prism.js и терминальная разметка
   */
  function initCodeHighlighting() {
    // Нормализуем языковые классы для Prism Autoloader
    document.querySelectorAll('.markdown-content pre code').forEach(code => {
      if (code.className.includes('language-dockerfile')) {
        code.classList.add('language-docker');
      }
      if (code.className.includes('language-sh')) {
        code.classList.add('language-bash');
      }
      if (code.className.includes('language-yml')) {
        code.classList.add('language-yaml');
      }
    });

    if (window.Prism) {
      if (window.Prism.plugins && window.Prism.plugins.autoloader) {
        window.Prism.plugins.autoloader.languages_path = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/';
      }
      window.Prism.highlightAll();
    }
  }

  /**
   * 5. Интерактивная векторная отрисовка диаграмм Mermaid
   */
  function initMermaidDiagrams() {
    const mermaidElements = document.querySelectorAll('.markdown-content code.language-mermaid');
    if (!mermaidElements.length || !window.mermaid) return;

    const isDarkMode = document.documentElement.getAttribute('data-bs-theme') === 'dark' ||
                       document.documentElement.getAttribute('data-theme') === 'dark';

    window.mermaid.initialize({
      startOnLoad: false,
      theme: isDarkMode ? 'dark' : 'neutral',
      securityLevel: 'loose',
      fontFamily: 'Inter, system-ui, sans-serif'
    });

    mermaidElements.forEach((codeEl, index) => {
      const preEl = codeEl.closest('pre');
      if (!preEl) return;

      const graphDefinition = codeEl.textContent.trim();
      const container = document.createElement('div');
      container.className = 'mermaid-chart-container text-center my-4 p-3 rounded-4 border bg-body-tertiary shadow-sm overflow-x-auto';
      container.id = 'mermaid-chart-' + (index + 1);

      preEl.parentNode.insertBefore(container, preEl);
      preEl.remove();

      try {
        window.mermaid.render('mermaid-svg-' + (index + 1), graphDefinition).then(({ svg }) => {
          container.innerHTML = svg;
        }).catch((err) => {
          console.warn('Ошибка рендеринга Mermaid:', err);
          container.innerHTML = '<pre class="text-danger small">' + graphDefinition + '</pre>';
        });
      } catch (err) {
        console.warn('Ошибка вызова Mermaid:', err);
        container.innerHTML = '<pre class="text-danger small">' + graphDefinition + '</pre>';
      }
    });
  }
})();
