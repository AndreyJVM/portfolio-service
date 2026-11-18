---
name: lesson-content-and-lab-creator
description: Standardized playbook for authoring, validating, and publishing interactive DevOps learning tracks, terminal-first lab sessions, code snippets, and SEO metadata in Primary Cloud Lab.
version: 1.0.0
tags: [devops, education, golang, content-pipeline, agentic-skill]
---

# Skill Playbook: Lesson Content & Lab Creator

## 1. Overview & Capability

This skill provides an automated, repeatable playbook for AI agents and platform engineers to design, build, validate, and integrate new interactive DevOps courses and lessons into the **Primary Cloud Lab** ecosystem.

### Key Capabilities
* Authors production-grade educational content in GitHub Flavored Markdown (GFM).
* Constructs terminal-first hands-on labs with realistic Linux/DevOps shell commands.
* Registers new tracks and lessons with thread-safe caching and search indexing.
* Validates zero-downtime integration and regression-free test passes with Go's race detector.

---

## 2. Input Requirements & Metadata Schema

Before executing this skill, collect or generate the following inputs:

| Field | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `courseSlug` | `string` | Lowercase identifier matching track directory | `linux`, `podman`, `cicd` |
| `orderPrefix` | `string` | 2-digit zero-padded index for sequential ordering | `01`, `02`, `06` |
| `lessonSlug` | `string` | Kebab-case descriptive slug | `quadlet-systemd-integration` |
| `title` | `string` | Human-readable title in Russian/English | `Интеграция Podman Quadlet с systemd` |
| `duration` | `string` | Estimated hands-on completion time | `15-20 мин` |
| `level` | `string` | Difficulty tier | `Базовый`, `Middle`, `Senior` |
| `tags` | `[]string` | Technology keywords for search indexing | `["Podman", "Quadlet", "systemd"]` |

---

## 3. Step-by-Step Execution Protocol

### Step 1: Author the Lesson Markdown File
Create the markdown file at `content/courses/<courseSlug>/<orderPrefix>-<lessonSlug>.md`.

#### Lesson Structure Standard:
Every lesson file must strictly follow this 4-section pedagogical layout:
```markdown
# [Lesson Title matching H1]

[Brief introductory hook explaining WHY this topic matters in high-load production environments.]

---

### 1. Концептуальные основы и теория (Core Architecture & Theory)
* Explain low-level mechanics (e.g. kernel namespaces, systemd socket activation, multi-stage caching).
* Use ASCII or Mermaid diagrams where appropriate.

---

### 2. Конфигурация и анатомия файлов (Configuration & File Anatomy)
* Provide production-ready configuration examples (e.g., unit files, Containerfiles, Ansible tasks).
* Annotate every parameter and directive.

---

### 3. Лабораторная практика в терминале (Hands-on Terminal Lab)
* Provide sequenced, copy-paste ready commands with explanatory comments:
```bash
# Step 1: Create isolated user directory
mkdir -p ~/.config/containers/systemd/

# Step 2: Write quadlet container definition
cat << 'EOF' > ~/.config/containers/systemd/web.container
[Container]
Image=docker.io/library/nginx:alpine
PublishPort=8080:80
[Service]
Restart=always
[Install]
WantedBy=default.target
EOF

# Step 3: Reload user systemd daemon and start service
systemctl --user daemon-reload
systemctl --user start web
```

---

### 4. Проверка и отладка (Verification & Diagnostics)
* Provide diagnostic commands to inspect status, logs, and sockets:
```bash
systemctl --user status web
podman ps -a
journalctl --user -u web -n 50 --no-pager
```
```

---

### Step 2: Register New Tracks (If Creating a New Course)
If introducing a brand new course track rather than adding a lesson to an existing course:

1. Update `GetAllCourses()` in [`internal/service/course_service.go`](file:///home/it/IdeaProjects/portfolio-service/internal/service/course_service.go):
   ```go
   {
       ID:          "newtrack",
       Slug:        "newtrack",
       Title:       "Track Title",
       Icon:        "fas fa-icon-name",
       Color:       "#colorHex",
       Level:       "Middle → Senior",
       Description: "Track description...",
       Tags:        []string{"Tag1", "Tag2"},
   },
   ```
2. Update the sidebar navigation in [`web/templates/layout.html`](file:///home/it/IdeaProjects/portfolio-service/web/templates/layout.html):
   ```html
   <a class="nav-link {{if eq .Meta.ActiveCourse "newtrack"}}active{{end}}" href="/courses/newtrack">Track Title</a>
   ```

---

### Step 3: Automated Quality Gate & Verification

Execute the verification commands to ensure zero regressions:

```bash
# 1. Verify Go syntax and format
gofmt -l cmd/ internal/

# 2. Run race-detector unit and integration tests
go test -count=1 -race ./...

# 3. Ensure the browser title invariant is preserved
grep -rn "<title>" web/templates/
# Expected output must strictly be:
# web/templates/layout.html:10:  <title>Primary Cloud Lab</title>
```

---

### Step 4: Commit and Publish

Commit following the Conventional Commits format:
```bash
git add content/courses/<courseSlug>/<orderPrefix>-<lessonSlug>.md
git commit -m "feat(courses): add lesson <lessonSlug> to <courseSlug> track"
git push origin master
```

---

## 4. Best Practices for Terminal Labs

1. **Deterministic & Safe**: Commands should avoid destructive actions without explicit user caution (`rm -rf /` or modifying system partitions).
2. **Flag Explanation**: Never provide cryptic shell flags (e.g. `tar -xzvf`) without documenting what each flag stands for.
3. **No Placeholders**: Avoid `<YOUR_IP_HERE>` when possible; prefer `127.0.0.1` or exportable variables (`export MY_IP=$(hostname -I | awk '{print $1}')`).
4. **Self-Contained Verification**: Every lab step must have a corresponding verification step to confirm success before proceeding.
