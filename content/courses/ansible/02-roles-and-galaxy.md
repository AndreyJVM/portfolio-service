# Модульные роли и публикация в Ansible Galaxy

Роли в Ansible позволяют структурировать код плейбуков, разделяя задачи, шаблоны, переменные и обработчики на независимые переиспользуемые модули.

---

### 1. Стандартная структура роли

Инициализация новой роли с помощью `ansible-galaxy`:
```bash
ansible-galaxy role init docker_host
```

Дерево каталогов роли:
```
docker_host/
├── defaults/
│   └── main.yml        # Переменные по умолчанию (низкий приоритет)
├── vars/
│   └── main.yml        # Строгие переменные роли
├── tasks/
│   └── main.yml        # Список исполняемых задач
├── handlers/
│   └── main.yml        # Обработчики событий (триггеры notify)
├── templates/
│   └── daemon.json.j2  # Jinja2 шаблоны конфигураций
├── files/
│   └── docker.gpg      # Статические файлы
└── meta/
    └── main.yml        # Метаданные роли для публикации в Galaxy
```

---

### 2. Пример роли: Настройка Docker Daemon и Hardening

В файле `tasks/main.yml`:
```yaml
---
- name: Установка необходимых пакетов для репозитория Docker
  ansible.builtin.apt:
    name:
      - ca-certificates
      - curl
      - gnupg
    state: present

- name: Добавить официальный GPG-ключ Docker
  ansible.builtin.get_url:
    url: https://download.docker.com/linux/ubuntu/gpg
    dest: /etc/apt/keyrings/docker.asc
    mode: '0644'

- name: Настроить параметры Docker daemon (ротация логов и cgroups)
  ansible.builtin.template:
    src: daemon.json.j2
    dest: /etc/docker/daemon.json
    mode: '0644'
  notify: Restart Docker
```

Шаблон `templates/daemon.json.j2`:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "no-new-privileges": true,
  "live-restore": true
}
```

---

### 3. Шифрование чувствительных данных: Ansible Vault

Для хранения паролей, приватных ключей и токенов используется встроенное симметричное шифрование AES-256:

```bash
# Создание зашифрованного файла с переменными
ansible-vault create vars/secrets.yml

# Запуск плейбука с запросом пароля от Vault
ansible-playbook -i inventory.ini site.yml --ask-vault-pass
```

---

### 4. Публикация и переиспользование из Ansible Galaxy

Для совместного использования ролей команды публикуют их в репозиторий [Ansible Galaxy](https://galaxy.ansible.com/):

```yaml
# requirements.yml
roles:
  - name: andreyvorobev.docker_hardened
    src: https://github.com/AndreyJVM/ansible-role-docker
    version: 1.2.0
```

Установка зависимостей:
```bash
ansible-galaxy install -r requirements.yml
```
