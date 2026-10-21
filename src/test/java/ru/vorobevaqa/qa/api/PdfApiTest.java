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
import static org.assertj.core.api.Assertions.assertThat;
import static org.hamcrest.Matchers.containsString;
import static org.hamcrest.Matchers.equalTo;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import io.restassured.response.Response;
import java.nio.charset.StandardCharsets;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

@Epic("REST API")
@Feature("PDF Resume Generator")
@Story("GET /api/pdf - Генерация и скачивание динамического PDF резюме")
@Tag("api")
public class PdfApiTest extends BaseApiTest {

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("API-PDF-01: Успешное получение PDF резюме (200 OK, inline, magic bytes %PDF-)")
  @Description(
      "Проверка статус-кода 200, Content-Type application/pdf, заголовков Content-Disposition и"
          + " сигнатуры PDF")
  void shouldReturnPdfResumeSuccessfully() {
    Response response =
        given()
            .when()
            .get("/api/pdf/resume")
            .then()
            .statusCode(200)
            .contentType("application/pdf")
            .header(
                "Content-Disposition",
                containsString("inline; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\""))
            .header("Cache-Control", containsString("public, max-age=3600"))
            .extract()
            .response();

    byte[] body = response.getBody().asByteArray();
    assertThat(body).isNotNull();
    assertThat(body.length).isGreaterThan(200);

    String magic = new String(body, 0, 5, StandardCharsets.US_ASCII);
    assertThat(magic).isEqualTo("%PDF-");
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("API-PDF-02: Получение PDF резюме с параметром download=true (attachment)")
  @Description(
      "Проверка формирования заголовка Content-Disposition: attachment при запросе на"
          + " принудительное скачивание")
  void shouldReturnAttachmentPdfWhenDownloadRequested() {
    given()
        .queryParam("download", "true")
        .when()
        .get("/api/pdf/resume")
        .then()
        .statusCode(200)
        .contentType("application/pdf")
        .header(
            "Content-Disposition",
            containsString("attachment; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\""));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-PDF-03: Проверка эндпоинта здоровья микросервиса PDF /api/pdf/health")
  @Description("Проверка возврата JSON {status: UP, service: pdf-service}")
  void shouldReturnPdfServiceHealth() {
    given()
        .when()
        .get("/api/pdf/health")
        .then()
        .statusCode(200)
        .contentType("application/json")
        .body("status", equalTo("UP"))
        .body("service", equalTo("pdf-service"));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-PDF-04: Метод POST на /api/pdf/resume возвращает 405 Method Not Allowed")
  @Description("Проверка протокольных ограничений HTTP методов")
  void shouldRejectPostMethodOnPdfResume() {
    given().when().post("/api/pdf/resume").then().statusCode(405);
  }
}
