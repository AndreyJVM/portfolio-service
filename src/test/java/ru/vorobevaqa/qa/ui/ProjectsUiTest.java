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
package ru.vorobevaqa.qa.ui;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import ru.vorobevaqa.qa.pages.ProjectsPage;

@Epic("Web Application")
@Feature("Projects & Solutions")
@Story("Просмотр инженерных проектов и репозиториев")
@Tag("ui")
public class ProjectsUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("UI-PRJ-01: Отображение страницы проектов и заголовка")
  @Description("Проверка рендеринга страницы /projects, заголовка и быстрых кнопок действий")
  void shouldDisplayProjectsPageHeader() {
    new ProjectsPage().openPage().shouldSeeProjectsTitle();
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("UI-PRJ-02: Переход со страницы проектов в QR Playground по кнопке в шапке")
  @Description("Проверка перехода по клику на кнопку QR Playground")
  void shouldNavigateToQrPlaygroundFromProjectsHeader() {
    new ProjectsPage()
        .openPage()
        .clickQrPlaygroundButton()
        .shouldHaveUrlInputValue("https://vorobevaqa.ru");
  }
}
