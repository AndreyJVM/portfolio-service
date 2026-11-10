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
package ru.vorobevaqa.dto;

import java.util.Collections;
import java.util.List;

/**
 * Агрегированные данные витрины проектов для страницы /projects.
 */
public record ProjectsPageData(
    List<GitHubRepoDto> githubRepos,
    List<DockerHubRepoDto> dockerRepos,
    List<AnsibleRoleDto> ansibleRoles) {

  public ProjectsPageData {
    githubRepos = githubRepos != null ? Collections.unmodifiableList(githubRepos) : List.of();
    dockerRepos = dockerRepos != null ? Collections.unmodifiableList(dockerRepos) : List.of();
    ansibleRoles = ansibleRoles != null ? Collections.unmodifiableList(ansibleRoles) : List.of();
  }

  public static ProjectsPageData empty() {
    return new ProjectsPageData(List.of(), List.of(), List.of());
  }
}
