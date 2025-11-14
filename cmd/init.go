package cmd

import (
	"fmt"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the .issues directory structure",
	Long: `Initialize the .issues directory with the following structure:
  .issues/
  ├── open/       # Open issues
  ├── closed/     # Closed issues
  ├── .counter    # Issue ID counter
  └── template.md # Template for new issues`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if .issues already exists
	if pkg.RepoExists() {
		return fmt.Errorf(".issues directory already exists. Use 'git-issue list' to see existing issues")
	}

	// Initialize the repository
	if err := pkg.InitializeRepo(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Display success message
	fmt.Println("✓ Initialized .issues directory structure:")
	fmt.Println()
	fmt.Println("  .issues/")
	fmt.Println("  ├── open/       # Open issues")
	fmt.Println("  ├── closed/     # Closed issues")
	fmt.Println("  ├── .counter    # Issue ID counter (initialized to 1)")
	fmt.Println("  └── template.md # Template for new issues")
	fmt.Println()
	fmt.Println("You can now create issues with 'git-issue create <title>'")

	return nil
}
