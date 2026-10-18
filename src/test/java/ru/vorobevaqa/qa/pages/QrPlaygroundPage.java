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
package ru.vorobevaqa.qa.pages;

import static com.codeborne.selenide.Condition.attribute;
import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.value;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

public class QrPlaygroundPage extends BasePage<QrPlaygroundPage> {

  private final SelenideElement pageTitle = $("h1");
  private final SelenideElement urlInput = $("#url-input");
  private final SelenideElement generateButton = $("#generate-btn");
  private final SelenideElement qrResultContainer = $("#result");
  private final SelenideElement qrCodeImage = $("#qr-code-image");
  private final SelenideElement errorContainer = $("#error");
  private final SelenideElement errorText = $("#error-text");
  private final SelenideElement metaSize = $("#qrMetaSize");
  private final SelenideElement metaTime = $("#qrMetaTime");
  private final SelenideElement downloadButton = $("button[onclick*='downloadQR']");
  private final SelenideElement copyButton = $("#copyImgBtn");

  @Step("Открыть страницу QR Playground")
  public QrPlaygroundPage openPage() {
    open("/qr");
    pageTitle.shouldBe(visible);
    return this;
  }

  @Step("Ввести данные/URL: {url}")
  public QrPlaygroundPage enterUrl(String url) {
    urlInput.shouldBe(visible).clear();
    urlInput.setValue(url);
    return this;
  }

  @Step("Выбрать быстрый шаблон: {presetUrl}")
  public QrPlaygroundPage selectPreset(String presetUrl) {
    $(".qr-preset-pill[data-preset='" + presetUrl + "']").shouldBe(visible).click();
    return this;
  }

  @Step("Выбрать цветовой пресет: {hexColor}")
  public QrPlaygroundPage selectColor(String hexColor) {
    $(".color-swatch[data-color='" + hexColor + "']").shouldBe(visible).click();
    return this;
  }

  @Step("Нажать кнопку генерации QR-кода")
  public QrPlaygroundPage clickGenerate() {
    generateButton.shouldBe(visible).click();
    return this;
  }

  @Step("Проверить, что QR-код успешно сгенерирован и отображается")
  public QrPlaygroundPage shouldSeeQrImage() {
    qrResultContainer.shouldBe(visible);
    qrCodeImage.shouldBe(visible).shouldHave(attribute("src"));
    metaSize.shouldBe(visible).shouldHave(text("300x300 px"));
    return this;
  }

  @Step("Проверить отображение ошибки: {expectedError}")
  public QrPlaygroundPage shouldSeeError(String expectedError) {
    errorContainer.shouldBe(visible);
    errorText.shouldHave(text(expectedError));
    return this;
  }

  @Step("Проверить значение в поле ввода URL: {expectedUrl}")
  public QrPlaygroundPage shouldHaveUrlInputValue(String expectedUrl) {
    urlInput.shouldHave(value(expectedUrl));
    return this;
  }

  @Step("Проверить доступность кнопок скачивания и копирования")
  public QrPlaygroundPage shouldHaveExportButtons() {
    downloadButton.shouldBe(visible);
    copyButton.shouldBe(visible);
    return this;
  }
}
