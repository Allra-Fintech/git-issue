package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-issue",
	Short: "A lightweight CLI tool for managing issues as Markdown files",
	Long: `git-issue is a CLI tool for managing issues as Markdown files in your git repository.
It provides AI agents and developers direct access to issue context without external integrations.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
}
