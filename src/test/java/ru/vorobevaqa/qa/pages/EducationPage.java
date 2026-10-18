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
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.$$;
import static com.codeborne.selenide.Selenide.open;

import com.codeborne.selenide.ElementsCollection;
import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

public class EducationPage extends BasePage<EducationPage> {

  private final SelenideElement pageTitle = $("h1");
  private final SelenideElement printButton = $("button[onclick*='window.print']");
  private final ElementsCollection diplomaCards = $$(".card");
  private final SelenideElement certModal = $("#certModal");
  private final SelenideElement modalImage = $("#certModalImage");

  @Step("Открыть страницу Образование и сертификаты")
  public EducationPage openPage() {
    open("/education");
    pageTitle.shouldBe(visible);
    return this;
  }

  @Step("Проверить заголовок и кнопку печати резюме (CV)")
  public EducationPage shouldSeeEducationHeader() {
    pageTitle.shouldBe(visible).shouldHave(text("Квалификация и"));
    printButton.shouldBe(visible);
    return this;
  }

  @Step("Проверить наличие карточек образования и дипломов")
  public EducationPage shouldHaveDiplomaCards() {
    diplomaCards.shouldHave(com.codeborne.selenide.CollectionCondition.sizeGreaterThan(0));
    return this;
  }

  @Step("Открыть модальное окно просмотра сертификата")
  public EducationPage openFirstCertificateModal() {
    $(".cursor-pointer").shouldBe(visible).click();
    certModal.shouldBe(visible);
    modalImage.shouldBe(visible).shouldHave(attribute("src"));
    return this;
  }
}
