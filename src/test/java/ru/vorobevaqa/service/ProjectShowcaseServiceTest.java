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

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.when;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.util.List;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import ru.vorobevaqa.dto.AnsibleRoleDto;
import ru.vorobevaqa.dto.DockerHubRepoDto;
import ru.vorobevaqa.dto.GitHubRepoDto;
import ru.vorobevaqa.dto.ProjectsPageData;

@Epic("Service Layer")
@Feature("Projects Aggregation")
@Story("Параллельная загрузка данных для витрины проектов")
@Tag("unit")
@Tag("service")
@ExtendWith(MockitoExtension.class)
class ProjectShowcaseServiceTest {

  @Mock private GitHubService gitHubService;
  @Mock private DockerHubService dockerHubService;
  @Mock private AnsibleGalaxyService ansibleGalaxyService;

  @InjectMocks private ProjectShowcaseService projectShowcaseService;

  @Test
  @DisplayName("SVC-01: Успешная параллельная агрегация всех интеграций")
  @Description("Проверка объединения списков репозиториев и ролей в единый объект ProjectsPageData")
  void shouldAggregateAllProjectSourcesSuccessfully() {
    GitHubRepoDto gitHubDto = new GitHubRepoDto();
    gitHubDto.setName("test-repo");

    DockerHubRepoDto dockerDto = DockerHubRepoDto.builder().name("test-image").build();

    AnsibleRoleDto ansibleDto = new AnsibleRoleDto();
    ansibleDto.setName("test-role");

    when(gitHubService.getRecentRepositories()).thenReturn(List.of(gitHubDto));
    when(dockerHubService.getRecentRepositories()).thenReturn(List.of(dockerDto));
    when(ansibleGalaxyService.getRecentRoles()).thenReturn(List.of(ansibleDto));

    ProjectsPageData result = projectShowcaseService.getProjectsData();

    assertThat(result).isNotNull();
    assertThat(result.githubRepos()).containsExactly(gitHubDto);
    assertThat(result.dockerRepos()).containsExactly(dockerDto);
    assertThat(result.ansibleRoles()).containsExactly(ansibleDto);
  }

  @Test
  @DisplayName("SVC-02: Отказоустойчивость при сбое одного из источников")
  @Description("При возникновении исключения в одном сервисе остальные данные должны успешно вернуться")
  void shouldHandleExceptionsGracefully() {
    GitHubRepoDto gitHubDto = new GitHubRepoDto();
    gitHubDto.setName("test-repo");

    when(gitHubService.getRecentRepositories()).thenReturn(List.of(gitHubDto));
    when(dockerHubService.getRecentRepositories()).thenThrow(new RuntimeException("DockerHub Timeout"));
    when(ansibleGalaxyService.getRecentRoles()).thenReturn(List.of());

    ProjectsPageData result = projectShowcaseService.getProjectsData();

    assertThat(result).isNotNull();
    assertThat(result.githubRepos()).containsExactly(gitHubDto);
    assertThat(result.dockerRepos()).isEmpty();
    assertThat(result.ansibleRoles()).isEmpty();
  }
}
