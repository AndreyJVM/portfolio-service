/**
 * Основной скрипт глобальной инициализации приложения.
 */
document.addEventListener('DOMContentLoaded', () => {
  // Глобальная инициализация тултипов Bootstrap при необходимости
  const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
  tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl);
  });
});
