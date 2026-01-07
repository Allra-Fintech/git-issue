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
		return fmt.Errorf(".issues directory already exists. Use 'gi list' to see existing issues")
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
	fmt.Println("You can now create issues with 'gi create <title>'")
	fmt.Println()

	// Display AI agent instruction guidance
	printAIAgentInstructions()

	return nil
}

func printAIAgentInstructions() {
	fmt.Println("💡 AI Agent Integration")
	fmt.Println()
	fmt.Println("To help AI coding agents (Claude Code, Cursor, etc.) work with your issues,")
	fmt.Println("add the following instructions to your CLAUDE.md, AGENTS.md, or .cursorrules file:")
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("## Issue Management")
	fmt.Println()
	fmt.Println("This project uses git-issue for managing issues as Markdown files.")
	fmt.Println()
	fmt.Println("### Finding Issues")
	fmt.Println()
	fmt.Println("When a user references an issue like \"#001\" or \"issue 001\":")
	fmt.Println()
	fmt.Println("1. Search for the file in `.issues/open/` or `.issues/closed/`")
	fmt.Println("2. File naming pattern: `{id}-{slug}.md` (e.g., `001-implement-feature.md`)")
	fmt.Println("3. Use glob pattern: `.issues/open/001-*.md` or `.issues/closed/001-*.md`")
	fmt.Println("4. Read the full issue file including YAML frontmatter and Markdown body")
	fmt.Println()
	fmt.Println("### Issue File Format")
	fmt.Println()
	fmt.Println("```markdown")
	fmt.Println("---")
	fmt.Println("id: \"001\"")
	fmt.Println("assignee: username")
	fmt.Println("labels: [bug, backend]")
	fmt.Println("created: 2025-11-14T10:30:00Z")
	fmt.Println("updated: 2025-11-14T14:20:00Z")
	fmt.Println("---")
	fmt.Println()
	fmt.Println("# Issue Title")
	fmt.Println()
	fmt.Println("## Description")
	fmt.Println()
	fmt.Println("Full issue description...")
	fmt.Println("```")
	fmt.Println()
	fmt.Println("**Important:** Status is determined by directory location (`.issues/open/`")
	fmt.Println("= open, `.issues/closed/` = closed), NOT by a field in the YAML frontmatter.")
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("For more details, see: https://github.com/Allra-Fintech/git-issue")
}
