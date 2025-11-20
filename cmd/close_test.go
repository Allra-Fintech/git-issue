package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func TestRunCloseMovesIssueToClosed(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Close flow works"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	issueBefore, dir, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to load issue: %v", err)
	}
	if dir != pkg.OpenDir {
		t.Fatalf("issue should start in open dir, got %s", dir)
	}

	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("runClose() failed: %v", err)
	}

	issueAfter, dirAfter, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to reload issue: %v", err)
	}

	if dirAfter != pkg.ClosedDir {
		t.Fatalf("issue should be moved to closed dir, got %s", dirAfter)
	}

	if issueAfter.Updated.Before(issueBefore.Updated) {
		t.Errorf("issue updated timestamp was not refreshed")
	}
}

func TestRunCloseAlreadyClosed(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Close twice"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("first close should succeed: %v", err)
	}

	err := runClose(nil, []string{"001"})
	if err == nil {
		t.Fatal("runClose() should fail when issue already closed")
	}
	if !strings.Contains(err.Error(), "already closed") {
		t.Fatalf("expected already closed error, got %v", err)
	}
}

func TestRunCloseMissingIssue(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	err := runClose(nil, []string{"001"})
	if err == nil {
		t.Fatal("runClose() should fail when issue is missing")
	}
	if !strings.Contains(err.Error(), "failed to load issue") {
		t.Fatalf("expected load error, got %v", err)
	}
}

func TestRunCloseCommitWithoutGitRepo(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Needs git"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	closeCommit = true
	defer func() { closeCommit = false }()

	err := runClose(nil, []string{"001"})
	if err == nil {
		t.Fatal("runClose() should fail when --commit used outside git repo")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("expected git repository error, got %v", err)
	}
}

func TestRunCloseCommitCreatesGitCommit(t *testing.T) {
	repoDir, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	initGitRepository(t, repoDir)

	if err := runCreate(nil, []string{"Git integration"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	closeCommit = true
	defer func() { closeCommit = false }()

	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("runClose() failed: %v", err)
	}

	lastMessage := gitLastCommitMessage(t, repoDir)
	if lastMessage != "Close issue #001" {
		t.Fatalf("unexpected commit message %q", lastMessage)
	}
}

func TestRunClosePreservesFilenameWithKoreanTitle(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	// Create issue with English title
	if err := runCreate(nil, []string{"Initial English Title"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	// Get original filename
	originalPath, _, err := pkg.FindIssueFile("001")
	if err != nil {
		t.Fatalf("failed to find issue: %v", err)
	}

	// Edit issue to have Korean title
	issue, _, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to load issue: %v", err)
	}
	issue.Title = "한글 제목으로 변경"
	content, err := pkg.SerializeIssue(issue)
	if err != nil {
		t.Fatalf("failed to serialize issue: %v", err)
	}
	if err := os.WriteFile(originalPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write issue: %v", err)
	}

	// Close the issue
	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("runClose() failed: %v", err)
	}

	// Verify issue is in closed directory with original filename
	closedPath, dir, err := pkg.FindIssueFile("001")
	if err != nil {
		t.Fatalf("failed to find closed issue: %v", err)
	}
	if dir != pkg.ClosedDir {
		t.Fatalf("issue should be in closed dir, got %s", dir)
	}

	// Verify filename was preserved (contains "initial-english-title")
	if !strings.Contains(closedPath, "initial-english-title") {
		t.Errorf("filename should be preserved, got %s", closedPath)
	}

	// Verify no malformed files like "001-.md" were created
	openFiles, _ := pkg.ListIssues(pkg.OpenDir)
	closedFiles, _ := pkg.ListIssues(pkg.ClosedDir)

	if len(openFiles) != 0 {
		t.Errorf("open directory should be empty, found %d files", len(openFiles))
	}
	if len(closedFiles) != 1 {
		t.Errorf("closed directory should have exactly 1 file, found %d", len(closedFiles))
	}
}

func TestRunCloseAfterEditDoesNotCreateNewIssueFile(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	// Create issue
	if err := runCreate(nil, []string{"Original Title"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	// Simulate editing the issue to change the title (mimics `gi edit`)
	issue, dir, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to load issue: %v", err)
	}
	issue.Title = "Edited Title That Changes The Slug"
	if err := pkg.SaveIssue(issue, dir); err != nil {
		t.Fatalf("failed to save edited issue: %v", err)
	}

	// Close the issue
	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("runClose() failed: %v", err)
	}

	// Verify no new files were created in open/
	openFiles, _ := pkg.ListIssues(pkg.OpenDir)
	if len(openFiles) != 0 {
		t.Fatalf("open directory should be empty after closing, found %d issues", len(openFiles))
	}

	// Verify closed file uses the original filename
	closedPath, dirAfter, err := pkg.FindIssueFile("001")
	if err != nil {
		t.Fatalf("failed to find closed issue: %v", err)
	}
	if dirAfter != pkg.ClosedDir {
		t.Fatalf("issue should be in closed dir, got %s", dirAfter)
	}
	if !strings.Contains(closedPath, "original-title") {
		t.Fatalf("expected filename to retain original slug, got %s", closedPath)
	}
}

func TestRunClosePreservesFilenameWhenTitleModified(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	// Create issue
	if err := runCreate(nil, []string{"Original Title"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	// Get original filename
	originalPath, _, err := pkg.FindIssueFile("001")
	if err != nil {
		t.Fatalf("failed to find issue: %v", err)
	}

	// Edit issue to have completely different title
	issue, _, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to load issue: %v", err)
	}
	issue.Title = "Completely Different Modified Title"
	content, err := pkg.SerializeIssue(issue)
	if err != nil {
		t.Fatalf("failed to serialize issue: %v", err)
	}
	if err := os.WriteFile(originalPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write issue: %v", err)
	}

	// Close the issue
	if err := runClose(nil, []string{"001"}); err != nil {
		t.Fatalf("runClose() failed: %v", err)
	}

	// Verify issue is in closed directory with original filename
	closedPath, dir, err := pkg.FindIssueFile("001")
	if err != nil {
		t.Fatalf("failed to find closed issue: %v", err)
	}
	if dir != pkg.ClosedDir {
		t.Fatalf("issue should be in closed dir, got %s", dir)
	}

	// Verify filename was preserved (contains "original-title", not "completely-different")
	if !strings.Contains(closedPath, "original-title") {
		t.Errorf("original filename should be preserved, got %s", closedPath)
	}
	if strings.Contains(closedPath, "completely-different") {
		t.Errorf("filename should not change to modified title, got %s", closedPath)
	}

	// Verify no duplicate files were created
	openFiles, _ := pkg.ListIssues(pkg.OpenDir)
	closedFiles, _ := pkg.ListIssues(pkg.ClosedDir)

	if len(openFiles) != 0 {
		t.Errorf("open directory should be empty, found %d files", len(openFiles))
	}
	if len(closedFiles) != 1 {
		t.Errorf("closed directory should have exactly 1 file, found %d", len(closedFiles))
	}
}
