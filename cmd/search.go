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
	searchStatus   string
	searchAssignee string
	searchLabel    string
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search issues by text",
	Long: `Search for issues by text in title and body.

The search is case-insensitive and searches both the issue title and description.

Examples:
  gi search "Redis"
  gi search "authentication" --status open
  gi search "bug" --label backend --assignee john`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchStatus, "status", "", "Filter by status (open/closed)")
	searchCmd.Flags().StringVar(&searchAssignee, "assignee", "", "Filter by assignee")
	searchCmd.Flags().StringVar(&searchLabel, "label", "", "Filter by label")
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Check if repository is initialized
	if !pkg.RepoExists() {
		return fmt.Errorf(".issues directory not found. Run 'gi init' first")
	}

	// Join all args as search query (in case query has spaces and wasn't quoted)
	query := strings.TrimSpace(strings.Join(args, " "))

	// Validate query is not empty
	if query == "" {
		return fmt.Errorf("search query cannot be empty")
	}

	// Convert query to lowercase for case-insensitive search
	queryLower := strings.ToLower(query)

	// Determine which directories to search
	var dirsToSearch []string
	if searchStatus != "" {
		// Filter by specific status
		switch searchStatus {
		case "open":
			dirsToSearch = []string{pkg.OpenDir}
		case "closed":
			dirsToSearch = []string{pkg.ClosedDir}
		default:
			return fmt.Errorf("invalid status: %s (must be 'open' or 'closed')", searchStatus)
		}
	} else {
		// Default: search all issues
		dirsToSearch = []string{pkg.OpenDir, pkg.ClosedDir}
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

	// Search and filter issues
	var matchedIssues []issueWithStatus
	for _, item := range allIssues {
		// Search in title and body (case-insensitive)
		titleLower := strings.ToLower(item.issue.Title)
		bodyLower := strings.ToLower(item.issue.Body)

		if !strings.Contains(titleLower, queryLower) && !strings.Contains(bodyLower, queryLower) {
			continue
		}

		// Filter by assignee
		if searchAssignee != "" && item.issue.Assignee != searchAssignee {
			continue
		}

		// Filter by label
		if searchLabel != "" && !item.issue.HasLabel(searchLabel) {
			continue
		}

		matchedIssues = append(matchedIssues, item)
	}

	// Display results
	if len(matchedIssues) == 0 {
		fmt.Printf("No issues found matching '%s'.\n", query)
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

	for _, item := range matchedIssues {
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
	fmt.Printf("\nFound %d issue(s) matching '%s'\n", len(matchedIssues), query)

	return nil
}
