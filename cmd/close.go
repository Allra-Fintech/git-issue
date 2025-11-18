package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/spf13/cobra"
)

var closeCommit bool

var closeCmd = &cobra.Command{
	Use:   "close <issue-id>",
	Short: "Close an issue",
	Long:  `Close an issue by moving it from .issues/open/ to .issues/closed/`,
	Args:  cobra.ExactArgs(1),
	RunE:  runClose,
}

func init() {
	rootCmd.AddCommand(closeCmd)
	closeCmd.Flags().BoolVarP(&closeCommit, "commit", "c", false, "Auto-commit the change to git")
}

func runClose(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	// Load the issue to check its status
	_, currentDir, err := pkg.LoadIssue(issueID)
	if err != nil {
		return fmt.Errorf("failed to load issue: %w", err)
	}

	// Check if issue is already closed
	if currentDir == pkg.ClosedDir {
		return fmt.Errorf("issue #%s is already closed", issueID)
	}

	// Move issue from open to closed (timestamp update included)
	if err := pkg.MoveIssue(issueID, pkg.OpenDir, pkg.ClosedDir); err != nil {
		return fmt.Errorf("failed to move issue: %w", err)
	}

	fmt.Printf("✓ Closed issue #%s\n", issueID)

	// Handle git commit if requested
	if closeCommit {
		if err := gitCommitChanges(fmt.Sprintf("Close issue #%s", issueID)); err != nil {
			return fmt.Errorf("failed to commit changes: %w", err)
		}
		fmt.Println("✓ Changes committed to git")
	}

	return nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run() == nil
}

// gitCommitChanges stages and commits changes to .issues/
func gitCommitChanges(message string) error {
	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	// Stage changes
	stageCmd := exec.Command("git", "add", pkg.GetIssuesPath())
	stageCmd.Stdout = os.Stdout
	stageCmd.Stderr = os.Stderr
	if err := stageCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Commit changes
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
