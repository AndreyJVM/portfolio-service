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
	ActiveCourse  string
	CanonicalURL  string
}

// Course represents an educational track (e.g., Linux, Docker).
type Course struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Icon        string   `json:"icon"`
	Color       string   `json:"color"`
	Description string   `json:"description"`
	Level       string   `json:"level"`
	LessonCount int      `json:"lessonCount"`
	Tags        []string `json:"tags"`
	Lessons     []Lesson `json:"lessons,omitempty"`
}

// Lesson represents an individual course lesson.
type Lesson struct {
	ID          string `json:"id"`
	CourseSlug  string `json:"courseSlug"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Order       int    `json:"order"`
	Duration    string `json:"duration"`
	Description string `json:"description"`
	ContentHTML string `json:"contentHtml,omitempty"`
}

// SearchItem represents a search index item for fast lookup.
type SearchItem struct {
	Type        string   `json:"type"`        // "course" or "lesson"
	Title       string   `json:"title"`       // Lesson or course title
	CourseTitle string   `json:"courseTitle"` // Course title (for lessons)
	CourseSlug  string   `json:"courseSlug"`  // Course slug
	LessonSlug  string   `json:"lessonSlug"`  // Lesson slug (empty for course)
	URL         string   `json:"url"`         // Relative URL for navigation
	Description string   `json:"description"` // Short description or text snippet
	Tags        []string `json:"tags"`        // Technology tags
	Snippet     string   `json:"snippet"`     // Sanitized text snippet
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

// DockerHubRepoDto represents a Docker Hub image.
type DockerHubRepoDto struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	PullCount   int    `json:"pull_count"`
	StarCount   int    `json:"star_count"`
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
	DockerRepos  []DockerHubRepoDto
	AnsibleRoles []AnsibleRoleDto
}
