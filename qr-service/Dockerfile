# Stage 1: Build static binary
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates && \
    addgroup -S -g 10001 appgroup && \
    adduser -S -u 10001 -G appgroup -s /sbin/nologin appuser

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-s -w" -o /app/qr-service .

# Stage 2: Minimal scratch image (~10MB)
FROM scratch

LABEL maintainer="Primary Cloud Lab"
LABEL description="Go QR-Code Generation Microservice"

# Copy passwd, group, and TLS certificates from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder --chown=10001:10001 /app/qr-service /app/qr-service

# Run as non-root user
USER 10001:10001

EXPOSE 8081

ENTRYPOINT ["/app/qr-service"]
