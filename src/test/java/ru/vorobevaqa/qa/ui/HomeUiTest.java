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
import ru.vorobevaqa.qa.pages.HomePage;

@Epic("Web Application")
@Feature("Home Page & Presentation")
@Story("Презентация и навигация на главной странице")
@Tag("ui")
public class HomeUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("UI-HOME-01: Отображение главной страницы, Hero-секции и ключевых компетенций")
  @Description(
      "Проверка рендеринга Hero-секции, аватара, описания профиля и карточек стека технологий")
  void shouldDisplayHomePageHeroAndBentoGrid() {
    new HomePage()
        .openPage()
        .shouldSeeHeroSection()
        .shouldHaveBentoCards()
        .shouldHaveValidContactLink();
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-HOME-02: Переход в портфолио проектов по кнопке в Hero-секции")
  @Description("Проверка перехода на страницу проектов при нажатии на кнопку 'Смотреть портфолио'")
  void shouldNavigateToProjectsFromHeroButton() {
    new HomePage()
        .openPage()
        .clickViewPortfolio()
        .shouldSeeProjectsHeader()
        .shouldHaveProjectCards();
  }
}
