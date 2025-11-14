package cmd

import (
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
