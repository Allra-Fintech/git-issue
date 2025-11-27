package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	IssuesDir     = ".issues"
	OpenDir       = "open"
	ClosedDir     = "closed"
	CounterFile   = ".counter"
	TemplateFile  = "template.md"
	DefaultEditor = "vim"
)

// InitializeRepo creates the .issues/ directory structure
func InitializeRepo() error {
	// Create main directory
	if err := os.MkdirAll(IssuesDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", IssuesDir, err)
	}

	// Create open and closed subdirectories
	openPath := filepath.Join(IssuesDir, OpenDir)
	if err := os.MkdirAll(openPath, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", openPath, err)
	}

	closedPath := filepath.Join(IssuesDir, ClosedDir)
	if err := os.MkdirAll(closedPath, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", closedPath, err)
	}

	// Create .keep files to ensure directories are tracked in git
	openKeepPath := filepath.Join(openPath, ".keep")
	if _, err := os.Stat(openKeepPath); os.IsNotExist(err) {
		if err := os.WriteFile(openKeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .keep file in open directory: %w", err)
		}
	}

	closedKeepPath := filepath.Join(closedPath, ".keep")
	if _, err := os.Stat(closedKeepPath); os.IsNotExist(err) {
		if err := os.WriteFile(closedKeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .keep file in closed directory: %w", err)
		}
	}

	// Initialize counter file
	counterPath := filepath.Join(IssuesDir, CounterFile)
	if _, err := os.Stat(counterPath); os.IsNotExist(err) {
		if err := os.WriteFile(counterPath, []byte("1\n"), 0644); err != nil {
			return fmt.Errorf("failed to create counter file: %w", err)
		}
	}

	// Create template file
	templatePath := filepath.Join(IssuesDir, TemplateFile)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		template := `---
id: ""
assignee: ""
labels: []
created:
updated:
---

# Issue Title

## Description

Describe the issue here...

## Requirements

- Requirement 1
- Requirement 2

## Success Criteria

- [ ] Criterion 1
- [ ] Criterion 2
`
		if err := os.WriteFile(templatePath, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to create template file: %w", err)
		}
	}

	return nil
}

// GetNextID reads and increments the counter, skipping any IDs that already exist
func GetNextID() (int, error) {
	counterPath := filepath.Join(IssuesDir, CounterFile)

	// Read current counter value
	data, err := os.ReadFile(counterPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read counter: %w", err)
	}

	currentID, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("invalid counter value: %w", err)
	}

	// Find the next available ID by checking if current ID exists
	availableID := currentID
	for {
		formattedID := FormatID(availableID)
		_, _, err := FindIssueFile(formattedID)
		if err != nil {
			// ID not found, so it's available
			break
		}
		// ID exists (in either open or closed), try next one
		availableID++
	}

	// Write the next ID after the one we're returning
	nextID := availableID + 1
	if err := os.WriteFile(counterPath, []byte(fmt.Sprintf("%d\n", nextID)), 0644); err != nil {
		return 0, fmt.Errorf("failed to write counter: %w", err)
	}

	return availableID, nil
}

// SaveIssue writes an issue to the specified directory (open or closed)
func SaveIssue(issue *Issue, dir string) error {
	var path string

	// If the issue already exists in the target directory, preserve its existing filename
	if existingPath, existingDir, err := FindIssueFile(issue.ID); err == nil {
		if existingDir != dir {
			return fmt.Errorf("issue %s exists in %s directory, cannot save to %s", issue.ID, existingDir, dir)
		}
		path = existingPath
	} else {
		// Generate a new filename only when the issue doesn't exist yet
		slug := GenerateSlug(issue.Title)
		filename := fmt.Sprintf("%s-%s.md", issue.ID, slug)
		path = filepath.Join(IssuesDir, dir, filename)
	}

	// Serialize issue
	content, err := SerializeIssue(issue)
	if err != nil {
		return fmt.Errorf("failed to serialize issue: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write issue file: %w", err)
	}

	return nil
}

