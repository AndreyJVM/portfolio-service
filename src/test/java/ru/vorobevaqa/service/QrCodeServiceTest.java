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

import static org.junit.jupiter.api.Assertions.*;

import java.util.Base64;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

class QrCodeServiceTest {

  private QrCodeService qrCodeService;

  @BeforeEach
  void setUp() {
    qrCodeService = new QrCodeService();
  }

  @Test
  @DisplayName("Успешная генерация валидного Base64 PNG изображения")
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
  @DisplayName("Выброс IllegalArgumentException при пустой строке или null")
  void shouldThrowExceptionWhenTextIsBlankOrNull() {
    assertThrows(IllegalArgumentException.class, () -> qrCodeService.generateQrCodeBase64(null));
    assertThrows(IllegalArgumentException.class, () -> qrCodeService.generateQrCodeBase64("   "));
  }
}
