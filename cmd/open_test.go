package cmd

import (
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func forceMoveIssueToClosed(t *testing.T, id string) *pkg.Issue {
	t.Helper()

	issue, dir, err := pkg.LoadIssue(id)
	if err != nil {
		t.Fatalf("failed to load issue %s: %v", id, err)
	}

	if dir == pkg.ClosedDir {
		return issue
	}

	if err := pkg.MoveIssue(id, pkg.OpenDir, pkg.ClosedDir); err != nil {
		t.Fatalf("failed to move issue %s to closed dir: %v", id, err)
	}

	updated, dirAfter, err := pkg.LoadIssue(id)
	if err != nil {
		t.Fatalf("failed to reload issue %s: %v", id, err)
	}
	if dirAfter != pkg.ClosedDir {
		t.Fatalf("issue %s should be in closed dir, got %s", id, dirAfter)
	}

	return updated
}

func TestRunOpenMovesIssueToOpen(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Reopen flow works"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	issueBefore := forceMoveIssueToClosed(t, "001")

	if err := runOpen(nil, []string{"001"}); err != nil {
		t.Fatalf("runOpen() failed: %v", err)
	}

	issueAfter, dir, err := pkg.LoadIssue("001")
	if err != nil {
		t.Fatalf("failed to reload issue: %v", err)
	}
	if dir != pkg.OpenDir {
		t.Fatalf("issue should be moved to open dir, got %s", dir)
	}
	if issueAfter.Updated.Before(issueBefore.Updated) {
		t.Errorf("issue updated timestamp was not refreshed")
	}
}

func TestRunOpenAlreadyOpen(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Stay open"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	err := runOpen(nil, []string{"001"})
	if err == nil {
		t.Fatal("runOpen() should fail when issue already open")
	}
	if !strings.Contains(err.Error(), "already open") {
		t.Fatalf("expected already open error, got %v", err)
	}
}

func TestRunOpenMissingIssue(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	err := runOpen(nil, []string{"001"})
	if err == nil {
		t.Fatal("runOpen() should fail when issue is missing")
	}
	if !strings.Contains(err.Error(), "failed to load issue") {
		t.Fatalf("expected load error, got %v", err)
	}
}

func TestRunOpenCommitWithoutGitRepo(t *testing.T) {
	_, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	if err := runCreate(nil, []string{"Need git reopen"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	forceMoveIssueToClosed(t, "001")

	openCommit = true
	defer func() { openCommit = false }()

	err := runOpen(nil, []string{"001"})
	if err == nil {
		t.Fatal("runOpen() should fail when --commit is used outside git repo")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("expected git repository error, got %v", err)
	}
}

func TestRunOpenCommitCreatesGitCommit(t *testing.T) {
	repoDir, cleanup := setupCommandTestRepo(t)
	defer cleanup()

	initGitRepository(t, repoDir)

	if err := runCreate(nil, []string{"Git reopen"}); err != nil {
		t.Fatalf("runCreate() failed: %v", err)
	}

	forceMoveIssueToClosed(t, "001")

	openCommit = true
	defer func() { openCommit = false }()

	if err := runOpen(nil, []string{"001"}); err != nil {
		t.Fatalf("runOpen() failed: %v", err)
	}

	lastMessage := gitLastCommitMessage(t, repoDir)
	if lastMessage != "Reopen issue #001" {
		t.Fatalf("unexpected commit message %q", lastMessage)
	}
}
