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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import ru.vorobevaqa.config.CacheConfig;
import ru.vorobevaqa.config.HttpConstants;
import ru.vorobevaqa.dto.DockerHubRepoDto;
import ru.vorobevaqa.dto.DockerHubResponseDto;

@Service
public class DockerHubService {

  private static final Logger log = LoggerFactory.getLogger(DockerHubService.class);
  private static final Duration CONNECT_TIMEOUT = Duration.ofSeconds(2);
  private static final Duration READ_TIMEOUT = Duration.ofSeconds(3);

  private final RestClient restClient;
  private final String dockerhubUsername;

  @Autowired
  public DockerHubService(
      @Value("${dockerhub.api-url:https://hub.docker.com/v2}") String apiUrl,
      @Value("${dockerhub.username:andreyvorobevaqa}") String dockerhubUsername) {
    this(defaultRestClient(apiUrl), dockerhubUsername);
  }

  public DockerHubService(RestClient restClient, String dockerhubUsername) {
    this.restClient = restClient;
    this.dockerhubUsername = dockerhubUsername;
  }

  private static RestClient defaultRestClient(String baseUrl) {
    SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
    requestFactory.setConnectTimeout(CONNECT_TIMEOUT);
    requestFactory.setReadTimeout(READ_TIMEOUT);

    return RestClient.builder()
        .requestFactory(requestFactory)
        .baseUrl(baseUrl)
        .defaultHeader(HttpHeaders.USER_AGENT, HttpConstants.DEFAULT_USER_AGENT)
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

      if (response == null || response.results() == null) {
        return Collections.emptyList();
      }

      return response.results();
    } catch (Exception e) {
      log.error("Failed to fetch repositories from DockerHub API: {}", e.getMessage());
      return Collections.emptyList();
    }
  }
}
