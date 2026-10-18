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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import java.util.Base64;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

@Epic("Core Services")
@Feature("QR Code Generator")
@Story("QrCodeService - Локальная генерация QR-кодов через ZXing")
@Tag("unit")
@Tag("service")
class QrCodeServiceTest {

  private QrCodeService qrCodeService;

  @BeforeEach
  void setUp() {
    qrCodeService = new QrCodeService();
  }

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("SRV-QR-01: Успешная генерация валидного Base64 PNG изображения")
  @Description("Проверка создания непустой строки Base64 и корректности сигнатуры PNG формата")
  void shouldGenerateValidBase64QrCode() {
    String testUrl = "https://vorobevaqa.ru";
    String base64Image = qrCodeService.generateQrCodeBase64(testUrl);

    assertNotNull(base64Image);
    assertFalse(base64Image.isBlank());

    byte[] decodedBytes = Base64.getDecoder().decode(base64Image);
    assertTrue(decodedBytes.length > 0);

    // Проверка заголовка PNG
    assertEquals((byte) 0x89, decodedBytes[0]);
    assertEquals((byte) 'P', decodedBytes[1]);
    assertEquals((byte) 'N', decodedBytes[2]);
    assertEquals((byte) 'G', decodedBytes[3]);
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-QR-02: Выброс IllegalArgumentException при пустой строке или null")
  @Description("Проверка валидации входных данных перед генерацией матрицы QR")
  void shouldThrowExceptionWhenTextIsBlankOrNull() {
    assertThrows(IllegalArgumentException.class, () -> qrCodeService.generateQrCodeBase64(null));
    assertThrows(IllegalArgumentException.class, () -> qrCodeService.generateQrCodeBase64("   "));
  }
}
