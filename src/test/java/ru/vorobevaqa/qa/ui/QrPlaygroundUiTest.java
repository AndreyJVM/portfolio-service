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
import ru.vorobevaqa.qa.pages.QrPlaygroundPage;

@Epic("Web Application")
@Feature("QR Code Playground")
@Story("Генерация и кастомизация QR-кодов")
@Tag("ui")
public class QrPlaygroundUiTest extends BaseUiTest {

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("UI-QR-01: Успешная генерация QR-кода по умолчанию")
  @Description("Проверка начального состояния страницы QR Playground и генерации QR-кода")
  void shouldGenerateQrForDefaultUrl() {
    new QrPlaygroundPage()
        .openPage()
        .shouldHaveUrlInputValue("https://vorobevaqa.ru")
        .clickGenerate()
        .shouldSeeQrImage()
        .shouldHaveExportButtons();
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("UI-QR-02: Подстановка URL через быстрые шаблоны (Presets)")
  @Description("Проверка выбора шаблонов (GitHub, Telegram) и обновления поля ввода")
  void shouldGenerateQrWhenPresetSelected() {
    new QrPlaygroundPage()
        .openPage()
        .selectPreset("https://github.com/AndreyJVM")
        .shouldHaveUrlInputValue("https://github.com/AndreyJVM")
        .clickGenerate()
        .shouldSeeQrImage();
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-QR-03: Выбор цветового стиля QR-кода")
  @Description("Проверка выбора неонового зеленого и синего цветов для генерации QR-кода")
  void shouldChangeColorPreset() {
    new QrPlaygroundPage().openPage().selectColor("#22c55e").clickGenerate().shouldSeeQrImage();
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("UI-QR-04: Клиентская валидация протокола URL")
  @Description("Проверка отображения ошибки при вводе URL без протокола http:// или https://")
  void shouldShowErrorForInvalidUrl() {
    new QrPlaygroundPage()
        .openPage()
        .enterUrl("not-a-valid-url")
        .clickGenerate()
        .shouldSeeError("Ссылка должна начинаться с http:// или https://");
  }
}
