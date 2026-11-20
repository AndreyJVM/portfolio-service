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
	svc := NewIntegrationsService("test-user")

	mockData := model.ProjectsPageData{
		GitHubRepos: []model.GitHubRepoDto{
			{Name: "test-repo", Description: "Test Description"},
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
	if len(data.AnsibleRoles) != 1 || data.AnsibleRoles[0].Name != "test-role" {
		t.Errorf("Ожидались кэшированные данные Ansible, получено: %+v", data.AnsibleRoles)
	}
}

func TestIntegrationsService_PreloadDefaults(t *testing.T) {
	svc := NewIntegrationsService("test-user")
	// Intercept client calls so the background refresh does not fail
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("prevent external network call")
	})

	svc.PreloadDefaults()

	data := svc.GetProjectsData()
	if len(data.GitHubRepos) == 0 {
		t.Errorf("PreloadDefaults должен заполнять fallback список репозиториев GitHub")
	}
	if len(data.AnsibleRoles) == 0 {
		t.Errorf("PreloadDefaults должен заполнять роли Ansible")
	}
}

func TestIntegrationsService_StaleWhileRevalidate(t *testing.T) {
	svc := NewIntegrationsService("test-user")

	refreshed := make(chan bool, 1)
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		select {
		case refreshed <- true:
		default:
		}
		return nil, errors.New("network error")
	})

	// Set expired cache (time in the past)
	staleData := model.ProjectsPageData{
		GitHubRepos: []model.GitHubRepoDto{{Name: "stale-repo"}},
	}
	svc.cacheMutex.Lock()
	svc.cachedData = &staleData
	svc.cacheUntil = time.Now().Add(-1 * time.Minute)
	svc.cacheMutex.Unlock()

	// Call should return staleData immediately
	returned := svc.GetProjectsData()
	if len(returned.GitHubRepos) != 1 || returned.GitHubRepos[0].Name != "stale-repo" {
		t.Errorf("Ожидался мгновенный возврат устаревших данных, получено: %+v", returned)
	}

	// And background Refresh goroutine should be launched
	select {
	case <-refreshed:
		// Background refresh successfully triggered
	case <-time.After(500 * time.Millisecond):
		t.Errorf("Фоновое обновление не было запущено при истечении TTL")
	}
}

func TestIntegrationsService_MockAPIs_Success(t *testing.T) {
	svc := NewIntegrationsService("test-gh")

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
	if len(data.AnsibleRoles) == 0 {
		t.Errorf("Ожидались статические роли Ansible, получено 0")
	}
}

func TestIntegrationsService_APIErrors_Fallback(t *testing.T) {
	svc := NewIntegrationsService("test-gh")

	// Mock 500 Internal Server Error response
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("Server Error")),
			Header:     make(http.Header),
		}, nil
	})

	data := svc.GetProjectsData()

	// On server error fallback should be used
	if len(data.GitHubRepos) == 0 {
		t.Errorf("Ожидался fallback список репозиториев GitHub")
	}
}

func TestIntegrationsService_NetworkError_Fallback(t *testing.T) {
	svc := NewIntegrationsService("test-gh")

	// Mock client network failure
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection timeout or network unreachable")
	})

	data := svc.GetProjectsData()

	if len(data.GitHubRepos) == 0 {
		t.Errorf("Ожидался fallback список GitHub при ошибке сети")
	}
}

func TestIntegrationsService_MalformedJSON_And_Empty(t *testing.T) {
	svc := NewIntegrationsService("test-gh")

	// Mock 200 OK responses, but with malformed JSON
	svc.client.Transport = mockRoundTripper(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.Host, "api.github.com") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`[{"name": "broken-unclosed"`)),
				Header:     make(http.Header),
			}, nil
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	data := svc.GetProjectsData()

	if len(data.GitHubRepos) == 0 {
		t.Errorf("Ожидался fallback список GitHub при битом JSON")
	}
}
