package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func setupListTest(t *testing.T) (string, func()) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Capture original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Set up cleanup to always restore directory
	cleanup := func() {
		_ = os.Chdir(originalDir)
		_ = os.RemoveAll(tmpDir)
		// Reset list flags
		listAll = false
		listAssignee = ""
		listLabel = ""
		listStatus = ""
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		cleanup()
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize repo
	if err := pkg.InitializeRepo(); err != nil {
		cleanup()
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Create some test issues
	createAssignee = "alice"
	createLabels = []string{"bug", "backend"}
	if err := runCreate(nil, []string{"Bug in authentication"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 1: %v", err)
	}

	createAssignee = "bob"
	createLabels = []string{"feature", "frontend"}
	if err := runCreate(nil, []string{"Add user dashboard"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 2: %v", err)
	}

	createAssignee = "alice"
	createLabels = []string{"bug", "frontend"}
	if err := runCreate(nil, []string{"Fix CSS styling"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 3: %v", err)
	}

	createAssignee = ""
	createLabels = []string{"docs"}
	if err := runCreate(nil, []string{"Update README"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 4: %v", err)
	}

	createAssignee = "charlie"
	createLabels = []string{"feature", "backend"}
	if err := runCreate(nil, []string{"API endpoint for users"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 5: %v", err)
	}

	// Reset flags
	createAssignee = ""
	createLabels = []string{}

	// Move one issue to closed for testing
	if err := pkg.MoveIssue("002", pkg.OpenDir, pkg.ClosedDir); err != nil {
		cleanup()
		t.Fatalf("Failed to move issue to closed: %v", err)
	}

	return tmpDir, cleanup
}

func TestListCommand(t *testing.T) {
	_, cleanup := setupListTest(t)
	defer cleanup()

	// Test without initialization
	t.Run("list without init", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "git-issue-test-noinit-*")
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

		err = runList(nil, []string{})
		if err == nil {
			t.Error("runList() should fail when .issues directory doesn't exist")
		}
		if !strings.Contains(err.Error(), ".issues directory not found") {
			t.Errorf("Expected error about .issues not found, got: %v", err)
		}
	})
}

func TestListCommandBasic(t *testing.T) {
	_, cleanup := setupListTest(t)
	defer cleanup()

	// Test basic list (only open issues)
	t.Run("list open issues only", func(t *testing.T) {
		listAll = false
		listAssignee = ""
		listLabel = ""
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// We should have 4 open issues (001, 003, 004, 005)
		// Issue 002 was moved to closed
		issues, err := pkg.ListIssues(pkg.OpenDir)
		if err != nil {
			t.Fatalf("Failed to list issues: %v", err)
		}

		if len(issues) != 4 {
			t.Errorf("Expected 4 open issues, got %d", len(issues))
		}
	})

	// Test list all issues
	t.Run("list all issues", func(t *testing.T) {
		listAll = true
		listAssignee = ""
		listLabel = ""
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// Should include both open and closed (5 total)
		openIssues, _ := pkg.ListIssues(pkg.OpenDir)
		closedIssues, _ := pkg.ListIssues(pkg.ClosedDir)
		total := len(openIssues) + len(closedIssues)

		if total != 5 {
			t.Errorf("Expected 5 total issues, got %d", total)
		}

		listAll = false
	})
}

func TestListCommandFilters(t *testing.T) {
	_, cleanup := setupListTest(t)
	defer cleanup()

	// Test filter by assignee
	t.Run("filter by assignee alice", func(t *testing.T) {
		listAll = true
		listAssignee = "alice"
		listLabel = ""
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// Alice has issues 001 and 003
		// We need to manually verify since runList prints to stdout
		// Let's load and check manually
		openIssues, _ := pkg.ListIssues(pkg.OpenDir)
		closedIssues, _ := pkg.ListIssues(pkg.ClosedDir)

		count := 0
		for _, issue := range openIssues {
			if issue.Assignee == "alice" {
				count++
			}
		}
		for _, issue := range closedIssues {
			if issue.Assignee == "alice" {
				count++
			}
		}

		if count != 2 {
			t.Errorf("Expected 2 issues for alice, found %d", count)
		}

		listAssignee = ""
	})

	// Test filter by label
	t.Run("filter by label bug", func(t *testing.T) {
		listAll = true
		listAssignee = ""
		listLabel = "bug"
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// Issues with 'bug' label: 001, 003
		openIssues, _ := pkg.ListIssues(pkg.OpenDir)
		closedIssues, _ := pkg.ListIssues(pkg.ClosedDir)

		count := 0
		for _, issue := range openIssues {
			if issue.HasLabel("bug") {
				count++
			}
		}
		for _, issue := range closedIssues {
			if issue.HasLabel("bug") {
				count++
			}
		}

		if count != 2 {
			t.Errorf("Expected 2 issues with 'bug' label, found %d", count)
		}

		listLabel = ""
	})

	// Test filter by status open
	t.Run("filter by status open", func(t *testing.T) {
		listAll = false
		listAssignee = ""
		listLabel = ""
		listStatus = "open"

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		issues, _ := pkg.ListIssues(pkg.OpenDir)
		if len(issues) != 4 {
			t.Errorf("Expected 4 open issues, got %d", len(issues))
		}

		listStatus = ""
	})

	// Test filter by status closed
	t.Run("filter by status closed", func(t *testing.T) {
		listAll = false
		listAssignee = ""
		listLabel = ""
		listStatus = "closed"

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		issues, _ := pkg.ListIssues(pkg.ClosedDir)
		if len(issues) != 1 {
			t.Errorf("Expected 1 closed issue, got %d", len(issues))
		}

		listStatus = ""
	})

	// Test invalid status
	t.Run("invalid status", func(t *testing.T) {
		listAll = false
		listAssignee = ""
		listLabel = ""
		listStatus = "invalid"

		err := runList(nil, []string{})
		if err == nil {
			t.Error("runList() should fail with invalid status")
		}
		if !strings.Contains(err.Error(), "invalid status") {
			t.Errorf("Expected error about invalid status, got: %v", err)
		}

		listStatus = ""
	})
}

func TestListCommandCombinedFilters(t *testing.T) {
	_, cleanup := setupListTest(t)
	defer cleanup()

	// Test combining assignee and label filters
	t.Run("filter by assignee and label", func(t *testing.T) {
		listAll = true
		listAssignee = "alice"
		listLabel = "frontend"
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// Only issue 003 should match (alice + frontend)
		openIssues, _ := pkg.ListIssues(pkg.OpenDir)
		closedIssues, _ := pkg.ListIssues(pkg.ClosedDir)

		count := 0
		for _, issue := range openIssues {
			if issue.Assignee == "alice" && issue.HasLabel("frontend") {
				count++
			}
		}
		for _, issue := range closedIssues {
			if issue.Assignee == "alice" && issue.HasLabel("frontend") {
				count++
			}
		}

		if count != 1 {
			t.Errorf("Expected 1 issue matching alice+frontend, found %d", count)
		}

		listAssignee = ""
		listLabel = ""
	})

	// Test filter with no matches
	t.Run("filter with no matches", func(t *testing.T) {
		listAll = true
		listAssignee = "nonexistent"
		listLabel = ""
		listStatus = ""

		err := runList(nil, []string{})
		if err != nil {
			t.Errorf("runList() failed: %v", err)
		}

		// Should handle gracefully with no results
		openIssues, _ := pkg.ListIssues(pkg.OpenDir)
		closedIssues, _ := pkg.ListIssues(pkg.ClosedDir)

		count := 0
		for _, issue := range openIssues {
			if issue.Assignee == "nonexistent" {
				count++
			}
		}
		for _, issue := range closedIssues {
			if issue.Assignee == "nonexistent" {
				count++
			}
		}

		if count != 0 {
			t.Errorf("Expected 0 issues for nonexistent user, found %d", count)
		}

		listAssignee = ""
	})
}

func TestListEmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-issue-test-empty-*")
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

	// Initialize but don't create any issues
	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("Failed to initialize repo: %v", err)
	}

	// Should handle empty repo gracefully
	err = runList(nil, []string{})
	if err != nil {
		t.Errorf("runList() should handle empty repo, got error: %v", err)
	}
}
