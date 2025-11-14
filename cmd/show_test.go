package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func setupShowTest(t *testing.T) (string, func()) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	_ = os.Chdir(tmpDir)

	// Initialize repo
	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Create test issues
	createAssignee = "alice"
	createLabels = []string{"bug", "backend"}
	_ = runCreate(nil, []string{"Fix authentication bug"})

	createAssignee = "bob"
	createLabels = []string{"feature"}
	_ = runCreate(nil, []string{"Add user dashboard"})

	createAssignee = ""
	createLabels = []string{}
	_ = runCreate(nil, []string{"Update documentation"})

	// Reset flags
	createAssignee = ""
	createLabels = []string{}

	// Move one issue to closed
	_ = pkg.MoveIssue("002", pkg.OpenDir, pkg.ClosedDir)

	cleanup := func() {
		_ = os.Chdir(originalDir)
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestShowCommand(t *testing.T) {
	_, cleanup := setupShowTest(t)
	defer cleanup()

	// Test without initialization
	t.Run("show without init", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "git-issue-test-noinit-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(tmpDir) }()

		originalDir, _ := os.Getwd()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(originalDir) }()

		err = runShow(nil, []string{"001"})
		if err == nil {
			t.Error("runShow() should fail when .issues directory doesn't exist")
		}
		if !strings.Contains(err.Error(), ".issues directory not found") {
			t.Errorf("Expected error about .issues not found, got: %v", err)
		}
	})
}

func TestShowCommandBasic(t *testing.T) {
	_, cleanup := setupShowTest(t)
	defer cleanup()

	// Test showing an open issue
	t.Run("show open issue", func(t *testing.T) {
		err := runShow(nil, []string{"001"})
		if err != nil {
			t.Errorf("runShow() failed: %v", err)
		}

		// Verify issue exists and can be loaded
		issue, dir, err := pkg.LoadIssue("001")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.ID != "001" {
			t.Errorf("Expected ID 001, got %s", issue.ID)
		}

		if dir != pkg.OpenDir {
			t.Errorf("Expected issue in open directory, got %s", dir)
		}

		if issue.Title != "Fix authentication bug" {
			t.Errorf("Expected title 'Fix authentication bug', got %q", issue.Title)
		}

		if issue.Assignee != "alice" {
			t.Errorf("Expected assignee 'alice', got %q", issue.Assignee)
		}

		if len(issue.Labels) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(issue.Labels))
		}
	})

	// Test showing a closed issue
	t.Run("show closed issue", func(t *testing.T) {
		err := runShow(nil, []string{"002"})
		if err != nil {
			t.Errorf("runShow() failed: %v", err)
		}

		// Verify issue exists in closed directory
		issue, dir, err := pkg.LoadIssue("002")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.ID != "002" {
			t.Errorf("Expected ID 002, got %s", issue.ID)
		}

		if dir != pkg.ClosedDir {
			t.Errorf("Expected issue in closed directory, got %s", dir)
		}
	})

	// Test showing issue with no assignee or labels
	t.Run("show issue without assignee/labels", func(t *testing.T) {
		err := runShow(nil, []string{"003"})
		if err != nil {
			t.Errorf("runShow() failed: %v", err)
		}

		issue, _, err := pkg.LoadIssue("003")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.Assignee != "" {
			t.Errorf("Expected no assignee, got %q", issue.Assignee)
		}

		if len(issue.Labels) != 0 {
			t.Errorf("Expected no labels, got %d", len(issue.Labels))
		}
	})
}

func TestShowCommandErrors(t *testing.T) {
	_, cleanup := setupShowTest(t)
	defer cleanup()

	// Test showing non-existent issue
	t.Run("show non-existent issue", func(t *testing.T) {
		err := runShow(nil, []string{"999"})
		if err == nil {
			t.Error("runShow() should fail for non-existent issue")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected error about issue not found, got: %v", err)
		}
	})

	// Test showing with invalid ID
	t.Run("show with invalid ID", func(t *testing.T) {
		err := runShow(nil, []string{"abc"})
		if err == nil {
			t.Error("runShow() should fail for invalid issue ID")
		}
	})
}

func TestShowCommandIDPadding(t *testing.T) {
	_, cleanup := setupShowTest(t)
	defer cleanup()

	// Test showing issue with unpadded ID
	t.Run("show with unpadded ID", func(t *testing.T) {
		// Should work with "1" instead of "001"
		err := runShow(nil, []string{"1"})
		if err != nil {
			t.Errorf("runShow() should handle unpadded ID, got error: %v", err)
		}

		// Verify it loaded the correct issue
		issue, _, err := pkg.LoadIssue("001")
		if err != nil {
			t.Fatalf("Failed to load issue: %v", err)
		}

		if issue.ID != "001" {
			t.Errorf("Expected ID 001, got %s", issue.ID)
		}
	})

	// Test showing issue with various ID formats
	t.Run("show with different ID formats", func(t *testing.T) {
		testCases := []string{"1", "01", "001"}

		for _, id := range testCases {
			err := runShow(nil, []string{id})
			if err != nil {
				t.Errorf("runShow() should handle ID format %q, got error: %v", id, err)
			}
		}
	})
}

func TestShowIssueWithBody(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-issue-test-body-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Create an issue
	createAssignee = "test"
	createLabels = []string{"test"}
	defer func() {
		createAssignee = ""
		createLabels = []string{}
	}()

	err = runCreate(nil, []string{"Test Issue"})
	if err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	// Load the issue and add a body
	issue, _, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("Failed to load issue: %v", err)
	}

	issue.Body = `## Description

This is a detailed description of the issue.

## Steps to Reproduce

1. Step one
2. Step two
3. Step three

## Expected Behavior

The system should work correctly.`

	// Save the updated issue
	if err := pkg.SaveIssue(issue, pkg.OpenDir); err != nil {
		t.Fatalf("Failed to save issue: %v", err)
	}

	// Show the issue
	err = runShow(nil, []string{"001"})
	if err != nil {
		t.Errorf("runShow() failed: %v", err)
	}

	// Verify the body is present
	loadedIssue, _, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("Failed to reload issue: %v", err)
	}

	if loadedIssue.Body == "" {
		t.Error("Issue body should not be empty")
	}

	if !strings.Contains(loadedIssue.Body, "Steps to Reproduce") {
		t.Error("Issue body should contain 'Steps to Reproduce'")
	}

	if !strings.Contains(loadedIssue.Body, "Expected Behavior") {
		t.Error("Issue body should contain 'Expected Behavior'")
	}
}

func TestShowMultipleIssues(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-issue-test-multi-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Create multiple issues
	for i := 1; i <= 10; i++ {
		createAssignee = ""
		createLabels = []string{}
		err := runCreate(nil, []string{"Issue", fmt.Sprintf("%d", i)})
		if err != nil {
			t.Fatalf("Failed to create issue %d: %v", i, err)
		}
	}

	// Test showing each issue
	for i := 1; i <= 10; i++ {
		id := pkg.FormatID(i)
		err := runShow(nil, []string{id})
		if err != nil {
			t.Errorf("Failed to show issue %s: %v", id, err)
		}
	}

	// Verify all issues exist
	issues, err := pkg.ListIssues(pkg.OpenDir)
	if err != nil {
		t.Fatalf("Failed to list issues: %v", err)
	}

	if len(issues) != 10 {
		t.Errorf("Expected 10 issues, got %d", len(issues))
	}
}
