/**
 * Модуль интерактивного просмотра архитектуры проектов.
 * Динамически заполняет модальное окно схемой, компонентами и техническими деталями.
 */
const PROJECT_ARCH_DATA = {
  'portfolio': {
    title: 'Portfolio & Production Stack',
    badge: 'Production Active',
    badgeClass: 'bg-success-subtle text-success border border-success-subtle',
    summary: 'Отказоустойчивая микросервисная архитектура с обратным прокси, SSL-терминацией и полностью автоматизированным CI/CD пайплайном.',
    diagram: `
+--------------------------------------------------------------+
|                         INTERNET                             |
|       (HTTPS vorobevaqa.ru:443 / HTTP:80 Auto-Redirect)      |
+------------------------------+-------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|               NGINX REVERSE PROXY & SSL (VDS)                |
|  - Let's Encrypt Automated Certbot (TLSv1.3, HSTS, Gzip)     |
|  - Rate Limiting & Security Headers (CSP, X-Frame-Options)   |
|  - Reverse Proxy -> http://portfolio-app:8080                |
+------------------------------+-------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|              SPRING BOOT 3 APP CONTAINER (JVM 21)            |
|  [Spring MVC / Thymeleaf]  <-->  [REST Controller (/api/qr)] |
|  [Spring Boot Actuator]    <-->  [Health & Metrics Probe]    |
|  [ZXing 3.5.3 Engine]      <-->  [BitMatrix PNG Encoder]     |
+--------------------------------------------------------------+
                               ^
                               | (Deploy Trigger)
+--------------------------------------------------------------+
|              CI/CD PIPELINE (GITHUB ACTIONS)                 |
|  git push main -> mvn test -> Docker Build -> Docker Hub     |
|  -> SSH Deploy Script -> docker compose up -d (Zero-Downtime)|
+--------------------------------------------------------------+
    `.trim(),
    highlights: [
      { icon: 'fa-shield-halved', title: 'Безопасность', text: 'TLS 1.3, HSTS, строгие CSP-политики и изолированный запуск от непривилегированного пользователя.' },
      { icon: 'fa-gauge-high', title: 'Производительность', text: 'Java 21 Virtual Threads ready, кэширование статических ассетов с ETag и gzip сжатие на Nginx.' },
      { icon: 'fa-infinity', title: 'CI/CD Автоматизация', text: 'Многоэтапный Dockerfile (Multi-stage build), запуск JUnit тестов на каждый коммит и авто-деплой.' },
      { icon: 'fa-heart-pulse', title: 'Мониторинг', text: 'Spring Boot Actuator health endpoint с живым опросом состояния компонентов из фронтенда.' }
    ],
    stack: ['Java 21', 'Spring Boot 3.2', 'Thymeleaf', 'Docker & Compose', 'Nginx SSL', 'GitHub Actions', 'Actuator']
  },

  'qr': {
    title: 'QR Code Generator Engine',
    badge: 'Microservice API',
    badgeClass: 'bg-primary-subtle text-primary border border-primary-subtle',
    summary: 'Высокопроизводительный сервис генерации динамических QR-кодов по REST API с буферизацией и клиентом интерактивной кастомизации.',
    diagram: `
+--------------------------------------------------------------+
|                   BROWSER CLIENT / API CONSUMER              |
|  [Input URL] -> [Debounce 450ms] -> [Fetch POST /api/qr]     |
|  [Canvas Engine] <- (Render Hex Color Tint & Export PNG/Blob)|
+------------------------------+-------------------------------+
                               |
                               v (JSON: {"url": "..."})
+--------------------------------------------------------------+
|                  SPRING BOOT REST CONTROLLER                 |
|  - QrRestController.generateQrCode(@Valid @RequestBody)      |
|  - Bean Validation (@NotBlank, URL pattern check)            |
|  - GlobalExceptionHandler (RFC 7807 Problem Detail mapping)  |
+------------------------------+-------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|                     QR CODE SERVICE (ZXing)                  |
|  - QRCodeWriter.encode(url, BarcodeFormat.QR_CODE, w, h)     |
|  - CharacterSet UTF-8 & Margin Hint Optimization             |
|  - MatrixToImageWriter -> ByteArrayOutputStream (PNG)        |
|  - Base64.getEncoder().encodeToString(bytes)                 |
+--------------------------------------------------------------+
    `.trim(),
    highlights: [
      { icon: 'fa-bolt', title: 'Субмиллисекундная обработка', text: 'Генерация матричного растра за ~10-15 мс без сохранения промежуточных файлов на диск.' },
      { icon: 'fa-palette', title: 'Клиентская стилизация', text: 'Манипуляция пикселями через HTML5 Canvas Pixel Array без повторных нагрузок на сервер.' },
      { icon: 'fa-check-double', title: 'Строгая валидация', text: 'Контроль протоколов http/https, защита от XSS и инъекций на уровне контроллера.' },
      { icon: 'fa-cubes', title: 'Масштабируемость', text: 'Stateless-сервис без состояния сессий, легко масштабируемый горизонтально.' }
    ],
    stack: ['Java 21', 'Spring Boot 3', 'Google ZXing', 'REST API', 'Jakarta Validation', 'HTML5 Canvas API']
  },

  'samba': {
    title: 'Samba Web UI Manager',
    badge: 'Infrastructure Tool',
    badgeClass: 'bg-info-subtle text-info border border-info-subtle',
    summary: 'Панель администрирования файловых шар Linux/Samba с управлением учетными записями, ACL-правами и мониторингом подключений.',
    diagram: `
+--------------------------------------------------------------+
|                      ADMINISTRATOR GUI                       |
|  Web Dashboard: Share Management, Active Sessions, Quotas    |
+------------------------------+-------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|               DAEMON & API SERVICE LAYER                     |
|  - Auth & Role-based Access Control (Admin / Read-only)      |
|  - smb.conf AST Parser & atomic file writer (safe reloads)   |
|  - Subprocess runner: smbstatus, pdbedit, smbpasswd          |
+------------------------------+-------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|                   HOST LINUX OS / SAMBA                      |
|  - Samba Daemon (smbd / nmbd) serving SMB2/SMB3 clients      |
|  - Local Linux Filesystem (ext4/zfs) with POSIX ACLs         |
+--------------------------------------------------------------+
    `.trim(),
    highlights: [
      { icon: 'fa-user-lock', title: 'Управление доступом', text: 'Создание пользователей Samba с изоляцией домашних каталогов и ролевыми правами.' },
      { icon: 'fa-file-shield', title: 'Атомарные правки конфигурации', text: 'Проверка синтаксиса перед применением smb.conf исключает сбои файлового сервера.' },
      { icon: 'fa-chart-line', title: 'Мониторинг сессий', text: 'Инспекция активных клиентских блокировок и открытых дескрипторов в реальном времени.' },
      { icon: 'fa-box', title: 'Контейнерное развёртывание', text: 'Запуск сервиса в Docker-контейнере с доступом к хостовому сокету Samba.' }
    ],
    stack: ['Linux', 'Samba (SMB/CIFS)', 'Docker', 'Shell / Bash', 'REST / Web UI', 'POSIX ACL']
  }
};

