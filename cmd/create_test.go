package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func TestCreateCommand(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test without initialization
	t.Run("create without init", func(t *testing.T) {
		err := runCreate(nil, []string{"Test issue"})
		if err == nil {
			t.Error("runCreate() should fail when .issues directory doesn't exist")
		}
		if !strings.Contains(err.Error(), ".issues directory not found") {
			t.Errorf("Expected error about .issues not found, got: %v", err)
		}
	})

	// Initialize repo for subsequent tests
	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Test basic issue creation
	t.Run("create basic issue", func(t *testing.T) {
		// Reset flags
		createAssignee = ""
		createLabels = []string{}

		err := runCreate(nil, []string{"Fix", "authentication", "bug"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		// Verify issue file exists
		files, err := os.ReadDir(pkg.GetOpenPath())
		if err != nil {
			t.Fatalf("Failed to read open directory: %v", err)
		}

		// Count only .md files (skip .keep)
		var issueFiles []os.DirEntry
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".md") {
				issueFiles = append(issueFiles, f)
			}
		}

		if len(issueFiles) != 1 {
			t.Errorf("Expected 1 issue file, found %d", len(issueFiles))
		}

		// Verify filename format
		expectedPrefix := "001-"
		if !strings.HasPrefix(issueFiles[0].Name(), expectedPrefix) {
			t.Errorf("Expected filename to start with %q, got %q", expectedPrefix, issueFiles[0].Name())
		}

		// Load and verify issue
		issue, _, err := pkg.LoadIssue("001")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.ID != "001" {
			t.Errorf("Expected ID 001, got %s", issue.ID)
		}

		expectedTitle := "Fix authentication bug"
		if issue.Title != expectedTitle {
			t.Errorf("Expected title %q, got %q", expectedTitle, issue.Title)
		}

		if issue.Assignee != "" {
			t.Errorf("Expected no assignee, got %q", issue.Assignee)
		}

		if len(issue.Labels) != 0 {
			t.Errorf("Expected no labels, got %v", issue.Labels)
		}
	})

	// Test issue creation with assignee
	t.Run("create with assignee", func(t *testing.T) {
		createAssignee = "john"
		createLabels = []string{}

		err := runCreate(nil, []string{"Add user profile"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		issue, _, err := pkg.LoadIssue("002")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.Assignee != "john" {
			t.Errorf("Expected assignee 'john', got %q", issue.Assignee)
		}

		// Reset flag
		createAssignee = ""
	})

	// Test issue creation with labels
	t.Run("create with labels", func(t *testing.T) {
		createAssignee = ""
		createLabels = []string{"bug", "backend", "urgent"}

		err := runCreate(nil, []string{"Database connection issue"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		issue, _, err := pkg.LoadIssue("003")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if len(issue.Labels) != 3 {
			t.Errorf("Expected 3 labels, got %d", len(issue.Labels))
		}

		expectedLabels := []string{"bug", "backend", "urgent"}
		for i, expected := range expectedLabels {
			if i >= len(issue.Labels) || issue.Labels[i] != expected {
				t.Errorf("Expected label[%d] to be %q, got %q", i, expected, issue.Labels[i])
			}
		}

		// Reset flag
		createLabels = []string{}
	})

	// Test issue creation with both assignee and labels
	t.Run("create with assignee and labels", func(t *testing.T) {
		createAssignee = "jane"
		createLabels = []string{"feature", "frontend"}

		err := runCreate(nil, []string{"Implement dark mode"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		issue, _, err := pkg.LoadIssue("004")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.Assignee != "jane" {
			t.Errorf("Expected assignee 'jane', got %q", issue.Assignee)
		}

		if len(issue.Labels) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(issue.Labels))
		}

		// Reset flags
		createAssignee = ""
		createLabels = []string{}
	})

	// Test ID increment
	t.Run("verify ID increment", func(t *testing.T) {
		createAssignee = ""
		createLabels = []string{}

		err := runCreate(nil, []string{"Fifth issue"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		issue, _, err := pkg.LoadIssue("005")
		if err != nil {
			t.Fatalf("Failed to load issue 005: %v", err)
		}

		if issue.ID != "005" {
			t.Errorf("Expected ID 005, got %s", issue.ID)
		}
	})

	// Test filename slug generation
	t.Run("verify filename slug", func(t *testing.T) {
		createAssignee = ""
		createLabels = []string{}

		err := runCreate(nil, []string{"Fix: Special & Characters!! Test"})
		if err != nil {
			t.Errorf("runCreate() failed: %v", err)
		}

		// Find the file
		files, err := os.ReadDir(pkg.GetOpenPath())
		if err != nil {
			t.Fatalf("Failed to read open directory: %v", err)
		}

		var foundFile string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "006-") {
				foundFile = file.Name()
				break
			}
		}

		if foundFile == "" {
			t.Fatal("Could not find issue 006 file")
		}

		// Verify slug is URL-safe
		expectedSlug := "006-fix-special-characters-test.md"
		if foundFile != expectedSlug {
			t.Errorf("Expected filename %q, got %q", expectedSlug, foundFile)
		}
	})
}

func TestCreateCommandNoArgs(t *testing.T) {
	// This tests the Cobra validation
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize repo
	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Test with empty args
	err = runCreate(nil, []string{})
	if err == nil {
		t.Error("runCreate() should fail with empty title")
	}
	if !strings.Contains(err.Error(), "title cannot be empty") {
		t.Errorf("Expected error about empty title, got: %v", err)
	}

	// Test with whitespace-only title
	err = runCreate(nil, []string{"   ", "  "})
	if err == nil {
		t.Error("runCreate() should fail with whitespace-only title")
	}
	if !strings.Contains(err.Error(), "title cannot be empty") {
		t.Errorf("Expected error about empty title, got: %v", err)
	}
}

func TestCreateIssueFile(t *testing.T) {
	// Test that created files are valid markdown with YAML frontmatter
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	createAssignee = "testuser"
	createLabels = []string{"test"}
	defer func() {
		createAssignee = ""
		createLabels = []string{}
	}()

	err = runCreate(nil, []string{"Test Issue"})
	if err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	// Read the file directly
	filePath := filepath.Join(pkg.GetOpenPath(), "001-test-issue.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read issue file: %v", err)
	}

	contentStr := string(content)

	// Verify YAML frontmatter
	if !strings.HasPrefix(contentStr, "---\n") {
		t.Error("Issue file should start with YAML frontmatter delimiter")
	}

	if !strings.Contains(contentStr, "id: \"001\"") {
		t.Error("Issue file should contain ID in frontmatter")
	}

	if !strings.Contains(contentStr, "assignee: testuser") {
		t.Error("Issue file should contain assignee in frontmatter")
	}

	if !strings.Contains(contentStr, "# Test Issue") {
		t.Error("Issue file should contain title as heading")
	}
}
