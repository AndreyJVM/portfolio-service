---
title: "Рецепты эталонных Production Dockerfile (Go, Node.js, Python, Rust, Nginx)"
description: "Коллекция готовых, проверенных multi-stage Dockerfile для коммерческих микросервисов с соблюдением стандартов Non-Root безопасности"
duration: "20 мин"
level: "Справочник • Все уровни"
order: 2
tags: ["Dockerfile", "Recipes", "Go", "Node.js", "Python", "Rust", "Nginx", "Security"]
---
# Рецепты эталонных Production Dockerfile (Go, Node.js, Python, Rust, Nginx)

Готовые шаблоны `Dockerfile` для продакшна. Каждый пример реализует:
- Многоэтапную сборку (**Multi-stage build**);
- Минимизацию размера (минус 80–95% лишних мегабайт);
- Запуск от непривилегированного пользователя (**Non-Root**);
- Корректную обработку сигналов PID 1 (`SIGTERM`).

---

## 1. Go (Golang) — Образ Scratch (~12 МБ)

```dockerfile
# Этап 1: Сборка
FROM golang:1.26-alpine AS builder

WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# -w -s вырезает отладочную информацию и таблицу символов (минус 40% размера)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Этап 2: Финальный образ (Scratch — нулевой размер базовой ОС)
FROM scratch

# Копируем сертификаты SSL и часовые пояса
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /server

# Запуск от непривилегированного пользователя (UID 10001)
USER 10001:10001

EXPOSE 8080
ENTRYPOINT ["/server"]
```

---

## 2. Node.js / TypeScript — Образ Alpine (~120 МБ)

```dockerfile
# Этап 1: Сборка TypeScript
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json tsconfig.json ./
RUN npm ci
COPY src/ ./src/
RUN npm run build

# Этап 2: Финальный продакшн-образ
FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production

# dumb-init необходим для корректной передачи SIGTERM в процесс Node.js
RUN apk add --no-cache dumb-init

COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

COPY --from=builder /app/dist ./dist

# Встроенный пользователь node (UID 1000)
USER node

EXPOSE 3000
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", "dist/index.js"]
```

---

## 3. Python / FastAPI — Debian Slim (~90 МБ)

```dockerfile
FROM python:3.12-slim AS builder

WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev \
 && rm -rf /var/lib/apt/lists/*

RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Финальный образ
FROM python:3.12-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
 && rm -rf /var/lib/apt/lists/* \
 && groupadd -g 10001 appgroup \
 && useradd -u 10001 -g appgroup -s /bin/sh appuser

COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY . .
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

---

## 4. Статический Nginx (Single Page App / Frontend)

```dockerfile
# Этап 1: Сборка фронтенда (Vite / React / Vue)
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Этап 2: Непривилегированный Nginx (Non-Root)
FROM nginxinc/nginx-unprivileged:alpine-slim

# Копируем оптимизированный конфиг Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем скомпилированную статику
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
```
