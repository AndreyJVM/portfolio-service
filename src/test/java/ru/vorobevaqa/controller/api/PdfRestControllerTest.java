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

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.header;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.nio.charset.StandardCharsets;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import ru.vorobevaqa.service.PdfResumeService;

@Epic("REST API")
@Feature("PDF Resume Generator")
@Story("PdfRestController - Spring MVC обработка запросов на скачивание и генерацию PDF резюме")
@Tag("unit")
@Tag("controller")
@WebMvcTest(PdfRestController.class)
class PdfRestControllerTest {

  @Autowired private MockMvc mockMvc;

  @MockBean private PdfResumeService pdfResumeService;

  private static final byte[] MOCK_PDF =
      "%PDF-1.4 mock pdf content %%EOF".getBytes(StandardCharsets.US_ASCII);

  @Nested
  @DisplayName("1. Успешная отдача PDF резюме (200 OK)")
  class ResumeSuccessScenarios {

    @Test
    @DisplayName("GET /api/pdf/resume без параметров возвращает inline PDF с заголовками")
    void shouldReturnInlinePdf() throws Exception {
      when(pdfResumeService.getResumePdf()).thenReturn(MOCK_PDF);

      mockMvc
          .perform(get("/api/pdf/resume"))
          .andExpect(status().isOk())
          .andExpect(content().contentType(MediaType.APPLICATION_PDF_VALUE))
          .andExpect(
              header()
                  .string(
                      HttpHeaders.CONTENT_DISPOSITION,
                      "inline; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\""))
          .andExpect(header().string(HttpHeaders.CACHE_CONTROL, "public, max-age=3600"))
          .andExpect(content().bytes(MOCK_PDF));
    }

    @Test
    @DisplayName(
        "GET /api/pdf/resume?download=true возвращает attachment PDF для принудительного"
            + " скачивания")
    void shouldReturnAttachmentPdfWhenDownloadRequested() throws Exception {
      when(pdfResumeService.getResumePdf()).thenReturn(MOCK_PDF);

      mockMvc
          .perform(get("/api/pdf/resume").param("download", "true"))
          .andExpect(status().isOk())
          .andExpect(content().contentType(MediaType.APPLICATION_PDF_VALUE))
          .andExpect(
              header()
                  .string(
                      HttpHeaders.CONTENT_DISPOSITION,
                      "attachment; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\""))
          .andExpect(content().bytes(MOCK_PDF));
    }
  }

  @Nested
  @DisplayName("2. Проверка здоровья (Health Check)")
  class HealthScenarios {

    @Test
    @DisplayName("GET /api/pdf/health возвращает статус UP и имя сервиса")
    void shouldReturnHealthStatus() throws Exception {
      mockMvc
          .perform(get("/api/pdf/health"))
          .andExpect(status().isOk())
          .andExpect(content().contentType(MediaType.APPLICATION_JSON_VALUE))
          .andExpect(jsonPath("$.status").value("UP"))
          .andExpect(jsonPath("$.service").value("pdf-service"));
    }
  }

  @Nested
  @DisplayName("3. Протокольные ограничения")
  class ProtocolScenarios {

    @Test
    @DisplayName("POST /api/pdf/resume возвращает 405 Method Not Allowed")
    void shouldReturn405OnPost() throws Exception {
      mockMvc.perform(post("/api/pdf/resume")).andExpect(status().isMethodNotAllowed());
    }
  }
}
