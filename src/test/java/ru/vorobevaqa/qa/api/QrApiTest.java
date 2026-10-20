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
import static org.hamcrest.Matchers.emptyOrNullString;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.not;
import static org.hamcrest.Matchers.notNullValue;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import io.restassured.http.ContentType;
import java.util.Map;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

@Epic("REST API")
@Feature("QR Code Generator")
@Story("POST /api/qr - Генерация QR-кода")
@Tag("api")
public class QrApiTest extends BaseApiTest {

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("API-QR-01: Успешная генерация QR-кода для валидного URL")
  @Description("Проверка статус-кода 200 OK, валидного Base64 и соответствия входного URL")
  void shouldGenerateQrCodeSuccessfully() {
    given()
        .contentType(ContentType.JSON)
        .body(Map.of("url", "https://vorobevaqa.ru"))
        .when()
        .post("/api/qr")
        .then()
        .statusCode(200)
        .contentType(ContentType.JSON)
        .body("url", equalTo("https://vorobevaqa.ru"))
        .body("qrCode", not(emptyOrNullString()));
  }

  @ParameterizedTest
  @ValueSource(strings = {"ftp://site.com", "not-a-url", "javascript:alert(1)"})
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("API-QR-02: Ошибка валидации для невалидного протокола URL")
  @Description("Проверка возврата 400 Bad Request с сообщением валидации шаблона URL")
  void shouldReturnBadRequestWhenUrlIsInvalid(String invalidUrl) {
    given()
        .contentType(ContentType.JSON)
        .body(Map.of("url", invalidUrl))
        .when()
        .post("/api/qr")
        .then()
        .statusCode(400)
        .contentType(ContentType.JSON)
        .body("error", equalTo("URL must start with http:// or https://"));
  }

  @ParameterizedTest
  @ValueSource(strings = {"", "   "})
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-QR-03: Ошибка валидации для пустого или пробельного URL")
  @Description("Проверка возврата 400 Bad Request при передаче пустого URL")
  void shouldReturnBadRequestWhenUrlIsBlank(String blankUrl) {
    given()
        .contentType(ContentType.JSON)
        .body(Map.of("url", blankUrl))
        .when()
        .post("/api/qr")
        .then()
        .statusCode(400)
        .contentType(ContentType.JSON)
        .body("error", notNullValue());
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-QR-04: Ошибка валидации при отсутствии поля URL (null)")
  @Description("Проверка возврата 400 Bad Request при передаче JSON с null URL")
  void shouldReturnBadRequestWhenUrlIsNull() {
    given()
        .contentType(ContentType.JSON)
        .body("{\"url\": null}")
        .when()
        .post("/api/qr")
        .then()
        .statusCode(400)
        .contentType(ContentType.JSON)
        .body("error", equalTo("URL is required"));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-QR-05: Ошибка 405 Method Not Allowed при вызове GET /api/qr")
  @Description("Проверка возврата статуса 405 Method Not Allowed при неподдерживаемом HTTP методе")
  void shouldReturnMethodNotAllowedOnGetRequest() {
    given().when().get("/api/qr").then().statusCode(405);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-QR-06: Ошибка 415 Unsupported Media Type при неверном Content-Type")
  @Description(
      "Проверка возврата 415 при передаче тела с Content-Type text/plain вместо application/json")
  void shouldReturnUnsupportedMediaTypeOnPlainText() {
    given()
        .contentType(ContentType.TEXT)
        .body("https://vorobevaqa.ru")
        .when()
        .post("/api/qr")
        .then()
        .statusCode(415);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-QR-07: Ошибка валидации 400 при передаче пустого JSON {}")
  @Description("Проверка обязательности параметров и возврата 400 Bad Request")
  void shouldReturnBadRequestOnEmptyJsonBody() {
    given()
        .contentType(ContentType.JSON)
        .body("{}")
        .when()
        .post("/api/qr")
        .then()
        .statusCode(400)
        .contentType(ContentType.JSON)
        .body("error", equalTo("URL is required"));
  }
}
