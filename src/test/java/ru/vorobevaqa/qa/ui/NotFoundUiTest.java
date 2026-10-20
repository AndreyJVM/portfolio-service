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

import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

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
@Feature("Error Handling")
@Story("Отображение кастомной страницы 404 Not Found")
@Tag("ui")
public class NotFoundUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-ERR-01: Отображение кастомной 404 страницы для несуществующего маршрута")
  @Description(
      "Проверка рендеринга фирменного шаблона 404 с кнопками возврата при обращении к"
          + " несуществующему URL")
  void shouldDisplayCustom404Page() {
    open("/non-existing-page-xyz");
    $("h1").shouldBe(visible).shouldHave(text("404 — Страница не найдена"));
    $("a[href='/']").shouldBe(visible).click();
    new HomePage().shouldSeeHeroSection();
  }
}
