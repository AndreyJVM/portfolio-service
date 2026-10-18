# Stage 1: Build binary
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o qr-service .

# Stage 2: Minimal scratch image (~10MB)
FROM scratch

WORKDIR /app
COPY --from=builder /app/qr-service /app/qr-service

EXPOSE 8081

ENTRYPOINT ["/app/qr-service"]
