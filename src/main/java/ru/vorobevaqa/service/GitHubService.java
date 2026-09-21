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

import java.time.Duration;
import java.util.Collections;
import java.util.List;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import ru.vorobevaqa.config.CacheConfig;
import ru.vorobevaqa.dto.GitHubRepoDto;

@Slf4j
@Service
public class GitHubService {

  private final RestClient restClient;
  private final String githubUsername;

  public GitHubService(@Value("${github.username:AndreyJVM}") String githubUsername) {
    this.githubUsername = githubUsername;

    SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
    requestFactory.setConnectTimeout(Duration.ofSeconds(2));
    requestFactory.setReadTimeout(Duration.ofSeconds(3));

    this.restClient =
        RestClient.builder()
            .requestFactory(requestFactory)
            .baseUrl("https://api.github.com")
            .defaultHeader(HttpHeaders.USER_AGENT, "Spring-Boot-Portfolio-App")
            .defaultHeader(HttpHeaders.ACCEPT, "application/vnd.github.v3+json")
            .build();
  }

  @Cacheable(value = CacheConfig.GITHUB_REPOS_CACHE, unless = "#result.isEmpty()")
  public List<GitHubRepoDto> getRecentRepositories() {
    try {
      log.info("Fetching recent repositories from GitHub API for user: {}", githubUsername);
      List<GitHubRepoDto> repos =
          restClient
              .get()
              .uri("/users/{username}/repos?sort=updated&per_page=6&type=owner", githubUsername)
              .retrieve()
              .body(new ParameterizedTypeReference<List<GitHubRepoDto>>() {});

      if (repos == null) {
        return Collections.emptyList();
      }

      return repos.stream().filter(repo -> !repo.isFork()).limit(6).toList();
    } catch (Exception e) {
      log.error("Failed to fetch repositories from GitHub API: {}", e.getMessage());
      return Collections.emptyList();
    }
  }
}
