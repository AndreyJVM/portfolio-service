package model

// MetaData contains SEO and meta information for pages.
type MetaData struct {
	Title         string
	Description   string
	Keywords      string
	OGTitle       string
	OGDescription string
	OGImage       string
	OGURL         string
	ActiveTab     string
	CanonicalURL  string
}

// GitHubRepoDto represents a GitHub repository.
type GitHubRepoDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
}

// AnsibleRoleDto represents an Ansible Galaxy role.
type AnsibleRoleDto struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	GithubUser    string `json:"github_user"`
	DownloadCount int    `json:"download_count"`
	URL           string `json:"url"`
}

// ProjectsPageData aggregates all external integrations.
type ProjectsPageData struct {
	GitHubRepos  []GitHubRepoDto
	AnsibleRoles []AnsibleRoleDto
}
