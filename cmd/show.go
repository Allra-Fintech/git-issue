package cmd

import (
	"fmt"
	"strings"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [issue-id]",
	Short: "Show detailed information about an issue",
	Long: `Show detailed information about a specific issue.

Examples:
  git-issue show 001
  git-issue show 42`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	// Check if repository is initialized
	if !pkg.RepoExists() {
		return fmt.Errorf(".issues directory not found. Run 'git-issue init' first")
	}

	// Get issue ID
	issueID := args[0]

	// Pad ID if needed (e.g., "1" -> "001")
	if len(issueID) < 3 {
		issueID = fmt.Sprintf("%03s", issueID)
	}

	// Load issue
	issue, dir, err := pkg.LoadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue #%s not found", issueID)
	}

	// Determine status from directory
	status := "open"
	if dir == pkg.ClosedDir {
		status = "closed"
	}

	// Display issue details
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	// Header
	fmt.Printf("%s %s\n", bold("Issue"), bold("#"+issue.ID))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Title
	fmt.Printf("%s %s\n", bold("Title:"), issue.Title)
	fmt.Println()

	// Status
	var statusStr string
	if status == "open" {
		statusStr = green(status)
	} else {
		statusStr = red(status)
	}
	fmt.Printf("%s %s\n", bold("Status:"), statusStr)

	// Assignee
	if issue.Assignee != "" {
		fmt.Printf("%s %s\n", bold("Assignee:"), issue.Assignee)
	}

	// Labels
	if len(issue.Labels) > 0 {
		fmt.Printf("%s %s\n", bold("Labels:"), strings.Join(issue.Labels, ", "))
	}

	// Timestamps
	fmt.Printf("%s %s\n", bold("Created:"), issue.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s %s\n", bold("Updated:"), issue.Updated.Format("2006-01-02 15:04:05"))

	// Body
	if issue.Body != "" {
		fmt.Println()
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println()
		fmt.Println(issue.Body)
	}

	return nil
}
