package pkg

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Fix Bug", "fix-bug"},
		{"Add User Authentication", "add-user-authentication"},
		{"Fix  Multiple   Spaces", "fix-multiple-spaces"},
		{"Special!@#$%Characters", "specialcharacters"},
		{"Lowercase-Already", "lowercase-already"},
		{"Title With (Parentheses)", "title-with-parentheses"},
		{"Title/With/Slashes", "titlewithslashes"},
		{"  Leading and Trailing Spaces  ", "leading-and-trailing-spaces"},
		{"Multiple---Hyphens", "multiple-hyphens"},
		{strings.Repeat("A", 100), strings.Repeat("a", 50)}, // Test length limit
	}

	for _, tt := range tests {
		result := GenerateSlug(tt.input)
		if result != tt.expected {
			t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatID(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1, "001"},
		{10, "010"},
		{100, "100"},
		{999, "999"},
	}

	for _, tt := range tests {
		result := FormatID(tt.input)
		if result != tt.expected {
			t.Errorf("FormatID(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseMarkdown(t *testing.T) {
	validIssue := `---
id: "001"
assignee: "jonghun"
labels: [bug, backend]
created: 2025-11-14T10:30:00Z
updated: 2025-11-14T14:20:00Z
---

# Fix Redis Connection Timeout

## Description

Users are experiencing intermittent authentication failures.

## Steps to Reproduce

1. Open app
2. Try to login
`

	issue, err := ParseMarkdown(validIssue)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if issue.ID != "001" {
		t.Errorf("ID = %q, want %q", issue.ID, "001")
	}
	if issue.Assignee != "jonghun" {
		t.Errorf("Assignee = %q, want %q", issue.Assignee, "jonghun")
	}
	if len(issue.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(issue.Labels))
	}
	if issue.Title != "Fix Redis Connection Timeout" {
		t.Errorf("Title = %q, want %q", issue.Title, "Fix Redis Connection Timeout")
	}
	if !strings.Contains(issue.Body, "Users are experiencing") {
		t.Errorf("Body should contain description")
	}
}

func TestParseMarkdown_MissingFrontmatter(t *testing.T) {
	invalidIssue := `# Just a title

Some content here.
`

	_, err := ParseMarkdown(invalidIssue)
	if err == nil {
		t.Error("ParseMarkdown() should return error for missing frontmatter")
	}
}

func TestParseMarkdown_MissingTitle(t *testing.T) {
	invalidIssue := `---
id: "001"
assignee: ""
labels: []
created: 2025-11-14T10:30:00Z
updated: 2025-11-14T14:20:00Z
---

Some content without a title heading.
`

	_, err := ParseMarkdown(invalidIssue)
	if err == nil {
		t.Error("ParseMarkdown() should return error for missing title")
	}
}

func TestSerializeIssue(t *testing.T) {
	now := time.Date(2025, 11, 14, 10, 30, 0, 0, time.UTC)
	issue := &Issue{
		ID:       "001",
		Assignee: "jonghun",
		Labels:   []string{"bug", "backend"},
		Created:  now,
		Updated:  now,
		Title:    "Fix Redis Connection",
		Body:     "## Description\n\nFix the timeout issue.",
	}

	content, err := SerializeIssue(issue)
	if err != nil {
		t.Fatalf("SerializeIssue() error = %v", err)
	}

	// Verify frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") {
		t.Error("Content should start with ---")
	}

	// Verify title
	if !strings.Contains(content, "# Fix Redis Connection") {
		t.Error("Content should contain title")
	}

	// Verify body
	if !strings.Contains(content, "## Description") {
		t.Error("Content should contain body")
	}

	// Verify it can be parsed back
	parsed, err := ParseMarkdown(content)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if parsed.ID != issue.ID {
		t.Errorf("Round-trip ID = %q, want %q", parsed.ID, issue.ID)
	}
	if parsed.Title != issue.Title {
		t.Errorf("Round-trip Title = %q, want %q", parsed.Title, issue.Title)
	}
}

func TestParseSerializeRoundTrip(t *testing.T) {
	original := `---
id: "042"
assignee: "alice"
labels: [feature, frontend]
created: 2025-11-14T10:30:00Z
updated: 2025-11-14T14:20:00Z
---

# Implement Dark Mode

## Requirements

- Toggle button in settings
- Persist preference
- Update all components
`

	// Parse
	issue, err := ParseMarkdown(original)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	// Serialize
	serialized, err := SerializeIssue(issue)
	if err != nil {
		t.Fatalf("SerializeIssue() error = %v", err)
	}

	// Parse again
	reparsed, err := ParseMarkdown(serialized)
	if err != nil {
		t.Fatalf("ParseMarkdown() second pass error = %v", err)
	}

	// Verify key fields match
	if reparsed.ID != issue.ID {
		t.Errorf("ID mismatch: %q != %q", reparsed.ID, issue.ID)
	}
	if reparsed.Title != issue.Title {
		t.Errorf("Title mismatch: %q != %q", reparsed.Title, issue.Title)
	}
	if reparsed.Assignee != issue.Assignee {
		t.Errorf("Assignee mismatch: %q != %q", reparsed.Assignee, issue.Assignee)
	}
}
