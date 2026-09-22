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

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import ru.vorobevaqa.service.GitHubService;

@Controller
@RequiredArgsConstructor
public class PageController {

  private final GitHubService gitHubService;

  @GetMapping("/")
  public String index() {
    return "pages/index";
  }

  @GetMapping({"/about", "/education"})
  public String about() {
    return "pages/about";
  }

  @GetMapping("/projects")
  public String projects(Model model) {
    model.addAttribute("repos", gitHubService.getRecentRepositories());
    return "pages/projects";
  }

  @GetMapping("/qr")
  public String qrPage() {
    return "pages/qr";
  }
}
