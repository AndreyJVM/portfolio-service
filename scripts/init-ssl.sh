#!/usr/bin/env bash
# ==============================================================================
# Script: scripts/init-ssl.sh
# Purpose: Первичный выпуск и обновление SSL сертификатов Let's Encrypt через Certbot
# Usage: ./scripts/init-ssl.sh [email]
# ==============================================================================

set -euo pipefail

DOMAIN="vorobevaqa.ru"
EMAIL="${1:-andrey.vorobev.aqa@gmail.com}"

# Цвета для логов
CLR_RESET="\033[0m"
CLR_INFO="\033[1;34m"
CLR_SUCCESS="\033[1;32m"
CLR_WARN="\033[1;33m"
CLR_ERROR="\033[1;31m"

log_info()    { echo -e "${CLR_INFO}[INFO]${CLR_RESET} $*"; }
log_success() { echo -e "${CLR_SUCCESS}[SUCCESS]${CLR_RESET} $*"; }
log_warn()    { echo -e "${CLR_WARN}[WARN]${CLR_RESET} $*"; }
log_error()   { echo -e "${CLR_ERROR}[ERROR]${CLR_RESET} $*"; }

if command -v docker &>/dev/null && docker compose version &>/dev/null; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose &>/dev/null; then
  COMPOSE_CMD="docker-compose"
else
  log_error "Команда docker compose или docker-compose не найдена!"
  exit 1
fi

log_info "1. Подготовка каталогов certbot и ssl..."
mkdir -p certbot/conf certbot/www ssl

log_info "2. Проверка и пересоздание Nginx с поддержкой webroot certbot..."
$COMPOSE_CMD up -d --force-recreate nginx

log_info "3. Запрос SSL сертификата у Let's Encrypt для $DOMAIN..."
docker run --rm \
  -v "$(pwd)/certbot/conf:/etc/letsencrypt" \
  -v "$(pwd)/certbot/www:/var/www/certbot" \
  certbot/certbot:latest certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email "$EMAIL" \
  --agree-tos \
  --no-eff-email \
  --force-renewal \
  -d "$DOMAIN"

log_info "4. Копирование актуальных сертификатов в ./ssl/... (симлинки разыменовываются)"
cp -L "certbot/conf/live/$DOMAIN/fullchain.pem" ssl/public.pem
cp -L "certbot/conf/live/$DOMAIN/privkey.pem" ssl/private.pem
chmod 644 ssl/public.pem
chmod 600 ssl/private.pem

log_info "5. Перезагрузка конфигурации Nginx..."
$COMPOSE_CMD exec nginx nginx -s reload || $COMPOSE_CMD restart nginx

log_info "6. Запуск фонового контейнера автопродления Certbot..."
$COMPOSE_CMD up -d certbot

log_success "🎉 SSL сертификат успешно выпущен и установлен!"
log_success "Автопродление активно в фоновом контейнере certbot."
