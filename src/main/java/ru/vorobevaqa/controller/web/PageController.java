package ru.vorobevaqa.controller.web;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import ru.vorobevaqa.service.DockerHubService;
import ru.vorobevaqa.service.GitHubService;

@Controller
@RequiredArgsConstructor
public class PageController {

  private final GitHubService gitHubService;
  private final DockerHubService dockerHubService;

  @GetMapping("/")
  public String index() {
    return "pages/index";
  }

  @GetMapping({"/about", "/education"})
  public String education() {
    return "pages/education";
  }

  @GetMapping("/projects")
  public String projects(Model model) {
    model.addAttribute("githubRepos", gitHubService.getRecentRepositories());
    model.addAttribute("dockerRepos", dockerHubService.getRecentRepositories());
    return "pages/projects";
  }

  @GetMapping("/qr")
  public String qrPage() {
    return "pages/qr";
  }
}
