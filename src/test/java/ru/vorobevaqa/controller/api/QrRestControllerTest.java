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
package ru.vorobevaqa.controller.api;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.delete;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.put;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import ru.vorobevaqa.dto.QrRequest;
import ru.vorobevaqa.exception.GlobalExceptionHandler;
import ru.vorobevaqa.service.QrCodeService;

@Epic("REST API")
@Feature("QR Code Generator")
@Story("QrRestController - Spring MVC обработка запросов на генерацию QR")
@Tag("unit")
@Tag("controller")
@WebMvcTest(QrRestController.class)
@Import(GlobalExceptionHandler.class)
class QrRestControllerTest {

  @Autowired private MockMvc mockMvc;

  @Autowired private ObjectMapper objectMapper;

  @MockBean private QrCodeService qrCodeService;

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("API-REST-01: POST /api/qr - успешная генерация возвращает 200 и валидный JSON")
  @Description("Проверка возврата HTTP 200 OK, JSON-структуры с url и base64 qrCode")
  void shouldReturnQrCodeOnValidRequest() throws Exception {
    QrRequest request = new QrRequest();
    request.setUrl("https://vorobevaqa.ru");

    String fakeBase64 = "fakeBase64String";
    when(qrCodeService.generateQrCodeBase64(anyString())).thenReturn(fakeBase64);

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isOk())
        .andExpect(content().contentType(MediaType.APPLICATION_JSON))
        .andExpect(jsonPath("$.url").value("https://vorobevaqa.ru"))
        .andExpect(jsonPath("$.qrCode").value(fakeBase64));

