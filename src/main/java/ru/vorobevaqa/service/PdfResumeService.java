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

import java.io.IOException;
import java.io.InputStream;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;

@Service
public class PdfResumeService {

  private static final Logger log = LoggerFactory.getLogger(PdfResumeService.class);
  private static final String RESUME_RESOURCE_PATH = "pdf/resume.pdf";

  public byte[] getResumePdf() {
    ClassPathResource resource = new ClassPathResource(RESUME_RESOURCE_PATH);
    try (InputStream is = resource.getInputStream()) {
      log.debug("Successfully loaded resume PDF from classpath: {}", RESUME_RESOURCE_PATH);
      return is.readAllBytes();
    } catch (IOException e) {
      log.error("Failed to load resume PDF from classpath: {}", RESUME_RESOURCE_PATH, e);
      return new byte[0];
    }
  }
}