// LoadIssue reads an issue from file system by ID (searches both open/ and closed/)
func LoadIssue(id string) (*Issue, string, error) {
	// Try to find the issue file
	path, dir, err := FindIssueFile(id)
	if err != nil {
		return nil, "", err
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read issue file: %w", err)
	}

	// Parse issue
	issue, err := ParseMarkdown(string(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse issue: %w", err)
	}

	return issue, dir, nil
}

// MoveIssue moves an issue file between directories and updates the timestamp
func MoveIssue(id string, fromDir, toDir string) error {
	// Find the issue file
	oldPath, currentDir, err := FindIssueFile(id)
	if err != nil {
		return err
	}

	// Verify it's in the expected source directory
	if currentDir != fromDir {
		return fmt.Errorf("issue %s is in %s, not %s", id, currentDir, fromDir)
	}

	// Generate new path (filename stays the same)
	filename := filepath.Base(oldPath)
	newPath := filepath.Join(IssuesDir, toDir, filename)

	// Move file atomically
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move issue file: %w", err)
	}

	// Update timestamp after successful move
	data, err := os.ReadFile(newPath)
	if err != nil {
		return fmt.Errorf("failed to read issue file: %w", err)
	}

	issue, err := ParseMarkdown(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse issue: %w", err)
	}

	issue.Updated = time.Now()

	content, err := SerializeIssue(issue)
	if err != nil {
		return fmt.Errorf("failed to serialize issue: %w", err)
	}

	if err := os.WriteFile(newPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write issue file: %w", err)
	}

	return nil
}

// ListIssues gets all issues from a directory
func ListIssues(dir string) ([]*Issue, error) {
	dirPath := filepath.Join(IssuesDir, dir)

	// Read directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var issues []*Issue
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Read and parse issue
		path := filepath.Join(dirPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue // Skip files we can't read
		}

		issue, err := ParseMarkdown(string(data))
		if err != nil {
			continue // Skip files we can't parse
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// FindIssueFile searches for an issue file by ID pattern in both open/ and closed/
// Returns the full path and the directory name (open or closed)
func FindIssueFile(id string) (string, string, error) {
	// Search in open directory first
	openPath := filepath.Join(IssuesDir, OpenDir)
	if path, err := findInDirectory(openPath, id); err == nil {
		return path, OpenDir, nil
	}

	// Search in closed directory
	closedPath := filepath.Join(IssuesDir, ClosedDir)
	if path, err := findInDirectory(closedPath, id); err == nil {
		return path, ClosedDir, nil
	}

	return "", "", fmt.Errorf("issue %s not found", id)
}

// findInDirectory searches for a file matching the ID pattern in a specific directory
func findInDirectory(dir, id string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	pattern := fmt.Sprintf("%s-", id)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), pattern) && strings.HasSuffix(entry.Name(), ".md") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("not found")
}

// DeleteIssue removes an issue file (for cleanup/testing)
func DeleteIssue(id string) error {
	path, _, err := FindIssueFile(id)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete issue file: %w", err)
	}

	return nil
}

// RepoExists checks if the .issues directory exists
func RepoExists() bool {
	_, err := os.Stat(IssuesDir)
	return err == nil
}

// GetIssuesPath returns the path to the .issues directory
func GetIssuesPath() string {
	return IssuesDir
}

// GetOpenPath returns the path to the open issues directory
func GetOpenPath() string {
	return filepath.Join(IssuesDir, OpenDir)
}

// GetClosedPath returns the path to the closed issues directory
func GetClosedPath() string {
	return filepath.Join(IssuesDir, ClosedDir)
}

// LoadTemplateBody reads the template file and extracts the body content (after title heading)
func LoadTemplateBody() string {
	templatePath := filepath.Join(IssuesDir, TemplateFile)

	// Read template file
	data, err := os.ReadFile(templatePath)
	if err != nil {
		// If template doesn't exist, return empty body
		return ""
	}

	// Parse the template to extract body
	content := string(data)
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		// Invalid template format, return empty
		return ""
	}

	// Get everything after frontmatter
	body := strings.TrimSpace(parts[2])

	// Skip the title line (first # heading) and get everything after
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			// Return everything after the title line
			if i+1 < len(lines) {
				return strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			}
			return ""
		}
	}

	// If no title found, return the whole body
	return body
}

// NewIssue creates a new Issue with default values
func NewIssue(id int, title, assignee string, labels []string) *Issue {
	now := time.Now()

	// Load template body
	templateBody := LoadTemplateBody()

	return &Issue{
		ID:       FormatID(id),
		Assignee: assignee,
		Labels:   labels,
		Created:  now,
		Updated:  now,
		Title:    title,
		Body:     templateBody,
	}
}
