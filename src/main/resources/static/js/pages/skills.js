document.addEventListener('DOMContentLoaded', () => {
    const filterButtons = document.querySelectorAll('.filter-btn');
    const skillCards = document.querySelectorAll('.skill-card-item');

    if (!filterButtons.length || !skillCards.length) {
        return;
    }

    filterButtons.forEach(button => {
        button.addEventListener('click', () => {
            // Переключаем активный класс у кнопок
            filterButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            const selectedCategory = button.getAttribute('data-category');

            // Фильтруем карточки навыков
            skillCards.forEach(card => {
                const cardCategory = card.getAttribute('data-category');
                if (selectedCategory === 'ALL' || cardCategory === selectedCategory) {
                    card.style.display = 'block';
                } else {
                    card.style.display = 'none';
                }
            });
        });
    });
});