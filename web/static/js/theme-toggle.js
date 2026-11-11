/**
 * Управление темой оформления (Light / Dark).
 * Синхронизирует кастомные CSS-переменные и Bootstrap 5.3 (data-bs-theme).
 */
(function () {
    const STORAGE_KEY = 'theme';
    const THEME_DARK = 'dark';
    const THEME_LIGHT = 'light';

    function getPreferredTheme() {
        const savedTheme = localStorage.getItem(STORAGE_KEY);
        if (savedTheme) {
            return savedTheme;
        }
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? THEME_DARK : THEME_LIGHT;
    }

    function applyTheme(theme) {
        document.documentElement.setAttribute('data-bs-theme', theme);
        document.documentElement.setAttribute('data-theme', theme);

        const icon = document.getElementById('themeIcon');
        if (icon) {
            icon.className = theme === THEME_DARK ? 'fas fa-sun' : 'fas fa-moon';
        }
    }

    // Применяем тему немедленно при разборе <head>, чтобы избежать белого мерцания
    const initialTheme = getPreferredTheme();
    applyTheme(initialTheme);

    // Подключаем обработчик клика после построения DOM-дерева
    document.addEventListener('DOMContentLoaded', function () {
        applyTheme(getPreferredTheme());

        const toggleBtn = document.getElementById('themeToggle');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', function () {
                const currentTheme = document.documentElement.getAttribute('data-bs-theme') || THEME_LIGHT;
                const nextTheme = currentTheme === THEME_DARK ? THEME_LIGHT : THEME_DARK;

                localStorage.setItem(STORAGE_KEY, nextTheme);
                applyTheme(nextTheme);
            });
        }
    });
})();