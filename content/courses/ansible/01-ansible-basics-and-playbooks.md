# Основы Ansible: Архитектура, Inventory и первый Playbook

Ansible — это инструмент автоматизации конфигураций и оркестрации с архитектурой **Agentless** (работает через SSH без установки агентов на целевые ноды).

---

### 1. Архитектура и ключевые концепции

- **Control Node:** Машина, с которой запускаются команды (Linux/macOS c установленным Python и `ansible-core`).
- **Managed Nodes:** Целевые серверы (VDS, bare-metal), на которых Ansible выполняет модули через временные скрипты Python.
- **Идемпотентность (Idempotence):** Повторный запуск плейбука не меняет состояние системы, если оно уже соответствует желаемому.

---

### 2. Описание инфраструктуры: Inventory

Файл инвентаря описывает хосты, группы и переменные подключения:

```ini
# inventory.ini
[webservers]
vds-primary ansible_host=194.87.x.x ansible_user=root ansible_ssh_private_key_file=~/.ssh/id_ed25519

[databases]
db-primary ansible_host=10.0.0.5 ansible_user=ubuntu

[production:children]
webservers
databases
```

Проверка доступности хостов с помощью ad-hoc команды:
```bash
ansible -i inventory.ini webservers -m ping
```

---

### 3. Первый идемпотентный Playbook

Создадим `setup-nginx.yml` для установки и настройки веб-сервера:

```yaml
---
- name: Развертывание и конфигурация веб-сервера Nginx
  hosts: webservers
  become: true
  vars:
    server_domain: "vorobevaqa.ru"
    app_port: 8080

  tasks:
    - name: Обновить apt-кэш и установить Nginx
      ansible.builtin.apt:
        name: nginx
        state: present
        update_cache: true
        cache_valid_time: 3600

    - name: Создать конфигурацию VirtualHost для приложения
      ansible.builtin.template:
        src: templates/nginx-site.conf.j2
        dest: /etc/nginx/sites-available/{{ server_domain }}.conf
        mode: '0644'
      notify: Reload Nginx

    - name: Активировать сайт через symlink
      ansible.builtin.file:
        src: /etc/nginx/sites-available/{{ server_domain }}.conf
        dest: /etc/nginx/sites-enabled/{{ server_domain }}.conf
        state: link
      notify: Reload Nginx

    - name: Убедиться, что служба Nginx включена и запущена
      ansible.builtin.service:
        name: nginx
        state: started
        enabled: true

  handlers:
    - name: Reload Nginx
      ansible.builtin.service:
        name: nginx
        state: reloaded
```

---

### 4. Запуск и отладка

```bash
# Проверка синтаксиса
ansible-playbook -i inventory.ini setup-nginx.yml --syntax-check

# Dry-run (тестовый прогон без реальных изменений)
ansible-playbook -i inventory.ini setup-nginx.yml --check --diff

# Боевое применение
ansible-playbook -i inventory.ini setup-nginx.yml
```
