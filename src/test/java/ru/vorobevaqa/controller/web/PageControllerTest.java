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
package ru.vorobevaqa.controller.web;

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.model;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.view;

import io.qameta.allure.Description;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Severity;
import io.qameta.allure.SeverityLevel;
import io.qameta.allure.Story;
import java.util.List;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import ru.vorobevaqa.dto.GitHubRepoDto;
import ru.vorobevaqa.service.AnsibleGalaxyService;
import ru.vorobevaqa.service.DockerHubService;
import ru.vorobevaqa.service.GitHubService;

@Epic("Web Layer")
@Feature("Page Navigation & Routing")
@Story("PageController - Рендеринг HTML-страниц Thymeleaf")
@Tag("unit")
@Tag("controller")
@WebMvcTest(PageController.class)
class PageControllerTest {

  @Autowired private MockMvc mockMvc;

  @MockBean private GitHubService gitHubService;
  @MockBean private DockerHubService dockerHubService;
  @MockBean private AnsibleGalaxyService ansibleGalaxyService;

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("WEB-01: GET / - Рендеринг главной страницы index")
  @Description("Проверка возврата HTTP 200 и шаблона pages/index")
  void shouldReturnIndexPage() throws Exception {
    mockMvc.perform(get("/")).andExpect(status().isOk()).andExpect(view().name("pages/index"));
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("WEB-02: GET /education - Рендеринг страницы образования")
  @Description("Проверка возврата HTTP 200 и шаблона pages/education")
  void shouldReturnEducationPage() throws Exception {
    mockMvc
        .perform(get("/education"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/education"));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("WEB-03: GET /about - Проверка легаси алиаса для страницы образования")
  @Description("Проверка маппинга /about на шаблон pages/education")
  void shouldReturnAboutPage() throws Exception {
    mockMvc
        .perform(get("/about"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/education"));
  }

  @Test
  @Severity(SeverityLevel.CRITICAL)
  @DisplayName("WEB-04: GET /projects - Рендеринг страницы проектов со списком репозиториев")
  @Description("Проверка передачи модели githubRepos в представление pages/projects")
  void shouldReturnProjectsPageWithRepos() throws Exception {
    GitHubRepoDto mockRepo = new GitHubRepoDto();
    mockRepo.setName("portfolio-service");
    mockRepo.setDescription("Portfolio service with Spring Boot");
    mockRepo.setHtmlUrl("https://github.com/AndreyJVM/portfolio-service");
    mockRepo.setLanguage("Java");
    mockRepo.setStars(5);
    mockRepo.setForks(2);

    when(gitHubService.getRecentRepositories()).thenReturn(List.of(mockRepo));

    mockMvc
        .perform(get("/projects"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/projects"))
        .andExpect(model().attributeExists("githubRepos"))
        .andExpect(model().attribute("githubRepos", List.of(mockRepo)));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("WEB-05: GET /projects - Рендеринг страницы проектов при пустом списке репозиториев")
  @Description("Проверка отображения страницы при отсутствии репозиториев")
  void shouldReturnProjectsPageWithEmptyRepos() throws Exception {
    when(gitHubService.getRecentRepositories()).thenReturn(List.of());

    mockMvc
        .perform(get("/projects"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/projects"))
        .andExpect(model().attributeExists("githubRepos"))
        .andExpect(model().attribute("githubRepos", List.of()));
  }

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("WEB-06: GET /qr - Рендеринг страницы QR Playground")
  @Description("Проверка возврата HTTP 200 и шаблона pages/qr")
  void shouldReturnQrPage() throws Exception {
    mockMvc.perform(get("/qr")).andExpect(status().isOk()).andExpect(view().name("pages/qr"));
  }
}
