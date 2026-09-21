/*
 * Copyright 2024-2026 Andrey Vorobev
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
package ru.vorobevaqa.service;

import java.util.Arrays;
import java.util.List;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.vorobevaqa.dto.SkillDto;
import ru.vorobevaqa.entity.Skill;
import ru.vorobevaqa.entity.SkillCategory;
import ru.vorobevaqa.entity.SkillLevel;
import ru.vorobevaqa.repository.SkillRepository;

@Slf4j
@Service
@RequiredArgsConstructor
public class SkillService implements CommandLineRunner {

  private final SkillRepository skillRepository;

  @Transactional(readOnly = true)
  public List<SkillDto> getAllSkills() {
    return skillRepository.findAllByOrderByDisplayOrderAscIdAsc().stream()
        .map(SkillDto::fromEntity)
        .toList();
  }

  @Transactional(readOnly = true)
  public List<SkillDto> getSkillsByCategory(SkillCategory category) {
    return skillRepository.findByCategoryOrderByDisplayOrderAscIdAsc(category).stream()
        .map(SkillDto::fromEntity)
        .toList();
  }

  @Override
  @Transactional
  public void run(String... args) {
    if (skillRepository.count() == 0) {
      log.info("No skills found in database. Initializing default skills matrix...");
      seedInitialSkills();
      log.info("Skills matrix successfully initialized with {} records.", skillRepository.count());
    }
  }

  private void seedInitialSkills() {
    List<Skill> defaultSkills =
        Arrays.asList(
            // QA Automation
            Skill.builder()
                .name("Автоматизация API и UI")
                .category(SkillCategory.QA_AUTOMATION)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Проектирование E2E и интеграционных автотестов, разработка тестовых"
                        + " фреймворков.")
                .keyTools("JUnit 5, RestAssured, Selenium/Playwright, Allure Report, Postman")
                .iconClass("fas fa-vial")
                .displayOrder(1)
                .build(),
            Skill.builder()
                .name("Нагрузочное и стресс-тестирование")
                .category(SkillCategory.QA_AUTOMATION)
                .level(SkillLevel.INTERMEDIATE)
                .description(
                    "Создание сценариев нагрузки API, сбор метрик времени отклика и утилизации"
                        + " ресурсов.")
                .keyTools("JMeter, k6, Grafana, Locust")
                .iconClass("fas fa-chart-line")
                .displayOrder(2)
                .build(),

            // DevOps & CI/CD
            Skill.builder()
                .name("CI/CD пайплайны")
                .category(SkillCategory.DEVOPS_CICD)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Построение автоматических сборок, прогона тестов, сборки образов и деплоя на"
                        + " VDS.")
                .keyTools("GitHub Actions, GitLab CI, Docker, Docker Compose")
                .iconClass("fas fa-infinity")
                .displayOrder(3)
                .build(),
            Skill.builder()
                .name("Контейнеризация и инфраструктура")
                .category(SkillCategory.DEVOPS_CICD)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Управление multi-stage Dockerfiles, Docker Compose, оптимизация размера слоёв"
                        + " и безопасность образов.")
                .keyTools("Docker, Docker Compose, Registry/Docker Hub, Nginx Reverse Proxy")
                .iconClass("fab fa-docker")
                .displayOrder(4)
                .build(),

            // Linux & Administration
            Skill.builder()
                .name("Linux System Administration")
                .category(SkillCategory.LINUX_INFRA)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Настройка и поддержка серверов (Debian/Ubuntu), systemd сервисы, мониторинг"
                        + " ресурсов, bash-скриптинг.")
                .keyTools("Bash, Systemd, SSH, UFW/iptables, Journalctl, Cron")
                .iconClass("fab fa-linux")
                .displayOrder(5)
                .build(),
            Skill.builder()
                .name("Сетевые протоколы и сервисы")
                .category(SkillCategory.LINUX_INFRA)
                .level(SkillLevel.INTERMEDIATE)
                .description(
                    "Конфигурация сетевого взаимодействия, DNS, TLS/SSL сертификаты (Certbot),"
                        + " Samba shares, роутинг.")
                .keyTools("TCP/IP, HTTP/HTTPS, SSL/TLS, Let's Encrypt, Samba, Nginx")
                .iconClass("fas fa-network-wired")
                .displayOrder(6)
                .build(),

            // Backend & DB
            Skill.builder()
                .name("Java & Spring Boot")
                .category(SkillCategory.BACKEND_DB)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Разработка REST API сервисов, MVC веб-приложений, конфигурирование Spring"
                        + " Security и Data JPA.")
                .keyTools("Java 17, Spring Boot 3, Spring Data JPA, Maven, Lombok")
                .iconClass("fab fa-java")
                .displayOrder(7)
                .build(),
            Skill.builder()
                .name("Реляционные БД и SQL")
                .category(SkillCategory.BACKEND_DB)
                .level(SkillLevel.PROFICIENT)
                .description(
                    "Проектирование схем, написание сложных SQL-запросов, транзакционность,"
                        + " миграции.")
                .keyTools("PostgreSQL, DBeaver, JDBC, Hibernate")
                .iconClass("fas fa-database")
                .displayOrder(8)
                .build());

    skillRepository.saveAll(defaultSkills);
  }
}
