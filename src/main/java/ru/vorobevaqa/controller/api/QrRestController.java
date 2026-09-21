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

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ru.vorobevaqa.dto.QrRequest;
import ru.vorobevaqa.dto.QrResponse;
import ru.vorobevaqa.service.QrCodeService;

@Slf4j
@RestController
@RequestMapping("/api/qr")
@RequiredArgsConstructor
public class QrRestController {

  private final QrCodeService qrCodeService;

  @PostMapping
  public ResponseEntity<QrResponse> generateQrCode(@Valid @RequestBody QrRequest request) {
    log.info("Request to generate QR code for URL: {}", request.getUrl());
    String qrCodeBase64 = qrCodeService.generateQrCodeBase64(request.getUrl());
    return ResponseEntity.ok(new QrResponse(request.getUrl(), qrCodeBase64));
  }
}
