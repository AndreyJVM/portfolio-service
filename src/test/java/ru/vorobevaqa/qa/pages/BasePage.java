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
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

public abstract class BasePage<T extends BasePage<T>> {

  protected final SelenideElement themeToggle = $("#themeToggle");
  protected final SelenideElement themeIcon = $("#themeIcon");
  protected final SelenideElement navHome = $("a.nav-link[href='/']");
  protected final SelenideElement navEducation = $("a.nav-link[href='/education']");
  protected final SelenideElement navProjects = $("a.nav-link[href='/projects']");
  protected final SelenideElement navQr = $("a.nav-link[href='/qr']");
  protected final SelenideElement brandLogo = $("a.navbar-brand");

  @Step("Переключить тему оформления (Dark / Light)")
  @SuppressWarnings("unchecked")
  public T toggleTheme() {
    themeToggle.shouldBe(visible).click();
    return (T) this;
  }

  @Step("Проверить, что установлена тема оформления: {expectedTheme}")
  @SuppressWarnings("unchecked")
  public T shouldHaveTheme(String expectedTheme) {
    $("html").shouldHave(attribute("data-bs-theme", expectedTheme));
    return (T) this;
  }

  @Step("Перейти в раздел QR Playground через навигацию")
  public QrPlaygroundPage goToQrPlayground() {
    navQr.shouldBe(visible).click();
    return new QrPlaygroundPage();
  }

  @Step("Перейти на Главную через навигацию")
  public HomePage goToHome() {
    navHome.shouldBe(visible).click();
    return new HomePage();
  }

  @Step("Перейти в Проекты через навигацию")
  public ProjectsPage goToProjects() {
    navProjects.shouldBe(visible).click();
    return new ProjectsPage();
  }
}
