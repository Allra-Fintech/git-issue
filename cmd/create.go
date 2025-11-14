package cmd

import (
	"fmt"
	"strings"

	"github.com/Allra-Fintech/git-issue/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	createAssignee string
	createLabels   []string
)

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new issue",
	Long: `Create a new issue with the specified title.

Examples:
  gi create "Fix authentication bug"
  gi create "Add user profile" --assignee john --label feature --label backend`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createAssignee, "assignee", "", "Assign the issue to a user")
	createCmd.Flags().StringSliceVar(&createLabels, "label", []string{}, "Add labels to the issue (can be specified multiple times)")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Check if repository is initialized
	if !pkg.RepoExists() {
		return fmt.Errorf(".issues directory not found. Run 'gi init' first")
	}

	// Join all args as title (in case title has spaces and wasn't quoted)
	title := strings.TrimSpace(strings.Join(args, " "))

	// Validate title is not empty
	if title == "" {
		return fmt.Errorf("issue title cannot be empty")
	}

	// Get next ID
	id, err := pkg.GetNextID()
	if err != nil {
		return fmt.Errorf("failed to get next issue ID: %w", err)
	}

	// Create new issue
	issue := pkg.NewIssue(id, title, createAssignee, createLabels)

	// Save issue to open directory
	if err := pkg.SaveIssue(issue, pkg.OpenDir); err != nil {
		return fmt.Errorf("failed to save issue: %w", err)
	}

	// Display success message
	green := color.New(color.FgGreen, color.Bold)
	fmt.Print("✓ Created issue ")
	_, _ = green.Printf("#%s", issue.ID)
	fmt.Printf(": %s\n", issue.Title)
	fmt.Println()

	// Display issue details
	fmt.Printf("  ID:       %s\n", issue.ID)
	fmt.Printf("  Title:    %s\n", issue.Title)
	fmt.Printf("  Status:   open\n")
	if issue.Assignee != "" {
		fmt.Printf("  Assignee: %s\n", issue.Assignee)
	}
	if len(issue.Labels) > 0 {
		fmt.Printf("  Labels:   %s\n", strings.Join(issue.Labels, ", "))
	}
	fmt.Printf("  Created:  %s\n", issue.Created.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Show file path
	slug := pkg.GenerateSlug(issue.Title)
	filename := fmt.Sprintf("%s-%s.md", issue.ID, slug)
	fmt.Printf("Issue saved to: .issues/open/%s\n", filename)
	fmt.Printf("Edit the file to add a detailed description.\n")

	return nil
}
