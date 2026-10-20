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
import ru.vorobevaqa.dto.AnsibleRoleDto;
import ru.vorobevaqa.dto.DockerHubRepoDto;
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
  @DisplayName("WEB-04: GET /projects - Рендеринг страницы проектов со всеми интеграциями")
  @Description(
      "Проверка передачи моделей githubRepos, dockerRepos и ansibleRoles в представление"
          + " pages/projects")
  void shouldReturnProjectsPageWithRepos() throws Exception {
    GitHubRepoDto mockGitHub = new GitHubRepoDto();
    mockGitHub.setName("portfolio-service");
    mockGitHub.setDescription("Portfolio service with Spring Boot");
    mockGitHub.setHtmlUrl("https://github.com/AndreyJVM/portfolio-service");
    mockGitHub.setLanguage("Java");
    mockGitHub.setStars(5);
    mockGitHub.setForks(2);

    DockerHubRepoDto mockDocker =
        DockerHubRepoDto.builder()
            .name("portfolio-service")
            .namespace("andreyjvm")
            .description("Docker image for portfolio")
            .pullCount(100)
            .starCount(3)
            .build();

    AnsibleRoleDto mockAnsible = new AnsibleRoleDto();
    mockAnsible.setName("security-baseline");
    mockAnsible.setDescription("Ansible hardening role");
    mockAnsible.setGithubUser("AndreyJVM");
    mockAnsible.setDownloadCount(50);

    when(gitHubService.getRecentRepositories()).thenReturn(List.of(mockGitHub));
    when(dockerHubService.getRecentRepositories()).thenReturn(List.of(mockDocker));
    when(ansibleGalaxyService.getRecentRoles()).thenReturn(List.of(mockAnsible));

    mockMvc
        .perform(get("/projects"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/projects"))
        .andExpect(model().attributeExists("githubRepos"))
        .andExpect(model().attribute("githubRepos", List.of(mockGitHub)))
        .andExpect(model().attributeExists("dockerRepos"))
        .andExpect(model().attribute("dockerRepos", List.of(mockDocker)))
        .andExpect(model().attributeExists("ansibleRoles"))
        .andExpect(model().attribute("ansibleRoles", List.of(mockAnsible)));
  }

  @Test
  @Severity(SeverityLevel.NORMAL)
  @DisplayName("WEB-05: GET /projects - Рендеринг страницы проектов при пустых ответах интеграций")
  @Description(
      "Проверка корректного рендеринга страницы при пустых списках проектов из внешних API")
  void shouldReturnProjectsPageWithEmptyRepos() throws Exception {
    when(gitHubService.getRecentRepositories()).thenReturn(List.of());
    when(dockerHubService.getRecentRepositories()).thenReturn(List.of());
    when(ansibleGalaxyService.getRecentRoles()).thenReturn(List.of());

    mockMvc
        .perform(get("/projects"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/projects"))
        .andExpect(model().attributeExists("githubRepos"))
        .andExpect(model().attribute("githubRepos", List.of()))
        .andExpect(model().attributeExists("dockerRepos"))
        .andExpect(model().attribute("dockerRepos", List.of()))
        .andExpect(model().attributeExists("ansibleRoles"))
        .andExpect(model().attribute("ansibleRoles", List.of()));
  }

  @Test
  @Severity(SeverityLevel.BLOCKER)
  @DisplayName("WEB-06: GET /qr - Рендеринг страницы QR Playground")
  @Description("Проверка возврата HTTP 200 и шаблона pages/qr")
  void shouldReturnQrPage() throws Exception {
    mockMvc.perform(get("/qr")).andExpect(status().isOk()).andExpect(view().name("pages/qr"));
  }
}
