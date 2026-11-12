# --- Этап 1: Сборка статического бинарника Go ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем описание модулей и вендорированные зависимости
COPY go.mod go.sum ./
COPY vendor ./vendor

# Копируем исходный код
COPY cmd ./cmd
COPY internal ./internal

# Статическая компиляция бинарника со стриппингом отладочных символов (-w -s)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-w -s" -o /app/devops-platform ./cmd/server

# --- Этап 2: Минимальный production runtime ---
FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl tzdata && \
    addgroup -S -g 10001 appgroup && \
    adduser -S -u 10001 -G appgroup -s /sbin/nologin appuser

WORKDIR /app

COPY --from=builder --chown=10001:10001 /app/devops-platform .
COPY --chown=10001:10001 content ./content
COPY --chown=10001:10001 web ./web

USER 10001:10001

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["./devops-platform"]
