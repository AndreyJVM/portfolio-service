/**
 * Основной скрипт глобальной инициализации приложения:
 * - Инициализация компонентов Bootstrap (тултипы, поповеры)
 * - Опрос Spring Boot Actuator (/actuator/health) для живого статуса системы
 * - Глобальная система Toast-уведомлений (window.showToast)
 * - Микро-анимации появления элементов при скролле (Scroll Reveal / IntersectionObserver)
 */
document.addEventListener('DOMContentLoaded', () => {
  // 1. Инициализация тултипов Bootstrap
  const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
  tooltipTriggerList.forEach((el) => new bootstrap.Tooltip(el));

  // 2. Инициализация поповеров Bootstrap
  const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
  popoverTriggerList.forEach((el) => new bootstrap.Popover(el, { html: true }));

  // 3. Мониторинг доступности и здоровья системы (Actuator Health)
  initSystemHealthBadge();

  // 4. Плавное появление контента при прокрутке
  initScrollReveal();
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

function initSystemHealthBadge() {
  const badge = document.getElementById('systemStatusBadge');
  const statusPing = document.getElementById('statusPing');
  const statusText = document.getElementById('systemStatusText');

  if (!badge || !statusPing || !statusText) {
    return;
  }

  const startTime = performance.now();

  fetch('/actuator/health', { cache: 'no-store' } || {})
    .then(async (response) => {
      const latency = Math.round(performance.now() - startTime);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const data = await response.json();
      const isUp = data && data.status === 'UP';

      if (isUp) {
        statusText.textContent = 'Systems Online';
        statusPing.className = 'status-ping me-2';

        const dbStatus = data.components?.db?.status || 'UP';
        const diskStatus = data.components?.diskSpace?.status || 'UP';

        const popoverContent = `
          <div class="p-1" style="min-width: 200px;">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <span class="text-body-secondary small"><i class="fas fa-server text-primary me-1"></i>Backend:</span>
              <span class="badge bg-success-subtle text-success border border-success-subtle">UP</span>
            </div>
            <div class="d-flex justify-content-between align-items-center mb-1">
              <span class="text-body-secondary small"><i class="fas fa-database text-primary me-1"></i>Database:</span>
              <span class="badge bg-success-subtle text-success border border-success-subtle">${dbStatus}</span>
            </div>
            <div class="d-flex justify-content-between align-items-center mb-1">
              <span class="text-body-secondary small"><i class="fas fa-hard-drive text-primary me-1"></i>Disk:</span>
              <span class="badge bg-success-subtle text-success border border-success-subtle">${diskStatus}</span>
            </div>
            <div class="d-flex justify-content-between align-items-center mb-2">
              <span class="text-body-secondary small"><i class="fas fa-bolt text-warning me-1"></i>Latency:</span>
              <span class="text-primary fw-semibold small">${latency} ms</span>
            </div>
            <div class="border-top pt-1 text-center text-body-secondary" style="font-size: 0.7rem;">
              Spring Boot Actuator &bull; Real-time
            </div>
          </div>
        `;

        updateBadgePopover(badge, popoverContent);
      } else {
        setDegradedStatus(badge, statusPing, statusText, 'Degraded', latency);
      }
    })
    .catch((err) => {
      const latency = Math.round(performance.now() - startTime);
      setDegradedStatus(badge, statusPing, statusText, 'Offline', latency);
    });
}

function setDegradedStatus(badge, statusPing, statusText, label, latency) {
  statusText.textContent = label;
  statusPing.className = 'status-ping status-ping--down me-2';

  const popoverContent = `
    <div class="p-1 text-danger small">
      <div><strong>Внимание:</strong> Сервисы временно недоступны или работают с задержкой.</div>
      <div class="text-muted mt-1" style="font-size: 0.75rem;">Задержка: ${latency} ms</div>
    </div>
  `;
  updateBadgePopover(badge, popoverContent);
}

function updateBadgePopover(element, content) {
  element.setAttribute('data-bs-content', content);
  const existingPopover = bootstrap.Popover.getInstance(element);
  if (existingPopover) {
    existingPopover.setContent({
      '.popover-body': content
    });
  }
}

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
