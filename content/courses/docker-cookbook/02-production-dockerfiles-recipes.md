---
title: "Рецепты эталонных Production Dockerfile (Go, Node.js, Python, Java, Rust, Nginx)"
description: "Коллекция готовых, проверенных multi-stage Dockerfile для коммерческих микросервисов с соблюдением стандартов Non-Root безопасности, кэширования слоев и правильной обработки сигналов"
duration: "25 мин"
level: "Справочник • Все уровни"
order: 2
tags: ["Dockerfile", "Recipes", "Go", "Node.js", "Python", "Java", "Rust", "Nginx", "Security", "Multi-stage"]
---
# Рецепты эталонных Production Dockerfile (Go, Node.js, Python, Java, Rust, Nginx)

Готовые и проверенные шаблоны `Dockerfile` для промышленной эксплуатации. Каждый пример строго реализует:
- Многоэтапную сборку (**Multi-stage build**);
- Минимизацию размера (минус 80–95% лишних мегабайт без компиляторов и утилит сборки);
- Запуск от непривилегированного пользователя (**Non-Root**);
- Корректную обработку сигналов PID 1 (`SIGTERM` и грациозную остановку);
- Оптимизацию кэширования зависимостей (изменение кода не инвалидирует кэш библиотек).

---

## 1. Go (Golang) — Образ Scratch (~12 МБ)

Идеальный паттерн для компилируемых Go-микросервисов. Базовый образ `scratch` не содержит операционной системы, шелла и утилит, что делает RCE-эксплойты бесполезными.

```dockerfile
# Этап 1: Сборка бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Необходимые системные сертификаты для HTTPS-запросов и часовые пояса
RUN apk add --no-cache ca-certificates tzdata

# Сначала копируем только файлы модулей для эффективного кэширования слоев
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем 100% статичный бинарник без CGO:
# -ldflags="-w -s" вырезает отладочную DWARF-информацию и таблицу символов (экономит ~40% размера)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/server ./cmd/server

# Этап 2: Финальный образ (Scratch — нулевой оверхед ОС)
FROM scratch

# Копируем доверенные корневые SSL-сертификаты и базу временных зон
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Копируем скомпилированный бинарный файл
COPY --from=builder /app/server /server

# Запуск от выделенного непривилегированного пользователя (UID 10001)
USER 10001:10001

EXPOSE 8080

# Используем exec-форму для прямой доставки сигналов ядра в приложение
ENTRYPOINT ["/server"]
```

---

## 2. Node.js / TypeScript — Образ Alpine (~120 МБ)

```dockerfile
# Этап 1: Компиляция TypeScript и сборка
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

# dumb-init необходим для корректной диспетчеризации сигналов SIGTERM процессу Node.js (PID 1)
RUN apk add --no-cache dumb-init

# Устанавливаем только production-зависимости (без typescript, eslint, jest)
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# Копируем скомпилированный JavaScript из первого этапа
COPY --from=builder /app/dist ./dist

# Встроенный в официальный образ Alpine непривилегированный пользователь node (UID 1000)
USER node

EXPOSE 3000

# dumb-init запускается как PID 1 и грациозно транслирует SIGTERM в node
ENTRYPOINT ["dumb-init", "--"]
CMD ["node", "dist/index.js"]
```

---

## 3. Python / FastAPI (Debian Slim & uv) (~95 МБ)

Современный шаблон с использованием менеджера `uv` (в 10–50 раз быстрее `pip`), виртуального окружения `venv` и изоляцией непривилегированного пользователя.

```dockerfile
# Этап 1: Сборка виртуального окружения
FROM python:3.12-slim AS builder

WORKDIR /app

# Устанавливаем сверхбыстрый установщик uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

# Системные библиотеки, необходимые для компиляции C-расширений (psycopg2, ujson)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc libpq-dev \
 && rm -rf /var/lib/apt/lists/*

# Создаем виртуальное окружение
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY requirements.txt .
# uv pip install использует кэш и компилирует колеса без раздувания образа
RUN uv pip install --no-cache -r requirements.txt

# Этап 2: Финальный минималистичный рантайм
FROM python:3.12-slim

WORKDIR /app
ENV PATH="/opt/venv/bin:$PATH" \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Только рантайм-библиотеки (без gcc и заголовков)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
 && rm -rf /var/lib/apt/lists/* \
 && groupadd -g 10001 appgroup \
 && useradd -u 10001 -g appgroup -s /bin/sh appuser

# Копируем готовое виртуальное окружение
COPY --from=builder /opt/venv /opt/venv

# Копируем исходный код приложения
COPY --chown=appuser:appgroup . /app

USER appuser

EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "2"]
```

---

## 4. Java / Spring Boot — Layered JAR Pattern (~180 МБ)

Классический Dockerfile для Spring Boot копирует `fat-jar` целиком (50–100 МБ). При любом изменении 1 строки кода Docker вынужден перекачивать весь JAR-файл. Паттерн **Layered JAR** делит архив на 4 слоя: сторонние зависимости кэшируются намертво, а пересобирается только слой кода приложения (<1 МБ)!

