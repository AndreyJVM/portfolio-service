/**
 * Модуль фильтрации и поиска навыков в матрице компетенций.
 * Поддерживает одновременную фильтрацию по категории и полнотекстовый поиск
 * по названию навыка, описанию и списку используемых инструментов.
 */
document.addEventListener('DOMContentLoaded', () => {
  const filterButtons = document.querySelectorAll('.filter-btn');
  const skillCards = document.querySelectorAll('.skill-card-item');
  const searchInput = document.getElementById('skillsSearchInput');
  const visibleCountEl = document.getElementById('skillsVisibleCount');
  const totalCountEl = document.getElementById('skillsTotalCount');
  const noSkillsFoundEl = document.getElementById('noSkillsFound');
  const resetBtn = document.getElementById('skillsResetBtn');

  if (!skillCards.length) {
    return;
  }

  let currentCategory = 'ALL';
  let currentSearchQuery = '';

  const totalCards = skillCards.length;
  if (totalCountEl) totalCountEl.textContent = totalCards;
  if (visibleCountEl) visibleCountEl.textContent = totalCards;

  function filterSkills() {
    let visibleCount = 0;
    const query = currentSearchQuery.toLowerCase().trim();

    skillCards.forEach((card) => {
      const cardCategory = card.getAttribute('data-category') || '';
      const cardText = card.textContent.toLowerCase();

      const matchesCategory = currentCategory === 'ALL' || cardCategory === currentCategory;
      const matchesSearch = !query || cardText.includes(query);

      if (matchesCategory && matchesSearch) {
        card.style.display = 'block';
        visibleCount++;
      } else {
        card.style.display = 'none';
      }
    });

    if (visibleCountEl) {
      visibleCountEl.textContent = visibleCount;
    }

    if (noSkillsFoundEl) {
      noSkillsFoundEl.style.display = visibleCount === 0 ? 'block' : 'none';
    }

    if (resetBtn) {
      const isFiltered = currentCategory !== 'ALL' || query.length > 0;
      resetBtn.style.display = isFiltered ? 'inline-block' : 'none';
    }
  }

  // Обработчики кнопок категорий
  filterButtons.forEach((button) => {
    button.addEventListener('click', () => {
      filterButtons.forEach((btn) => btn.classList.remove('active'));
      button.classList.add('active');

      currentCategory = button.getAttribute('data-category') || 'ALL';
      filterSkills();
    });
  });

  // Обработчик строки поиска с мгновенным откликом
  if (searchInput) {
    searchInput.addEventListener('input', () => {
      currentSearchQuery = searchInput.value;
      filterSkills();
    });
  }

  // Кнопка сброса фильтров
  if (resetBtn) {
    resetBtn.addEventListener('click', () => {
      currentCategory = 'ALL';
      currentSearchQuery = '';
      if (searchInput) searchInput.value = '';

      filterButtons.forEach((btn) => {
        if (btn.getAttribute('data-category') === 'ALL') {
          btn.classList.add('active');
        } else {
          btn.classList.remove('active');
        }
      });

      filterSkills();
    });
  }

  // Первоначальный подсчет
  filterSkills();
});
