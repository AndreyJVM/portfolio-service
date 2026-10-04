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
import ru.vorobevaqa.dto.DockerHubRepoDto;
import ru.vorobevaqa.dto.DockerHubResponseDto;

@Slf4j
@Service
public class DockerHubService {

  private final RestClient restClient;
  private final String dockerhubUsername;

  public DockerHubService(
      @Value("${dockerhub.username:andreyvorobevaqa}") String dockerhubUsername) {
    this.dockerhubUsername = dockerhubUsername;

    SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
    requestFactory.setConnectTimeout(Duration.ofSeconds(2));
    requestFactory.setReadTimeout(Duration.ofSeconds(3));

    this.restClient =
        RestClient.builder()
            .requestFactory(requestFactory)
            .baseUrl("https://hub.docker.com/v2")
            .defaultHeader(HttpHeaders.USER_AGENT, "Spring-Boot-Portfolio-App")
            .build();
  }

  @Cacheable(value = CacheConfig.DOCKER_REPOS_CACHE, unless = "#result.isEmpty()")
  public List<DockerHubRepoDto> getRecentRepositories() {
    try {
      log.info(
          "Fetching recent repositories from DockerHub API for namespace: {}", dockerhubUsername);
      DockerHubResponseDto response =
          restClient
              .get()
              .uri(
                  "/repositories/{username}/?page_size=10&ordering=-last_updated",
                  dockerhubUsername)
              .retrieve()
              .body(DockerHubResponseDto.class);

      if (response == null || response.getResults() == null) {
        return Collections.emptyList();
      }

      return response.getResults();
    } catch (Exception e) {
      log.error("Failed to fetch repositories from DockerHub API: {}", e.getMessage());
      return Collections.emptyList();
    }
  }
}
