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
import ru.vorobevaqa.dto.AnsibleRoleDto;

@Epic("Core Services")
@Feature("External Integrations")
@Story("AnsibleGalaxyService - Получение списка ролей Ansible")
@Tag("unit")
@Tag("service")
class AnsibleGalaxyServiceTest {

  private AnsibleGalaxyService ansibleGalaxyService;
  private MockRestServiceServer mockServer;

  @BeforeEach
  void setUp() {
    RestClient.Builder builder = RestClient.builder().baseUrl("https://galaxy.ansible.com/api/v1");
    mockServer = MockRestServiceServer.bindTo(builder).build();
    ansibleGalaxyService = new AnsibleGalaxyService(builder.build(), "testuser");
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-ANS-01: Успешный парсинг списка ролей из ответа Ansible Galaxy API")
  @Description("Проверка десериализации JSON в список AnsibleRoleDto при статусе 200 OK")
  void shouldReturnRolesOnSuccessfulApiResponse() {
    String mockJson =
        """
        {
          "count": 1,
          "results": [
            {
              "id": 101,
              "name": "docker-role",
              "description": "Ansible role for Docker installation",
              "download_count": 42
            }
          ]
        }
        """;

    mockServer
        .expect(requestTo("https://galaxy.ansible.com/api/v1/roles/?github_user=testuser"))
        .andRespond(withSuccess(mockJson, MediaType.APPLICATION_JSON));

    List<AnsibleRoleDto> roles = ansibleGalaxyService.getRecentRoles();

    mockServer.verify();
    assertNotNull(roles);
    assertFalse(roles.isEmpty());
    assertEquals(1, roles.size());
    assertEquals("docker-role", roles.getFirst().getName());
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("SRV-ANS-02: Корректная обработка пустого списка результатов")
  @Description("Проверка возврата пустого списка, если у пользователя нет ролей в results")
  void shouldReturnEmptyListWhenResultsEmpty() {
    String mockJson =
        """
        {
          "count": 0,
          "results": []
        }
        """;

    mockServer
        .expect(requestTo("https://galaxy.ansible.com/api/v1/roles/?github_user=testuser"))
        .andRespond(withSuccess(mockJson, MediaType.APPLICATION_JSON));

    List<AnsibleRoleDto> roles = ansibleGalaxyService.getRecentRoles();

    mockServer.verify();
    assertNotNull(roles);
    assertTrue(roles.isEmpty());
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("SRV-ANS-03: Graceful fallback на пустой список при ошибке 500 от внешнего API")
  @Description("Проверка отсутствия необработанных исключений при HTTP 5xx от Ansible Galaxy")
  void shouldReturnEmptyListOnServerError() {
    mockServer
        .expect(requestTo("https://galaxy.ansible.com/api/v1/roles/?github_user=testuser"))
        .andRespond(withServerError());

    List<AnsibleRoleDto> roles = ansibleGalaxyService.getRecentRoles();

    mockServer.verify();
    assertNotNull(roles);
    assertTrue(roles.isEmpty());
  }
}
