---
title: "Шпаргалка команд Docker и Docker Compose (Cheat Sheet)"
description: "Исчерпывающий справочник консольных команд для повседневной работы: управление контейнерами, сетями, томами, логами и форматирование inspect"
duration: "15 мин"
level: "Справочник • Все уровни"
order: 1
tags: ["Cheat Sheet", "CLI", "Commands", "Docker Compose", "Inspect"]
---
# Шпаргалка команд Docker и Docker Compose (Cheat Sheet)

Консольная шпаргалка с самыми важными и эффективными командами для повседневной работы инженера. Без воды — только флаги, синтаксис и практические примеры.

---

## 1. Управление контейнерами (Lifecycle)

```bash
# Запустить контейнер в фоне с автоперезапуском и пробросом порта на localhost
docker run -d --name web --restart unless-stopped -p 127.0.0.1:8080:80 nginx:alpine

# Интерактивный вход в запущенный контейнер с командной оболочкой bash или sh
docker exec -it web sh

# Посмотреть активные контейнеры (с размером слоев)
docker ps -s

# Посмотреть ВСЕ контейнеры (включая упавшие и остановленные)
docker ps -a

# Быстро остановить и сразу удалить контейнер
docker rm -f web

# Мониторинг потребления ресурсов (CPU, RAM, Net I/O) в реальном времени
docker stats --no-stream
```

---

## 2. Диагностика и Логи

```bash
# Просмотр последних 100 строк логов и подключение к живому потоку (follow)
docker logs -f --tail 100 web

# Логи с таймстемпами (временем возникновения)
docker logs -t --tail 50 web

# Узнать точную причину падения (ExitCode и OOMKilled)
docker inspect web --format 'ExitCode: {{.State.ExitCode}}, OOMKilled: {{.State.OOMKilled}}'

# Узнать IP-адрес контейнера во внутренней сети Docker
docker inspect web --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

# Посмотреть изменения файлов внутри контейнера относительно базового образа
docker diff web
```

---

## 3. Образы и Сборка

```bash
# Сборка с явным именем тега и без использования кэша
docker build --no-cache -t myapp:v1.0 .

# Сборка через BuildKit с передачей секрета (секрет не попадает в слои образа)
DOCKER_BUILDKIT=1 docker build --secret id=mytoken,src=.token -t myapp:secure .

# Сканирование образа на уязвимости через встроенный сканер
docker scout quickview myapp:v1.0

# История слоев образа (что занимает больше всего места)
docker history --human myapp:v1.0
```

---

## 4. Сети и Тома

```bash
# Создать изолированную пользовательскую сеть-мост
docker network create --driver bridge app-network

# Подключить уже работающий контейнер к сети
docker network connect app-network web

# Список всех томов и инспекция конкретного тома
docker volume ls
docker volume inspect my_postgres_data

# Создать резервную копию тома в tar-архив на хосте
docker run --rm -v my_postgres_data:/data -v $(pwd):/backup alpine tar czvf /backup/db_backup.tar.gz -C /data .
```

---

## 5. Docker Compose v2

```bash
# Сборка и запуск сервисов в фоне с пересборкой образов
docker compose up -d --build

# Просмотр статуса здоровья (Healthy / Unhealthy)
docker compose ps

# Перезапуск одного конкретного сервиса без простоя остальных
docker compose restart backend

# Выполнить одноразовую команду (например, миграции БД)
docker compose run --rm backend ./migrate.sh

# Полная остановка с удалением контейнеров и сетей (тома сохраняются!)
docker compose down

# Полная остановка с удалением томов (⚠️ уничтожит базы данных!)
docker compose down -v
```
