# Сетевой стек Linux: порты, сокеты и фаервол UFW

Понимание сетевого стека Linux необходимо для диагностики сетевых сбоев и защиты серверов.

---

### 1. Проверка открытых портов и сокетов (ss / netstat)

Утилита `ss` (Socket Statistics) — современная и быстрая замена устаревшему `netstat`.

```bash
# Посмотреть все слушающие TCP порты с процессами:
sudo ss -tulpn

# Отфильтровать порты 80 и 443:
sudo ss -tulpn | grep -E ':(80|443)'
```

Флаги `ss`:
- `-t`: TCP сокеты
- `-u`: UDP сокеты
- `-l`: Только слушающие (Listening)
- `-p`: Имя процесса и PID (требует `sudo`)
- `-n`: Числовые порты и IP без DNS-резолвинга

---

### 2. Защита периметра: Uncomplicated Firewall (UFW)

UFW — удобная обертка над подсистемой iptables/nftables в ядре Linux.

#### Базовый чек-лист настройки на боевом VDS:
```bash
# 1. Запретить весь входящий трафик по умолчанию
sudo ufw default deny incoming

# 2. Разрешить весь исходящий трафик
sudo ufw default allow outgoing

# 3. КРИТИЧЕСКИ ВАЖНО: открыть порт SSH перед включением!
sudo ufw allow 22/tcp

# 4. Открыть порты для Web-сервера (Nginx)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 5. Активировать фаервол и проверить правила
sudo ufw enable
sudo ufw status verbose
```
