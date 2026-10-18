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

import static com.codeborne.selenide.Selenide.$;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import ru.vorobevaqa.qa.pages.HomePage;

@Epic("Web Application")
@Feature("Theme Management")
@Story("Переключение темы (Dark / Light mode)")
@Tag("ui")
public class ThemeToggleUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-TH-01: Переключение темы оформления в шапке сайта")
  @Description("Проверка корректной смены data-bs-theme при клике на переключатель темы")
  void shouldToggleThemeBetweenLightAndDark() {
    HomePage homePage = new HomePage().openPage();

    String currentTheme = $("html").getAttribute("data-bs-theme");
    String nextTheme = "dark".equalsIgnoreCase(currentTheme) ? "light" : "dark";

    homePage.toggleTheme().shouldHaveTheme(nextTheme);

    // Обратное переключение
    homePage.toggleTheme().shouldHaveTheme(currentTheme);
  }
}
