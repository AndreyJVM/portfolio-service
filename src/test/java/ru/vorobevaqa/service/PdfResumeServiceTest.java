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
package ru.vorobevaqa.service;

import static org.assertj.core.api.Assertions.assertThat;

import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.nio.charset.StandardCharsets;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

@Epic("Core Services")
@Feature("PDF Resume Generator")
@Story("PdfResumeService - Генерация и фалбэк PDF резюме")
@Tag("unit")
@Tag("service")
class PdfResumeServiceTest {

  private PdfResumeService pdfResumeService;

  @BeforeEach
  void setUp() {
    // Port 65530 simulates an unreachable Go service
    pdfResumeService = new PdfResumeService("http://127.0.0.1:65530");
  }

  @Test
  @DisplayName("generateFallbackPdf возвращает валидный PDF с магическими байтами %PDF-1.4")
  void shouldGenerateFallbackPdfWithCorrectMagicBytes() {
    byte[] pdf = pdfResumeService.generateFallbackPdf();

    assertThat(pdf).isNotNull();
    assertThat(pdf.length).isGreaterThan(200);

    String pdfHeader = new String(pdf, 0, 8, StandardCharsets.US_ASCII);
    assertThat(pdfHeader).startsWith("%PDF-1.");

    String pdfString = new String(pdf, StandardCharsets.US_ASCII);
    assertThat(pdfString).contains("Andrey Vorobev").contains("%%EOF");
  }

  @Test
  @DisplayName("getResumePdf возвращает fallback PDF при недоступности внешнего сервиса")
  void shouldReturnFallbackWhenServiceUnavailable() {
    byte[] pdf = pdfResumeService.getResumePdf();

    assertThat(pdf).isNotNull();
    assertThat(pdf.length).isGreaterThan(200);
    assertThat(new String(pdf, 0, 8, StandardCharsets.US_ASCII)).startsWith("%PDF-1.");
  }
}
