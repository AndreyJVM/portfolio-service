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

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class PdfResumeService {

  private static final Logger log = LoggerFactory.getLogger(PdfResumeService.class);

  private final String pdfServiceUrl;
  private final HttpClient httpClient;

  public PdfResumeService(@Value("${pdf.service.url:http://localhost:8082}") String pdfServiceUrl) {
    this.pdfServiceUrl = pdfServiceUrl;
    this.httpClient = HttpClient.newBuilder().connectTimeout(Duration.ofMillis(800)).build();
  }

  public byte[] getResumePdf() {
    try {
      HttpRequest request =
          HttpRequest.newBuilder()
              .uri(URI.create(pdfServiceUrl + "/api/pdf/resume"))
              .timeout(Duration.ofMillis(1500))
              .GET()
              .build();

      HttpResponse<byte[]> response =
          httpClient.send(request, HttpResponse.BodyHandlers.ofByteArray());
      if (response.statusCode() == 200 && response.body() != null && response.body().length > 0) {
        log.debug("Successfully fetched resume PDF from Go microservice ({})", pdfServiceUrl);
        return response.body();
      }
    } catch (Exception e) {
      log.debug(
          "Go PDF microservice unavailable at {}, using embedded fallback: {}",
          pdfServiceUrl,
          e.getMessage());
    }

    return generateFallbackPdf();
  }

  public byte[] generateFallbackPdf() {
    String content =
        "BT\n"
            + "/F1 20 Tf\n"
            + "40 790 Td\n"
            + "(Andrey Vorobev) Tj\n"
            + "0 -24 Td\n"
            + "/F1 12 Tf\n"
            + "(DevOps & Linux Infrastructure Engineer) Tj\n"
            + "0 -20 Td\n"
            + "/F2 9 Tf\n"
            + "(Website: https://vorobevaqa.ru | Telegram: @AndreyAQA | GitHub:"
            + " github.com/AndreyJVM) Tj\n"
            + "0 -35 Td\n"
            + "/F1 11 Tf\n"
            + "(Core Skills: Linux, Bash, Docker, Docker Compose, Nginx, Ansible, CI/CD, Java) Tj\n"
            + "ET\n";

    byte[] contentBytes = content.getBytes(StandardCharsets.US_ASCII);
    int streamLen = contentBytes.length;

    String header =
        "%PDF-1.4\n"
            + "1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj\n"
            + "2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj\n"
            + "3 0 obj <</Type /Page /Parent 2 0 R /Resources 4 0 R /MediaBox [0 0 595 842]"
            + " /Contents 5 0 R>> endobj\n"
            + "4 0 obj <</Font <</F1 <</Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold>> /F2"
            + " <</Type /Font /Subtype /Type1 /BaseFont /Helvetica>>>>>> endobj\n"
            + "5 0 obj <</Length "
            + streamLen
            + ">> stream\n";

    String footer =
        "endstream\n"
            + "endobj\n"
            + "xref\n"
            + "0 6\n"
            + "0000000000 65535 f \n"
            + "0000000009 00000 n \n"
            + "0000000058 00000 n \n"
            + "0000000115 00000 n \n"
            + "0000000219 00000 n \n"
            + "0000000350 00000 n \n"
            + "trailer <</Size 6 /Root 1 0 R>>\n"
            + "startxref\n"
            + "750\n"
            + "%%EOF\n";

    return (header + content + footer).getBytes(StandardCharsets.US_ASCII);
  }
}
