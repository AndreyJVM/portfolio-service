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
func (s *IntegrationsService) GetProjectsData() model.ProjectsPageData {
	s.cacheMutex.RLock()
	if s.cachedData != nil && time.Now().Before(s.cacheUntil) {
		data := *s.cachedData
		s.cacheMutex.RUnlock()
		return data
	}
	s.cacheMutex.RUnlock()

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

	s.cacheMutex.Lock()
	s.cachedData = &data
	s.cacheUntil = time.Now().Add(10 * time.Minute)
	s.cacheMutex.Unlock()

	return data
}

func (s *IntegrationsService) fetchGitHubRepos() []model.GitHubRepoDto {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=6", s.githubUser)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []model.GitHubRepoDto{}
	}
	req.Header.Set("User-Agent", "DevOpsPlatform-Go")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return []model.GitHubRepoDto{}
	}
	defer resp.Body.Close()

	var repos []model.GitHubRepoDto
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return []model.GitHubRepoDto{}
	}
	return repos
}

func (s *IntegrationsService) fetchDockerHubRepos() []model.DockerHubRepoDto {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/?page_size=6", s.dockerUser)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []model.DockerHubRepoDto{}
	}
	req.Header.Set("User-Agent", "DevOpsPlatform-Go")

	resp, err := s.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return []model.DockerHubRepoDto{}
	}
	defer resp.Body.Close()

	var result struct {
		Results []model.DockerHubRepoDto `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []model.DockerHubRepoDto{}
	}
	return result.Results
}

func (s *IntegrationsService) getStaticAnsibleRoles() []model.AnsibleRoleDto {
	return []model.AnsibleRoleDto{
		{
			Name:          "security-baseline",
			Description:   "Ansible роль для базового харденинга Ubuntu/Debian: UFW, SSH, Fail2ban.",
			GithubUser:    s.githubUser,
			DownloadCount: 42,
			URL:           "https://galaxy.ansible.com/ui/standalone/namespaces/24255/",
		},
		{
			Name:          "docker-environment",
			Description:   "Автоматизированная установка Docker Engine, Compose v2 и daemon.json.",
			GithubUser:    s.githubUser,
			DownloadCount: 38,
			URL:           "https://galaxy.ansible.com/ui/standalone/namespaces/24255/",
		},
	}
}
