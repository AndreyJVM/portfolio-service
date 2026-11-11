# Проектирование CI/CD в GitHub Actions с кэшированием BuildKit

Автоматизация тестирования и сборки контейнеров — основа практики Continuous Integration.

---

### 1. Архитектура Workflow в GitHub Actions

Пайплайн описывается в `.github/workflows/ci.yml`. Основные сущности:
- **Events (Triggers):** события, запускающие выполнение (`push`, `pull_request`, `workflow_dispatch`).
- **Jobs:** изолированные наборы шагов, по умолчанию выполняющиеся параллельно на виртуальных машинах runners.
- **Steps:** последовательные команды или переиспользуемые Action-блоки.

---

### 2. Оптимизация сборки: Кэширование слоев Docker (gha)

Обычная команда `docker build` на GitHub Actions каждый раз собирает образ с нуля из-за чистой виртуальной машины runner. С помощью драйвера `docker-container` и кэша GitHub Actions (`type=gha`) сборка ускоряется с 5 минут до 20 секунд:

```yaml
name: CI Pipeline

on:
  push:
    branches: [ master ]

jobs:
  test-and-build:
    name: Test & Build Go Application
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'
          cache-dependency-path: go.sum

      - name: Run Unit Tests
        run: go test -v -race ./...

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and Push with BuildKit Cache
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ${{ secrets.DOCKER_USERNAME }}/devops-app:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

### 3. Concurrency Lock: Защита от конфликтующих деплоев

Если два разработчика пушат коммиты с разницей в секунды, деплои могут перемешаться. Для исключения гонок используется группа `concurrency`:

```yaml
concurrency:
  group: production-deploy
  cancel-in-progress: false
```
Параметр `cancel-in-progress: false` гарантирует, что текущий релиз штатно завершится, а следующий встанет в очередь.
