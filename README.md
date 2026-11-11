<p align="center">
  <a href="https://ubuntu.com/"><img src="https://img.shields.io/badge/Linux-Ubuntu_24.04_LTS-E95420?logo=ubuntu&logoColor=white" alt="Linux"></a>
  <a href="https://hub.docker.com/r/andreyvorobevaqa/portfolio-service"><img src="https://img.shields.io/badge/Docker-Compose_%26_Swarm-2496ED?logo=docker&logoColor=white" alt="Docker"></a>
  <a href="https://galaxy.ansible.com/ui/standalone/namespaces/24255/"><img src="https://img.shields.io/badge/Ansible-Galaxy_Roles-EE0000?logo=ansible&logoColor=white" alt="Ansible"></a>
  <a href="https://github.com/AndreyJVM"><img src="https://img.shields.io/badge/CI%2FCD-GitHub_Actions_%26_GitLab-2088FF?logo=githubactions&logoColor=white" alt="CI/CD"></a>
  <a href="https://nginx.org/"><img src="https://img.shields.io/badge/Nginx-API_Gateway_%26_SSL-009639?logo=nginx&logoColor=white" alt="Nginx"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26_Platform-00ADD8?logo=go&logoColor=white" alt="Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
</p>

# Primary Cloud Lab • vorobevaqa.ru

Открытая интерактивная образовательная платформа по **DevOps, системному администрированию Linux, контейнеризации в Docker и автоматизации релизов**.

Переписана с монолитной Java/Spring Boot архитектуры на высокопроизводительный, легковесный микростек на **Go (Golang)**:
- 🚀 **Потребление RAM:** сократилось с ~450 МБ до **~13 МБ** (экономия >95% памяти на VDS);
- ⚡ **Время холодного старта:** снизилось с ~15 секунд до **<5 миллисекунд**;
- 📦 **Размер Docker-образа:** уменьшился с ~550 МБ до **~33 МБ**;
- 📚 **Курсы как код (Courses as Code):** уроки хранятся в Git в Markdown-формате и компилируются на лету.

