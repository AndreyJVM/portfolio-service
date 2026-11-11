# Multi-stage сборка Dockerfile и лимиты ресурсов (cgroups)

Оптимизация размеров Docker-образов и защита хоста от утечек памяти с помощью Control Groups.

---

### 1. Почему важен Multi-stage сборщик?

В традиционном Dockerfile SDK, компиляторы и утилиты сборщика остаются в финальном образе, раздувая его до 1-2 ГБ и создавая векторы атак.

**Паттерн Multi-stage в Go:**
```dockerfile
# Этап 1: Сборка статического бинарника
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Этап 2: Финальный легковесный runtime (Scratch / Alpine)
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
USER 10001
EXPOSE 8080
ENTRYPOINT ["./server"]
```
*Результат:* Размер финального контейнера снижается **с 800 МБ до 18 МБ!**

---

### 2. Защита от OOM (Out Of Memory) через cgroups в Docker Compose

Если один контейнер начнет течь по памяти, Linux OOM Killer может убить критичные процессы на хосте (например, sshd или nginx).

```yaml
services:
  app:
    image: my-app:latest
    deploy:
      resources:
        limits:
          cpus: "0.50"
          memory: 128M
        reservations:
          cpus: "0.05"
          memory: 32M
```