window.openArchModal = function(projectKey) {
  const data = PROJECT_ARCH_DATA[projectKey];
  if (!data) return;

  const titleEl = document.getElementById('archModalTitle');
  const badgeEl = document.getElementById('archModalBadge');
  const summaryEl = document.getElementById('archModalSummary');
  const diagramEl = document.getElementById('archModalDiagram');
  const highlightsEl = document.getElementById('archModalHighlights');
  const stackEl = document.getElementById('archModalStack');

  if (titleEl) titleEl.textContent = data.title;
  if (badgeEl) {
    badgeEl.textContent = data.badge;
    badgeEl.className = 'badge ' + data.badgeClass;
  }
  if (summaryEl) summaryEl.textContent = data.summary;
  if (diagramEl) diagramEl.textContent = data.diagram;

  if (highlightsEl) {
    highlightsEl.innerHTML = data.highlights.map(h => `
      <div class="col-md-6">
        <div class="p-3 rounded-3 bg-body-tertiary border h-100">
          <div class="d-flex align-items-center mb-2">
            <i class="fas ${h.icon} text-primary me-2 fa-lg"></i>
            <h6 class="fw-bold mb-0">${h.title}</h6>
          </div>
          <p class="text-body-secondary small mb-0">${h.text}</p>
        </div>
      </div>
    `).join('');
  }

  if (stackEl) {
    stackEl.innerHTML = data.stack.map(s => `
      <span class="badge tech-badge">${s}</span>
    `).join('');
  }

  const modalEl = document.getElementById('projectArchModal');
  if (modalEl) {
    const modal = bootstrap.Modal.getOrCreateInstance(modalEl);
    modal.show();
  }
};
