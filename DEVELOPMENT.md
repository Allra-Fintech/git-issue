# Development Guide

## Prerequisites

- Go 1.21 or later

## Build

```bash
# Install dependencies (including dev tools)
go mod download

# Build for current platform
make build

# Build for all platforms (macOS ARM64/AMD64, Linux AMD64)
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter (uses golangci-lint as a dev dependency, no global install needed)
make lint

# Format code
make fmt
```

**Note:** `golangci-lint` is tracked as a dev dependency in `tools.go` and doesn't need to be installed globally. The Makefile automatically runs it via `go run`.

## Project Structure

```
git-issue/
├── .github/
│   └── workflows/
│       └── ci.yml       # GitHub Actions CI workflow
├── cmd/
│   ├── root.go          # Root command and global flags
│   ├── init.go          # Initialize command
│   ├── create.go        # Create command
│   ├── list.go          # List command
│   ├── show.go          # Show command
│   ├── close.go         # Close command
│   ├── open.go          # Open command
│   ├── edit.go          # Edit command
│   └── search.go        # Search command
├── pkg/
│   ├── issue.go         # Issue struct and operations
│   ├── storage.go       # File system operations
│   └── parser.go        # Markdown/YAML parsing
├── main.go
├── tools.go             # Dev tool dependencies
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Dependencies

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

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

- Place test files alongside the code they test (e.g., `storage.go` → `storage_test.go`)
- Use table-driven tests where appropriate
- Test both success and error cases
- Ensure proper cleanup in test setup/teardown functions
- Use descriptive test names that explain what is being tested

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write clear comments for exported functions
- Keep functions focused and single-purpose

## Release Process

1. Update version in `cmd/root.go`
2. Create a git tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
3. Push tag: `git push origin v0.1.0`
4. GitHub Actions will automatically build and create a release

## Architecture Notes

### Issue Storage

- Issues are stored as Markdown files with YAML frontmatter
- Status is determined by directory location (`.issues/open/` or `.issues/closed/`)
- File naming pattern: `{id}-{slug}.md`
- Counter file (`.issues/.counter`) tracks the next issue ID

### Parser Design

- YAML frontmatter is delimited by `---`
- Markdown body is preserved exactly as written
- Slug generation: lowercase, spaces/special chars replaced with hyphens

### Git Integration

- Git operations are optional (via `--commit` flag)
- Always check if directory is a git repository before git operations
- Never use `--force` or destructive git commands
- Commit messages follow pattern: "Close issue #001" or "Reopen issue #001"

## Troubleshooting

### Tests Failing

- Ensure you're in the project root directory
- Run `go mod tidy` to ensure dependencies are correct
- Check that temporary directories are being cleaned up properly

### Build Issues

- Verify Go version: `go version` (should be 1.21+)
- Clear build cache: `go clean -cache`
- Update dependencies: `go get -u ./...`

### Linter Errors

- Run `make fmt` to auto-format code
- Run `make lint` to see all linting issues
- Address issues one at a time
