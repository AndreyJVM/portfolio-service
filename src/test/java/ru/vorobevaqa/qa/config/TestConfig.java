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

  private TestConfig() {}

  public static String getBaseUrl() {
    return getPropertyOrEnv("test.base.url", "TEST_BASE_URL")
        .or(() -> getPropertyOrEnv("app.base-url", "APP_BASE_URL"))
        .orElse("http://localhost:8080");
  }

  public static String getBrowser() {
    return getPropertyOrEnv("selenium.browser", "SELENIUM_BROWSER").orElse("chrome");
  }

  public static String getBrowserVersion() {
    return getPropertyOrEnv("selenium.browser.version", "SELENIUM_BROWSER_VERSION")
        .orElse("latest");
  }

  public static String getBrowserSize() {
    return getPropertyOrEnv("selenium.browser.size", "SELENIUM_BROWSER_SIZE").orElse("1920x1080");
  }

  public static boolean isHeadless() {
    return Boolean.parseBoolean(
        getPropertyOrEnv("selenium.headless", "SELENIUM_HEADLESS").orElse("true"));
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
}