```dockerfile
# Этап 1: Распаковка слоев Spring Boot Layertools
FROM eclipse-temurin:21-jre-alpine AS extractor

WORKDIR /builder
COPY target/*.jar app.jar
# Извлекаем слои: dependencies, spring-boot-loader, snapshot-dependencies, application
RUN java -Djarmode=layertools -jar app.jar extract

# Этап 2: Финальный легковесный образ JRE
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем слои от наименее изменяемых к наиболее частым
COPY --from=extractor --chown=appuser:appgroup /builder/dependencies/ ./
COPY --from=extractor --chown=appuser:appgroup /builder/spring-boot-loader/ ./
COPY --from=extractor --chown=appuser:appgroup /builder/snapshot-dependencies/ ./
COPY --from=extractor --chown=appuser:appgroup /builder/application/ ./

USER appuser

# Оптимизация памяти JVM под контейнерные лимиты Cgroups v2
ENV JAVA_TOOL_OPTIONS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0 -XX:+ExitOnOutOfMemoryError"

EXPOSE 8080
ENTRYPOINT ["java", "org.springframework.boot.loader.launch.JarLauncher"]
```

---

## 5. Rust — Образ Scratch (~15 МБ)

```dockerfile
# Этап 1: Сборка с кросс-компиляцией через musl
FROM rust:1.78-alpine AS builder

WORKDIR /src
RUN apk add --no-cache musl-dev

# Копируем манифесты для кэширования крейтов
COPY Cargo.toml Cargo.lock ./
# Создаем заглушку, чтобы скомпилировать зависимости отдельно от бизнес-логики
RUN mkdir src && echo "fn main() {}" > src/main.rs && cargo build --release && rm -rf src

# Копируем реальный исходный код и собираем бинарник
COPY src ./src
RUN touch src/main.rs && cargo build --release --target x86_64-unknown-linux-musl
RUN strip target/x86_64-unknown-linux-musl/release/my-service

# Этап 2: Финальный образ Scratch
FROM scratch

COPY --from=builder /src/target/x86_64-unknown-linux-musl/release/my-service /my-service

USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/my-service"]
```

---

## 6. Статический Nginx (Single Page App / Vite / React / Vue) (~25 МБ)

```dockerfile
# Этап 1: Сборка фронтенда
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Этап 2: Официальный непривилегированный Nginx (Non-Root UID 101)
FROM nginxinc/nginx-unprivileged:alpine-slim

# Копируем оптимизированный конфиг Nginx (сжатие gzip, кэширование статики)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем скомпилированную статику из первого этапа
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
```

---

## 7. Сводная таблица параметров и размеров образов

| Стек / Язык | Базовый рантайм | Размер образа | Пользователь | Передача SIGTERM |
| :--- | :--- | :--- | :--- | :--- |
| **Go** | `scratch` | **~12 МБ** | `10001:10001` (Non-Root) | Нативная в бинарнике |
| **Node.js** | `node:20-alpine` | **~120 МБ** | `node (UID 1000)` | Через `dumb-init` |
| **Python** | `python:3.12-slim` | **~95 МБ** | `appuser (UID 10001)` | Нативная в Uvicorn |
| **Java** | `eclipse-temurin:21-jre-alpine`| **~180 МБ** | `appuser (Non-Root)` | Флаг `-XX:+ExitOnOutOfMemoryError`|
| **Rust** | `scratch` | **~15 МБ** | `10001:10001` (Non-Root) | Нативная в бинарнике |
| **Nginx SPA**| `nginx-unprivileged:alpine-slim`| **~25 МБ** | `nginx (UID 101)` | Нативная в Nginx |

---

## 8. Лабораторный вызов: Проверь себя

<details>
<summary><b>Вопрос 1:</b> Зачем в Java Spring Boot приложении распаковывать JAR через layertools вместо простого COPY app.jar?</summary>

> **Ответ:** Зависимости (библиотеки Spring, драйверы СУБД) весят 90% размера JAR и практически не меняются между коммитами. Распаковка через `layertools` изолирует их в отдельный неизменяемый Docker-слой, который навсегда кэшируется на сервере. При сборке нового релиза пересобирается только слой `application` размером в пару сотен килобайт, что ускоряет деплой в CI/CD в 10–20 раз.
</details>

<details>
<summary><b>Вопрос 2:</b> Почему для сборки Go и Rust бинарников под образ scratch необходимо отключать CGO и использовать musl?</summary>

> **Ответ:** По умолчанию компиляторы связывают код с динамической библиотекой C (`glibc`), которая ожидает наличие системных файлов `/lib/ld-linux-x86-64.so.2` на хосте. В образе `scratch` нет ни одного файла. Если скомпилировать бинарник с динамической линковкой, контейнер упадет с кодом 127 `not found`. Статическая линковка (`CGO_ENABLED=0` или `target=musl`) вшивает все функции внутрь бинарника, позволяя ему работать автономно.
</details>
