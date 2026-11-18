# AGENTS.md — AI Agent Operational Instructions & Repository Guardrails

> This document establishes the operational parameters, engineering standards, non-negotiable architectural invariants, and verification procedures for AI coding agents (Claude Code, Cursor, Copilot, Antigravity, Devin, Windsurf) working within the **Primary Cloud Lab** codebase.

---

## 1. Project Overview & Architecture

**Primary Cloud Lab** is an educational engineering platform and cloud showcase built for DevOps, SRE, and Systems Engineers. It teaches real-world Linux system administration, container internals (Docker, Podman, LXC/Incus, Cgroups v2, Namespaces), Infrastructure as Code (Ansible), and zero-downtime CI/CD pipelines.

### Architectural Blueprint
* **Core Web Server (`cmd/server`)**: High-performance HTTP application written in idiomatic Go standard library (`net/http`) without heavy external web frameworks.
* **Content Engine (`internal/service/course_service.go`)**: GFM-compliant markdown-to-HTML pipeline using `goldmark`, backed by thread-safe in-memory caching (`sync.RWMutex`) and automated full-text search indexing.
* **Microservices Ecosystem (`qr-service/`)**: Standalone, decoupled Go microservice responsible for real-time QR code generation and image serving.
* **Frontend Philosophy**: Server-Side Rendering (Go `html/template`), vanilla JavaScript (ES6+), modern CSS variables, and Bootstrap 5 for grid layout. **No heavy client-side frameworks** (React/Vue/Angular) and **no Node.js build steps**.
* **Edge & Ingress Gateway (`nginx/`)**: Nginx reverse proxy with TLS/SSL termination, HTTP/2, strict Content Security Policy (CSP), rate limiting, and gzip/brotli compression.

---

## 2. Non-Negotiable Invariants & Guardrails

When modifying any part of this repository, agents **MUST** strictly adhere to the following invariants:

### 1. Browser Tab Title Invariant
* The HTML `<title>` tag in `web/templates/layout.html` **MUST ALWAYS REMAIN EXACTLY**:
  ```html
  <title>Primary Cloud Lab</title>
  ```
* Under **no circumstances** should dynamic page titles, lesson names, or suffixes be injected into the browser tab `<title>`. Dynamic page context is communicated via on-page breadcrumbs and `<h1>` headings.

### 2. Code Formatting, Style & Comments
* Every Go file must strictly conform to official `gofmt` specifications before any commit:
  ```bash
  gofmt -l cmd/ internal/ qr-service/
  ```
* Avoid unnecessary external dependencies. Always favor the Go standard library (`net/http`, `sync`, `os`, `context`, `regexp`).
* **English Only for Code Comments**: All code comments, docstrings, and inline annotations across all programming languages in the project must be written exclusively in English.

### 3. Test Suite & Concurrency Verification
* All code changes must pass tests with the Go race detector enabled. Zero regressions allowed:
  ```bash
  # Root service
  go test -count=1 -race ./...

  # QR microservice
  cd qr-service && go test -count=1 -race ./...
  ```
* Thread safety: Any shared mutable state (e.g. caches, metrics) must use explicit synchronization (`sync.RWMutex` or `sync.Mutex`). Never introduce data races.

### 4. Content Security Policy (CSP) & Frontend Rules
* **No inline JavaScript**: Never add inline `<script>...</script>` blocks inside HTML templates.
* All scripts must reside in `web/static/js/pages/` and be loaded via external `<script src="..." defer></script>`.
* Event handlers must be bound unobtrusively with `addEventListener`, not inline HTML attributes (`onclick=...`).

### 5. Input Validation & Directory Traversal Protection
* Path parameters (e.g., `courseSlug`, `lessonSlug`) must be strictly sanitized against directory traversal attacks (`..`, `\`, `/`). Always verify clean paths before calling filesystem APIs.

---

## 3. Repository Topology

```
.
├── .agents/
│   ├── english-code-comments/
│   │   └── SKILL.md           # Standardized rule for English-only code comments
│   └── lesson-creator/
│       └── SKILL.md           # Modular playbook for authoring courses & interactive labs
├── AGENTS.md                  # Operational directives for AI coding agents
├── README.md                  # Project overview, architecture diagrams, and quickstart (EN/RU)
├── cmd/
│   └── server/
│       └── main.go            # Monolith entrypoint (flags, router assembly, graceful shutdown)
├── content/
│   └── courses/               # Educational Markdown tracks (linux, docker, podman, lxc, internals, ansible, cicd)
├── nginx/                     # Production Nginx reverse proxy configuration & TLS
├── scripts/                   # Deployment and SSL lifecycle management scripts
├── internal/
│   ├── config/                # Environment variable parsing and application config
│   ├── handler/               # HTTP handlers (courses, pages, search, SEO, QR integration)
│   ├── model/                 # Shared data structures (Course, Lesson, SearchItem, MetaData)
│   └── service/               # Core domain logic (course loading, search indexing, integrations)
├── qr-service/                # Independent microservice with its own go.mod and test suite
├── web/
│   ├── static/                # Static assets (CSS, modular JS, SVG icons, images)
│   └── templates/             # Go html/template layouts, pages, and components
└── go.mod                     # Go module definition
```

---

## 4. Development & Verification Commands

Agents must use the following standard commands when building, testing, and verifying the application:

```bash
# 1. Format check
gofmt -l cmd/ internal/ qr-service/

# 2. Run unit and integration tests with Race Detector
go test -count=1 -race ./...
cd qr-service && go test -count=1 -race ./... && cd ..

# 3. Build binaries locally
go build -v -o bin/server ./cmd/server
(cd qr-service && go build -v -o ../bin/qr-service .)

# 4. Run application locally
./bin/server -port=8080 -content=./content/courses -templates=./web/templates -static=./web/static
```

---

## 5. Git Commit Protocol

* All commit messages must follow the **Conventional Commits** specification:
  * `feat:` — New user-facing feature or educational track
  * `fix:` — Bug fix or broken link correction
  * `refactor:` — Code restructuring without behavior changes
  * `perf:` — Performance optimization (e.g. caching, allocations)
  * `chore:` — Maintenance, dependency updates, cleanup
  * `docs:` — Documentation and specification updates
  * `test:` — Test suite additions or modifications
* Verify working directory cleanliness (`git status`) before committing.
* Never commit secrets, temporary binaries (`bin/`), or operating system artifacts.
