---
title: "Шпаргалка команд Docker и Docker Compose (Cheat Sheet)"
description: "Исчерпывающий справочник консольных команд для повседневной работы: управление контейнерами, сетями, томами, логами, инспекция Go-шаблонами и отладочные контейнеры"
duration: "20 мин"
level: "Справочник • Все уровни"
order: 1
tags: ["Cheat Sheet", "CLI", "Commands", "Docker Compose", "Inspect", "Netshoot", "Troubleshooting"]
---
# Шпаргалка команд Docker и Docker Compose (Cheat Sheet)

Консольная шпаргалка с самыми важными и эффективными командами для повседневной работы DevOps-инженера и разработчика. Без воды — только флаги, синтаксис, фильтры и практические примеры.

---

## 1. Управление контейнерами (Lifecycle)

```bash
# Запустить контейнер в фоне с автоперезапуском и пробросом порта на localhost
docker run -d --name web --restart unless-stopped -p 127.0.0.1:8080:80 nginx:alpine

# Интерактивный вход в запущенный контейнер с командной оболочкой bash или sh
docker exec -it web sh

# Посмотреть активные контейнеры (с отображением реального размера слоя diff)
docker ps -s

# Посмотреть ВСЕ контейнеры (включая завершившиеся аварийно со статусом Exited)
docker ps -a --filter "status=exited"

# Быстро и принудительно удалить контейнер
docker rm -f web

# Мониторинг потребления ресурсов (CPU, RAM, Net I/O, PIDs) в реальном времени
docker stats --no-stream
```

---

## 2. Диагностика и Логи

```bash
# Просмотр последних 100 строк логов с живым стримингом (follow)
docker logs -f --tail 100 web

# Логи с отображением точного времени (RFC3339 timestamps)
docker logs -t --tail 50 web

# Просмотр изменений файловой системы контейнера относительно исходного образа
# (A - добавлен, C - изменен, D - удален)
docker diff web

# Просмотр процессов, запущенных внутри контейнера на хосте
docker top web aux

# Системные события демона Docker в реальном времени (die, kill, oom, restart)
docker events --filter 'event=die' --filter 'event=oom'
```

---

## 3. Продвинутая инспекция через шаблоны Go (`docker inspect`)

Используйте флаг `--format` или `-f`, чтобы извлекать точные поля конфигурации без обработки гигантского JSON-вывода:

```bash
# 1. Узнать статус OOM Killer и ExitCode упавшего контейнера
docker inspect web --format 'ExitCode: {{.State.ExitCode}}, OOMKilled: {{.State.OOMKilled}}'

# 2. Узнать точный внутренний IP-адрес контейнера в сети
docker inspect web -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

# 3. Узнать статус готовности Healthcheck и последний вывод пробы
docker inspect web -f '{{.State.Health.Status}} (вывод: {{range .State.Health.Log}}{{.Output}}{{end}})'

# 4. Просмотреть все точки монтирования (Volumes и Bind Mounts)
docker inspect web -f '{{range .Mounts}}{{.Type}} -> {{.Source}} : {{.Destination}} ({{.Mode}}){{"\n"}}{{end}}'

# 5. Узнать PID процесса контейнера в главном пространстве имен хоста Linux
docker inspect web -f '{{.State.Pid}}'
```

---

## 4. Образы и Движок BuildKit

```bash
# Сборка образа с использованием кэша и явным тегом
docker build -t myapp:v1.2.0 .

# Сборка без кэша (чистая сборка с нуля)
docker build --no-cache -t myapp:v1.2.0 .

# Безопасная сборка с передачей секрета BuildKit (секрет не сохраняется в слоях!)
DOCKER_BUILDKIT=1 docker build --secret id=npmrc,src=$HOME/.npmrc -t myapp:secure .

# Сборка с ограничением памяти сборщика (защита от падения сборки на слабых серверах)
docker build --memory 2g -t myapp:v1 .

# Просмотр дерева слоев и веса каждой команды
docker history --human myapp:v1.2.0

# Просмотр объема кэша сборщика BuildKit
docker buildx du

# Очистка дискового кэша сборщика BuildKit
docker buildx prune -f
```

---

## 5. Сети и Постоянство Данных (Volumes)

