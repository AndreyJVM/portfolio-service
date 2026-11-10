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

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.servlet.view.RedirectView;
import ru.vorobevaqa.dto.ProjectsPageData;
import ru.vorobevaqa.service.ProjectShowcaseService;

/**
 * Контроллер веб-страниц Thymeleaf платформы.
 */
@Controller
public class PageController {

  public static final String VIEW_INDEX = "pages/index";
  public static final String VIEW_EDUCATION = "pages/education";
  public static final String VIEW_PROJECTS = "pages/projects";
  public static final String VIEW_QR = "pages/qr";
  public static final String VIEW_MUSTANG = "pages/mustang";

  public static final String ATTR_PROJECTS_DATA = "projectsData";
  public static final String ATTR_GITHUB_REPOS = "githubRepos";
  public static final String ATTR_DOCKER_REPOS = "dockerRepos";
  public static final String ATTR_ANSIBLE_ROLES = "ansibleRoles";

  private final ProjectShowcaseService projectShowcaseService;

  public PageController(ProjectShowcaseService projectShowcaseService) {
    this.projectShowcaseService = projectShowcaseService;
  }

  @GetMapping("/")
  public String index() {
    return VIEW_INDEX;
  }

  @GetMapping("/education")
  public String education() {
    return VIEW_EDUCATION;
  }

  /**
   * Каноничный HTTP 301 Permanent Redirect для предотвращения дублирования контента в SEO.
   */
  @GetMapping("/about")
  public RedirectView redirectAbout() {
    RedirectView redirectView = new RedirectView("/education", true);
    redirectView.setStatusCode(HttpStatus.MOVED_PERMANENTLY);
    return redirectView;
  }

  @GetMapping("/projects")
  public String projects(Model model) {
    ProjectsPageData projectsData = projectShowcaseService.getProjectsData();
    model.addAttribute(ATTR_PROJECTS_DATA, projectsData);
    // Обратная совместимость с текущими переменными в шаблоне
    model.addAttribute(ATTR_GITHUB_REPOS, projectsData.githubRepos());
    model.addAttribute(ATTR_DOCKER_REPOS, projectsData.dockerRepos());
    model.addAttribute(ATTR_ANSIBLE_ROLES, projectsData.ansibleRoles());
    return VIEW_PROJECTS;
  }

  @GetMapping("/qr")
  public String qrPage() {
    return VIEW_QR;
  }

  @GetMapping("/mustang")
  public String mustang() {
    return VIEW_MUSTANG;
  }
}
