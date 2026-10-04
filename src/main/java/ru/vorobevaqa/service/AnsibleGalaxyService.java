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
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import ru.vorobevaqa.config.CacheConfig;
import ru.vorobevaqa.dto.AnsibleGalaxyResponseDto;
import ru.vorobevaqa.dto.AnsibleRoleDto;

@Slf4j
@Service
public class AnsibleGalaxyService {

  private final RestClient restClient;
  private final String githubUsername;

  public AnsibleGalaxyService(@Value("${github.username:AndreyJVM}") String githubUsername) {
    this.githubUsername = githubUsername;

    SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
    requestFactory.setConnectTimeout(Duration.ofSeconds(5));
    requestFactory.setReadTimeout(Duration.ofSeconds(10));

    this.restClient =
        RestClient.builder()
            .requestFactory(requestFactory)
            .baseUrl("https://galaxy.ansible.com/api/v1")
            .defaultHeader(HttpHeaders.USER_AGENT, "Spring-Boot-Portfolio-App")
            .build();
  }

  @Cacheable(value = CacheConfig.ANSIBLE_REPOS_CACHE, unless = "#result.isEmpty()")
  public List<AnsibleRoleDto> getRecentRoles() {
    try {
      log.info("Fetching recent roles from Ansible Galaxy API for user: {}", githubUsername);
      AnsibleGalaxyResponseDto response =
          restClient
              .get()
              .uri("/roles/?github_user={username}", githubUsername)
              .retrieve()
              .body(AnsibleGalaxyResponseDto.class);

      if (response == null || response.getResults() == null) {
        return Collections.emptyList();
      }

      return response.getResults();
    } catch (Exception e) {
      log.error("Failed to fetch roles from Ansible Galaxy API: {}", e.getMessage());
      return Collections.emptyList();
    }
  }
}
