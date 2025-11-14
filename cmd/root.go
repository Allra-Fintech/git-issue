package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "git-issue",
	Short: "A lightweight CLI tool for managing issues as Markdown files",
	Long: `git-issue is a CLI tool for managing issues as Markdown files in your git repository.
It provides AI agents and developers direct access to issue context without external integrations.`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	setVersionTemplate()

	// Enable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

// SetVersion overrides the default CLI version (useful for ldflags injection).
func SetVersion(v string) {
	if v == "" {
		return
	}
	version = v
	rootCmd.Version = version
	setVersionTemplate()
}

func setVersionTemplate() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("git-issue version %s\n", version))
}
