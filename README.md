# Git Issue

A lightweight CLI tool for managing issues as Markdown files in your git repository, giving AI agents and developers direct access to issue context without external integrations.

## Features

- 📝 Create, list, and manage issues as Markdown files
- 🏷️ Support for labels, assignees, and status tracking
- 🔍 Search and filter issues
- 📊 Simple statistics and reporting
- 🤖 AI-friendly format with structured frontmatter
- 🔄 Git-native workflow - all issues version controlled

## Installation

### From Release (Recommended)

Download the latest binary for your platform from the [releases page](https://github.com/Allra-Fintech/git-issue/releases):

```bash
# macOS (ARM)
curl -L https://github.com/Allra-Fintech/git-issue/releases/latest/download/git-issue-darwin-arm64 -o git-issue
chmod +x git-issue
sudo mv git-issue /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/Allra-Fintech/git-issue/releases/latest/download/git-issue-darwin-amd64 -o git-issue
chmod +x git-issue
sudo mv git-issue /usr/local/bin/

# Linux
curl -L https://github.com/Allra-Fintech/git-issue/releases/latest/download/git-issue-linux-amd64 -o git-issue
chmod +x git-issue
sudo mv git-issue /usr/local/bin/

# Windows
# Download git-issue-windows-amd64.exe and add to PATH
```

### From Source

```bash
# Requires Go 1.21 or later
go install github.com/Allra-Fintech/git-issue@latest
```

### Build from Source

```bash
git clone https://github.com/Allra-Fintech/git-issue.git
cd git-issue
go build -o git-issue
```

## Directory Structure

```
.issues/
├── open/
│   ├── 001-user-auth-bug.md
│   └── 002-performance-improvement.md
├── closed/
│   └── 000-initial-setup.md
├── .counter
└── template.md
```

## Usage

### Initialize issue tracking in your repository

```bash
git-issue init
```

This creates the `.issues/` directory structure in your current repository.

### Create a new issue

```bash
git-issue create "Fix Redis connection timeout"
git-issue create "Fix Redis connection timeout" --assignee jonghun --label bug --label backend
```

Interactive mode:

```bash
git-issue create
# Prompts for title, description, assignee, labels
```

### List issues

```bash
# List all open issues
git-issue list

# List all issues including closed
git-issue list --all

# Filter by assignee
git-issue list --assignee jonghun

# Filter by label
git-issue list --label bug

# Combine filters
git-issue list --assignee jonghun --label backend --status open
```

### View an issue

```bash
git-issue show 001
git-issue show 001 --full  # Show with full description and comments
```

### Close an issue

```bash
git-issue close 001
git-issue close 001 --comment "Fixed in commit abc123"
```

### Reopen an issue

```bash
git-issue open 001
```

### Add a comment

```bash
git-issue comment 001 "Found the root cause in Redis session handling"
```

### Edit an issue

```bash
git-issue edit 001
# Opens the issue file in $EDITOR (defaults to vim)
```

### Search issues

```bash
git-issue search "Redis"
git-issue search "authentication" --status open
```

### Update labels or assignee

```bash
git-issue update 001 --assignee jonghun
git-issue update 001 --add-label urgent --remove-label low-priority
```

### Statistics

```bash
git-issue stats
# Shows: Total issues, Open/Closed count, Issues by label, Issues by assignee
```

## AI Integration

This format is designed to be easily readable by AI assistants:

- **Claude/ChatGPT**: Can read issue files directly from the repository
- **GitHub Copilot**: Has context of open issues while coding
- **Custom AI agents**: Can parse YAML frontmatter and Markdown content

Example AI query:

```
"Look at the open issues in .issues/open/ and suggest which one I should work on next based on urgency and my recent commits"
```

## Git Workflow Integration

```bash
# Create issue for current work
git-issue create "Implement user profile API"

# Reference in commits
git commit -m "Add profile endpoint (issue #005)"

# Close when done
git-issue close 005 --comment "Completed in PR #42"
git add .issues/
git commit -m "Close issue #005"
```

## Commands Reference

| Command               | Description                                     |
| --------------------- | ----------------------------------------------- |
| `init`                | Initialize issue tracking in current repository |
| `create [title]`      | Create a new issue                              |
| `list`                | List issues with optional filters               |
| `show <id>`           | Show issue details                              |
| `close <id>`          | Close an issue                                  |
| `open <id>`           | Reopen a closed issue                           |
| `comment <id> <text>` | Add a comment to an issue                       |
| `edit <id>`           | Edit an issue in your editor                    |
| `update <id>`         | Update issue metadata (assignee, labels)        |
| `search <query>`      | Search issues by text                           |
| `stats`               | Show issue statistics                           |

## Global Flags

- `--config <path>` - Path to .issues directory (default: `./.issues`)
- `--no-color` - Disable colored output
- `-h, --help` - Show help for any command

## Command-Specific Options

### create

- `--assignee <name>` - Assign to user
- `--label <label>` - Add label (can be used multiple times)
- `--interactive, -i` - Interactive mode with prompts

### list

- `--assignee <name>` - Filter by assignee
- `--label <label>` - Filter by label
- `--status <status>` - Filter by status (open/closed)
- `--all, -a` - Include closed issues
- `--format <format>` - Output format (table/json/simple)

### show

- `--full, -f` - Show full content including description and comments

### close/open

- `--comment <text>` - Add a comment when changing status

### update

- `--assignee <name>` - Update assignee
- `--add-label <label>` - Add a label
- `--remove-label <label>` - Remove a label

### search

- `--status <status>` - Filter by status
- `--assignee <name>` - Filter by assignee
- `--label <label>` - Filter by label

## Development

### Prerequisites

- Go 1.21 or later

### Build

```bash
# Build for current platform
go build -o git-issue

# Build for all platforms
make build-all

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run
```

### Project Structure

```
git-issue-tracker/
├── cmd/
│   ├── root.go          # Root command and global flags
│   ├── init.go          # Initialize command
│   ├── create.go        # Create command
│   ├── list.go          # List command
│   ├── show.go          # Show command
│   ├── close.go         # Close command
│   ├── open.go          # Open command
│   ├── comment.go       # Comment command
│   ├── edit.go          # Edit command
│   ├── update.go        # Update command
│   ├── search.go        # Search command
│   └── stats.go         # Stats command
├── pkg/
│   ├── issue/
│   │   ├── issue.go     # Issue struct and operations
│   │   ├── storage.go   # File system operations
│   │   └── parser.go    # Markdown/YAML parsing
│   ├── config/
│   │   └── config.go    # Configuration management
│   └── ui/
│       └── format.go    # Output formatting and colors
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Dependencies

```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0
    gopkg.in/yaml.v3 v3.0.1
    github.com/fatih/color v1.16.0
    github.com/olekukonko/tablewriter v0.0.5
)
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Why Git-based Issue Tracking?

