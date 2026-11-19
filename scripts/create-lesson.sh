#!/usr/bin/env bash
# ==============================================================================
# Script: create-lesson.sh
# Purpose: Scaffold a new production-ready lesson for Docker Library with frontmatter
# Usage: ./scripts/create-lesson.sh <module-slug> <lesson-slug> "<Title>" [duration] [level] [tags]
# Example: ./scripts/create-lesson.sh docker-storage 03-permissions-uid-gid "Права доступа UID и GID в томах" "15 мин" "Middle" "Storage,Permissions"
# ==============================================================================

set -Eeuo pipefail

if [ "$#" -lt 3 ]; then
  echo "Использование: $0 <module-slug> <lesson-slug> \"<Заголовок>\" [длительность] [уровень] [теги]"
  echo "Пример: $0 docker-storage 03-permissions-uid-gid \"Права доступа UID и GID в томах\" \"15 мин\" \"Middle\" \"Storage,Permissions\""
  exit 1
fi

MODULE_SLUG="$1"
LESSON_SLUG="$2"
TITLE="$3"
DURATION="${4:-15-20 мин}"
LEVEL="${5:-Middle}"
TAGS="${6:-Docker}"

TARGET_DIR="content/courses/${MODULE_SLUG}"
TARGET_FILE="${TARGET_DIR}/${LESSON_SLUG}.md"

mkdir -p "${TARGET_DIR}"

if [ -f "${TARGET_FILE}" ]; then
  echo "❌ Ошибка: Файл '${TARGET_FILE}' уже существует!"
  exit 1
fi

# Extract order from slug prefix (e.g. "03-..." -> 3)
ORDER=$(echo "${LESSON_SLUG}" | grep -o '^[0-9]\+' || echo "1")
ORDER=$((10#$ORDER))

# Format tags as JSON array
FORMATTED_TAGS=$(echo "${TAGS}" | awk -F',' '{
  printf "["
  for (i=1; i<=NF; i++) {
    gsub(/^[ \t]+|[ \t]+$/, "", $i);
    printf "\"%s\"", $i;
    if (i<NF) printf ", ";
  }
  printf "]"
}')

cat <<EOF > "${TARGET_FILE}"
---
title: "${TITLE}"
description: "Практическое руководство и лабораторный разбор темы: ${TITLE}"
duration: "${DURATION}"
level: "${LEVEL}"
order: ${ORDER}
tags: ${FORMATTED_TAGS}
---
# ${TITLE}

Краткое введение в проблему: почему эта тема критически важна в production и какие задачи решает.

\`\`\`mermaid
flowchart TD
    A["Входные данные / Проблема"] --> B["Архитектурный механизм"]
    B --> C["Результат / Безопасное состояние"]
\`\`\`

---

## 1. Архитектура и принципы работы

Подробное описание физики процесса без воды: системные вызовы, взаимодействие с ядром Linux или демоном Docker.

---

## 2. Практическая конфигурация (Production Standard)

\`\`\`bash
# Пример команды для терминала
docker run -d --name service myimage:latest
\`\`\`

### Фрагмент compose.yaml:
\`\`\`yaml
services:
  app:
    image: myimage:latest
    restart: unless-stopped
\`\`\`

---

## 3. Частые ошибки и грабли в продакшене

- **Ошибка 1:** Описание проблемы, почему она возникает и как ее избежать.
- **Ошибка 2:** Описание проблемы и способ решения.

---

## 4. Лабораторный вызов: Проверь себя

<details>
<summary><b>Вопрос:</b> Контрольный вопрос для самопроверки инженера?</summary>

> **Ответ:** Развернутый и точный ответ с пояснением логики.
</details>
EOF

echo "✅ Урок успешно создан: ${TARGET_FILE}"
echo "📝 Заголовок: ${TITLE}"
echo "⏱ Длительность: ${DURATION} | Уровень: ${LEVEL} | Порядок: ${ORDER}"
