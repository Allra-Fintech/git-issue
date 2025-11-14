---
id: "001"
assignee: ""
labels: [feature, enhancement]
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T00:00:00Z
---

# Implement git-issue CLI Tool

## Description

Build a lightweight CLI tool for managing issues as Markdown files in git repositories, providing AI agents and developers direct access to issue context without external integrations.

## Requirements

### Phase 1: Project Setup & Core Infrastructure

**1.1 Initialize Go Project**
- Create Go module: `go mod init github.com/Allra-Fintech/git-issue`
- Set up project structure:
  - `main.go` - Entry point
  - `cmd/` - CLI commands using Cobra
  - `pkg/issue/` - Core issue management logic
- Install dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `gopkg.in/yaml.v3` - YAML parsing for frontmatter
  - `github.com/fatih/color` - Terminal colors
  - `github.com/olekukonko/tablewriter` - Table formatting

**1.2 Define Core Data Structures (pkg/issue/issue.go)**
- `Issue` struct with fields:
  - `ID` (string)
  - `Title` (string)
  - `Status` (string: "open" or "closed")
  - `Assignee` (string)
  - `Labels` ([]string)
  - `Created` (time.Time)
  - `Updated` (time.Time)
  - `Body` (string) - Markdown content

### Phase 2: Storage & File System Operations

**2.1 Implement Storage Layer (pkg/issue/storage.go)**
- `InitializeRepo()` - Create `.issues/` directory structure
- `GetNextID()` - Read and increment counter
- `SaveIssue()` - Write issue to appropriate directory
- `LoadIssue(id)` - Read issue from file system
- `MoveIssue(id, status)` - Move between open/closed directories
- `ListIssues(status)` - Get all issues from a directory
- `FindIssueFile(id)` - Search for issue file by ID pattern

**2.2 Implement Parser (pkg/issue/parser.go)**
- `ParseMarkdown(content)` - Parse YAML frontmatter + Markdown body
- `SerializeIssue(issue)` - Convert Issue struct to Markdown with frontmatter
- Slug generation from title (e.g., "Fix Bug" → "fix-bug")

### Phase 3: CLI Commands Implementation

**Required Commands:**
1. `init` - Initialize issue tracking
2. `create [title]` - Create new issue with flags: `--assignee`, `--label`
3. `list` - List issues with flags: `--all`, `--assignee`, `--label`, `--status`
4. `show <id>` - Display issue details
5. `close <id>` - Close issue with flag: `--commit`
6. `open <id>` - Reopen issue with flag: `--commit`
7. `edit <id>` - Edit issue in $EDITOR
8. `search <query>` - Search issues with flags: `--status`, `--assignee`, `--label`

### Phase 4: Git Integration

- Check if current directory is a git repository
- Implement `--commit` flag for close/open commands:
  - `git add .issues/`
  - `git commit -m "message"`
- Error handling for non-git directories

### Phase 5: Testing & Quality

- Unit tests for storage operations
- Unit tests for parser
- Integration tests for full workflows
- Error handling (missing files, invalid IDs, non-git repos)

### Phase 6: Build & Release

**Makefile targets:**
- `make build` - Build for current platform
- `make build-all` - Cross-compile for macOS (ARM64/AMD64), Linux (AMD64)
- `make test` - Run all tests
- `make lint` - Run golangci-lint

## Technical Considerations

1. **File Naming:** Use pattern `{id}-{slug}.md` where slug is URL-safe title
2. **Atomic Operations:** Ensure file moves are atomic
3. **Concurrency:** Handle `.counter` file race conditions
4. **Cross-platform:** Test on macOS and Linux (Windows not currently supported)
5. **Editor Integration:** Support various `$EDITOR` values
6. **Git Safety:** Never force operations, check git status

## Directory Structure

```
.issues/
├── open/
│   └── 001-implement-git-issue-cli-tool.md
├── closed/
├── .counter
└── template.md
```

## Expected Output Format

Issue files should have YAML frontmatter with Markdown body:

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

Issue description here...
```

## Success Criteria

- [ ] All 8 commands implemented and working
- [ ] YAML frontmatter parsing works correctly
- [ ] File system operations are atomic and safe
- [ ] Git integration with --commit flag works
- [ ] Interactive mode for create command works
- [ ] Search and filtering work across all commands
- [ ] Cross-platform builds available
- [ ] Tests pass with >80% coverage
- [ ] README documentation is complete

## Implementation Order

1. **Critical Path:** Project setup, core data structures, storage layer, parser, commands (init, create, list, show)
2. **Secondary:** Commands (close, open, edit), search functionality, git integration
3. **Polish:** Filtering options, table formatting, error handling, tests
4. **Release:** Build system, CI/CD, documentation
