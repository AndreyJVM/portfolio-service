package model

import (
	"encoding/json"
	"testing"
)

func TestProjectsPageDataJSONSerialization(t *testing.T) {
	data := ProjectsPageData{
		GitHubRepos: []GitHubRepoDto{
			{
				Name:        "portfolio-service",
				Description: "High-performance Go service",
				HTMLURL:     "https://github.com/AndreyJVM/portfolio-service",
				Language:    "Go",
				Stargazers:  10,
				Forks:       2,
			},
		},
		DockerRepos: []DockerHubRepoDto{
			{
				Name:        "portfolio-service",
				Namespace:   "andreyvorobevaqa",
				Description: "Production image",
				PullCount:   100,
				StarCount:   5,
			},
		},
		AnsibleRoles: []AnsibleRoleDto{
			{
				Name:          "docker_host",
				Description:   "Ansible role for Docker host",
				GithubUser:    "AndreyJVM",
				DownloadCount: 50,
				URL:           "https://galaxy.ansible.com",
			},
		},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Ошибка сериализации ProjectsPageData в JSON: %v", err)
	}

	var parsed ProjectsPageData
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("Ошибка десериализации ProjectsPageData из JSON: %v", err)
	}

	if len(parsed.GitHubRepos) != 1 || parsed.GitHubRepos[0].Name != "portfolio-service" {
		t.Errorf("Некорректная десериализация GitHubRepos: %+v", parsed.GitHubRepos)
	}
	if len(parsed.DockerRepos) != 1 || parsed.DockerRepos[0].PullCount != 100 {
		t.Errorf("Некорректная десериализация DockerRepos: %+v", parsed.DockerRepos)
	}
	if len(parsed.AnsibleRoles) != 1 || parsed.AnsibleRoles[0].Name != "docker_host" {
		t.Errorf("Некорректная десериализация AnsibleRoles: %+v", parsed.AnsibleRoles)
	}
}

func TestMetaData(t *testing.T) {
	meta := MetaData{
		Title:        "Об авторе • Андрей Воробьев",
		ActiveTab:    "about",
		CanonicalURL: "https://vorobevaqa.ru/",
	}
	if meta.ActiveTab != "about" {
		t.Errorf("Ожидался ActiveTab 'about', получено '%s'", meta.ActiveTab)
	}
	if meta.CanonicalURL != "https://vorobevaqa.ru/" {
		t.Errorf("Ожидался CanonicalURL 'https://vorobevaqa.ru/', получено '%s'", meta.CanonicalURL)
	}
}
