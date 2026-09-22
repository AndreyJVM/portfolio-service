/**
 * Интерактивный CLI-виджет (Терминал) на главной странице.
 * Эмулирует Linux Bash окружение, поддерживает быстрые кнопки,
 * ручной ввод команд, историю команд и вывод с подсветкой.
 */
document.addEventListener('DOMContentLoaded', () => {
  const terminalInput = document.getElementById('terminalInput');
  const terminalHistory = document.getElementById('terminalHistory');
  const terminalBody = document.getElementById('terminalBody');
  const terminalClearBtn = document.getElementById('terminalClearBtn');
  const terminalCopyBtn = document.getElementById('terminalCopyBtn');
  const quickBtns = document.querySelectorAll('.terminal-quick-btn');

  if (!terminalInput || !terminalHistory) {
    return;
  }

  const COMMAND_REGISTRY = {
    'whoami': () => `
<span class="terminal-text-cyan">{</span>
  <span class="terminal-text-yellow">"developer"</span>: <span class="terminal-text-green">"Андрей Воробьев"</span>,
  <span class="terminal-text-yellow">"role"</span>: <span class="terminal-text-green">"Java Backend / QA Automation / DevOps Engineer"</span>,
  <span class="terminal-text-yellow">"location"</span>: <span class="terminal-text-green">"Moscow / Remote"</span>,
  <span class="terminal-text-yellow">"github"</span>: <span class="terminal-text-green">"https://github.com/AndreyJVM"</span>,
  <span class="terminal-text-yellow">"telegram"</span>: <span class="terminal-text-green">"@AndreyAQA"</span>,
  <span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"Ready for production challenges & high-load systems"</span>
<span class="terminal-text-cyan">}</span>`,

    'docker ps': () => `
<span class="terminal-text-muted">CONTAINER ID   IMAGE                                      STATUS         PORTS</span>
<span class="terminal-text-cyan">7b8e1f92a3c4</span>   andreyvorobevaqa/portfolio-service:latest  Up (healthy)   0.0.0.0:8080-&gt;8080/tcp
<span class="terminal-text-cyan">3c4d5e6f7a8b</span>   postgres:17-alpine                         Up (healthy)   0.0.0.0:5432-&gt;5432/tcp
<span class="terminal-text-cyan">9f8a7b6c5d4e</span>   nginx:alpine (Let's Encrypt SSL)           Up (running)   0.0.0.0:80-&gt;80/tcp, 443-&gt;443/tcp`,

    'curl -s /actuator/health': () => `
<span class="terminal-text-cyan">{</span>
  <span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"UP"</span>,
  <span class="terminal-text-yellow">"components"</span>: {
    <span class="terminal-text-yellow">"db"</span>: {<span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"UP"</span>, <span class="terminal-text-yellow">"details"</span>: {<span class="terminal-text-yellow">"database"</span>: <span class="terminal-text-green">"PostgreSQL 17"</span>}},
    <span class="terminal-text-yellow">"caffeineCache"</span>: {<span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"UP"</span>, <span class="terminal-text-yellow">"details"</span>: {<span class="terminal-text-yellow">"ttl"</span>: <span class="terminal-text-green">"60m"</span>}},
    <span class="terminal-text-yellow">"diskSpace"</span>: {<span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"UP"</span>, <span class="terminal-text-yellow">"details"</span>: {<span class="terminal-text-yellow">"free"</span>: <span class="terminal-text-green">"38.4 GB"</span>}},
    <span class="terminal-text-yellow">"ping"</span>: {<span class="terminal-text-yellow">"status"</span>: <span class="terminal-text-green">"UP"</span>}
  }
<span class="terminal-text-cyan">}</span>`,

    'curl health': () => COMMAND_REGISTRY['curl -s /actuator/health'](),

    'mvn test': () => `
[INFO] Scanning for projects...
[INFO] Building Portfolio &amp; QR Code Generator Service 1.0.0
[INFO] -------------------------------------------------------
[INFO]  T E S T S
[INFO] -------------------------------------------------------
[INFO] Running ru.vorobevaqa.controller.api.QrRestControllerTest  ... <span class="terminal-text-green">PASSED</span>
[INFO] Running ru.vorobevaqa.controller.web.PageControllerTest    ... <span class="terminal-text-green">PASSED</span>
[INFO] Running ru.vorobevaqa.service.QrCodeServiceTest            ... <span class="terminal-text-green">PASSED</span>
[INFO] Results: Tests run: 7, Failures: 0, Errors: 0, Skipped: 0
[INFO] -------------------------------------------------------
[INFO] <span class="terminal-text-green fw-bold">BUILD SUCCESS</span>`,

    'skills': () => `
<span class="terminal-text-yellow fw-bold">Backend:</span>      Java 17/21, Spring Boot 3, Spring Data JPA, Hibernate, PostgreSQL, H2
<span class="terminal-text-yellow fw-bold">DevOps:</span>       Docker, Docker Compose, Linux (Ubuntu/Debian), GitHub Actions, Nginx SSL
<span class="terminal-text-yellow fw-bold">QA Automation:</span> REST Assured, JUnit 5, TestNG, Selenide, Allure Report, Postman`,

    'contact': () => `
<span class="terminal-text-green">Telegram:</span>  <a href="https://t.me/AndreyAQA" target="_blank" rel="noopener noreferrer" class="text-info">https://t.me/AndreyAQA</a>
<span class="terminal-text-green">GitHub:</span>    <a href="https://github.com/AndreyJVM" target="_blank" rel="noopener noreferrer" class="text-info">https://github.com/AndreyJVM</a>
<span class="terminal-text-green">Email:</span>     <a href="mailto:andrey.vorobev.qa@yandex.ru" class="text-info">andrey.vorobev.qa@yandex.ru</a>
<span class="terminal-text-green">Website:</span>   <a href="https://vorobevaqa.ru" class="text-info">https://vorobevaqa.ru</a>`,

    'help': () => `
<span class="terminal-text-yellow">Доступные команды:</span>
  <span class="terminal-text-green">whoami</span>             - информация об инженере и специализации
  <span class="terminal-text-green">docker ps</span>          - активные контейнеры продакшн-окружения
  <span class="terminal-text-green">curl health</span>        - онлайн проверка состояния БД, кэша и диска
  <span class="terminal-text-green">mvn test</span>           - прогон набора unit и integration тестов
  <span class="terminal-text-green">skills</span>             - технологический стек и компетенции
  <span class="terminal-text-green">contact</span>            - прямые каналы связи
  <span class="terminal-text-green">clear</span>              - очистить экран терминала`
  };

  function executeCommand(rawCommand) {
    const cmd = rawCommand.trim();
    if (!cmd) return;

    if (cmd.toLowerCase() === 'clear') {
      terminalHistory.innerHTML = '';
      terminalInput.value = '';
      return;
    }

    const entryDiv = document.createElement('div');
    entryDiv.className = 'terminal-entry mb-2';

    const promptLine = document.createElement('div');
    promptLine.className = 'terminal-prompt-line';
    promptLine.innerHTML = `<span class="terminal-prompt text-success fw-bold">guest@vorobevaqa:~$</span> <span class="terminal-cmd text-white">${escapeHtml(cmd)}</span>`;
    entryDiv.appendChild(promptLine);

    const outputDiv = document.createElement('div');
    outputDiv.className = 'terminal-output mt-1 text-light';

    const normalizedCmd = cmd.toLowerCase();
    const handler = COMMAND_REGISTRY[normalizedCmd];

    if (handler) {
      outputDiv.innerHTML = handler();
    } else {
      outputDiv.innerHTML = `<span class="terminal-text-red">bash: command not found: ${escapeHtml(cmd)}</span>. Введите <span class="terminal-text-yellow">'help'</span> для списка команд.`;
    }

    entryDiv.appendChild(outputDiv);
    terminalHistory.appendChild(entryDiv);

    terminalInput.value = '';
    scrollToBottom();
  }

  function escapeHtml(text) {
    const map = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, (m) => map[m]);
  }

  function scrollToBottom() {
    terminalBody.scrollTop = terminalBody.scrollHeight;
  }

  terminalInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      executeCommand(terminalInput.value);
    }
  });

  quickBtns.forEach((btn) => {
    btn.addEventListener('click', () => {
      const cmd = btn.getAttribute('data-cmd');
      if (cmd) {
        executeCommand(cmd);
        terminalInput.focus();
      }
    });
  });

  if (terminalClearBtn) {
    terminalClearBtn.addEventListener('click', () => {
      terminalHistory.innerHTML = '';
      terminalInput.value = '';
      terminalInput.focus();
    });
  }

  if (terminalCopyBtn) {
    terminalCopyBtn.addEventListener('click', () => {
      const text = terminalHistory.innerText;
      if (text) {
        navigator.clipboard.writeText(text).then(() => {
          terminalCopyBtn.innerHTML = '<i class="fas fa-check text-success"></i>';
          setTimeout(() => {
            terminalCopyBtn.innerHTML = '<i class="fas fa-copy"></i>';
          }, 1500);
        });
      }
    });
  }

  // Автоматический запуск whoami при первой загрузке страницы для вовлечения
  setTimeout(() => {
    executeCommand('whoami');
  }, 400);
});
