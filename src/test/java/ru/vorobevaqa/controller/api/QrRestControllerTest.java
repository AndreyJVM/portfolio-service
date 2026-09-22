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
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.DisplayName;
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

@WebMvcTest(QrRestController.class)
@Import(GlobalExceptionHandler.class)
class QrRestControllerTest {

  @Autowired private MockMvc mockMvc;

  @Autowired private ObjectMapper objectMapper;

  @MockBean private QrCodeService qrCodeService;

  @Test
  @DisplayName("POST /api/qr - успешная генерация возвращает 200 и валидный JSON")
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
  @DisplayName("POST /api/qr - ошибка валидации (невалидный формат URL) возвращает 400 Bad Request")
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
  @DisplayName("POST /api/qr - пустой или null URL возвращает 400 Bad Request")
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
