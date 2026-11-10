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

import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import ru.vorobevaqa.dto.AnsibleRoleDto;
import ru.vorobevaqa.dto.DockerHubRepoDto;
import ru.vorobevaqa.dto.GitHubRepoDto;
import ru.vorobevaqa.dto.ProjectsPageData;

/**
 * Фасадный сервис для параллельной асинхронной агрегации витрины проектов.
 */
@Service
public class ProjectShowcaseService {

  private static final Logger log = LoggerFactory.getLogger(ProjectShowcaseService.class);

  private final GitHubService gitHubService;
  private final DockerHubService dockerHubService;
  private final AnsibleGalaxyService ansibleGalaxyService;

  public ProjectShowcaseService(
      GitHubService gitHubService,
      DockerHubService dockerHubService,
      AnsibleGalaxyService ansibleGalaxyService) {
    this.gitHubService = gitHubService;
    this.dockerHubService = dockerHubService;
    this.ansibleGalaxyService = ansibleGalaxyService;
  }

  /**
   * Параллельно запрашивает данные из GitHub, DockerHub и Ansible Galaxy.
   *
   * @return агрегированные данные для витрины проектов
   */
  public ProjectsPageData getProjectsData() {
    CompletableFuture<List<GitHubRepoDto>> githubFuture =
        CompletableFuture.supplyAsync(gitHubService::getRecentRepositories)
            .exceptionally(ex -> {
              log.error("Не удалось получить репозитории GitHub: {}", ex.getMessage());
              return Collections.emptyList();
            });

    CompletableFuture<List<DockerHubRepoDto>> dockerFuture =
        CompletableFuture.supplyAsync(dockerHubService::getRecentRepositories)
            .exceptionally(ex -> {
              log.error("Не удалось получить репозитории DockerHub: {}", ex.getMessage());
              return Collections.emptyList();
            });

    CompletableFuture<List<AnsibleRoleDto>> ansibleFuture =
        CompletableFuture.supplyAsync(ansibleGalaxyService::getRecentRoles)
            .exceptionally(ex -> {
              log.error("Не удалось получить роли Ansible Galaxy: {}", ex.getMessage());
              return Collections.emptyList();
            });

    CompletableFuture.allOf(githubFuture, dockerFuture, ansibleFuture).join();

    return new ProjectsPageData(
        githubFuture.join(),
        dockerFuture.join(),
        ansibleFuture.join()
    );
  }
}
