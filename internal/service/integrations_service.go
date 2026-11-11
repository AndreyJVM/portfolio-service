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
	dockerUser string

	cacheMutex sync.RWMutex
	cachedData *model.ProjectsPageData
	cacheUntil time.Time
}

func NewIntegrationsService(githubUser, dockerUser string) *IntegrationsService {
	return &IntegrationsService{
		client: &http.Client{
			Timeout: 4 * time.Second,
		},
		githubUser: githubUser,
		dockerUser: dockerUser,
	}
}

// GetProjectsData параллельно запрашивает GitHub, DockerHub и отдает кэшированный результат (TTL: 10 минут).
// Защищен от Thundering Herd (двойная проверка блокировки).
func (s *IntegrationsService) GetProjectsData() model.ProjectsPageData {
	s.cacheMutex.RLock()
	if s.cachedData != nil && time.Now().Before(s.cacheUntil) {
		data := *s.cachedData
		s.cacheMutex.RUnlock()
		return data
	}
	s.cacheMutex.RUnlock()

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Double-Checked Locking: проверяем, не успела ли другая горутина обновить кэш
	if s.cachedData != nil && time.Now().Before(s.cacheUntil) {
		return *s.cachedData
	}

	var wg sync.WaitGroup
	var ghRepos []model.GitHubRepoDto
	var dhRepos []model.DockerHubRepoDto
	var ansibleRoles []model.AnsibleRoleDto

	wg.Add(3)

	go func() {
		defer wg.Done()
		ghRepos = s.fetchGitHubRepos()
	}()

	go func() {
		defer wg.Done()
		dhRepos = s.fetchDockerHubRepos()
	}()

	go func() {
		defer wg.Done()
		ansibleRoles = s.getStaticAnsibleRoles()
	}()

	wg.Wait()

	data := model.ProjectsPageData{
		GitHubRepos:  ghRepos,
		DockerRepos:  dhRepos,
		AnsibleRoles: ansibleRoles,
	}

	s.cachedData = &data
	s.cacheUntil = time.Now().Add(10 * time.Minute)

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
			Description: "Primary Cloud Lab • Открытая образовательная платформа по DevOps, Docker и CI/CD на Go.",
			HTMLURL:     "https://github.com/AndreyJVM/portfolio-service",
			Language:    "Go",
			Stargazers:  12,
			Forks:       2,
		},
		{
			Name:        "ansible-role-docker",
			Description: "Ansible Galaxy роль автоматизированной установки и hardening Docker daemon на Ubuntu 24.04.",
			HTMLURL:     "https://github.com/AndreyJVM/ansible-role-docker",
			Language:    "YAML",
			Stargazers:  8,
			Forks:       1,
		},
		{
			Name:        "qr-microservice",
			Description: "Высокопроизводительный микросервис генерации QR-кодов на Go (Alpine/Scratch, RAM ~8MB).",
			HTMLURL:     "https://github.com/AndreyJVM/portfolio-service/tree/master/qr-service",
			Language:    "Go",
			Stargazers:  5,
			Forks:       0,
		},
		{
			Name:        "infrastructure-live",
			Description: "Production Nginx Gateway, Let's Encrypt SSL Certbot и Blue-Green деплой скрипты.",
			HTMLURL:     "https://github.com/AndreyJVM/portfolio-service/tree/master/nginx",
			Language:    "Shell",
			Stargazers:  4,
			Forks:       0,
		},
	}
}

func (s *IntegrationsService) fetchDockerHubRepos() []model.DockerHubRepoDto {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/?page_size=6", s.dockerUser)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return s.fallbackDockerHubRepos()
	}

	req.Header.Set("User-Agent", "DevOpsPlatform-Go")
	resp, err := s.client.Do(req)
	if err != nil {
		return s.fallbackDockerHubRepos()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.fallbackDockerHubRepos()
	}

	var result struct {
		Results []model.DockerHubRepoDto `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && len(result.Results) > 0 {
		for i := range result.Results {
			if result.Results[i].Namespace == "" {
				result.Results[i].Namespace = s.dockerUser
			}
		}
		return result.Results
	}

	return s.fallbackDockerHubRepos()
}

func (s *IntegrationsService) fallbackDockerHubRepos() []model.DockerHubRepoDto {
	return []model.DockerHubRepoDto{
		{
			Name:        "portfolio-service",
			Namespace:   s.dockerUser,
			Description: "DevOps Learning Platform • Легковесный multi-stage образ на Alpine Linux (~33MB).",
			PullCount:   180,
			StarCount:   3,
		},
		{
			Name:        "qr-service",
			Namespace:   s.dockerUser,
			Description: "Go Scratch QR generator microservice с in-binary healthcheck.",
			PullCount:   120,
			StarCount:   2,
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
