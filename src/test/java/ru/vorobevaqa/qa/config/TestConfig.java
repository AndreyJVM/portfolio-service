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

import java.util.Optional;

public final class TestConfig {

  public static final String DEFAULT_BASE_URL = "http://localhost:8080";
  public static final String DEFAULT_BROWSER = "chrome";
  public static final String DEFAULT_BROWSER_VERSION = "latest";
  public static final String DEFAULT_BROWSER_SIZE = "1920x1080";
  public static final boolean DEFAULT_HEADLESS = true;
  public static final long DEFAULT_TIMEOUT = 10000L;
  public static final long DEFAULT_PAGE_LOAD_TIMEOUT = 60000L;
  public static final String DEFAULT_PAGE_LOAD_STRATEGY = "eager";
  public static final long DEFAULT_POLLING_INTERVAL = 200L;
  public static final boolean DEFAULT_HOLD_BROWSER_OPEN = false;
  public static final boolean DEFAULT_CLICK_VIA_JS = false;
  public static final boolean DEFAULT_SCREENSHOTS = true;
  public static final boolean DEFAULT_SAVE_PAGE_SOURCE = true;

  private TestConfig() {}

  public static String getBaseUrl() {
    String url =
        getPropertyOrEnv("test.base.url", "TEST_BASE_URL")
            .or(() -> getPropertyOrEnv("app.base-url", "APP_BASE_URL"))
            .orElse(DEFAULT_BASE_URL);
    return stripTrailingSlash(url);
  }

  public static boolean isCustomBaseUrlConfigured() {
    return getPropertyOrEnv("test.base.url", "TEST_BASE_URL")
        .or(() -> getPropertyOrEnv("app.base-url", "APP_BASE_URL"))
        .map(TestConfig::stripTrailingSlash)
        .filter(
            url ->
                !url.equals(DEFAULT_BASE_URL)
                    && !url.contains("localhost:8080")
                    && !url.contains("127.0.0.1:8080"))
        .isPresent();
  }

  public static String getBrowser() {
    return getPropertyOrEnv("selenium.browser", "SELENIUM_BROWSER")
        .map(String::toLowerCase)
        .orElse(DEFAULT_BROWSER);
  }

  public static String getBrowserVersion() {
    return getPropertyOrEnv("selenium.browser.version", "SELENIUM_BROWSER_VERSION")
        .orElse(DEFAULT_BROWSER_VERSION);
  }

  public static String getBrowserSize() {
    return getPropertyOrEnv("selenium.browser.size", "SELENIUM_BROWSER_SIZE")
        .orElse(DEFAULT_BROWSER_SIZE);
  }

  public static boolean isHeadless() {
    return getBooleanPropertyOrEnv("selenium.headless", "SELENIUM_HEADLESS", DEFAULT_HEADLESS);
  }

  public static long getTimeout() {
    return getLongPropertyOrEnv("selenide.timeout", "SELENIDE_TIMEOUT", DEFAULT_TIMEOUT);
  }

  public static long getPageLoadTimeout() {
    return getLongPropertyOrEnv(
        "selenide.page-load-timeout", "SELENIDE_PAGE_LOAD_TIMEOUT", DEFAULT_PAGE_LOAD_TIMEOUT);
  }

  public static String getPageLoadStrategy() {
    return getPropertyOrEnv("selenide.page-load-strategy", "SELENIDE_PAGE_LOAD_STRATEGY")
        .orElse(DEFAULT_PAGE_LOAD_STRATEGY);
  }

  public static long getPollingInterval() {
    return getLongPropertyOrEnv(
        "selenide.polling-interval", "SELENIDE_POLLING_INTERVAL", DEFAULT_POLLING_INTERVAL);
  }

  public static boolean isHoldBrowserOpen() {
    return getBooleanPropertyOrEnv(
        "selenide.hold-browser-open", "SELENIDE_HOLD_BROWSER_OPEN", DEFAULT_HOLD_BROWSER_OPEN);
  }

  public static boolean isClickViaJs() {
    return getBooleanPropertyOrEnv(
        "selenide.click-via-js", "SELENIDE_CLICK_VIA_JS", DEFAULT_CLICK_VIA_JS);
  }

  public static boolean isScreenshotsEnabled() {
    return getBooleanPropertyOrEnv("allure.screenshots", "ALLURE_SCREENSHOTS", DEFAULT_SCREENSHOTS);
  }

  public static boolean isSavePageSourceEnabled() {
    return getBooleanPropertyOrEnv(
        "allure.save-page-source", "ALLURE_SAVE_PAGE_SOURCE", DEFAULT_SAVE_PAGE_SOURCE);
  }

  public static Optional<String> getRemoteUrl() {
    return getPropertyOrEnv("selenium.remote.url", "SELENIUM_REMOTE_URL");
  }

  public static boolean isRemote() {
    return getRemoteUrl().isPresent();
  }

  public static boolean isCi() {
    return getBooleanPropertyOrEnv("ci", "CI", false)
        || getBooleanPropertyOrEnv("github.actions", "GITHUB_ACTIONS", false);
  }

  private static String stripTrailingSlash(String url) {
    if (url != null && url.endsWith("/")) {
      return url.replaceAll("/+$", "");
    }
    return url;
  }

  private static Optional<String> getPropertyOrEnv(String propKey, String envKey) {
    String sysProp = System.getProperty(propKey);
    if (sysProp != null && !sysProp.isBlank()) {
      return Optional.of(sysProp.trim());
    }
    String envVal = System.getenv(envKey);
    if (envVal != null && !envVal.isBlank()) {
      return Optional.of(envVal.trim());
    }
    return Optional.empty();
  }

  private static boolean getBooleanPropertyOrEnv(
      String propKey, String envKey, boolean defaultValue) {
    return getPropertyOrEnv(propKey, envKey).map(Boolean::parseBoolean).orElse(defaultValue);
  }

  private static long getLongPropertyOrEnv(String propKey, String envKey, long defaultValue) {
    return getPropertyOrEnv(propKey, envKey)
        .map(
            val -> {
              try {
                long parsed = Long.parseLong(val);
                return parsed > 0 ? parsed : defaultValue;
              } catch (NumberFormatException e) {
                return defaultValue;
              }
            })
        .orElse(defaultValue);
  }
}
