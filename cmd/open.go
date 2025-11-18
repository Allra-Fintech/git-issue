package cmd

import (
	"fmt"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/spf13/cobra"
)

var openCommit bool

var openCmd = &cobra.Command{
	Use:   "open <issue-id>",
	Short: "Reopen a closed issue",
	Long:  `Reopen a closed issue by moving it from .issues/closed/ to .issues/open/`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&openCommit, "commit", "c", false, "Auto-commit the change to git")
}

func runOpen(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	// Load the issue to check its status
	_, currentDir, err := pkg.LoadIssue(issueID)
	if err != nil {
		return fmt.Errorf("failed to load issue: %w", err)
	}

	// Check if issue is already open
	if currentDir == pkg.OpenDir {
		return fmt.Errorf("issue #%s is already open", issueID)
	}

	// Move issue from closed to open (timestamp update included)
	if err := pkg.MoveIssue(issueID, pkg.ClosedDir, pkg.OpenDir); err != nil {
		return fmt.Errorf("failed to move issue: %w", err)
	}

	fmt.Printf("✓ Reopened issue #%s\n", issueID)

	// Handle git commit if requested
	if openCommit {
		if err := gitCommitChanges(fmt.Sprintf("Reopen issue #%s", issueID)); err != nil {
			return fmt.Errorf("failed to commit changes: %w", err)
		}
		fmt.Println("✓ Changes committed to git")
	}

	return nil
}
