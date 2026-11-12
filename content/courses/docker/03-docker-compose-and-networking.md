# Сети в Docker Compose и изоляция микросервисов

Docker по умолчанию предоставляет встроенный DNS-сервер для обнаружения сервисов (Service Discovery).

---

### 1. Архитектура пользовательской сети Bridge

Когда сервисы объединены в пользовательскую сеть `bridge`:
1. Контейнеры могут обращаться друг к другу **по имени сервиса**, а не по динамическим IP-адресам.
2. Внутренние порты (например, `8080` у бэкенда или `5432` у базы данных) не нужно пробрасывать наружу в секции `ports:`. Достаточно открыть наружу только `80:80` и `443:443` на Nginx Gateway.

```yaml
services:
  app:
    image: my-app:latest
    networks:
      - internal-net

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    networks:
      - internal-net

networks:
  internal-net:
    driver: bridge
```

---

### 2. Практическая проверка резолвинга DNS внутри контейнера

```bash
# Выполнить DNS-запрос к соседнему сервису
docker compose exec nginx nslookup app
docker compose exec nginx ping -c 2 app
```
