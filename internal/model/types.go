package model

// MetaData содержит SEO и метаинформацию для страниц.
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

// Course представляет образовательный трек (например, Linux, Docker).
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

// Lesson представляет отдельный урок курса.
type Lesson struct {
	ID          string   `json:"id"`
	CourseSlug  string   `json:"courseSlug"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Order       int      `json:"order"`
	Duration    string   `json:"duration"`
	Description string   `json:"description"`
	ContentHTML string   `json:"contentHtml,omitempty"`
	Quiz        *Quiz    `json:"quiz,omitempty"`
}

// SearchItem представляет элемент поискового индекса для быстрого поиска.
type SearchItem struct {
	Type        string   `json:"type"`        // "course" или "lesson"
	Title       string   `json:"title"`       // Заголовок урока или курса
	CourseTitle string   `json:"courseTitle"` // Название курса (для уроков)
	CourseSlug  string   `json:"courseSlug"`  // Слаг курса
	LessonSlug  string   `json:"lessonSlug"`  // Слаг урока (пустой для курса)
	URL         string   `json:"url"`         // Относительный URL для перехода
	Description string   `json:"description"` // Краткое описание или фрагмент текста
	Tags        []string `json:"tags"`        // Теги технологий
	Snippet     string   `json:"snippet"`     // Очищенный текстовый сниппет
}

// Quiz представляет интерактивный тест к уроку для самопроверки.
type Quiz struct {
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	CorrectIdx  int      `json:"correctIndex"`
	Explanation string   `json:"explanation"`
}

// GitHubRepoDto представляет репозиторий GitHub.
type GitHubRepoDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
}

// DockerHubRepoDto представляет образ в Docker Hub.
type DockerHubRepoDto struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	PullCount   int    `json:"pull_count"`
	StarCount   int    `json:"star_count"`
}

// AnsibleRoleDto представляет роль Ansible Galaxy.
type AnsibleRoleDto struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	GithubUser    string `json:"github_user"`
	DownloadCount int    `json:"download_count"`
	URL           string `json:"url"`
}

// ProjectsPageData объединяет все внешние интеграции.
type ProjectsPageData struct {
	GitHubRepos []GitHubRepoDto
	DockerRepos []DockerHubRepoDto
	AnsibleRoles []AnsibleRoleDto
}
