package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	listAll      bool
	listAssignee string
	listLabel    string
	listStatus   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long: `List issues with optional filtering.

By default, only open issues are shown. Use --all to include closed issues.

Examples:
  git-issue list                           # List open issues
  git-issue list --all                     # List all issues
  git-issue list --assignee john           # List issues assigned to john
  git-issue list --label bug               # List issues with 'bug' label
  git-issue list --status closed           # List closed issues`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Include closed issues")
	listCmd.Flags().StringVar(&listAssignee, "assignee", "", "Filter by assignee")
	listCmd.Flags().StringVar(&listLabel, "label", "", "Filter by label")
	listCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (open/closed)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Check if repository is initialized
	if !pkg.RepoExists() {
		return fmt.Errorf(".issues directory not found. Run 'git-issue init' first")
	}

	// Determine which directories to search
	var dirsToSearch []string
	if listStatus != "" {
		// Filter by specific status
		switch listStatus {
		case "open":
			dirsToSearch = []string{pkg.OpenDir}
		case "closed":
			dirsToSearch = []string{pkg.ClosedDir}
		default:
			return fmt.Errorf("invalid status: %s (must be 'open' or 'closed')", listStatus)
		}
	} else if listAll {
		// Show all issues
		dirsToSearch = []string{pkg.OpenDir, pkg.ClosedDir}
	} else {
		// Default: only open issues
		dirsToSearch = []string{pkg.OpenDir}
	}

	// Collect issues from all directories
	type issueWithStatus struct {
		issue  *pkg.Issue
		status string
	}
	var allIssues []issueWithStatus

	for _, dir := range dirsToSearch {
		issues, err := pkg.ListIssues(dir)
		if err != nil {
			// If directory doesn't exist yet, just skip it
			continue
		}

		status := "open"
		if dir == pkg.ClosedDir {
			status = "closed"
		}

		for _, issue := range issues {
			allIssues = append(allIssues, issueWithStatus{
				issue:  issue,
				status: status,
			})
		}
	}

	// Apply filters
	var filteredIssues []issueWithStatus
	for _, item := range allIssues {
		// Filter by assignee
		if listAssignee != "" && item.issue.Assignee != listAssignee {
			continue
		}

		// Filter by label
		if listLabel != "" && !item.issue.HasLabel(listLabel) {
			continue
		}

		filteredIssues = append(filteredIssues, item)
	}

	// Display results
	if len(filteredIssues) == 0 {
		fmt.Println("No issues found.")
		return nil
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Status", "Assignee", "Labels"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	// Add rows
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, item := range filteredIssues {
		issue := item.issue
		status := item.status

		// Color-code status
		var statusStr string
		if status == "open" {
			statusStr = green(status)
		} else {
			statusStr = red(status)
		}

		// Format labels
		labelsStr := "-"
		if len(issue.Labels) > 0 {
			labelsStr = strings.Join(issue.Labels, ", ")
		}

		// Format assignee
		assigneeStr := "-"
		if issue.Assignee != "" {
			assigneeStr = issue.Assignee
		}

		table.Append([]string{
			"#" + issue.ID,
			issue.Title,
			statusStr,
			assigneeStr,
			labelsStr,
		})
	}

	table.Render()

	// Summary
	fmt.Printf("\nTotal: %d issue(s)\n", len(filteredIssues))

	return nil
}
