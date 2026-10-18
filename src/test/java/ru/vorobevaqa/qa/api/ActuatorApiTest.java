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

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.equalTo;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

@Epic("REST API")
@Feature("Actuator Endpoints")
@Story("GET /actuator/health - Мониторинг жизнеспособности сервиса")
@Tag("api")
public class ActuatorApiTest extends BaseApiTest {

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("API-ACT-01: Проверка доступности Spring Boot Actuator Health")
  @Description("Проверка статус-кода 200 OK и статуса UP для /actuator/health")
  void shouldReturnHealthStatusUp() {
    given()
        .accept(ContentType.JSON)
        .when()
        .get("/actuator/health")
        .then()
        .statusCode(200)
        .body("status", equalTo("UP"));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-ACT-02: Проверка эндпоинта Actuator Info")
  @Description("Проверка статус-кода 200 OK для эндпоинта /actuator/info")
  void shouldReturnInfoEndpointOk() {
    given().accept(ContentType.JSON).when().get("/actuator/info").then().statusCode(200);
  }
}
