package cmd

import (
	"os"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func setupSearchTest(t *testing.T) (string, func()) {
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
		// Reset search flags
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""
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

	// Create test issues with different content
	createAssignee = "alice"
	createLabels = []string{"bug", "backend"}
	if err := runCreate(nil, []string{"Fix Redis connection timeout"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 1: %v", err)
	}

	createAssignee = "bob"
	createLabels = []string{"feature", "frontend"}
	if err := runCreate(nil, []string{"Add user authentication"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 2: %v", err)
	}

	createAssignee = "alice"
	createLabels = []string{"bug", "frontend"}
	if err := runCreate(nil, []string{"Fix CSS styling issue"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 3: %v", err)
	}

	createAssignee = ""
	createLabels = []string{"docs"}
	if err := runCreate(nil, []string{"Update README documentation"}); err != nil {
		cleanup()
		t.Fatalf("Failed to create issue 4: %v", err)
	}

	// Add body content to some issues for searching
	issue1, _, err := pkg.LoadIssue("001")
	if err != nil {
		cleanup()
		t.Fatalf("Failed to load issue 1: %v", err)
	}
	issue1.Body = "The Redis connection keeps timing out after 30 seconds. We need to increase the timeout configuration."
	if err := pkg.SaveIssue(issue1, pkg.OpenDir); err != nil {
		cleanup()
		t.Fatalf("Failed to save issue 1: %v", err)
	}

	issue2, _, err := pkg.LoadIssue("002")
	if err != nil {
		cleanup()
		t.Fatalf("Failed to load issue 2: %v", err)
	}
	issue2.Body = "Implement user authentication using JWT tokens. Include login, logout, and session management."
	if err := pkg.SaveIssue(issue2, pkg.OpenDir); err != nil {
		cleanup()
		t.Fatalf("Failed to save issue 2: %v", err)
	}

	// Move one issue to closed for testing
	if err := pkg.MoveIssue("002", pkg.OpenDir, pkg.ClosedDir); err != nil {
		cleanup()
		t.Fatalf("Failed to move issue to closed: %v", err)
	}

	// Reset flags
	createAssignee = ""
	createLabels = []string{}

	return tmpDir, cleanup
}

func TestSearchCommand(t *testing.T) {
	_, cleanup := setupSearchTest(t)
	defer cleanup()

	// Test without initialization
	t.Run("search without init", func(t *testing.T) {
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

		err = runSearch(nil, []string{"test"})
		if err == nil {
			t.Error("runSearch() should fail when .issues directory doesn't exist")
		}
	})
}

func TestSearchCommandBasic(t *testing.T) {
	_, cleanup := setupSearchTest(t)
	defer cleanup()

	// Test search by title
	t.Run("search in title", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"Redis"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issue 001 which has "Redis" in title
	})

	// Test search by body
	t.Run("search in body", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"JWT"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issue 002 which has "JWT" in body
	})

	// Test case-insensitive search
	t.Run("case insensitive search", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"REDIS"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issue 001 (case-insensitive)
	})

	// Test no results
	t.Run("no results", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"nonexistent"})
		if err != nil {
			t.Errorf("runSearch() should handle no results gracefully: %v", err)
		}
	})
}

func TestSearchCommandFilters(t *testing.T) {
	_, cleanup := setupSearchTest(t)
	defer cleanup()

	// Test search with status filter
	t.Run("search with status filter open", func(t *testing.T) {
		searchStatus = "open"
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"Fix"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issues 001 and 003 (both open and have "Fix" in title)
		searchStatus = ""
	})

	t.Run("search with status filter closed", func(t *testing.T) {
		searchStatus = "closed"
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"authentication"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issue 002 (closed and has "authentication" in title)
		searchStatus = ""
	})

	// Test search with assignee filter
	t.Run("search with assignee filter", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = "alice"
		searchLabel = ""

		err := runSearch(nil, []string{"Fix"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issues 001 and 003 (assigned to alice and have "Fix")
		searchAssignee = ""
	})

	// Test search with label filter
	t.Run("search with label filter", func(t *testing.T) {
		searchStatus = ""
		searchAssignee = ""
		searchLabel = "bug"

		err := runSearch(nil, []string{"Fix"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find issues 001 and 003 (both have "bug" label)
		searchLabel = ""
	})

	// Test combined filters
	t.Run("search with combined filters", func(t *testing.T) {
		searchStatus = "open"
		searchAssignee = "alice"
		searchLabel = "backend"

		err := runSearch(nil, []string{"Redis"})
		if err != nil {
			t.Errorf("runSearch() failed: %v", err)
		}

		// Should find only issue 001
		searchStatus = ""
		searchAssignee = ""
		searchLabel = ""
	})

	// Test invalid status
	t.Run("invalid status filter", func(t *testing.T) {
		searchStatus = "invalid"
		searchAssignee = ""
		searchLabel = ""

		err := runSearch(nil, []string{"test"})
		if err == nil {
			t.Error("runSearch() should fail with invalid status")
		}

		searchStatus = ""
	})
}

func TestSearchCommandEmptyQuery(t *testing.T) {
	_, cleanup := setupSearchTest(t)
	defer cleanup()

	t.Run("empty query", func(t *testing.T) {
		err := runSearch(nil, []string{})
		if err == nil {
			t.Error("runSearch() should fail with empty query")
		}
	})

	t.Run("whitespace only query", func(t *testing.T) {
		err := runSearch(nil, []string{"   ", "  "})
		if err == nil {
			t.Error("runSearch() should fail with whitespace-only query")
		}
	})
}

func TestSearchEmptyRepo(t *testing.T) {
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
	err = runSearch(nil, []string{"test"})
	if err != nil {
		t.Errorf("runSearch() should handle empty repo, got error: %v", err)
	}
}
