---
id: "002"
assignee: ""
labels: [feature, setup]
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T00:00:00Z
---

# Project Setup and Core Infrastructure

**Parent Issue:** #001

## Description

Set up the Go project structure, initialize dependencies, and define core data structures for the git-issue CLI tool.

## Tasks

### Initialize Go Project

- [ ] Create Go module: `go mod init github.com/Allra-Fintech/git-issue`
- [ ] Create `main.go` entry point
- [ ] Set up directory structure:
  - `cmd/` - CLI commands
  - `pkg/issue/` - Core issue management logic
- [ ] Install dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `gopkg.in/yaml.v3` - YAML parsing for frontmatter
  - `github.com/fatih/color` - Terminal colors
  - `github.com/olekukonko/tablewriter` - Table formatting

### Define Core Data Structures

Create `pkg/issue/issue.go` with:

```go
type Issue struct {
    ID       string
    Title    string
    Status   string    // "open" or "closed"
    Assignee string
    Labels   []string
    Created  time.Time
    Updated  time.Time
    Body     string    // Markdown content
}
```

### Create Makefile

- [ ] `make build` - Build for current platform
- [ ] `make build-all` - Cross-compile for macOS (ARM64/AMD64), Linux (AMD64)
- [ ] `make test` - Run tests
- [ ] `make lint` - Run golangci-lint

## Success Criteria

- [ ] Go module initialized
- [ ] All dependencies installed
- [ ] Directory structure created
- [ ] Core `Issue` struct defined
- [ ] Makefile with all targets working
- [ ] Project builds successfully with `go build`

## Dependencies

None - this is the foundation for all other issues.
