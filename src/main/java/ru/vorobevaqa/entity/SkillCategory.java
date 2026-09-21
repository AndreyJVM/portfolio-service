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
package ru.vorobevaqa.entity;

public enum SkillCategory {
  QA_AUTOMATION("Тестирование и автоматизация", "fas fa-vial"),
  DEVOPS_CICD("DevOps & CI/CD", "fas fa-infinity"),
  LINUX_INFRA("Linux, Сети и Администрирование", "fas fa-server"),
  BACKEND_DB("Backend, Java & Базы данных", "fas fa-database");

  private final String displayName;
  private final String iconClass;

  SkillCategory(String displayName, String iconClass) {
    this.displayName = displayName;
    this.iconClass = iconClass;
  }

  public String getDisplayName() {
    return displayName;
  }

  public String getIconClass() {
    return iconClass;
  }
}
