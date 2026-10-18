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
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
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
  }
}
