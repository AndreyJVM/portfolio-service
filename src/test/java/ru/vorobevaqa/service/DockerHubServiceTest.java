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
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withServerError;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import java.util.List;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import org.springframework.web.client.RestClient;
import ru.vorobevaqa.dto.DockerHubRepoDto;

@Epic("Core Services")
@Feature("External Integrations")
@Story("DockerHubService - Получение списка репозиториев DockerHub")
@Tag("unit")
@Tag("service")
class DockerHubServiceTest {

  private DockerHubService dockerHubService;
  private MockRestServiceServer mockServer;

  @BeforeEach
  void setUp() {
    RestClient.Builder builder = RestClient.builder().baseUrl("https://hub.docker.com/v2");
    mockServer = MockRestServiceServer.bindTo(builder).build();
    dockerHubService = new DockerHubService(builder.build(), "testuser");
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-DOC-01: Успешный парсинг списка репозиториев из DockerHub API")
  @Description("Проверка десериализации JSON в список DockerHubRepoDto при статусе 200 OK")
  void shouldReturnRepositoriesOnSuccessfulApiResponse() {
    String mockJson =
        """
        {
          "count": 1,
          "results": [
            {
              "name": "portfolio-service",
              "description": "Portfolio and QR service docker image",
              "pull_count": 150,
              "star_count": 5
            }
          ]
        }
        """;

    mockServer
        .expect(
            requestTo(
                "https://hub.docker.com/v2/repositories/testuser/?page_size=10&ordering=-last_updated"))
        .andRespond(withSuccess(mockJson, MediaType.APPLICATION_JSON));

    List<DockerHubRepoDto> repos = dockerHubService.getRecentRepositories();

    mockServer.verify();
    assertNotNull(repos);
    assertFalse(repos.isEmpty());
    assertEquals(1, repos.size());
    assertEquals("portfolio-service", repos.getFirst().getName());
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-DOC-02: Graceful fallback на пустой список при ошибке 500 от DockerHub API")
  @Description("Проверка отсутствия исключений при HTTP 5xx от серверов DockerHub")
  void shouldReturnEmptyListOnServerError() {
    mockServer
        .expect(
            requestTo(
                "https://hub.docker.com/v2/repositories/testuser/?page_size=10&ordering=-last_updated"))
        .andRespond(withServerError());

    List<DockerHubRepoDto> repos = dockerHubService.getRecentRepositories();

    mockServer.verify();
    assertNotNull(repos);
    assertTrue(repos.isEmpty());
  }
}
