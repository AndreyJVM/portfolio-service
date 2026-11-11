package service

import (
	"testing"
	"time"

	"ru.vorobevaqa/devops-platform/internal/model"
)

func TestIntegrationsService_Cache(t *testing.T) {
	svc := NewIntegrationsService("nonexistent-user-123456", "nonexistent-user-123456")

	// Pre-fill cache
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

func TestIntegrationsService_Fallback(t *testing.T) {
	// При обращении к несуществующим профилям сервис должен вернуть fallback без паники
	svc := NewIntegrationsService("user-definitely-does-not-exist-9999", "user-definitely-does-not-exist-9999")

	data := svc.GetProjectsData()
	if len(data.AnsibleRoles) == 0 {
		t.Errorf("Ожидались дефолтные роли Ansible, получено 0")
	}
}
