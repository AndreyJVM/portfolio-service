/**
 * Full-Text Search Modal & Client-Side Fuzzy Search Engine
 * Primary Cloud Lab
 */
(function() {
  let searchData = null;
  let isLoading = false;

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

  function performSearch(query) {
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
    return text.replace(regex, '<mark class="p-0 bg-warning-subtle text-warning-emphasis fw-bold">$1</mark>');
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

    let html = '<div class="list-group list-group-flush">';
    items.forEach(item => {
      const isCourse = item.type === 'course';
      const typeBadge = isCourse
        ? '<span class="badge bg-primary-subtle text-primary border border-primary-subtle me-2">Курс</span>'
        : '<span class="badge bg-secondary-subtle text-secondary border border-secondary-subtle me-2">Урок</span>';

      const subtitle = isCourse
        ? (item.tags ? item.tags.join(' • ') : '')
        : `${item.courseTitle || ''}`;

      const textSnippet = item.snippet ? highlight(item.snippet, query) : '';

      html += `
        <a href="${item.url}" class="list-group-item list-group-item-action p-3 rounded-3 mb-2 border shadow-sm">
          <div class="d-flex w-100 justify-content-between align-items-center mb-1">
            <div class="d-flex align-items-center">
              ${typeBadge}
              <h6 class="mb-0 fw-bold">${highlight(item.title, query)}</h6>
            </div>
            <small class="text-body-secondary">${subtitle}</small>
          </div>
          <p class="mb-0 small text-body-secondary lh-sm">${textSnippet}</p>
        </a>
      `;
    });
    html += '</div>';

    searchResults.innerHTML = html;
  }
})();