```bash
# Создать изолированную пользовательскую сеть-мост
docker network create --driver bridge internal-net

# Создать изолированную сеть без выхода в интернет (internal)
docker network create --driver bridge --internal secure-backend-net

# Подключить уже работающий контейнер ко второй сети на лету
docker network connect internal-net web

# Отключить контейнер от сети
docker network disconnect bridge web

# Создать резервную копию Named Volume в tar.gz архив на хосте
docker run --rm -v pg_data:/data -v $(pwd):/backup alpine \
  tar czvf /backup/pg_data_$(date +%F).tar.gz -C /data .

# Восстановить Named Volume из архива
docker run --rm -v pg_data:/data -v $(pwd):/backup alpine \
  tar xzvf /backup/pg_data_backup.tar.gz -C /data
```

---

## 6. Трюки Production-отладки (Ephemeral Debug Containers)

Если контейнер собран на базе `scratch` или `distroless`, внутри нет `curl`, `netstat`, `ping`, `sh`. Как исследовать его сетевые и процессные проблемы без остановки сервиса?

Используйте совместное использование пространств имен Linux (**Namespace Sharing**):

```bash
# 1. Подключиться швейцарским ножом сетевой диагностики (Netshoot)
# в сетевой стек целевого контейнера 'web':
docker run --rm -it --net=container:web nicolaka/netshoot

# Внутри открытой консоли вам доступны:
# tcpdump, iftop, curl, dig, mtr, ss, ngrep, tshark!
# Все запросы выполняются прямо из сетевого интерфейса контейнера 'web'!

# 2. Инспекция файловой системы и процессов целевого контейнера:
# Находим PID целевого процесса:
PID=$(docker inspect -f '{{.State.Pid}}' web)

# Входим в сетевой стек через системную утилиту nsenter хоста:
sudo nsenter -t $PID -n ss -tulpen
```

---

## 7. Docker Compose v2 (Шпаргалка CLI)

```bash
# Запуск мультиконтейнерного стека в фоне с пересборкой измененных образов
docker compose up -d --build

# Просмотр статуса сервисов, проброшенных портов и проверок healthcheck
docker compose ps

# Чтение логов конкретного сервиса с таймстемпами
docker compose logs -f -t --tail 50 backend

# Перезапуск одного сервиса без влияния на соседние контейнеры
docker compose restart redis

# Выполнение одноразовой команды в изолированном контейнере (например, накатка миграций)
docker compose run --rm backend ./migrate -env=production

# Валидация синтаксиса compose.yaml и просмотр интерполированных переменных
docker compose config

# Грациозная остановка и удаление контейнеров и сетей (тома сохраняются!)
docker compose down

# Полное удаление стека вместе с томами (⚠️ УДАЛИТ ДАННЫЕ БАЗ ДАННЫХ!)
docker compose down -v
```

---

## 8. Безопасная очистка диска

```bash
# 1. Посмотреть детальную раскладку занятого места Docker:
docker system df -v

# 2. Безопасная чистка (удаляет остановленные контейнеры, повисшие образы <none> и неиспользуемые сети):
docker system prune -f

# 3. Полная очистка неиспользуемых образов (включая образы без запущенных контейнеров):
docker image prune -a -f --filter "until=168h"   # Только образы старше 7 дней

# ⚠️ ОПАСНО: Удаление ВСЕХ непримонтированных томов (никогда не использовать в cron!):
docker volume prune -f
```

---

## 9. Лабораторный вызов: Проверь себя

<details>
<summary><b>Вопрос 1:</b> Как проинспектировать сетевой трафик контейнера на базе scratch, если внутри нет утилит и шелла?</summary>

> **Ответ:** Запустить временный отладочный контейнер в общем сетевом пространстве имен целевого контейнера: `docker run --rm -it --net=container:<target_name> nicolaka/netshoot`. Внутри будут доступны все сетевые утилиты (`tcpdump`, `curl`, `ss`), разделяющие сетевой интерфейс и IP с целевым контейнером.
</details>

<details>
<summary><b>Вопрос 2:</b> Чем отличается docker compose down от docker compose down -v?</summary>

> **Ответ:** Команда `docker compose down` останавливает и удаляет только контейнеры и пользовательские сети, оставляя именованные тома (Volumes) нетронутыми. Флаг `-v` (`--volumes`) безвозвратно удаляет все объявленные тома, что приведет к полному уничтожению баз данных и постоянных данных.
</details>
