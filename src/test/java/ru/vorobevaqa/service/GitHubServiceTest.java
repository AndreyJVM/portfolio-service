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
import ru.vorobevaqa.dto.GitHubRepoDto;

@Epic("Core Services")
@Feature("External Integrations")
@Story("GitHubService - Получение списка репозиториев GitHub")
@Tag("unit")
@Tag("service")
class GitHubServiceTest {

  private GitHubService gitHubService;
  private MockRestServiceServer mockServer;

  @BeforeEach
  void setUp() {
    RestClient.Builder builder = RestClient.builder().baseUrl("https://api.github.com");
    mockServer = MockRestServiceServer.bindTo(builder).build();
    gitHubService = new GitHubService(builder.build(), "testuser");
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-GIT-01: Успешный парсинг списка репозиториев из GitHub API (без форков)")
  @Description("Проверка фильтрации форков и десериализации JSON в список GitHubRepoDto")
  void shouldReturnNonForkRepositoriesOnSuccessfulApiResponse() {
    String mockJson =
        """
        [
          {
            "name": "portfolio-service",
            "description": "My Portfolio Service",
            "html_url": "https://github.com/testuser/portfolio-service",
            "language": "Java",
            "stargazers_count": 10,
            "forks_count": 2,
            "fork": false
          },
          {
            "name": "forked-repo",
            "description": "Some forked project",
            "html_url": "https://github.com/testuser/forked-repo",
            "language": "Go",
            "stargazers_count": 0,
            "forks_count": 0,
            "fork": true
          }
        ]
        """;

    mockServer
        .expect(
            requestTo(
                "https://api.github.com/users/testuser/repos?sort=updated&per_page=6&type=owner"))
        .andRespond(withSuccess(mockJson, MediaType.APPLICATION_JSON));

    List<GitHubRepoDto> repos = gitHubService.getRecentRepositories();

    mockServer.verify();
    assertNotNull(repos);
    assertFalse(repos.isEmpty());
    assertEquals(1, repos.size());
    assertEquals("portfolio-service", repos.getFirst().getName());
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-GIT-02: Graceful fallback на пустой список при ошибке 500 от GitHub API")
  @Description("Проверка отсутствия исключений при падении GitHub API")
  void shouldReturnEmptyListOnServerError() {
    mockServer
        .expect(
            requestTo(
                "https://api.github.com/users/testuser/repos?sort=updated&per_page=6&type=owner"))
        .andRespond(withServerError());

    List<GitHubRepoDto> repos = gitHubService.getRecentRepositories();

    mockServer.verify();
    assertNotNull(repos);
    assertTrue(repos.isEmpty());
  }
}
