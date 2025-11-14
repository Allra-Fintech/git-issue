package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

// setupCommandTestRepo initializes a temporary repository and returns a cleanup function.
func setupCommandTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "git-issue-cmd-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	if err := pkg.InitializeRepo(); err != nil {
		t.Fatalf("failed to initialize repo: %v", err)
	}

	createAssignee = ""
	createLabels = []string{}

	cleanup := func() {
		_ = os.Chdir(originalDir)
		_ = os.RemoveAll(tmpDir)
		closeCommit = false
		openCommit = false
		createAssignee = ""
		createLabels = []string{}
	}

	return tmpDir, cleanup
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func initGitRepository(t *testing.T, dir string) {
	t.Helper()
	runGitCommand(t, dir, "init")
	runGitCommand(t, dir, "config", "user.email", "tests@example.com")
	runGitCommand(t, dir, "config", "user.name", "git-issue tests")
}

func gitLastCommitMessage(t *testing.T, dir string) string {
	t.Helper()

	cmd := exec.Command("git", "log", "-1", "--pretty=%s")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git log failed: %v\n%s", err, string(output))
	}

	return strings.TrimSpace(string(output))
}
