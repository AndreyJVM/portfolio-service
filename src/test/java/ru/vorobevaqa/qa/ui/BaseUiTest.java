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

import com.codeborne.selenide.Configuration;
import com.codeborne.selenide.Selenide;
import com.codeborne.selenide.logevents.SelenideLogger;
import io.qameta.allure.selenide.AllureSelenide;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.openqa.selenium.chrome.ChromeOptions;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import ru.vorobevaqa.qa.config.TestConfig;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public abstract class BaseUiTest {

  @LocalServerPort private int serverPort;

  @BeforeAll
  static void setUpSelenideConfig() {
    Configuration.browser = TestConfig.getBrowser();
    Configuration.browserSize = TestConfig.getBrowserSize();
    Configuration.headless = TestConfig.isHeadless();
    Configuration.timeout = 10000;
    Configuration.pageLoadTimeout = 30000;

    ChromeOptions options = new ChromeOptions();
    options.addArguments("--no-sandbox", "--disable-dev-shm-usage", "--disable-gpu");
    Configuration.browserCapabilities = options;

    if (!SelenideLogger.hasListener("AllureSelenide")) {
      SelenideLogger.addListener(
          "AllureSelenide",
          new AllureSelenide().screenshots(true).savePageSource(true).includeSelenideSteps(true));
    }
  }

  @BeforeEach
  void initBaseUrl() {
    String configuredUrl = TestConfig.getBaseUrl();
    if (configuredUrl != null && !configuredUrl.contains("localhost:8080")) {
      Configuration.baseUrl = configuredUrl;
    } else {
      Configuration.baseUrl = "http://localhost:" + serverPort;
    }
  }

  @AfterAll
  static void tearDown() {
    Selenide.closeWebDriver();
  }
}
