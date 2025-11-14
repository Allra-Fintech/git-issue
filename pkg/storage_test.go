package pkg

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setupTestRepo creates a temporary test repository
func setupTestRepo(t *testing.T) func() {
	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Return cleanup function
	return func() {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	}
}

func TestInitializeRepo(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	err := InitializeRepo()
	if err != nil {
		t.Fatalf("InitializeRepo() error = %v", err)
	}

	// Check directories exist
	if !dirExists(IssuesDir) {
		t.Errorf("%s directory not created", IssuesDir)
	}
	if !dirExists(filepath.Join(IssuesDir, OpenDir)) {
		t.Errorf("%s directory not created", OpenDir)
	}
	if !dirExists(filepath.Join(IssuesDir, ClosedDir)) {
		t.Errorf("%s directory not created", ClosedDir)
	}

	// Check counter file exists
	counterPath := filepath.Join(IssuesDir, CounterFile)
	if !fileExists(counterPath) {
		t.Errorf("Counter file not created")
	}

	// Check template file exists
	templatePath := filepath.Join(IssuesDir, TemplateFile)
	if !fileExists(templatePath) {
		t.Errorf("Template file not created")
	}

	// Verify counter initial value
	data, err := os.ReadFile(counterPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "1\n" {
		t.Errorf("Counter initial value = %q, want %q", string(data), "1\n")
	}
}

func TestGetNextID(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Get first ID
	id1, err := GetNextID()
	if err != nil {
		t.Fatalf("GetNextID() error = %v", err)
	}
	if id1 != 1 {
		t.Errorf("First ID = %d, want 1", id1)
	}

	// Get second ID
	id2, err := GetNextID()
	if err != nil {
		t.Fatalf("GetNextID() error = %v", err)
	}
	if id2 != 2 {
		t.Errorf("Second ID = %d, want 2", id2)
	}

	// Get third ID
	id3, err := GetNextID()
	if err != nil {
		t.Fatalf("GetNextID() error = %v", err)
	}
	if id3 != 3 {
		t.Errorf("Third ID = %d, want 3", id3)
	}
}

func TestSaveAndLoadIssue(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Create test issue
	now := time.Now()
	issue := &Issue{
		ID:       "001",
		Assignee: "alice",
		Labels:   []string{"bug", "urgent"},
		Created:  now,
		Updated:  now,
		Title:    "Test Issue",
		Body:     "This is a test issue.",
	}

	// Save issue
	err := SaveIssue(issue, OpenDir)
	if err != nil {
		t.Fatalf("SaveIssue() error = %v", err)
	}

	// Load issue
	loaded, dir, err := LoadIssue("001")
	if err != nil {
		t.Fatalf("LoadIssue() error = %v", err)
	}

	if dir != OpenDir {
		t.Errorf("LoadIssue() dir = %q, want %q", dir, OpenDir)
	}

	if loaded.ID != issue.ID {
		t.Errorf("ID = %q, want %q", loaded.ID, issue.ID)
	}
	if loaded.Title != issue.Title {
		t.Errorf("Title = %q, want %q", loaded.Title, issue.Title)
	}
	if loaded.Assignee != issue.Assignee {
		t.Errorf("Assignee = %q, want %q", loaded.Assignee, issue.Assignee)
	}
}

func TestMoveIssue(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Create and save issue in open directory
	issue := NewIssue(1, "Test Issue", "bob", []string{"feature"})
	if err := SaveIssue(issue, OpenDir); err != nil {
		t.Fatal(err)
	}

	// Move to closed
	err := MoveIssue("001", OpenDir, ClosedDir)
	if err != nil {
		t.Fatalf("MoveIssue() error = %v", err)
	}

	// Verify it's now in closed
	_, dir, err := LoadIssue("001")
	if err != nil {
		t.Fatalf("LoadIssue() error = %v", err)
	}
	if dir != ClosedDir {
		t.Errorf("Issue dir = %q, want %q", dir, ClosedDir)
	}

	// Try to move from open again (should fail)
	err = MoveIssue("001", OpenDir, ClosedDir)
	if err == nil {
		t.Error("MoveIssue() should fail when issue is not in source directory")
	}
}

func TestListIssues(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Create multiple issues
	issue1 := NewIssue(1, "First Issue", "alice", []string{"bug"})
	issue2 := NewIssue(2, "Second Issue", "bob", []string{"feature"})
	issue3 := NewIssue(3, "Third Issue", "charlie", []string{"enhancement"})

	// Save issues
	SaveIssue(issue1, OpenDir)
	SaveIssue(issue2, OpenDir)
	SaveIssue(issue3, ClosedDir)

	// List open issues
	openIssues, err := ListIssues(OpenDir)
	if err != nil {
		t.Fatalf("ListIssues() error = %v", err)
	}
	if len(openIssues) != 2 {
		t.Errorf("len(openIssues) = %d, want 2", len(openIssues))
	}

	// List closed issues
	closedIssues, err := ListIssues(ClosedDir)
	if err != nil {
		t.Fatalf("ListIssues() error = %v", err)
	}
	if len(closedIssues) != 1 {
		t.Errorf("len(closedIssues) = %d, want 1", len(closedIssues))
	}
}

func TestFindIssueFile(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Create and save issue
	issue := NewIssue(42, "Find Me", "", nil)
	if err := SaveIssue(issue, OpenDir); err != nil {
		t.Fatal(err)
	}

	// Find by ID
	path, dir, err := FindIssueFile("042")
	if err != nil {
		t.Fatalf("FindIssueFile() error = %v", err)
	}

	if dir != OpenDir {
		t.Errorf("dir = %q, want %q", dir, OpenDir)
	}

	if !fileExists(path) {
		t.Errorf("File not found at path: %s", path)
	}

	// Try to find non-existent issue
	_, _, err = FindIssueFile("999")
	if err == nil {
		t.Error("FindIssueFile() should return error for non-existent issue")
	}
}

func TestDeleteIssue(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Create and save issue
	issue := NewIssue(1, "To Delete", "", nil)
	if err := SaveIssue(issue, OpenDir); err != nil {
		t.Fatal(err)
	}

	// Verify it exists
	_, _, err := FindIssueFile("001")
	if err != nil {
		t.Fatal("Issue should exist before deletion")
	}

	// Delete issue
	err = DeleteIssue("001")
	if err != nil {
		t.Fatalf("DeleteIssue() error = %v", err)
	}

	// Verify it's gone
	_, _, err = FindIssueFile("001")
	if err == nil {
		t.Error("Issue should not exist after deletion")
	}
}

func TestRepoExists(t *testing.T) {
	cleanup := setupTestRepo(t)
	defer cleanup()

	// Should not exist initially
	if RepoExists() {
		t.Error("RepoExists() = true, want false")
	}

	// Initialize repo
	if err := InitializeRepo(); err != nil {
		t.Fatal(err)
	}

	// Should exist now
	if !RepoExists() {
		t.Error("RepoExists() = false, want true")
	}
}

func TestNewIssue(t *testing.T) {
	issue := NewIssue(5, "New Issue", "dave", []string{"bug", "urgent"})

	if issue.ID != "005" {
		t.Errorf("ID = %q, want %q", issue.ID, "005")
	}
	if issue.Title != "New Issue" {
		t.Errorf("Title = %q, want %q", issue.Title, "New Issue")
	}
	if issue.Assignee != "dave" {
		t.Errorf("Assignee = %q, want %q", issue.Assignee, "dave")
	}
	if len(issue.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(issue.Labels))
	}
	if issue.Created.IsZero() {
		t.Error("Created timestamp should not be zero")
	}
	if issue.Updated.IsZero() {
		t.Error("Updated timestamp should not be zero")
	}
}

// Helper functions
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
