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
	fmt.Println("# AI Agent Instructions")
	fmt.Println()
	fmt.Println("## Issue Management")
	fmt.Println()
	fmt.Println("This project uses [gi](https://github.com/Allra-Fintech/git-issue) for managing issues as Markdown files.")
	fmt.Println()
	fmt.Println("### Finding Issues")
	fmt.Println()
	fmt.Println("- **Open issues**: Located in `.issues/open/`")
	fmt.Println("- **Closed issues**: Located in `.issues/closed/`")
	fmt.Println("- **Issue file naming**: `{id}-{title-slug}.md` (e.g., `001-user-auth-bug.md`)")
	fmt.Println()
	fmt.Println("### When a user references an issue")
	fmt.Println()
	fmt.Println("If a user says \"implement #001\" or \"fix issue 001\":")
	fmt.Println()
	fmt.Println("1. Search for the file matching the issue ID in `.issues/open/` or `.issues/closed/`")
	fmt.Println("2. Read the entire issue file to understand requirements")
	fmt.Println("3. Parse the YAML frontmatter for metadata (assignee, labels)")
	fmt.Println("4. Note: Status is determined by directory location (open/ = open, closed/ = closed)")
	fmt.Println("5. Use the issue description and details to guide your implementation")
	fmt.Println()
	fmt.Println("Example: For \"#001\", look for `.issues/open/001-*.md`")
	fmt.Println()
	fmt.Println("### Creating Issues")
	fmt.Println()
	fmt.Println("Always use the `gi create` command instead of manually creating files:")
	fmt.Println()
	fmt.Println("```bash")
	fmt.Println("gi create \"Your issue title\"")
	fmt.Println("gi create \"Fix bug\" --assignee username --label bug")
	fmt.Println("```")
	fmt.Println()
	fmt.Println("**Never** create issue files directly - the command handles ID generation and formatting.")
	fmt.Println()
	fmt.Println("### Working with issues")
	fmt.Println()
	fmt.Println("- Always read the full issue before implementing")
	fmt.Println("- Reference the issue file path in your responses")
	fmt.Println("- Use `gi close <id>` and `gi open <id>` to change issue status")
	fmt.Println("- Maintain the YAML frontmatter structure when editing issues")
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────────────────")
}
