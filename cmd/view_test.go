package cmd

import (
	"testing"
)

func TestViewMissingArgs(t *testing.T) {
	cmd := viewCmd
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error for missing arguments")
	}
}

func TestViewNonExistentIssue(t *testing.T) {
	// runView should fail when the issue doesn't exist
	err := runView(viewCmd, []string{"999"})
	if err == nil {
		t.Error("expected error for non-existent issue")
	}
}
