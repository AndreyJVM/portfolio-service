# Stage 1: Build binary
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pdf-service .

# Stage 2: Minimal scratch image (~15MB)
FROM scratch

WORKDIR /app
COPY --from=builder /app/pdf-service /app/pdf-service

EXPOSE 8082

ENTRYPOINT ["/app/pdf-service"]
