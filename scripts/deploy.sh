#!/usr/bin/env bash
# ==============================================================================
# Script: scripts/deploy.sh
# Purpose: Zero-downtime Blue-Green deployment with automatic rollback
# Target: Ubuntu VDS running Docker Compose
# ==============================================================================

set -euo pipefail

# ANSI Color Codes
CLR_RESET="\033[0m"
CLR_INFO="\033[1;34m"
CLR_SUCCESS="\033[1;32m"
CLR_WARN="\033[1;33m"
CLR_ERROR="\033[1;31m"

log_info()    { echo -e "${CLR_INFO}[INFO]${CLR_RESET} $*"; }
log_success() { echo -e "${CLR_SUCCESS}[SUCCESS]${CLR_RESET} $*"; }
log_warn()    { echo -e "${CLR_WARN}[WARN]${CLR_RESET} $*"; }
log_error()   { echo -e "${CLR_ERROR}[ERROR]${CLR_RESET} $*"; }

# ------------------------------------------------------------------------------
# 0. Проверка рабочего каталога и Docker Compose
# ------------------------------------------------------------------------------
if command -v docker &>/dev/null && docker compose version &>/dev/null; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose &>/dev/null; then
  COMPOSE_CMD="docker-compose"
else
  log_error "Команда docker compose или docker-compose не найдена!"
  exit 1
fi

log_info "Используется утилита: $COMPOSE_CMD"

# ------------------------------------------------------------------------------
# 1. Загрузка актуальных образов из Docker Hub
# ------------------------------------------------------------------------------
log_info "📥 Шаг 1/5: Загрузка свежих Docker-образов..."
$COMPOSE_CMD pull app qr-service

# ------------------------------------------------------------------------------
# 2. Обновление Go микросервиса (мгновенный старт scratch ~10ms)
# ------------------------------------------------------------------------------
log_info "⚡ Шаг 2/5: Обновление Go микросервиса (qr-service)..."
$COMPOSE_CMD up -d --no-deps qr-service

# Ожидание готовности healthcheck для Go сервисов
for svc in qr-service; do
  log_info "Проверка готовности $svc..."
  SVC_READY=false
  for i in $(seq 1 15); do
    STATUS=$($COMPOSE_CMD ps "$svc" --format '{{.Health}}' 2>/dev/null || true)
    if [[ "$STATUS" == *"healthy"* ]]; then
      log_success "$svc готов и прошел healthcheck!"
      SVC_READY=true
      break
    fi
    sleep 1
  done
  if [ "$SVC_READY" = false ]; then
    log_warn "$svc не сообщил о healthy в течение 15s (продолжаем деплой)..."
  fi
done

# ------------------------------------------------------------------------------
# 3. Бесперебойный (Zero-Downtime) деплой Java приложения
# ------------------------------------------------------------------------------
log_info "🚀 Шаг 3/5: Развертывание новой версии Java приложения (app)..."

OLD_APP_CONTAINER=$($COMPOSE_CMD ps -q app 2>/dev/null | head -n 1 || true)

if [ -z "$OLD_APP_CONTAINER" ]; then
  log_info "Предыдущий контейнер app не найден. Выполняется первичный запуск..."
  $COMPOSE_CMD up -d --no-deps app
else
  log_info "Обнаружен работающий контейнер app: $OLD_APP_CONTAINER"
  log_info "Запуск нового экземпляра рядом со старым (scale app=2)..."
  $COMPOSE_CMD up -d --no-deps --scale app=2 --no-recreate app

  # Определение ID нового запущенного контейнера
  NEW_APP_CONTAINER=""
  for cid in $($COMPOSE_CMD ps -q app); do
    if [ "$cid" != "$OLD_APP_CONTAINER" ]; then
      NEW_APP_CONTAINER="$cid"
      break
    fi
  done

  log_info "Новый инстанс: ${NEW_APP_CONTAINER:-unknown}. Ожидание прохождения healthcheck (до 60 секунд)..."

  HEALTHCHECK_PASSED=false
  for i in $(seq 1 30); do
    HEALTHY_COUNT=$($COMPOSE_CMD ps app --format '{{.Health}}' 2>/dev/null | grep -c "healthy" || true)
    if [ "$HEALTHY_COUNT" -ge 2 ]; then
      log_success "Новый контейнер app успешно прогрелся и сообщает status: healthy!"
      HEALTHCHECK_PASSED=true
      break
    fi
    echo "  [Ожидание $i/30] Здоровых контейнеров app: $HEALTHY_COUNT/2..."
    sleep 2
  done

  # Автоматический откат (Rollback) в случае непрохождения проверки
  if [ "$HEALTHCHECK_PASSED" = false ]; then
    log_error "❌ Новый контейнер app НЕ прошел проверку здоровья за 60 секунд!"
    log_warn "🔄 Инициирован автоматический откат (ROLLBACK)..."

    if [ -n "$NEW_APP_CONTAINER" ]; then
      log_info "Остановка сбойного инстанса $NEW_APP_CONTAINER..."
      docker stop "$NEW_APP_CONTAINER" || true
      docker rm -f "$NEW_APP_CONTAINER" || true
    fi

    log_info "Восстановление масштабирования до 1 (старый рабочий контейнер)..."
    $COMPOSE_CMD up -d --no-deps --scale app=1 app || true

    log_error "Деплой отменен из-за непрохождения healthcheck. Старая версия продолжает работу без простоя."
    exit 1
  fi

  # ----------------------------------------------------------------------------
  # 4. Переключение Nginx и остановка старого контейнера
  # ----------------------------------------------------------------------------
  log_info "🔄 Шаг 4/5: Плавное переключение трафика Nginx Gateway..."
  $COMPOSE_CMD up -d --no-deps nginx
  $COMPOSE_CMD exec -T nginx nginx -s reload || true

  log_info "🧹 Остановка и удаление устаревшего контейнера app ($OLD_APP_CONTAINER)..."
  docker stop "$OLD_APP_CONTAINER" || true
  docker rm -f "$OLD_APP_CONTAINER" || true

  # Фиксация scale=1
  $COMPOSE_CMD up -d --no-deps --scale app=1 app
fi

# ------------------------------------------------------------------------------
# 5. Очистка неиспользуемых слоев и завершение
# ------------------------------------------------------------------------------
log_info "🗑️ Шаг 5/5: Очистка устаревших Docker-образов..."
docker image prune -f

log_success "🎉 Zero-Downtime деплой успешно завершен!"
