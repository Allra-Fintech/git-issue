# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`git-issue` is a lightweight CLI tool written in Go for managing issues as Markdown files within git repositories. The tool provides AI agents and developers direct access to issue context without external integrations.

## Architecture

### Project Structure

```
git-issue/
├── cmd/gi/
│   └── main.go          # Entry point
├── cmd/                 # CLI commands (Cobra-based)
│   ├── root.go          # Root command and global flags
│   ├── init.go          # Initialize .issues/ directory
│   ├── create.go        # Create new issues
│   ├── list.go          # List/filter issues
│   ├── show.go          # Display issue details
│   ├── close.go         # Close issues (with --commit flag)
│   ├── open.go          # Reopen issues (with --commit flag)
│   ├── edit.go          # Edit issues in $EDITOR
│   └── search.go        # Full-text search
├── pkg/
│   ├── issue.go         # Core Issue struct and operations
│   ├── storage.go       # File system operations (.counter, file I/O)
│   └── parser.go        # YAML frontmatter + Markdown parsing
└── .issues/             # Issue storage (NOT in version control for this repo)
    ├── open/            # Open issues: {id}-{slug}.md
    ├── closed/          # Closed issues
    ├── .counter         # Next issue ID
    └── template.md      # Template for new issues
```

### Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML frontmatter parsing
- `github.com/fatih/color` - Terminal colors
- `github.com/olekukonko/tablewriter` - Table output formatting

## Issue Management

This project uses its own issue tracking system (dogfooding).

### Finding Issues

When a user references an issue like "#001" or "issue 001":

1. Search for the file in `.issues/open/` or `.issues/closed/`
2. File naming pattern: `{id}-{slug}.md` (e.g., `001-implement-git-issue-cli-tool.md`)
3. Read the full issue file including YAML frontmatter and Markdown body

Example: For "#001", use glob pattern `.issues/open/001-*.md` or `.issues/closed/001-*.md`

### Issue File Format

```markdown
---
id: "001"
assignee: username
labels: [bug, backend]
created: 2025-11-14T10:30:00Z
updated: 2025-11-14T14:20:00Z
---

# Issue Title

## Description

Full issue description...
```

**Note:** Status is determined by directory location (`.issues/open/` = open, `.issues/closed/` = closed), not by a field in the YAML frontmatter.

## Development Workflow

### Build Commands

```bash
# Build for current platform
go build -o gi ./cmd/gi

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run
```

### Cross-Platform Builds

The Makefile should support:

```bash
make build           # Current platform
make build-all       # macOS (ARM64/AMD64), Linux (AMD64)
make test           # Run all tests
make lint           # Run golangci-lint
```

## Critical Implementation Details

### File Naming

- Pattern: `{id}-{slug}.md`
- ID is zero-padded 3 digits from `.counter` file
- Slug is URL-safe version of title (lowercase, hyphens, no special chars)

### Storage Layer (pkg/storage.go)

- **File moves**: File moves between `open/` and `closed/` use `os.Rename`
- **FindIssueFile**: Must search by ID prefix since filename includes slug which may not match

### Git Integration

- `--commit` flag on `close` and `open` commands auto-commits changes
- Must check if directory is a git repository before git operations
- Never use `--force` or destructive git commands
- Commit message format: "Close issue #001" or "Reopen issue #001"

### Parser (pkg/parser.go)

- YAML frontmatter is delimited by `---`
- Must preserve Markdown body exactly as written
- Slug generation: lowercase, replace spaces/special chars with hyphens
- Status field is NOT in YAML - status is determined by directory location

### Commands

All commands use Cobra. Global flag: `-h, --help`

**Command-specific flags:**

- `create`: `--assignee <name>`, `--label <label>` (repeatable)
- `list`: `--all`, `--assignee <name>`, `--label <label>`, `--status <status>`
- `close/open`: `--commit` (auto-commit to git)
- `search`: `--status <status>`, `--assignee <name>`, `--label <label>`

## Testing Requirements

- Unit tests for storage operations (file I/O, counter increment)
- Unit tests for YAML/Markdown parser
- Integration tests for command workflows
- Error cases: missing files, invalid IDs, non-git directories, malformed YAML
- Target: >80% code coverage
- run manual test by running the command in the `test/` directory which already gitignored

## Cross-Platform Considerations

- File path separators (use `filepath` package)
- Line endings (normalize to LF in files)
- `$EDITOR` environment variable fallback chain (vim as default)
- Test on macOS and Linux (Windows not currently supported)
