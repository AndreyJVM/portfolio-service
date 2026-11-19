package service

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// lessonMeta stores metadata extracted from frontmatter or filename.
type lessonMeta struct {
	Title       string
	Description string
	Duration    string
	Level       string
	Order       int
	Tags        []string
}

// parseFrontmatter extracts optional YAML-like frontmatter delimited by "---".
func parseFrontmatter(raw []byte) (lessonMeta, []byte) {
	meta := lessonMeta{}
	content := string(raw)

	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return meta, raw
	}

	var bodyStart int
	var frontmatterStr string
	if idx := strings.Index(content[4:], "\n---\n"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 5
	} else if idx := strings.Index(content[4:], "\r\n---\r\n"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 7
	} else if idx := strings.Index(content[4:], "\n---"); idx != -1 {
		frontmatterStr = content[4 : 4+idx]
		bodyStart = 4 + idx + 4
		if bodyStart < len(content) && content[bodyStart] == '\n' {
			bodyStart++
		}
	} else {
		return meta, raw
	}

	lines := strings.Split(frontmatterStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		switch key {
		case "title":
			meta.Title = val
		case "description":
			meta.Description = val
		case "duration":
			meta.Duration = val
		case "level":
			meta.Level = val
		case "order":
			if o, err := strconv.Atoi(val); err == nil {
				meta.Order = o
			}
		case "tags":
			val = strings.Trim(val, "[]")
			for _, tag := range strings.Split(val, ",") {
				t := strings.Trim(strings.TrimSpace(tag), "\"'")
				if t != "" {
					meta.Tags = append(meta.Tags, t)
				}
			}
		}
	}

	return meta, []byte(content[bodyStart:])
}

// estimateDuration calculates reading time based on word count.
func estimateDuration(body []byte) string {
	words := len(strings.Fields(string(body)))
	mins := words / 160
	if mins < 5 {
		mins = 10
	} else if mins > 45 {
		mins = 45
	}
	return fmt.Sprintf("%d мин", mins)
}

// extractOrderFromSlug parses numeric prefix (e.g., "01-basics" -> 1).
func extractOrderFromSlug(slug string, fallback int) int {
	parts := strings.SplitN(slug, "-", 2)
	if len(parts) > 0 {
		if num, err := strconv.Atoi(parts[0]); err == nil && num > 0 {
			return num
		}
	}
	return fallback
}

// extractTitle extracts the h1 title from the first markdown line or formats the fallback slug.
func extractTitle(content []byte, fallbackSlug string) string {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	parts := strings.SplitN(fallbackSlug, "-", 2)
	cleanedSlug := fallbackSlug
	if len(parts) == 2 {
		if _, err := strconv.Atoi(parts[0]); err == nil {
			cleanedSlug = parts[1]
		}
	}
	title := strings.ReplaceAll(cleanedSlug, "-", " ")
	return titleCase(title)
}

// titleCase capitalizes the first letter of each word (replacement for deprecated strings.Title).
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			r := []rune(w)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
