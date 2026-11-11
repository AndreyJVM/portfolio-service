# Права доступа в Linux и управление сервисами systemd

Управление правами доступа и жизненным циклом демонов — фундаментальный навык системного инженера.

---

### 1. Модель прав доступа: DAC (Discretionary Access Control)

Каждый файл и директория в Linux принадлежат:
- **Пользователю (User / Owner)**
- **Группе (Group)**
- **Остальным (Others)**

```bash
# Пример вывода ls -l
-rwxr-xr-- 1 devops webapps 4096 Nov 11 20:00 deploy.sh
```

- `r` (read) = 4
- `w` (write) = 2
- `x` (execute) = 1

#### Практические команды:
```bash
# Выставить права только владельцу на приватный ключ SSH:
chmod 600 ~/.ssh/id_ed25519

# Сменить владельца и группу рекурсивно:
sudo chown -R devops:webapps /var/www/app
```

---

### 2. Управление демонами через systemd

Systemd — стандартная система инициализации и менеджер служб в большинстве дистрибутивов (Ubuntu, Debian, RHEL).

#### Анатомия unit-файла (`/etc/systemd/system/myapp.service`):
```ini
[Unit]
Description=DevOps Platform Go Daemon
After=network.target

[Service]
Type=simple
User=appuser
Group=appgroup
WorkingDirectory=/opt/app
ExecStart=/opt/app/server
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

#### Основные команды:
```bash
sudo systemctl daemon-reload        # Перечитать конфигурацию
sudo systemctl enable --now myapp   # Включить в автозагрузку и запустить
sudo systemctl status myapp         # Проверить статус
journalctl -u myapp -f              # Читать логи сервиса в реальном времени
```
