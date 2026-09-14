# Portfolio & QR Code Generator Service

[![Java](https://img.shields.io/badge/Java-17-orange?logo=openjdk)](https://openjdk.org/projects/jdk/17/)
[![Spring Boot](https://img.shields.io/badge/Spring%20Boot-3-brightgreen?logo=springboot)](https://spring.io/projects/spring-boot)
[![Docker](https://img.shields.io/badge/Docker-ready-blue?logo=docker)](https://hub.docker.com/r/andreyvorobevaqa/portfolio-service)
[![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=githubactions&logoColor=white)](https://github.com/AndreyJVM)
[![License](https://img.shields.io/badge/license-MIT-lightgrey)](LICENSE)

Production-ready web application combining a personal portfolio and an interactive QR code generator microservice. Built with Java 17 and Spring Boot 3, containerized with Docker, and deployed via a fully automated CI/CD pipeline.

**Live demo:** [vorobevaqa.ru](https://vorobevaqa.ru) · **Docker Hub:** [`andreyvorobevaqa/portfolio-service`](https://hub.docker.com/r/andreyvorobevaqa/portfolio-service)


## Table of Contents

- [Tech Stack](#tech-stack)
- [Key Features](#key-features)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [CI/CD Pipeline](#cicd-pipeline)
- [Project Structure](#project-structure)
- [Author](#author)
- [License](#license)


## Tech Stack

| Layer | Technologies |
|---|---|
| **Backend** | Java 17, Spring Boot 3, Spring MVC, Thymeleaf |
| **QR Engine** | Google ZXing (Core & JavaSE 3.5.3) |
| **Frontend** | Bootstrap 5.3, Font Awesome 6, Vanilla JS (ES6 modules, theme management) |
| **Build & Tools** | Maven, Lombok |
| **Infrastructure** | Docker, Docker Compose, Nginx (reverse proxy, SSL, Gzip), Ubuntu VDS |
| **CI/CD** | GitHub Actions (build, test, Docker Hub push, zero-downtime SSH deploy) |

## Key Features

- **Personal Portfolio** — responsive showcase with dedicated sections for completed engineering projects and professional skills.
- **Dark / Light Theme** — native toggle synchronized with Bootstrap 5.3 color modes (`data-bs-theme`), no flicker on load.
- **QR Code REST API** — high-performance QR code generation endpoint returning Base64-encoded images.
- **Client Features** — one-click PNG download, Web Share API integration, instant clipboard copy.
- **Production-Tuned** — Nginx reverse proxy with HTTP/2, Let's Encrypt SSL, and Gzip compression; JVM container limits configured via `JAVA_OPTS`.

## API Reference

### Generate QR Code

| | |
|---|---|
| **Endpoint** | `POST /api/qr` |
| **Content-Type** | `application/json` |

#### Request Body

```json
{
  "text": "https://vorobevaqa.ru",
  "size": 300
}
```

#### Response `200 OK`

```json
{
  "image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
}
```

#### Response `400 Bad Request`

```json
{
  "error": "Field 'text' must not be blank"
}
```

> **Note:** confirm the exact field names, validation rules, and error format against your controller/DTO before publishing — the samples above are illustrative placeholders.

## CI/CD Pipeline

1. **Build & Test** — runs `mvn clean verify` on Temurin JDK 17 on every push/PR.
2. **Containerization** — builds an optimized Docker image and pushes tagged versions (`latest`, `${{ github.sha }}`) to Docker Hub.
3. **Automated Deployment** — deploys the target image to the remote VDS over SSH with configuration sync and zero-downtime container recreation.

```mermaid
flowchart LR
    A[Push / PR] --> B[Build & Test<br/>mvn clean verify]
    B --> C[Build Docker Image]
    C --> D[Push to Docker Hub]
    D --> E[SSH Deploy to VDS]
    E --> F[Zero-downtime Restart]
```

## Project Structure

```
portfolio/
├── src/
│   ├── main/
│   │   ├── java/            # Application source code
│   │   └── resources/       # Templates, static assets, config
│   └── test/                # Unit & integration tests
├── Dockerfile
├── docker-compose.yml
├── pom.xml
└── .github/workflows/       # CI/CD pipeline definitions
```