    verify(qrCodeService).generateQrCodeBase64("https://vorobevaqa.ru");
  }

  @ParameterizedTest
  @ValueSource(strings = {"invalid-url", "ftp://example.com", "javascript:alert(1)", "   "})
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("API-REST-02: POST /api/qr - ошибка валидации URL возвращает 400 Bad Request")
  @Description(
      "Проверка валидации протоколов http:// и https:// и возврата ошибки от"
          + " GlobalExceptionHandler")
  void shouldReturnBadRequestWhenUrlIsInvalid(String invalidUrl) throws Exception {
    QrRequest request = new QrRequest();
    request.setUrl(invalidUrl);

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isBadRequest())
        .andExpect(jsonPath("$.error").exists());

    verifyNoInteractions(qrCodeService);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-03: POST /api/qr - пустой или null URL возвращает 400 Bad Request")
  @Description("Проверка валидации @NotBlank для поля url в теле запроса")
  void shouldReturnBadRequestWhenUrlIsNull() throws Exception {
    QrRequest request = new QrRequest();
    request.setUrl(null);

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isBadRequest())
        .andExpect(jsonPath("$.error").exists());

    verifyNoInteractions(qrCodeService);
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName(
      "API-REST-04: POST /api/qr - обработка ошибки сервиса возвращает 500 Internal Server Error")
  @Description("Проверка обработки RuntimeException глобальным обработчиком GlobalExceptionHandler")
  void shouldReturnInternalServerErrorWhenServiceFails() throws Exception {
    QrRequest request = new QrRequest();
    request.setUrl("https://vorobevaqa.ru");

    when(qrCodeService.generateQrCodeBase64(anyString()))
        .thenThrow(new RuntimeException("ZXing encoding error"));

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isInternalServerError())
        .andExpect(jsonPath("$.error").value("Internal server error: ZXing encoding error"));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-05: POST /api/qr - пустой JSON объект {} возвращает 400 Bad Request")
  @Description("Проверка валидации обязательности поля url при передаче пустого JSON")
  void shouldReturnBadRequestWhenEmptyJsonProvided() throws Exception {
    mockMvc
        .perform(post("/api/qr").contentType(MediaType.APPLICATION_JSON).content("{}"))
        .andExpect(status().isBadRequest())
        .andExpect(jsonPath("$.error").value("URL is required"));

    verifyNoInteractions(qrCodeService);
  }

  @ParameterizedTest
  @ValueSource(
      strings = {
        "http://vorobevaqa.ru",
        "https://sub.domain.example.com:8443/api/v1/path?query=1&sort=asc#section",
        "https://мойдомен.рф/тест",
        "http://127.0.0.1:8080/ping"
      })
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-06: POST /api/qr - поддержка различных валидных структур и протоколов URL")
  @Description(
      "Проверка генерации QR для HTTP, HTTPS, нестандартных портов, query-параметров и кириллицы")
  void shouldAcceptVariousValidUrlFormats(String validUrl) throws Exception {
    QrRequest request = new QrRequest();
    request.setUrl(validUrl);

    when(qrCodeService.generateQrCodeBase64(validUrl)).thenReturn("mockQrBase64");

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isOk())
        .andExpect(jsonPath("$.url").value(validUrl))
        .andExpect(jsonPath("$.qrCode").value("mockQrBase64"));

    verify(qrCodeService).generateQrCodeBase64(validUrl);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-07: POST /api/qr - игнорирование неизвестных дополнительных полей в JSON")
  @Description("Проверка устойчивости контроллера при получении нештатных дополнительных полей")
  void shouldIgnoreUnknownFieldsInJson() throws Exception {
    String payloadWithExtra =
        "{\"url\": \"https://vorobevaqa.ru\", \"unknownField\": \"value\", \"id\": 999}";

    when(qrCodeService.generateQrCodeBase64("https://vorobevaqa.ru")).thenReturn("extraFieldQr");

    mockMvc
        .perform(post("/api/qr").contentType(MediaType.APPLICATION_JSON).content(payloadWithExtra))
        .andExpect(status().isOk())
        .andExpect(jsonPath("$.url").value("https://vorobevaqa.ru"))
        .andExpect(jsonPath("$.qrCode").value("extraFieldQr"));
  }

  @ParameterizedTest
  @ValueSource(strings = {MediaType.APPLICATION_XML_VALUE, MediaType.TEXT_PLAIN_VALUE})
  @Severity(SeverityLevel.NORMAL)
  @DisplayName(
      "API-REST-08: POST /api/qr - возврат 415 Unsupported Media Type при неверном Content-Type")
  @Description("Проверка возврата HTTP 415 при передаче XML или plain text вместо JSON")
  void shouldReturnUnsupportedMediaType(String contentType) throws Exception {
    mockMvc
        .perform(post("/api/qr").contentType(contentType).content("https://vorobevaqa.ru"))
        .andExpect(status().isUnsupportedMediaType());

    verifyNoInteractions(qrCodeService);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-09: GET /api/qr - возврат 405 Method Not Allowed при вызове GET")
  @Description("Проверка запрета метода GET для эндпоинта /api/qr")
  void shouldReturnMethodNotAllowedOnGet() throws Exception {
    mockMvc.perform(get("/api/qr")).andExpect(status().isMethodNotAllowed());
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-10: PUT /api/qr - возврат 405 Method Not Allowed при вызове PUT")
  @Description("Проверка запрета метода PUT для эндпоинта /api/qr")
  void shouldReturnMethodNotAllowedOnPut() throws Exception {
    mockMvc
        .perform(
            put("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content("{\"url\": \"https://vorobevaqa.ru\"}"))
        .andExpect(status().isMethodNotAllowed());
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-11: DELETE /api/qr - возврат 405 Method Not Allowed при вызове DELETE")
  @Description("Проверка запрета метода DELETE для эндпоинта /api/qr")
  void shouldReturnMethodNotAllowedOnDelete() throws Exception {
    mockMvc.perform(delete("/api/qr")).andExpect(status().isMethodNotAllowed());
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-12: POST /api/qr - поддержка длинного URL (>1000 символов)")
  @Description("Проверка корректной передачи длинного URL без искажений в бизнес-логику")
  void shouldHandleLongUrlSuccessfully() throws Exception {
    String longPath = "a".repeat(1000);
    String longUrl = "https://vorobevaqa.ru/" + longPath;

    QrRequest request = new QrRequest();
    request.setUrl(longUrl);

    when(qrCodeService.generateQrCodeBase64(longUrl)).thenReturn("longUrlQrBase64");

    mockMvc
        .perform(
            post("/api/qr")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
        .andExpect(status().isOk())
        .andExpect(jsonPath("$.url").value(longUrl))
        .andExpect(jsonPath("$.qrCode").value("longUrlQrBase64"));

    verify(qrCodeService).generateQrCodeBase64(longUrl);
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("API-REST-13: POST /api/qr - синтаксически некорректный JSON возвращает 400 или 500")
  @Description("Проверка реакции обработчика исключений при получении битого JSON (malformed JSON)")
  void shouldHandleMalformedJsonBody() throws Exception {
    String malformedJson = "{\"url\": \"https://vorobevaqa.ru\""; // незакрытая фигурная скобка

    mockMvc
        .perform(post("/api/qr").contentType(MediaType.APPLICATION_JSON).content(malformedJson))
        .andExpect(status().is5xxServerError()); // GlobalExceptionHandler handles RuntimeException
    // (HttpMessageNotReadableException)

    verifyNoInteractions(qrCodeService);
  }
}
