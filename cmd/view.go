package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <issue-id>",
	Short: "Open an issue in the default program",
	Long:  `Open an issue's markdown file in the system's default program (e.g., Typora, VS Code, Obsidian).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func runView(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	// Find the issue file
	path, _, err := pkg.FindIssueFile(issueID)
	if err != nil {
		return fmt.Errorf("failed to find issue: %w", err)
	}

	// Determine the OS-appropriate open command
	var opener string
	switch runtime.GOOS {
	case "darwin":
		opener = "open"
	case "linux":
		opener = "xdg-open"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Open the file in the default program (non-blocking)
	if err := exec.Command(opener, path).Start(); err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	color.New(color.FgGreen).Printf("✓ Opened issue #%s in default program\n", issueID)
	return nil
}
