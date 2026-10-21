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
  private static final String PROP_TIMEOUT = "selenide.timeout";
  private static final String PROP_PAGE_LOAD_TIMEOUT = "selenide.page-load-timeout";
  private static final String PROP_PAGE_LOAD_STRATEGY = "selenide.page-load-strategy";
  private static final String PROP_POLLING_INTERVAL = "selenide.polling-interval";
  private static final String PROP_HOLD_BROWSER_OPEN = "selenide.hold-browser-open";
  private static final String PROP_CLICK_VIA_JS = "selenide.click-via-js";
  private static final String PROP_SCREENSHOTS = "allure.screenshots";
  private static final String PROP_SAVE_PAGE_SOURCE = "allure.save-page-source";
  private static final String PROP_REMOTE_URL = "selenium.remote.url";
  private static final String PROP_CI = "ci";
  private static final String PROP_GITHUB_ACTIONS = "github.actions";

  @BeforeEach
  @AfterEach
  void cleanUpSystemProperties() {
    System.clearProperty(PROP_TEST_BASE_URL);
    System.clearProperty(PROP_APP_BASE_URL);
    System.clearProperty(PROP_BROWSER);
    System.clearProperty(PROP_BROWSER_VERSION);
    System.clearProperty(PROP_BROWSER_SIZE);
    System.clearProperty(PROP_HEADLESS);
    System.clearProperty(PROP_TIMEOUT);
    System.clearProperty(PROP_PAGE_LOAD_TIMEOUT);
    System.clearProperty(PROP_PAGE_LOAD_STRATEGY);
    System.clearProperty(PROP_POLLING_INTERVAL);
    System.clearProperty(PROP_HOLD_BROWSER_OPEN);
    System.clearProperty(PROP_CLICK_VIA_JS);
    System.clearProperty(PROP_SCREENSHOTS);
    System.clearProperty(PROP_SAVE_PAGE_SOURCE);
    System.clearProperty(PROP_REMOTE_URL);
    System.clearProperty(PROP_CI);
    System.clearProperty(PROP_GITHUB_ACTIONS);
  }

  @Nested
  @DisplayName("1. Базовый URL (getBaseUrl & isCustomBaseUrlConfigured)")
  class BaseUrlTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-01: Дефолтный Base URL возвращает http://localhost:8080")
    @Description(
        "Проверка значения по умолчанию при отсутствии системных свойств и переменных окружения")
    void shouldReturnDefaultBaseUrl() {
      assertThat(TestConfig.getBaseUrl()).isEqualTo("http://localhost:8080");
      assertThat(TestConfig.isCustomBaseUrlConfigured()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.CRITICAL)
    @DisplayName("CFG-02: Чтение Base URL из системного свойства test.base.url")
    @Description("Проверка приоритетного считывания test.base.url")
    void shouldReturnBaseUrlFromTestBaseUrlProperty() {
      System.setProperty(PROP_TEST_BASE_URL, "https://qa.vorobevaqa.ru");
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://qa.vorobevaqa.ru");
      assertThat(TestConfig.isCustomBaseUrlConfigured()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-03: Фолбэк на свойство app.base-url при отсутствии test.base.url")
    @Description("Проверка считывания альтернативного ключа конфигурации app.base-url")
    void shouldFallbackToAppBaseUrlProperty() {
      System.setProperty(PROP_APP_BASE_URL, "https://staging.vorobevaqa.ru");
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://staging.vorobevaqa.ru");
      assertThat(TestConfig.isCustomBaseUrlConfigured()).isTrue();
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
      assertThat(TestConfig.isCustomBaseUrlConfigured()).isFalse();
    }

    @ParameterizedTest
    @ValueSource(strings = {"https://vorobevaqa.ru/", "https://vorobevaqa.ru///"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName(
        "CFG-07: Удаление закрывающего слэша (trailing slash) для предотвращения двойных слэшей")
    @Description("Проверка нормализации URL и удаления завершающих слэшей")
    void shouldStripTrailingSlashesFromBaseUrl(String urlWithSlash) {
      System.setProperty(PROP_TEST_BASE_URL, urlWithSlash);
      assertThat(TestConfig.getBaseUrl()).isEqualTo("https://vorobevaqa.ru");
    }

    @ParameterizedTest
    @ValueSource(
        strings = {
          "http://localhost:8080",
          "http://localhost:8080/",
          "http://127.0.0.1:8080",
          "http://127.0.0.1:8080/"
        })
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-08: isCustomBaseUrlConfigured возвращает false для локальных хостов")
    @Description("Проверка, что localhost и 127.0.0.1 не считаются внешними стендами")
    void shouldReturnFalseForLocalhostAndLoopbackBaseUrl(String localUrl) {
      System.setProperty(PROP_TEST_BASE_URL, localUrl);
      assertThat(TestConfig.isCustomBaseUrlConfigured()).isFalse();
    }
  }

  @Nested
  @DisplayName("2. Параметры браузера (Browser, Version, Size)")
  class BrowserSettingsTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-09: Дефолтный браузер - chrome")
    @Description("Проверка значения по умолчанию для selenium.browser")
    void shouldReturnDefaultBrowser() {
      assertThat(TestConfig.getBrowser()).isEqualTo("chrome");
    }

    @ParameterizedTest
    @ValueSource(strings = {"firefox", "edge", "safari"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-10: Чтение кастомного браузера из свойства selenium.browser")
    @Description("Проверка переопределения браузера через системные свойства")
    void shouldReturnCustomBrowserFromProperty(String browser) {
      System.setProperty(PROP_BROWSER, browser);
      assertThat(TestConfig.getBrowser()).isEqualTo(browser);
    }

    @ParameterizedTest
    @ValueSource(strings = {"CHROME", "Firefox", "EdGe"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-11: Приведение имени браузера к нижнему регистру (toLowerCase)")
    @Description("Проверка нормализации регистра названий браузеров")
    void shouldNormalizeBrowserNameToLowerCase(String mixedCaseBrowser) {
      System.setProperty(PROP_BROWSER, mixedCaseBrowser);
      assertThat(TestConfig.getBrowser()).isEqualTo(mixedCaseBrowser.toLowerCase());
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-12: Дефолтная версия браузера - latest")
    @Description("Проверка значения по умолчанию для selenium.browser.version")
    void shouldReturnDefaultBrowserVersion() {
      assertThat(TestConfig.getBrowserVersion()).isEqualTo("latest");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-13: Чтение версии браузера из свойства selenium.browser.version")
    @Description("Проверка переопределения версии браузера")
    void shouldReturnCustomBrowserVersionFromProperty() {
      System.setProperty(PROP_BROWSER_VERSION, "130.0");
      assertThat(TestConfig.getBrowserVersion()).isEqualTo("130.0");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-14: Дефолтное разрешение экрана - 1920x1080")
    @Description("Проверка значения по умолчанию для selenium.browser.size")
    void shouldReturnDefaultBrowserSize() {
      assertThat(TestConfig.getBrowserSize()).isEqualTo("1920x1080");
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-15: Чтение разрешения экрана из свойства selenium.browser.size")
    @Description("Проверка установки кастомного разрешения окна браузера")
    void shouldReturnCustomBrowserSizeFromProperty() {
      System.setProperty(PROP_BROWSER_SIZE, "1366x768");
      assertThat(TestConfig.getBrowserSize()).isEqualTo("1366x768");
    }
  }

  @Nested
  @DisplayName("3. Таймауты, интервалы и стратегии загрузки (Timeouts & Strategies)")
  class TimeoutsAndStrategiesTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-16: Дефолтный таймаут ожидания элементов - 10000 мс")
    @Description("Проверка значения по умолчанию для selenide.timeout")
    void shouldReturnDefaultTimeout() {
      assertThat(TestConfig.getTimeout()).isEqualTo(10000L);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-17: Установка кастомного таймаута ожидания элементов")
    @Description("Проверка чтения значения из selenide.timeout")
    void shouldReturnCustomTimeout() {
      System.setProperty(PROP_TIMEOUT, "15000");
      assertThat(TestConfig.getTimeout()).isEqualTo(15000L);
    }

    @ParameterizedTest
    @ValueSource(strings = {"invalid-long", "-5000", "0"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName(
        "CFG-18: Фолбэк на дефолт при невалидном или отрицательном значении selenide.timeout")
    @Description("Проверка защиты от нечисловых, нулевых и отрицательных значений таймаута")
    void shouldFallbackToDefaultTimeoutOnInvalidOrNonPositive(String invalidVal) {
      System.setProperty(PROP_TIMEOUT, invalidVal);
      assertThat(TestConfig.getTimeout()).isEqualTo(10000L);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-19: Дефолтный таймаут загрузки страницы - 60000 мс")
    @Description("Проверка значения по умолчанию для selenide.page-load-timeout")
    void shouldReturnDefaultPageLoadTimeout() {
      assertThat(TestConfig.getPageLoadTimeout()).isEqualTo(60000L);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-20: Установка кастомного таймаута загрузки страницы")
    @Description("Проверка чтения значения из selenide.page-load-timeout")
    void shouldReturnCustomPageLoadTimeout() {
      System.setProperty(PROP_PAGE_LOAD_TIMEOUT, "45000");
      assertThat(TestConfig.getPageLoadTimeout()).isEqualTo(45000L);
    }

    @ParameterizedTest
    @ValueSource(strings = {"abc", "-1", "0"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-21: Фолбэк на дефолт при некорректном selenide.page-load-timeout")
    @Description("Проверка безопасной обработки отрицательного или некорректного таймаута загрузки")
    void shouldFallbackToDefaultPageLoadTimeoutOnInvalid(String invalidVal) {
      System.setProperty(PROP_PAGE_LOAD_TIMEOUT, invalidVal);
      assertThat(TestConfig.getPageLoadTimeout()).isEqualTo(60000L);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-22: Дефолтная стратегия загрузки страниц - eager")
    @Description("Проверка значения по умолчанию для selenide.page-load-strategy")
    void shouldReturnDefaultPageLoadStrategy() {
      assertThat(TestConfig.getPageLoadStrategy()).isEqualTo("eager");
    }

    @ParameterizedTest
    @ValueSource(strings = {"normal", "none"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-23: Установка кастомной стратегии загрузки страницы")
    @Description("Проверка чтения стратегий normal и none")
    void shouldReturnCustomPageLoadStrategy(String strategy) {
      System.setProperty(PROP_PAGE_LOAD_STRATEGY, strategy);
      assertThat(TestConfig.getPageLoadStrategy()).isEqualTo(strategy);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-24: Дефолтный интервал опроса элементов - 200 мс")
    @Description("Проверка значения по умолчанию для selenide.polling-interval")
    void shouldReturnDefaultPollingInterval() {
      assertThat(TestConfig.getPollingInterval()).isEqualTo(200L);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-25: Установка кастомного интервала опроса элементов")
    @Description("Проверка чтения значения из selenide.polling-interval")
    void shouldReturnCustomPollingInterval() {
      System.setProperty(PROP_POLLING_INTERVAL, "350");
      assertThat(TestConfig.getPollingInterval()).isEqualTo(350L);
    }

    @ParameterizedTest
    @ValueSource(strings = {"invalid-interval", "-20", "0"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-26: Фолбэк на дефолт при отрицательном или некорректном интервале опроса")
    @Description("Проверка безопасной обработки интервала опроса")
    void shouldFallbackToDefaultPollingIntervalOnInvalid(String invalidVal) {
      System.setProperty(PROP_POLLING_INTERVAL, invalidVal);
      assertThat(TestConfig.getPollingInterval()).isEqualTo(200L);
    }
  }

  @Nested
  @DisplayName("4. Расширенные параметры Selenide и Allure")
  class AdvancedSettingsTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-27: Дефолтный holdBrowserOpen - false")
    @Description("Проверка значения по умолчанию для selenide.hold-browser-open")
    void shouldReturnDefaultHoldBrowserOpen() {
      assertThat(TestConfig.isHoldBrowserOpen()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-28: Включение holdBrowserOpen через свойство selenide.hold-browser-open")
    @Description("Проверка удержания браузера открытым после тестов")
    void shouldEnableHoldBrowserOpen() {
      System.setProperty(PROP_HOLD_BROWSER_OPEN, "true");
      assertThat(TestConfig.isHoldBrowserOpen()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-29: Дефолтный clickViaJs - false")
    @Description("Проверка значения по умолчанию для selenide.click-via-js")
    void shouldReturnDefaultClickViaJs() {
      assertThat(TestConfig.isClickViaJs()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-30: Включение клика через JS через свойство selenide.click-via-js")
    @Description("Проверка включения эмуляции кликов через JavaScript")
    void shouldEnableClickViaJs() {
      System.setProperty(PROP_CLICK_VIA_JS, "true");
      assertThat(TestConfig.isClickViaJs()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-31: Дефолтное создание скриншотов Allure - true")
    @Description("Проверка значения по умолчанию для allure.screenshots")
    void shouldReturnDefaultScreenshotsEnabled() {
      assertThat(TestConfig.isScreenshotsEnabled()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-32: Отключение создания скриншотов Allure через свойство allure.screenshots")
    @Description("Проверка отключения автоматических скриншотов Allure")
    void shouldDisableScreenshots() {
      System.setProperty(PROP_SCREENSHOTS, "false");
      assertThat(TestConfig.isScreenshotsEnabled()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-33: Дефолтное сохранение Page Source в Allure - true")
    @Description("Проверка значения по умолчанию для allure.save-page-source")
    void shouldReturnDefaultSavePageSourceEnabled() {
      assertThat(TestConfig.isSavePageSourceEnabled()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-34: Отключение сохранения Page Source через свойство allure.save-page-source")
    @Description("Проверка отключения сохранения исходного HTML страниц при падениях")
    void shouldDisableSavePageSource() {
      System.setProperty(PROP_SAVE_PAGE_SOURCE, "false");
      assertThat(TestConfig.isSavePageSourceEnabled()).isFalse();
    }
  }

  @Nested
  @DisplayName("5. Headless, Remote Grid, CI и Рефлексия (Headless, Remote, CI, Reflection)")
  class HeadlessRemoteAndCiTests {

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-35: Дефолтный headless режим - true")
    @Description("Проверка значения по умолчанию для selenium.headless")
    void shouldReturnDefaultHeadlessTrue() {
      assertThat(TestConfig.isHeadless()).isTrue();
    }

    @ParameterizedTest
    @ValueSource(strings = {"false", "FALSE", "False"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-36: Отключение headless режима через свойство selenium.headless")
    @Description("Проверка отключения headless режима (false) независимо от регистра")
    void shouldReturnHeadlessFalseWhenConfigured(String headlessVal) {
      System.setProperty(PROP_HEADLESS, headlessVal);
      assertThat(TestConfig.isHeadless()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-37: По умолчанию Remote Grid не сконфигурирован")
    @Description("Проверка getRemoteUrl().isEmpty() и isRemote() == false")
    void shouldReturnRemoteDisabledByDefault() {
      assertThat(TestConfig.getRemoteUrl()).isEmpty();
      assertThat(TestConfig.isRemote()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.CRITICAL)
    @DisplayName("CFG-38: Конфигурация Remote Selenium Grid / Selenoid")
    @Description("Проверка считывания URL удалённого грида и флага isRemote()")
    void shouldConfigureRemoteGridUrl() {
      String remoteHub = "http://selenoid.local:4444/wd/hub";
      System.setProperty(PROP_REMOTE_URL, remoteHub);

      assertThat(TestConfig.getRemoteUrl()).contains(remoteHub);
      assertThat(TestConfig.isRemote()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName(
        "CFG-39: Дефолтное состояние isCi() в зависимости от окружения или явного отключения")
    @Description("Проверка выявления окружения CI/CD или дефолтного значения")
    void shouldReturnDefaultCiOrRespectExplicitFalse() {
      if (System.getenv("CI") != null || System.getenv("GITHUB_ACTIONS") != null) {
        assertThat(TestConfig.isCi()).isTrue();
      } else {
        assertThat(TestConfig.isCi()).isFalse();
      }

      System.setProperty(PROP_CI, "false");
      System.setProperty(PROP_GITHUB_ACTIONS, "false");
      assertThat(TestConfig.isCi()).isFalse();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-40: Определение CI окружения через свойство ci=true")
    @Description("Проверка активации флага CI через системное свойство ci")
    void shouldDetectCiViaCiProperty() {
      System.setProperty(PROP_CI, "true");
      assertThat(TestConfig.isCi()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("CFG-41: Определение CI окружения через свойство github.actions=true")
    @Description("Проверка активации флага CI через github.actions")
    void shouldDetectCiViaGitHubActionsProperty() {
      System.setProperty(PROP_GITHUB_ACTIONS, "true");
      assertThat(TestConfig.isCi()).isTrue();
    }

    @Test
    @Severity(SeverityLevel.MINOR)
    @DisplayName("CFG-42: Проверка приватного конструктора утилитного класса")
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
