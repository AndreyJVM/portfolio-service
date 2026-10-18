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
package ru.vorobevaqa.qa.api;

import io.qameta.allure.restassured.AllureRestAssured;
import io.restassured.RestAssured;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import ru.vorobevaqa.qa.config.TestConfig;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public abstract class BaseApiTest {

  @LocalServerPort private int serverPort;

  @BeforeEach
  void setUpRestAssured() {
    RestAssured.enableLoggingOfRequestAndResponseIfValidationFails();

    String configuredUrl = TestConfig.getBaseUrl();
    if (configuredUrl != null && !configuredUrl.contains("localhost:8080")) {
      RestAssured.baseURI = configuredUrl;
    } else {
      RestAssured.baseURI = "http://localhost";
      RestAssured.port = serverPort;
    }

    RestAssured.replaceFiltersWith(new AllureRestAssured());
  }
}
