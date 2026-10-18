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

import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

public class ProjectsPage extends BasePage<ProjectsPage> {

  private final SelenideElement pageTitle = $("h1");
  private final SelenideElement qrPlaygroundLink = $("a[href='/qr']");

  @Step("Открыть страницу Проектов")
  public ProjectsPage openPage() {
    open("/projects");
    pageTitle.shouldBe(visible);
    return this;
  }

  @Step("Проверить заголовок страницы проектов")
  public ProjectsPage shouldSeeProjectsTitle() {
    pageTitle.shouldBe(visible).shouldHave(text("Инженерные Решения"));
    qrPlaygroundLink.shouldBe(visible);
    return this;
  }
}
