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
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.request;
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
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.ResultActions;
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

  private static final String QR_ENDPOINT = "/api/qr";
  private static final String DEFAULT_VALID_URL = "https://vorobevaqa.ru";
  private static final String MOCK_BASE64_QR = "iVBORw0KGgoAAAANSUhEUgAA...";

  @Autowired private MockMvc mockMvc;
  @Autowired private ObjectMapper objectMapper;

  @MockitoBean private QrCodeService qrCodeService;

  // ==========================================
  // Helper Methods for clean, readable DSL
  // ==========================================

  private static QrRequest requestFor(String url) {
    QrRequest request = new QrRequest();
    request.setUrl(url);
    return request;
  }

  private ResultActions performPost(QrRequest request) throws Exception {
    return performPostRawJson(objectMapper.writeValueAsString(request));
  }

  private ResultActions performPostRawJson(String json) throws Exception {
    return mockMvc.perform(
        post(QR_ENDPOINT)
            .contentType(MediaType.APPLICATION_JSON)
            .accept(MediaType.APPLICATION_JSON)
            .content(json));
  }

  // ==========================================
  // 1. Успешные сценарии (200 OK)
  // ==========================================
  @Nested
  @DisplayName("1. Успешная генерация QR-кода (200 OK)")
  class SuccessScenarios {

    @Test
    @Severity(SeverityLevel.BLOCKER)
    @DisplayName("API-REST-01: Успешная генерация возвращает 200 OK и валидный JSON")
    @Description("Проверка возврата HTTP 200 OK, JSON-структуры с url и base64 qrCode")
    void shouldReturnQrCodeOnValidRequest() throws Exception {
      when(qrCodeService.generateQrCodeBase64(DEFAULT_VALID_URL)).thenReturn(MOCK_BASE64_QR);

      performPost(requestFor(DEFAULT_VALID_URL))
          .andExpect(status().isOk())
          .andExpect(content().contentType(MediaType.APPLICATION_JSON))
          .andExpect(jsonPath("$.url").value(DEFAULT_VALID_URL))
          .andExpect(jsonPath("$.qrCode").value(MOCK_BASE64_QR));

      verify(qrCodeService).generateQrCodeBase64(DEFAULT_VALID_URL);
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
    @DisplayName("API-REST-02: Поддержка различных валидных структур и протоколов URL")
    @Description(
        "Проверка генерации QR для HTTP, HTTPS, нестандартных портов, query-параметров и кириллицы")
    void shouldAcceptVariousValidUrlFormats(String validUrl) throws Exception {
      when(qrCodeService.generateQrCodeBase64(validUrl)).thenReturn(MOCK_BASE64_QR);

      performPost(requestFor(validUrl))
          .andExpect(status().isOk())
          .andExpect(jsonPath("$.url").value(validUrl))
          .andExpect(jsonPath("$.qrCode").value(MOCK_BASE64_QR));

      verify(qrCodeService).generateQrCodeBase64(validUrl);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-03: Корректная обработка длинного URL (>1000 символов)")
    @Description("Проверка передачи длинного URL без обрезания в бизнес-логику")
    void shouldHandleLongUrlSuccessfully() throws Exception {
      String longUrl = DEFAULT_VALID_URL + "/" + "a".repeat(1000);
      when(qrCodeService.generateQrCodeBase64(longUrl)).thenReturn(MOCK_BASE64_QR);

      performPost(requestFor(longUrl))
          .andExpect(status().isOk())
          .andExpect(jsonPath("$.url").value(longUrl))
          .andExpect(jsonPath("$.qrCode").value(MOCK_BASE64_QR));

      verify(qrCodeService).generateQrCodeBase64(longUrl);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-04: Игнорирование неизвестных дополнительных полей в JSON")
    @Description("Проверка устойчивости контроллера при получении нештатных дополнительных полей")
    void shouldIgnoreUnknownFieldsInJson() throws Exception {
      String payloadWithExtra =
          "{\"url\": \"" + DEFAULT_VALID_URL + "\", \"unknownField\": \"value\", \"id\": 999}";
      when(qrCodeService.generateQrCodeBase64(DEFAULT_VALID_URL)).thenReturn(MOCK_BASE64_QR);

      performPostRawJson(payloadWithExtra)
          .andExpect(status().isOk())
          .andExpect(jsonPath("$.url").value(DEFAULT_VALID_URL))
          .andExpect(jsonPath("$.qrCode").value(MOCK_BASE64_QR));

      verify(qrCodeService).generateQrCodeBase64(DEFAULT_VALID_URL);
    }
  }

  // ==========================================
  // 2. Валидация входных данных (400 Bad Request)
  // ==========================================
  @Nested
  @DisplayName("2. Валидация входных данных (400 Bad Request)")
  class ValidationScenarios {

    @ParameterizedTest
    @ValueSource(
        strings = {
          "invalid-url",
          "ftp://example.com",
          "file:///etc/passwd",
          "javascript:alert(1)",
          "www.vorobevaqa.ru",
          "   "
        })
    @Severity(SeverityLevel.CRITICAL)
    @DisplayName("API-REST-05: Невалидный формат URL возвращает 400 Bad Request")
    @Description("Проверка валидации протоколов http:// и https://")
    void shouldReturnBadRequestWhenUrlIsInvalid(String invalidUrl) throws Exception {
      performPost(requestFor(invalidUrl))
          .andExpect(status().isBadRequest())
          .andExpect(jsonPath("$.error").exists());

      verifyNoInteractions(qrCodeService);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-06: Null URL возвращает 400 Bad Request с понятным сообщением")
    @Description("Проверка срабатывания аннотации @NotBlank для null значения")
    void shouldReturnBadRequestWhenUrlIsNull() throws Exception {
      performPost(requestFor(null))
          .andExpect(status().isBadRequest())
          .andExpect(jsonPath("$.error").value("URL is required"));

      verifyNoInteractions(qrCodeService);
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-07: Пустой JSON объект {} возвращает 400 Bad Request")
    @Description("Проверка обязательности поля url при передаче пустого JSON")
    void shouldReturnBadRequestWhenEmptyJsonProvided() throws Exception {
      performPostRawJson("{}")
          .andExpect(status().isBadRequest())
          .andExpect(jsonPath("$.error").value("URL is required"));

      verifyNoInteractions(qrCodeService);
    }
  }

  // ==========================================
  // 3. Протокольные ограничения (405 & 415)
  // ==========================================
  @Nested
  @DisplayName("3. Протокольные ограничения (405 Method Not Allowed & 415 Unsupported Media Type)")
  class ProtocolScenarios {

    @ParameterizedTest
    @ValueSource(strings = {"GET", "PUT", "DELETE", "PATCH"})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-08: Неподдерживаемые HTTP методы возвращают 405 Method Not Allowed")
    @Description("Проверка запрета методов GET, PUT, DELETE, PATCH для эндпоинта /api/qr")
    void shouldReturnMethodNotAllowedForUnsupportedHttpMethods(String httpMethod) throws Exception {
      mockMvc
          .perform(
              request(HttpMethod.valueOf(httpMethod), QR_ENDPOINT)
                  .contentType(MediaType.APPLICATION_JSON)
                  .content(objectMapper.writeValueAsString(requestFor(DEFAULT_VALID_URL))))
          .andExpect(status().isMethodNotAllowed());

      verifyNoInteractions(qrCodeService);
    }

    @ParameterizedTest
    @ValueSource(strings = {MediaType.APPLICATION_XML_VALUE, MediaType.TEXT_PLAIN_VALUE})
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-09: Неподдерживаемый Content-Type возвращает 415 Unsupported Media Type")
    @Description("Проверка возврата HTTP 415 при передаче XML или plain text вместо JSON")
    void shouldReturnUnsupportedMediaType(String contentType) throws Exception {
      mockMvc
          .perform(post(QR_ENDPOINT).contentType(contentType).content(DEFAULT_VALID_URL))
          .andExpect(status().isUnsupportedMediaType());

      verifyNoInteractions(qrCodeService);
    }
  }

  // ==========================================
  // 4. Обработка серверных сбоев (500 Internal Server Error)
  // ==========================================
  @Nested
  @DisplayName("4. Обработка серверных сбоев (500 Internal Server Error)")
  class ErrorHandlingScenarios {

    @Test
    @Severity(SeverityLevel.CRITICAL)
    @DisplayName("API-REST-10: Ошибка генерации QR в сервисе возвращает 500 Internal Server Error")
    @Description(
        "Проверка перехвата RuntimeException глобальным обработчиком GlobalExceptionHandler")
    void shouldReturnInternalServerErrorWhenServiceFails() throws Exception {
      when(qrCodeService.generateQrCodeBase64(anyString()))
          .thenThrow(new RuntimeException("ZXing encoding error"));

      performPost(requestFor(DEFAULT_VALID_URL))
          .andExpect(status().isInternalServerError())
          .andExpect(jsonPath("$.error").value("Internal server error: ZXing encoding error"));
    }

    @Test
    @Severity(SeverityLevel.NORMAL)
    @DisplayName("API-REST-11: Синтаксически некорректный JSON возвращает 500")
    @Description("Проверка обработки синтаксически битого JSON через GlobalExceptionHandler")
    void shouldHandleMalformedJsonBody() throws Exception {
      String malformedJson = "{\"url\": \"https://vorobevaqa.ru\""; // незакрытая фигурная скобка

      performPostRawJson(malformedJson).andExpect(status().is5xxServerError());

      verifyNoInteractions(qrCodeService);
    }
  }
}
