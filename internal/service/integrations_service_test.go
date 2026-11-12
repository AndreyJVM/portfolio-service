package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"ru.vorobevaqa/devops-platform/internal/model"
)

type mockRoundTripper func(req *http.Request) (*http.Response, error)

func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

func TestIntegrationsService_Cache(t *testing.T) {
	svc := NewIntegrationsService("test-user", "test-user")

	mockData := model.ProjectsPageData{
		GitHubRepos: []model.GitHubRepoDto{
			{Name: "test-repo", Description: "Test Description"},
		},
		DockerRepos: []model.DockerHubRepoDto{
			{Name: "test-image", PullCount: 42},
		},
		AnsibleRoles: []model.AnsibleRoleDto{
			{Name: "test-role", Description: "Role Description"},
		},
	}

	svc.cachedData = &mockData
	svc.cacheUntil = time.Now().Add(10 * time.Minute)

	data := svc.GetProjectsData()
	if len(data.GitHubRepos) != 1 || data.GitHubRepos[0].Name != "test-repo" {
		t.Errorf("Ожидались кэшированные данные GitHub, получено: %+v", data.GitHubRepos)
	}
	if len(data.DockerRepos) != 1 || data.DockerRepos[0].PullCount != 42 {
		t.Errorf("Ожидались кэшированные данные Docker, получено: %+v", data.DockerRepos)
	}
	if len(data.AnsibleRoles) != 1 || data.AnsibleRoles[0].Name != "test-role" {
		t.Errorf("Ожидались кэшированные данные Ansible, получено: %+v", data.AnsibleRoles)
	}
}

func TestIntegrationsService_MockAPIs_Success(t *testing.T) {
	svc := NewIntegrationsService("test-gh", "test-dh")

	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.Host, "api.github.com") {
			body := `[
				{
					"name": "custom-gh-repo",
					"description": "Mocked GitHub Repo",
					"html_url": "https://github.com/test/custom-gh-repo",
					"language": "Go",
					"stargazers_count": 99,
					"forks_count": 5
				}
			]`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}, nil
		}

		if strings.Contains(req.URL.Host, "hub.docker.com") {
			body := `{
				"results": [
					{
						"name": "custom-docker-image",
						"description": "Mocked Docker Image",
						"pull_count": 500,
						"star_count": 10,
						"last_updated": "2026-11-12T10:00:00Z"
					}
				]
			}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString("not found")),
			Header:     make(http.Header),
		}, nil
	})

	data := svc.GetProjectsData()

	if len(data.GitHubRepos) != 1 || data.GitHubRepos[0].Name != "custom-gh-repo" {
		t.Errorf("Ожидался замоканный репозиторий GitHub, получено: %+v", data.GitHubRepos)
	}
	if len(data.DockerRepos) != 1 || data.DockerRepos[0].Name != "custom-docker-image" {
		t.Errorf("Ожидался замоканный репозиторий Docker, получено: %+v", data.DockerRepos)
	}
	if len(data.AnsibleRoles) == 0 {
		t.Errorf("Ожидались статические роли Ansible, получено 0")
	}
}

func TestIntegrationsService_APIErrors_Fallback(t *testing.T) {
	svc := NewIntegrationsService("test-gh", "test-dh")

	// Мокаем ответ с кодом 500 Internal Server Error
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("Server Error")),
			Header:     make(http.Header),
		}, nil
	})

	data := svc.GetProjectsData()

	// При ошибке сервера должен сработать fallback
	if len(data.GitHubRepos) == 0 {
		t.Errorf("Ожидался fallback список репозиториев GitHub")
	}
	if len(data.DockerRepos) == 0 {
		t.Errorf("Ожидался fallback список репозиториев Docker")
	}
}

func TestIntegrationsService_NetworkError_Fallback(t *testing.T) {
	svc := NewIntegrationsService("test-gh", "test-dh")

	// Мокаем сетевой сбой клиента
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection timeout or network unreachable")
	})

	data := svc.GetProjectsData()

	if len(data.GitHubRepos) == 0 {
		t.Errorf("Ожидался fallback список GitHub при ошибке сети")
	}
	if len(data.DockerRepos) == 0 {
		t.Errorf("Ожидался fallback список Docker при ошибке сети")
	}
}
