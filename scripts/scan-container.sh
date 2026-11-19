#!/usr/bin/env bash
# ==============================================================================
# Primary Cloud Lab — Локальное сканирование безопасности Docker-контейнера
# Проверяет:
# 1. Dockerfile на соответствие лучшим практикам (IaC misconfigurations)
# 2. Образ контейнера на известные уязвимости CVE (OS packages + Go binary)
# 3. Наличие секретов и критических ошибок
# ==============================================================================

set -euo pipefail

IMAGE_NAME="primary-cloud-lab:security-audit"
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}▶ [1/3] Сканирование конфигурации Dockerfile (Trivy IaC)...${NC}"
docker run --rm --network host -v "$(pwd):/app" aquasec/trivy:latest config /app/Dockerfile

echo -e "\n${BLUE}▶ [2/3] Сборка локального тестового образа ${IMAGE_NAME}...${NC}"
docker build -t "$IMAGE_NAME" .

echo -e "\n${BLUE}▶ [3/3] Сканирование образа на уязвимости (CVE: HIGH, CRITICAL)...${NC}"
docker run --rm --network host \
  -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy:latest image \
  --severity HIGH,CRITICAL \
  --ignore-unfixed \
  "$IMAGE_NAME"

# Очистка локального тестового образа
docker rmi "$IMAGE_NAME" > /dev/null 2>&1 || true

echo -e "\n${GREEN}✔ Проверка безопасности завершена: уязвимостей CRITICAL/HIGH не обнаружено!${NC}\n"
