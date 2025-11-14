package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <issue-id>",
	Short: "Edit an issue in your editor",
	Long:  `Edit an issue by opening it in your configured editor ($EDITOR, defaults to vim)`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	// Find the issue file
	path, dir, err := pkg.FindIssueFile(issueID)
	if err != nil {
		return fmt.Errorf("failed to find issue: %w", err)
	}

	// Get editor from environment, default to vim
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = pkg.DefaultEditor
	}

	// Open editor
	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Read and validate the edited file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read edited file: %w", err)
	}

	// Parse to validate YAML frontmatter
	issue, err := pkg.ParseMarkdown(string(data))
	if err != nil {
		return fmt.Errorf("invalid issue format after editing: %w", err)
	}

	// Update timestamp
	issue.Updated = time.Now()

	// Save the updated issue
	if err := pkg.SaveIssue(issue, dir); err != nil {
		return fmt.Errorf("failed to save issue: %w", err)
	}

	fmt.Printf("✓ Updated issue #%s\n", issueID)

	return nil
}
