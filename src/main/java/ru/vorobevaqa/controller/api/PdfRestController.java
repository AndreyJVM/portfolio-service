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

import java.util.Map;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import ru.vorobevaqa.service.PdfResumeService;

@Slf4j
@RestController
@RequestMapping("/api/pdf")
@RequiredArgsConstructor
public class PdfRestController {

  private final PdfResumeService pdfResumeService;

  @GetMapping(value = "/resume", produces = MediaType.APPLICATION_PDF_VALUE)
  public ResponseEntity<byte[]> getResume(
      @RequestParam(name = "download", defaultValue = "false") boolean download) {
    log.info("Request to generate/download resume PDF (download={})", download);
    byte[] pdfBytes = pdfResumeService.getResumePdf();
    String disposition = download ? "attachment" : "inline";

    return ResponseEntity.ok()
        .header(
            HttpHeaders.CONTENT_DISPOSITION,
            disposition + "; filename=\"Andrey_Vorobev_QA_DevOps_Resume.pdf\"")
        .header(HttpHeaders.CACHE_CONTROL, "public, max-age=3600")
        .body(pdfBytes);
  }

  @GetMapping(value = "/health", produces = MediaType.APPLICATION_JSON_VALUE)
  public ResponseEntity<Map<String, String>> health() {
    return ResponseEntity.ok(Map.of("status", "UP", "service", "pdf-service"));
  }
}
