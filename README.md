# git-issue

A lightweight CLI tool for managing issues as Markdown files in your git repository, giving AI agents and developers direct access to issue context without external integrations.

## Features

- 📝 Create, list, and manage issues as Markdown files
- 🏷️ Support for labels and assignees
- 🔍 Search and filter issues (status determined by directory: open/ or closed/)
- 🤖 AI-friendly format with structured frontmatter
- 🔄 Git-native workflow - all issues version controlled

## Usage

### Initialize issue tracking in your repository

```bash
git-issue init
```

This creates the `.issues/` directory structure in your current repository:

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

### Create a new issue

```bash
git-issue create "Fix Redis connection timeout"
git-issue create "Fix Redis connection timeout" --assignee jonghun --label bug --label backend
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
```

### Close an issue

```bash
git-issue close 001
```

### Reopen an issue

```bash
git-issue open 001
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

### Makefile Targets

The repository ships with a Makefile that wraps the common developer workflows:

| Target | Description |
| --- | --- |
| `make build` | Build the CLI for the current platform with embedded version info |
| `make build-all` | Cross-compile binaries for macOS (arm64/amd64) and Linux (amd64) |
| `make test` / `make test-coverage` | Run the test suite (optionally with coverage) |
| `make lint` | Run `golangci-lint` against the codebase |
| `make fmt` | Format the tree with `go fmt ./...` |
| `make clean` | Remove build artifacts and coverage reports |

Override `VERSION` to embed a specific version string:

```bash
VERSION=v1.2.3 make build
./git-issue --version
```

### Release Process

See [RELEASE.md](RELEASE.md) for the full release checklist, tagging instructions, and binary verification steps.

### Shell Completion

#### Zsh (macOS)
```bash
git-issue completion zsh > $(brew --prefix)/share/zsh/site-functions/_git-issue
```

#### Bash

**Linux:**
```bash
git-issue completion bash > /etc/bash_completion.d/git-issue
```

**macOS:**
```bash
git-issue completion bash > $(brew --prefix)/etc/bash_completion.d/git-issue
```

For other shells or custom setups, run `git-issue completion --help`.

## AI Integration

This format is designed to be easily readable by AI agents:

- **Claude/ChatGPT**: Can read issue files directly from the repository
- **GitHub Copilot**: Has context of open issues while coding
- **Custom AI agents**: Can parse YAML frontmatter and Markdown content

Example AI queries:

**Planning your work:**

```
"Look at the open issues in .issues/open/ and suggest which one I should work on next based on urgency and my recent commits"
```

**Getting implementation guidance:**

```
"Read issue .issues/open/003-add-user-authentication.md and provide a detailed implementation plan with:
1. Required dependencies and packages
2. Step-by-step implementation guide
3. Security best practices to follow
4. Test cases to cover
5. Potential edge cases to handle"
```

**Code review with context:**

```
"Review my changes in src/auth.js against issue .issues/open/003-add-user-authentication.md and check if all requirements are met"
```

### Setting up AI Agent Instructions

For optimal AI agent integration, create instruction files in your repository root to teach agents how to work with your issues:

**AGENTS.md** or **CLAUDE.md**:

```markdown
# AI Agent Instructions

## Issue Management

This project uses git-issue for managing issues as Markdown files.

### Finding Issues

- **Open issues**: Located in `.issues/open/`
- **Closed issues**: Located in `.issues/closed/`
- **Issue file naming**: `{id}-{title-slug}.md` (e.g., `001-user-auth-bug.md`)

### When a user references an issue

If a user says "implement #001" or "fix issue 001":

1. Search for the file matching the issue ID in `.issues/open/` or `.issues/closed/`
2. Read the entire issue file to understand requirements
3. Parse the YAML frontmatter for metadata (assignee, labels)
4. Note: Status is determined by directory location (open/ = open, closed/ = closed)
5. Use the issue description and details to guide your implementation

Example: For "#001", look for `.issues/open/001-*.md`

### Working with issues

- Always read the full issue before implementing
- Reference the issue file path in your responses
- Status is determined by directory: move files between open/ and closed/ to change status
- Maintain the YAML frontmatter structure when editing issues
```

**Example workflow:**

```bash
# User: "Give me a plan to implement #001"
# AI agent will:
# 1. Find .issues/open/001-*.md
# 2. Read the issue content
# 3. Provide implementation plan based on issue requirements
```

## Git Workflow Integration

```bash
# Create issue for current work
git-issue create "Implement user profile API"

# Work on the implementation
git commit -m "Add profile endpoint (issue #005)"

# Close issue and automatically commit the change
git-issue close 005 --commit

# Or manually manage the commit
git-issue close 005
git add .issues/
git commit -m "Close issue #005"

# Reopen issue with automatic commit
git-issue open 005 --commit
```

## Commands Reference

| Command          | Description                                     |
| ---------------- | ----------------------------------------------- |
| `init`           | Initialize issue tracking in current repository |
| `create [title]` | Create a new issue                              |
| `list`           | List issues with optional filters               |
| `show <id>`      | Show issue details                              |
| `close <id>`     | Close an issue                                  |
| `open <id>`      | Reopen a closed issue                           |
| `edit <id>`      | Edit an issue in your editor                    |
| `search <query>` | Search issues by text                           |

## Global Flags

- `-h, --help` - Show help for any command

## Command-Specific Options

### create

- `--assignee <name>` - Assign to user
- `--label <label>` - Add label (can be used multiple times)

### list

- `--assignee <name>` - Filter by assignee
- `--label <label>` - Filter by label
- `--status <status>` - Filter by status (open/closed)
- `--all, -a` - Include closed issues

### close/open

- `--commit, -c` - Commit the change to git

### search

- `--status <status>` - Filter by status
- `--assignee <name>` - Filter by assignee
- `--label <label>` - Filter by label

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed development guidelines, build instructions, and contribution workflow.

## Why Git-based Issue Tracking?

- ✅ **Offline-first**: Work without internet connection
- ✅ **Version controlled**: Full history of all changes
- ✅ **No vendor lock-in**: Just Markdown files
- ✅ **AI-friendly**: Direct access for AI agents
- ✅ **Simple**: No database or server required
- ✅ **Portable**: Easy to migrate or backup
- ✅ **Single binary**: No runtime dependencies

## Comparison with Other Tools

| Feature      | git-issue | GitHub Issues | Jira    | Linear   |
| ------------ | --------- | ------------- | ------- | -------- |
| Offline      | ✅        | ❌            | ❌      | ❌       |
| AI Context   | ✅        | ⚠️            | ⚠️      | ⚠️       |
| Setup Time   | < 1 min   | 5 min         | 30+ min | 10 min   |
| Dependencies | None      | GitHub        | Server  | Internet |
| Cost         | Free      | Free          | $$      | $$       |

## License

MIT

## Author

[Allra fintech](https://github.com/Allra-Fintech)

---

**Note**: This tool is designed to complement, not replace, full-featured issue trackers. For teams already using Jira/Linear/GitHub Issues, consider using this as a synced cache for AI context rather than the source of truth.
