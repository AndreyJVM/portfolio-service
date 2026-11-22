package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"ru.vorobevaqa/devops-platform/internal/model"
)

type IntegrationsService struct {
	client     *http.Client
	githubUser string

	cacheMutex sync.RWMutex
	cachedData *model.ProjectsPageData
	cacheUntil time.Time
}

func NewIntegrationsService(githubUser string) *IntegrationsService {
	return &IntegrationsService{
		client: &http.Client{
			Timeout: 4 * time.Second,
		},
		githubUser: githubUser,
	}
}

// PreloadDefaults immediately pre-populates the cache with verified static data
// for a guaranteed instantaneous server response (0ms cold start)
// and triggers an asynchronous background update from APIs.
func (s *IntegrationsService) PreloadDefaults() {
	s.cacheMutex.Lock()
	data := model.ProjectsPageData{
		GitHubRepos:  s.fallbackGitHubRepos(),
		AnsibleRoles: s.getStaticAnsibleRoles(),
	}
	s.cachedData = &data
	s.cacheUntil = time.Now().Add(10 * time.Minute)
	s.cacheMutex.Unlock()

	// Launch background refresh without blocking the main goroutine
	go s.Refresh()
}

// GetProjectsData serves the cached O(1) result from memory without network latency.
// When TTL expires, it triggers a background refresh (Stale-While-Revalidate) without delaying the user.
func (s *IntegrationsService) GetProjectsData() model.ProjectsPageData {
	s.cacheMutex.RLock()
	if s.cachedData != nil {
		data := *s.cachedData
		isExpired := time.Now().After(s.cacheUntil)
		s.cacheMutex.RUnlock()

		if isExpired {
			// Background refresh: return stale data instantly and perform update asynchronously
			go s.Refresh()
		}
		return data
	}
	s.cacheMutex.RUnlock()

	// If cache is not yet initialized (e.g. in unit tests), perform synchronous population
	return s.Refresh()
}

// Refresh queries GitHub in parallel with static roles and updates the in-memory cache.
func (s *IntegrationsService) Refresh() model.ProjectsPageData {
	var wg sync.WaitGroup
	var ghRepos []model.GitHubRepoDto
	var ansibleRoles []model.AnsibleRoleDto

	wg.Add(2)

	go func() {
		defer wg.Done()
		ghRepos = s.fetchGitHubRepos()
	}()

	go func() {
		defer wg.Done()
		ansibleRoles = s.getStaticAnsibleRoles()
	}()

	wg.Wait()

	data := model.ProjectsPageData{
		GitHubRepos:  ghRepos,
		AnsibleRoles: ansibleRoles,
	}

	s.cacheMutex.Lock()
	s.cachedData = &data
	s.cacheUntil = time.Now().Add(15 * time.Minute)
	s.cacheMutex.Unlock()

	return data
}

func (s *IntegrationsService) fetchGitHubRepos() []model.GitHubRepoDto {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=6", s.githubUser)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return s.fallbackGitHubRepos()
	}

	req.Header.Set("User-Agent", "DevOpsPlatform-Go")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return s.fallbackGitHubRepos()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.fallbackGitHubRepos()
	}

	var repos []model.GitHubRepoDto
	if err := json.NewDecoder(resp.Body).Decode(&repos); err == nil && len(repos) > 0 {
		return repos
	}

	return s.fallbackGitHubRepos()
}

func (s *IntegrationsService) fallbackGitHubRepos() []model.GitHubRepoDto {
	return []model.GitHubRepoDto{
		{
			Name:        "portfolio-service",
			Description: "Персональный сайт-визитка и инженерная платформа на Go с Docker, Nginx и Zero-Downtime CI/CD.",
			HTMLURL:     "https://github.com/AndreyJVM/portfolio-service",
			Language:    "Go",
			Stargazers:  12,
			Forks:       2,
		},
		{
			Name:        "container-ctl",
			Description: "Быстрый и легковесный терминальный TUI монитор для хоста, Docker, Compose и Swarm на Go.",
			HTMLURL:     "https://github.com/AndreyJVM/container-ctl",
			Language:    "Go",
			Stargazers:  8,
			Forks:       1,
		},
		{
			Name:        "ansible-role-selenium-grid",
			Description: "Ansible роль для развертывания масштабируемого Selenium Grid 4 через Docker Compose за Traefik с Let's Encrypt.",
			HTMLURL:     "https://github.com/AndreyJVM/ansible-role-selenium-grid",
			Language:    "Java",
			Stargazers:  5,
			Forks:       1,
		},
		{
			Name:        "bash_wiki",
			Description: "Практическая библиотека Bash и Shell скриптов для системного администрирования и автоматизации Linux.",
			HTMLURL:     "https://github.com/AndreyJVM/bash_wiki",
			Language:    "Shell",
			Stargazers:  4,
			Forks:       0,
		},
	}
}

func (s *IntegrationsService) getStaticAnsibleRoles() []model.AnsibleRoleDto {
	return []model.AnsibleRoleDto{
		{
			Name:          "andreyvorobev.docker_hardened",
			Description:   "Установка Docker CE, настройка cgroups v2, no-new-privileges и ротации логов.",
			GithubUser:    s.githubUser,
			DownloadCount: 420,
			URL:           "https://galaxy.ansible.com/ui/standalone/namespaces/24255/",
		},
		{
			Name:          "andreyvorobev.nginx_gateway",
			Description:   "Конфигурация Nginx Reverse Proxy с HTTP/2, Gzip, Let's Encrypt и защитой от перегрузок.",
			GithubUser:    s.githubUser,
			DownloadCount: 290,
			URL:           "https://galaxy.ansible.com/ui/standalone/namespaces/24255/",
		},
		{
			Name:          "andreyvorobev.security_ufw",
			Description:   "Hardening хоста: закрытие неиспользуемых портов, Fail2ban, SSH по ed25519 ключам.",
			GithubUser:    s.githubUser,
			DownloadCount: 175,
			URL:           "https://galaxy.ansible.com/ui/standalone/namespaces/24255/",
		},
	}
}
