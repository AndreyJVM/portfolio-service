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

import com.google.zxing.BarcodeFormat;
import com.google.zxing.EncodeHintType;
import com.google.zxing.WriterException;
import com.google.zxing.client.j2se.MatrixToImageWriter;
import com.google.zxing.common.BitMatrix;
import com.google.zxing.qrcode.QRCodeWriter;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.EnumMap;
import java.util.Map;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Slf4j
@Service
public class QrCodeService {

  private static final int DEFAULT_SIZE = 300;

  public String generateQrCodeBase64(String text) {
    return generateQrCodeBase64(text, DEFAULT_SIZE, DEFAULT_SIZE);
  }

  public String generateQrCodeBase64(String text, int width, int height) {
    if (text == null || text.isBlank()) {
      throw new IllegalArgumentException("Text or URL for QR code cannot be empty");
    }

    try {
      QRCodeWriter qrCodeWriter = new QRCodeWriter();

      // Настройки генерации: поддержка кириллицы и аккуратные границы
      Map<EncodeHintType, Object> hints = new EnumMap<>(EncodeHintType.class);
      hints.put(EncodeHintType.CHARACTER_SET, StandardCharsets.UTF_8.name());
      hints.put(EncodeHintType.MARGIN, 1);

      BitMatrix bitMatrix = qrCodeWriter.encode(text, BarcodeFormat.QR_CODE, width, height, hints);

      ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
      MatrixToImageWriter.writeToStream(bitMatrix, "PNG", outputStream);

      return Base64.getEncoder().encodeToString(outputStream.toByteArray());
    } catch (WriterException | IOException e) {
      log.error("Failed to generate QR code for input: {}", text, e);
      throw new RuntimeException("Failed to generate QR code", e);
    }
  }
}
