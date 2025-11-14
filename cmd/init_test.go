package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Allra-Fintech/git-issue/pkg"
)

func TestInitCommand(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "git-issue-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test successful initialization
	t.Run("successful init", func(t *testing.T) {
		err := runInit(nil, []string{})
		if err != nil {
			t.Errorf("runInit() failed: %v", err)
		}

		// Verify directory structure exists
		if !pkg.RepoExists() {
			t.Error(".issues directory was not created")
		}

		// Verify subdirectories
		openPath := filepath.Join(pkg.IssuesDir, pkg.OpenDir)
		if _, err := os.Stat(openPath); os.IsNotExist(err) {
			t.Error("open/ directory was not created")
		}

		closedPath := filepath.Join(pkg.IssuesDir, pkg.ClosedDir)
		if _, err := os.Stat(closedPath); os.IsNotExist(err) {
			t.Error("closed/ directory was not created")
		}

		// Verify counter file
		counterPath := filepath.Join(pkg.IssuesDir, pkg.CounterFile)
		if _, err := os.Stat(counterPath); os.IsNotExist(err) {
			t.Error(".counter file was not created")
		}

		// Verify counter value is 1
		data, err := os.ReadFile(counterPath)
		if err != nil {
			t.Errorf("Failed to read counter file: %v", err)
		}
		if string(data) != "1\n" {
			t.Errorf("Counter file contains %q, expected %q", string(data), "1\n")
		}

		// Verify template file
		templatePath := filepath.Join(pkg.IssuesDir, pkg.TemplateFile)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Error("template.md file was not created")
		}
	})

	// Test init when directory already exists
	t.Run("init already exists", func(t *testing.T) {
		err := runInit(nil, []string{})
		if err == nil {
			t.Error("runInit() should fail when .issues already exists")
		}
	})
}

func TestInitCommandIsolated(t *testing.T) {
	// Test in a fresh temporary directory
	tmpDir, err := os.MkdirTemp("", "git-issue-test-isolated-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Verify .issues doesn't exist initially
	if pkg.RepoExists() {
		t.Fatal(".issues directory should not exist initially")
	}

	// Run init
	err = runInit(nil, []string{})
	if err != nil {
		t.Fatalf("runInit() failed: %v", err)
	}

	// Verify it exists now
	if !pkg.RepoExists() {
		t.Error(".issues directory should exist after init")
	}
}