- ✅ **Offline-first**: Work without internet connection
- ✅ **Version controlled**: Full history of all changes
- ✅ **No vendor lock-in**: Just Markdown files
- ✅ **AI-friendly**: Direct access for AI agents
- ✅ **Simple**: No database or server required
- ✅ **Portable**: Easy to migrate or backup
- ✅ **Single binary**: No runtime dependencies

## Use Cases

- Small to medium teams (5-30 developers)
- Projects where context is important for AI assistance
- Teams wanting to reduce external dependencies
- Open source projects with simple tracking needs
- Personal projects and side projects
- Temporary issue tracking before migrating to full system

## Comparison with Other Tools

| Feature      | git-issue | GitHub Issues | Jira    | Linear   |
| ------------ | --------- | ------------- | ------- | -------- |
| Offline      | ✅        | ❌            | ❌      | ❌       |
| AI Context   | ✅        | ⚠️            | ⚠️      | ⚠️       |
| Setup Time   | < 1 min   | 5 min         | 30+ min | 10 min   |
| Dependencies | None      | GitHub        | Server  | Internet |
| Cost         | Free      | Free          | $$      | $$       |

## Roadmap

- [ ] Milestone support
- [ ] Issue templates
- [ ] Export to other formats (JSON, CSV)
- [ ] Import from GitHub/Jira/Linear
- [ ] Sync with external systems
- [ ] Web UI for browsing
- [ ] Git hooks integration
- [ ] Issue dependencies/blocking

## License

MIT

## Author

JongHun - Allra Engineering

---

**Note**: This tool is designed to complement, not replace, full-featured issue trackers. For teams already using Jira/Linear/GitHub Issues, consider using this as a synced cache for AI context rather than the source of truth.
Updated with:

Go-specific installation methods (binaries, go install)
Go development setup and project structure
More detailed command options and flags
Build instructions with Makefile
Go dependencies listed
Better examples for Go CLI patterns

Want me to start implementing the actual Go code now?RetryClaude can make mistakes. Please double-check responses.
