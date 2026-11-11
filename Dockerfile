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
FROM alpine:latest

RUN apk add --no-cache ca-certificates curl tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/devops-platform .
COPY content ./content
COPY web ./web

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["./devops-platform"]