**Продакшн-платформа:** [vorobevaqa.ru](https://vorobevaqa.ru) · **Docker Hub:** [`andreyvorobevaqa`](https://hub.docker.com/repositories/andreyvorobevaqa) · **Ansible Galaxy:** [роли автоматизации](https://galaxy.ansible.com/ui/standalone/namespaces/24255/)


## Содержание

- [Инфраструктура и Архитектура](#инфраструктура-и-архитектура)
- [Образовательные треки (Курсы)](#образовательные-треки-курсы)
- [DevOps & Hardening контейнеров](#devops--hardening-контейнеров)
- [CI/CD & Бесперебойный Blue-Green деплой](#cicd--бесперебойный-blue-green-деплой)
- [Тестирование и локальный запуск](#тестирование-и-локальный-запуск)
- [Структура проекта](#структура-проекта)
- [Автор и Контакты](#автор-и-контакты)
- [Лицензия](#лицензия)


## Инфраструктура и Архитектура

```mermaid
flowchart TD
    Client["Внешний клиент / Студент"] -->|HTTPS :443 / HTTP :80| Host["Linux VDS (Ubuntu 24.04 LTS)<br/>UFW Firewall • SSH ed25519 • Fail2ban • logrotate"]
    Host --> Nginx["Nginx Reverse Proxy & API Gateway<br/>Терминация SSL • HTTP/2 • Gzip • Fallback 50x"]
    Nginx -->|/courses, /, /projects, /mustang, /api/*| GoApp["Go Primary Cloud Lab (:8080)<br/>Chi Router • Goldmark Markdown • RAM ~13MB"]
    Nginx -->|/api/qr fallback| GoQR["Go QR Microservice (:8081)<br/>Scratch-образ ~10MB • cgroups квоты"]
```

| Уровень | Технологии & Решения |
|---|---|
| **Хост и ОС** | Linux (Ubuntu 24.04 LTS), Bash-скрипты автоматизации, systemd, cron, ротация logrotate |
| **Сетевая безопасность** | UFW firewall (ограничение SSH, 80/443), Fail2ban от брутфорса, SSH-аутентификация по ключам |
| **Шлюз и Прокси** | Nginx (Reverse Proxy, API Gateway, SSL Let's Encrypt / Certbot, HTTP/2, Gzip, кастомные 50x fallback) |
| **Контейнеризация** | Docker, Docker Compose, Docker Hub Registry, Alpine Linux runtime |
| **Бэкенд платформы** | Go 1.26 (`go-chi/chi/v5`, `yuin/goldmark`, `skip2/go-qrcode`), статическая сборка со стриппингом отладочных символов |
| **IaC & Оркестрация** | Ansible (идемпотентные плейбуки и роли в Ansible Galaxy), Blue-Green деплой скриптом `deploy.sh` |
| **CI/CD** | GitHub Actions (параллельная сборка, BuildKit gha-кэш, Zero-Downtime деплой с авто-откатом) |


## Образовательные треки (Курсы)

1. 🐧 **Системное администрирование Linux** (`/courses/linux`):
   - Права доступа (POSIX, SUID, SGID, Sticky bit), управление юнитами и таймерами systemd;
   - Сетевой стек Linux (iproute2, ss, iptables, ufw), анализ трафика и сокетов.
2. 🐳 **Контейнеризация и Docker** (`/courses/docker`):
   - Многоэтапные сборки (Multi-stage), Linux cgroups v2, неймспейсы (PID/NET/MNT);
   - Docker Compose v2, внутренняя DNS-маршрутизация и изоляция bridge-сетей.


## DevOps & Hardening контейнеров

Стек развертывается через `docker-compose.yml` в соответствии с production-стандартами:

1. **Минимальные права и non-root пользователь:**
   Контейнер работает от непривилегированного пользователя `appuser:appgroup` с `security_opt: [no-new-privileges:true]`.
2. **Встроенный Healthcheck:**
   ```yaml
   healthcheck:
     test: [ "CMD-SHELL", "curl -f http://localhost:8080/health || exit 1" ]
     interval: 5s
     timeout: 3s
     retries: 4
     start_period: 3s
   ```
3. **Лимиты ресурсов (cgroups):**
   ```yaml
   deploy:
     resources:
       limits:
         cpus: "0.5"
         memory: 128M
       reservations:
         cpus: "0.05"
         memory: 32M
   ```
4. **Ротация логов (Log Rotation):**
   ```yaml
   logging:
     driver: "json-file"
     options:
       max-size: "10m"
       max-file: "3"
   ```


## CI/CD & Бесперебойный Blue-Green деплой

```mermaid
flowchart LR
    A[Git Push master] --> B[detect-changes<br/>Paths Filter]
    B --> C[Build & Test Go<br/>100% Passing]
    C --> D[Docker Build & Push<br/>BuildKit Cache]
    D --> E[Zero-Downtime Deploy VDS<br/>scripts/deploy.sh]
    E --> F[Healthcheck Verify<br/>Nginx Reload & Prune]
```

1. **Path-based trigger:** Сборка запускается только для изменившихся компонентов (Go платформа, QR-сервис, конфигурация деплоя).
2. **Concurrency lock:** Группа `production-deploy` блокирует одновременные гонки релизов.
3. **Бесперебойное обновление:** `scripts/deploy.sh` запускает новый инстанс параллельно (`scale app=2`), ждет `healthy`, перезагружает Nginx и гасит старый инстанс. В случае сбоя срабатывает мгновенный автоматический откат (`rollback`).


## Тестирование и локальный запуск

```bash
# Запуск тестов платформы на Go
go test -v ./...

# Локальный запуск приложения Go
go run cmd/server/main.go

# Сборка и запуск в Docker
docker build -t andreyvorobevaqa/portfolio-service:test .
docker run -p 8080:8080 andreyvorobevaqa/portfolio-service:test
```


## Структура проекта

```
portfolio/
├── cmd/server/main.go       # Точка входа веб-сервера Go (Chi роутер, Graceful shutdown)
├── content/courses/         # Образовательные треки и уроки в Markdown (Linux, Docker)
├── internal/
│   ├── config/              # Конфигурация окружения
│   ├── handler/             # HTTP-обработчики страниц и REST API
│   ├── model/               # Модели данных (Courses, Lessons, Integrations)
│   └── service/             # Бизнес-логика каталога и интеграций (GitHub, DockerHub, Ansible)
├── web/
│   ├── static/              # CSS, JS, изображения сертификатов, PDF резюме
│   └── templates/           # HTML-шаблоны Go (layout, index, courses, projects, mustang)
├── nginx/
│   ├── nginx.conf           # Конфигурация Reverse Proxy, SSL Let's Encrypt и роутинга
│   └── 50x.html             # Кастомная страница аварийных ошибок
├── scripts/
│   └── deploy.sh            # Blue-Green развертывание с авто-откатом (Rollback)
├── docker-compose.yml       # Конфигурация продакшн-стека микросервисов
├── Dockerfile               # Multi-stage сборка статического бинарника Go на Alpine
├── go.mod / go.sum          # Модули Go 1.26
└── .github/workflows/       # Пайплайн GitHub Actions
```


## Автор и Контакты

**Андрей Воробьев**
- **Специализация:** DevOps & Linux Infrastructure Engineer
- **Веб-сайт:** [vorobevaqa.ru](https://vorobevaqa.ru)
- **GitHub:** [@AndreyJVM](https://github.com/AndreyJVM)
- **Docker Hub:** [andreyvorobevaqa](https://hub.docker.com/repositories/andreyvorobevaqa)
- **Ansible Galaxy:** [Профиль автора](https://galaxy.ansible.com/ui/standalone/namespaces/24255/)
- **Telegram:** [@AndreyAQA](https://t.me/AndreyAQA)

---

## Лицензия

Проект распространяется под свободной лицензией [Apache License 2.0](LICENSE).
