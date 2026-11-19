# --- Этап 1: Сборка статического бинарника Go ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем описание модулей и вендорированные зависимости для кэширования слоев
COPY go.mod go.sum ./
COPY vendor ./vendor

# Копируем исходный код
COPY cmd ./cmd
COPY internal ./internal

# Кроссплатформенная сборка со стриппингом отладочных символов (-w -s) и кэшированием компилятора BuildKit
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -mod=vendor -ldflags="-w -s" -o /app/devops-platform ./cmd/server

# --- Этап 2: Минимальный production runtime ---
FROM alpine:3.21

# OCI Metadata Labels
LABEL org.opencontainers.image.title="Primary Cloud Lab" \
      org.opencontainers.image.description="Docker & Cloud Engineering Knowledge Base" \
      org.opencontainers.image.source="https://github.com/AndreyJVM/portfolio-service" \
      org.opencontainers.image.licenses="MIT"

# Создаем изолированного системного пользователя офлайн (без внешних сетевых вызовов)
RUN addgroup -S -g 10001 appgroup && \
    adduser -S -u 10001 -G appgroup -s /sbin/nologin appuser

WORKDIR /app

# Копируем скомпилированный бинарник и статические ресурсы с правами пользователя
COPY --from=builder --chown=10001:10001 /app/devops-platform ./devops-platform
COPY --chown=10001:10001 content ./content
COPY --chown=10001:10001 web ./web

USER 10001:10001

EXPOSE 8080

# Легковесный healthcheck через нативный busybox wget (не требует curl и внешних сетевых пакетов)
HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=3 \
  CMD wget -q --spider http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["./devops-platform"]
