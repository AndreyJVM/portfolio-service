/**
 * Full-Text Search Modal & Client-Side Fuzzy Search Engine
 * Primary Cloud Lab
 */
(function() {
  let searchData = null;
  let isLoading = false;
  let selectedIndex = -1;

  const searchTriggerMobile = document.getElementById('searchTriggerMobile');
  const searchTriggerSidebar = document.getElementById('searchTriggerSidebar');
  const searchModalEl = document.getElementById('searchModal');
  const searchInput = document.getElementById('searchInput');
  const searchResults = document.getElementById('searchResults');
  const searchLoading = document.getElementById('searchLoading');
  const searchEmpty = document.getElementById('searchEmpty');

  if (!searchModalEl || !searchInput) return;

  const searchModal = new bootstrap.Modal(searchModalEl);

  // Клавиатурное сочетание Ctrl+K или Cmd+K
  document.addEventListener('keydown', function(e) {
    if ((e.ctrlKey || e.metaKey) && (e.key === 'k' || e.key === 'K')) {
      e.preventDefault();
      searchModal.show();
    }
  });

  if (searchTriggerMobile) {
    searchTriggerMobile.addEventListener('click', () => searchModal.show());
  }
  if (searchTriggerSidebar) {
    searchTriggerSidebar.addEventListener('click', () => searchModal.show());
  }

  searchModalEl.addEventListener('shown.bs.modal', function() {
    searchInput.focus();
    loadSearchIndex();
  });

  searchModalEl.addEventListener('hidden.bs.modal', function() {
    selectedIndex = -1;
  });

  async function loadSearchIndex() {
    if (searchData || isLoading) return;
    isLoading = true;
    if (searchLoading) searchLoading.classList.remove('d-none');

    try {
      const res = await fetch('/api/search');
      if (res.ok) {
        searchData = await res.json();
      }
    } catch (e) {
      console.warn('Failed to load search index:', e);
    } finally {
      isLoading = false;
      if (searchLoading) searchLoading.classList.add('d-none');
      performSearch(searchInput.value.trim());
    }
  }

  searchInput.addEventListener('input', function() {
    performSearch(this.value.trim());
  });

  searchInput.addEventListener('keydown', function(e) {
    const items = searchResults.querySelectorAll('.search-result-item');
    if (!items.length) return;

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      selectedIndex = (selectedIndex + 1) % items.length;
      updateActiveItem(items);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      selectedIndex = (selectedIndex - 1 + items.length) % items.length;
      updateActiveItem(items);
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (selectedIndex >= 0 && selectedIndex < items.length) {
        items[selectedIndex].click();
      } else if (items.length > 0) {
        items[0].click();
      }
    }
  });

  function updateActiveItem(items) {
    items.forEach((item, idx) => {
      if (idx === selectedIndex) {
        item.classList.add('active');
        item.setAttribute('aria-selected', 'true');
        item.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
      } else {
        item.classList.remove('active');
        item.removeAttribute('aria-selected');
      }
    });
  }

  function performSearch(query) {
    selectedIndex = -1;
    if (!searchData) return;

    if (!query) {
      searchResults.innerHTML = '';
      if (searchEmpty) searchEmpty.classList.add('d-none');
      return;
    }

    const q = query.toLowerCase();
    const matches = searchData.filter(item => {
      return (item.title && item.title.toLowerCase().includes(q)) ||
             (item.courseTitle && item.courseTitle.toLowerCase().includes(q)) ||
             (item.description && item.description.toLowerCase().includes(q)) ||
             (item.snippet && item.snippet.toLowerCase().includes(q)) ||
             (item.tags && item.tags.some(t => t.toLowerCase().includes(q)));
    });

    renderResults(matches, query);
  }

  function highlight(text, query) {
    if (!text || !query) return text || '';
    const regex = new RegExp(`(${escapeRegex(query)})`, 'gi');
    return text.replace(regex, '<mark class="p-0 bg-warning-subtle text-warning-emphasis fw-bold rounded px-1">$1</mark>');
  }

  function escapeRegex(string) {
    return string.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
  }

  function renderResults(items, query) {
    if (items.length === 0) {
      searchResults.innerHTML = '';
      if (searchEmpty) searchEmpty.classList.remove('d-none');
      return;
    }

    if (searchEmpty) searchEmpty.classList.add('d-none');

    let html = '<div class="list-group list-group-flush" role="listbox">';
    items.forEach((item, index) => {
      const isCourse = item.type === 'course';
      const icon = isCourse
        ? '<i class="fab fa-docker text-primary me-2"></i>'
        : '<i class="fas fa-file-code text-info me-2"></i>';
      const typeBadge = isCourse
        ? '<span class="badge bg-primary-subtle text-primary border border-primary-subtle me-2">Модуль</span>'
        : '<span class="badge bg-info-subtle text-info-emphasis border border-info-subtle me-2">Урок</span>';

      const subtitle = isCourse
        ? (item.tags ? item.tags.slice(0, 4).join(' • ') : '')
        : `${item.courseTitle || ''}`;

      const textSnippet = item.snippet ? highlight(item.snippet, query) : '';

      html += `
        <a href="${item.url}" class="list-group-item list-group-item-action p-3 rounded-3 mb-2 border shadow-sm search-result-item" data-index="${index}" role="option">
          <div class="d-flex w-100 justify-content-between align-items-center mb-1">
            <div class="d-flex align-items-center flex-grow-1 text-truncate pe-2">
              ${icon}
              ${typeBadge}
              <h6 class="mb-0 fw-bold text-truncate">${highlight(item.title, query)}</h6>
            </div>
            <div class="d-flex align-items-center gap-2 flex-shrink-0">
              <small class="text-body-secondary">${subtitle}</small>
              <i class="fas fa-arrow-right text-body-tertiary ms-1 search-arrow-indicator small"></i>
            </div>
          </div>
          ${textSnippet ? `<p class="mb-0 small text-body-secondary lh-sm text-truncate">${textSnippet}</p>` : ''}
        </a>
      `;
    });
    html += '</div>';

    searchResults.innerHTML = html;

    const resultElements = searchResults.querySelectorAll('.search-result-item');
    resultElements.forEach((el, idx) => {
      el.addEventListener('mouseenter', () => {
        selectedIndex = idx;
        updateActiveItem(resultElements);
      });
    });
  }
})();
