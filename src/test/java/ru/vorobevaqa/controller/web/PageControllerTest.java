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
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

import java.util.List;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import ru.vorobevaqa.dto.SkillDto;
import ru.vorobevaqa.service.GitHubService;
import ru.vorobevaqa.service.SkillService;

@WebMvcTest(PageController.class)
class PageControllerTest {

  @Autowired private MockMvc mockMvc;

  @MockBean private SkillService skillService;

  @MockBean private GitHubService gitHubService;

  @Test
  @DisplayName("GET /about - should return about view with skills and categories in model")
  void shouldReturnAboutPageWithSkills() throws Exception {
    // Подготавливаем тестовые данные
    SkillDto mockSkill =
        SkillDto.builder()
            .id(1L)
            .name("CI/CD пайплайны")
            .categoryKey("DEVOPS_CICD")
            .categoryDisplayName("DevOps & CI/CD")
            .levelTitle("Уверенный")
            .levelPercentage(75)
            .badgeColor("bg-primary")
            .description("Тестовое описание")
            .tools(List.of("Docker", "GitHub Actions"))
            .build();

    when(skillService.getAllSkills()).thenReturn(List.of(mockSkill));

    mockMvc
        .perform(get("/about"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/about"))
        .andExpect(model().attributeExists("skills"))
        .andExpect(model().attributeExists("categories"));
  }

  @Test
  @DisplayName("GET / - should return index view")
  void shouldReturnIndexPage() throws Exception {
    mockMvc.perform(get("/")).andExpect(status().isOk()).andExpect(view().name("pages/index"));
  }

  @Test
  @DisplayName("GET /projects - should return projects view with repos in model")
  void shouldReturnProjectsPageWithRepos() throws Exception {
    when(gitHubService.getRecentRepositories()).thenReturn(List.of());

    mockMvc
        .perform(get("/projects"))
        .andExpect(status().isOk())
        .andExpect(view().name("pages/projects"))
        .andExpect(model().attributeExists("repos"));
  }
}
