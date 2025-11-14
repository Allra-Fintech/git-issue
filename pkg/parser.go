package pkg

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseMarkdown parses a markdown file with YAML frontmatter into an Issue struct
func ParseMarkdown(content string) (*Issue, error) {
	// Split frontmatter and body
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid issue format: missing YAML frontmatter")
	}

	// Parse YAML frontmatter
	var issue Issue
	if err := yaml.Unmarshal([]byte(parts[1]), &issue); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	// Extract body (everything after second ---)
	body := strings.TrimSpace(parts[2])

	// Extract title from first # heading
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			issue.Title = strings.TrimPrefix(line, "# ")
			// Body is everything after the title line
			if i+1 < len(lines) {
				issue.Body = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			}
			break
		}
	}

	if issue.Title == "" {
		return nil, fmt.Errorf("issue missing title (# heading)")
	}

	return &issue, nil
}

// SerializeIssue converts an Issue struct to markdown format with YAML frontmatter
func SerializeIssue(issue *Issue) (string, error) {
	var buf bytes.Buffer

	// Write YAML frontmatter
	buf.WriteString("---\n")
	yamlData, err := yaml.Marshal(issue)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	buf.Write(yamlData)
	buf.WriteString("---\n\n")

	// Write title
	buf.WriteString("# ")
	buf.WriteString(issue.Title)
	buf.WriteString("\n\n")

	// Write body
	if issue.Body != "" {
		buf.WriteString(issue.Body)
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// GenerateSlug generates a URL-safe slug from a title
// Converts to lowercase, replaces spaces and special chars with hyphens
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove all non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Limit length to 50 characters
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// FormatID formats an issue ID as a zero-padded 3-digit string
func FormatID(id int) string {
	return fmt.Sprintf("%03d", id)
}

// ParseTimestamp parses a timestamp string, returns zero time if empty
func ParseTimestamp(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
