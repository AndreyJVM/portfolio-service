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
package ru.vorobevaqa.qa.config;

import static org.assertj.core.api.Assertions.assertThat;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Modifier;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

@Epic("QA Infrastructure")
@Feature("Test Configuration")
@Story("TestConfig - Управление параметрами тестовой инфраструктуры")
@Tag("unit")
@Tag("config")
class TestConfigTest {

  private static final String PROP_TEST_BASE_URL = "test.base.url";
  private static final String PROP_APP_BASE_URL = "app.base-url";
  private static final String PROP_BROWSER = "selenium.browser";
  private static final String PROP_BROWSER_VERSION = "selenium.browser.version";
  private static final String PROP_BROWSER_SIZE = "selenium.browser.size";
  private static final String PROP_HEADLESS = "selenium.headless";

  @BeforeEach
  @AfterEach
  void cleanUpSystemProperties() {
    System.clearProperty(PROP_TEST_BASE_URL);
    System.clearProperty(PROP_APP_BASE_URL);
    System.clearProperty(PROP_BROWSER);
    System.clearProperty(PROP_BROWSER_VERSION);
    System.clearProperty(PROP_BROWSER_SIZE);
    System.clearProperty(PROP_HEADLESS);
  }

  @Nested
  @DisplayName("1. Базовый URL (getBaseUrl)")
  class BaseUrlTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-01: Дефолтный Base URL возвращает http://localhost:8080")
    @Description(
        "Проверка значения по умолчанию при отсутствии системных свойств и переменных окружения")
    void shouldReturnDefaultBaseUrl() {
      assertThat(TestConfig.getBaseUrl()).isEqualTo("http://localhost:8080");
    }

    @Test
    @Severity(SeverityLevel.CRITICAL)
    @DisplayName("CFG-02: Чтение Base URL из системного свойства test.base.url")
    @Description("Проверка приоритетного считывания test.base.url")
    void shouldReturnBaseUrlFromTestBaseUrlProperty() {
      System.setProperty(PROP_TEST_BASE_URL, "https://qa.vorobevaqa.ru");
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://qa.vorobevaqa.ru");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-03: Фолбэк на свойство app.base-url при отсутствии test.base.url")
    @Description("Проверка считывания альтернативного ключа конфигурации app.base-url")
    void shouldFallbackToAppBaseUrlProperty() {
      System.setProperty(PROP_APP_BASE_URL, "https://staging.vorobevaqa.ru");
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://staging.vorobevaqa.ru");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-04: Приоритет test.base.url над app.base-url")
    @Description("Проверка порядка разрешения свойств при одновременной установке обоих ключей")
    void shouldPrioritizeTestBaseUrlOverAppBaseUrl() {
      System.setProperty(PROP_TEST_BASE_URL, "https://primary.vorobevaqa.ru");
      System.setProperty(PROP_APP_BASE_URL, "https://secondary.vorobevaqa.ru");

      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://primary.vorobevaqa.ru");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-05: Обрезка начальных и конечных пробелов в Base URL")
    @Description("Проверка trim() для строковых параметров конфигурации")
    void shouldTrimWhitespaceFromBaseUrl() {
      System.setProperty(PROP_TEST_BASE_URL, "   https://vorobevaqa.ru   ");
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://vorobevaqa.ru");
    }

    @ParameterizedTest
    @ValueSource(strings = {"", "   ", "\t", "\n"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-06: Игнорирование пустых/пробельных значений системного свойства")
    @Description("Проверка возврата дефолта при передаче пустых или пробельных строк")
    void shouldIgnoreBlankPropertyAndReturnDefault(String blankValue) {
      System.setProperty(PROP_TEST_BASE_URL, blankValue);
      assertThat(TestConfig.getBaseUrl()).isEqualTo("http://localhost:8080");
    }
  }

  @Nested
  @DisplayName("2. Параметры браузера (Browser, Version, Size)")
  class BrowserSettingsTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-07: Дефолтный браузер - chrome")
    @Description("Проверка значения по умолчанию для selenium.browser")
    void shouldReturnDefaultBrowser() {
      assertThat(TestConfig.getBrowser()).isEqualTo("chrome");
    }

    @ParameterizedTest
    @ValueSource(strings = {"firefox", "edge", "safari"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-08: Чтение кастомного браузера из свойства selenium.browser")
    @Description("Проверка переопределения браузера через системные свойства")
    void shouldReturnCustomBrowserFromProperty(String browser) {
      System.setProperty(PROP_BROWSER, browser);
      assertThat(TestConfig.getBrowser()).isEqualTo(browser);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-09: Дефолтная версия браузера - latest")
    @Description("Проверка значения по умолчанию для selenium.browser.version")
    void shouldReturnDefaultBrowserVersion() {
      assertThat(TestConfig.getBrowserVersion()).isEqualTo("latest");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-10: Чтение версии браузера из свойства selenium.browser.version")
    @Description("Проверка переопределения версии браузера")
    void shouldReturnCustomBrowserVersionFromProperty() {
      System.setProperty(PROP_BROWSER_VERSION, "130.0");
      assertThat(TestConfig.getBrowserVersion()).isEqualTo("130.0");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-11: Дефолтное разрешение экрана - 1920x1080")
    @Description("Проверка значения по умолчанию для selenium.browser.size")
    void shouldReturnDefaultBrowserSize() {
      assertThat(TestConfig.getBrowserSize()).isEqualTo("1920x1080");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-12: Чтение разрешения экрана из свойства selenium.browser.size")
    @Description("Проверка установки кастомного разрешения окна браузера")
    void shouldReturnCustomBrowserSizeFromProperty() {
      System.setProperty(PROP_BROWSER_SIZE, "1366x768");
      assertThat(TestConfig.getBrowserSize()).isEqualTo("1366x768");
    }
  }

  @Nested
  @DisplayName("3. Headless режим и рефлексия (isHeadless)")
  class HeadlessAndReflectionTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-13: Дефолтный headless режим - true")
    @Description("Проверка значения по умолчанию для selenium.headless")
    void shouldReturnDefaultHeadlessTrue() {
      assertThat(TestConfig.isHeadless()).isTrue();
    }

    @ParameterizedTest
    @ValueSource(strings = {"false", "FALSE", "False"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-14: Отключение headless режима через свойство selenium.headless")
    @Description("Проверка отключения headless режима (false) независимо от регистра")
    void shouldReturnHeadlessFalseWhenConfigured(String headlessVal) {
      System.setProperty(PROP_HEADLESS, headlessVal);
      assertThat(TestConfig.isHeadless()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.MINOR)
    @DisplayName("CFG-15: Проверка приватного конструктора утилитного класса")
    @Description("Проверка приватного модификатора конструктора через рефлексию")
    void shouldHavePrivateConstructor()
        throws NoSuchMethodException,
            InvocationTargetException,
            InstantiationException,
            IllegalAccessException {
      Constructor<TestConfig> constructor = TestConfig.class.getDeclaredConstructor();
      assertThat(Modifier.isPrivate(constructor.getModifiers())).isTrue();

      constructor.setAccessible(true);
      TestConfig instance = constructor.newInstance();
      assertThat(instance).isNotNull();
    }
  }
}
