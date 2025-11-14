---
id: "008"
assignee: ""
labels:
    - build
    - release
    - devops
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T17:13:30.625812+09:00
---

# Build System and Release

**Parent Issue:** #001

## Description

Set up the build system, cross-compilation support, and release automation for distributing git-issue binaries.

## Tasks

### Makefile Implementation

Create comprehensive Makefile with the following targets:

**Build Targets**
- [ ] `make build` - Build for current platform
  ```bash
  go build -o git-issue
  ```

- [ ] `make build-all` - Cross-compile for all supported platforms:
  - macOS ARM64: `git-issue-darwin-arm64`
  - macOS AMD64: `git-issue-darwin-amd64`
  - Linux AMD64: `git-issue-linux-amd64`

  ```bash
  GOOS=darwin GOARCH=arm64 go build -o git-issue-darwin-arm64
  GOOS=darwin GOARCH=amd64 go build -o git-issue-darwin-amd64
  GOOS=linux GOARCH=amd64 go build -o git-issue-linux-amd64
  ```

**Test Targets**
- [ ] `make test` - Run all tests
  ```bash
  go test ./...
  ```

- [ ] `make test-coverage` - Run tests with coverage
  ```bash
  go test -cover ./...
  ```


**Quality Targets**
- [ ] `make lint` - Run golangci-lint
  ```bash
  golangci-lint run
  ```

- [ ] `make fmt` - Format code
  ```bash
  go fmt ./...
  ```

**Utility Targets**
- [ ] `make clean` - Remove build artifacts
- [ ] `make install` - Install to /usr/local/bin (or $GOPATH/bin)
- [ ] `make help` - Show available targets

### GitHub Actions CI/CD (Optional)

**.github/workflows/test.yml**
- [ ] Run tests on push and PR
- [ ] Test on multiple platforms (Ubuntu, macOS)
- [ ] Check code coverage
- [ ] Run linter

**.github/workflows/release.yml**
- [ ] Trigger on git tags (e.g., `v1.0.0`)
- [ ] Build binaries for all platforms
- [ ] Create GitHub release
- [ ] Upload binaries as release assets
- [ ] Generate release notes

### Release Process Documentation

Create `RELEASE.md` with:
- [ ] Version numbering scheme (semantic versioning)
- [ ] Release checklist
- [ ] How to create a release
- [ ] How to test releases

### Build Configuration

**Version Information**
- [ ] Embed version info in binary using build flags:
  ```bash
  go build -ldflags "-X main.version=$(git describe --tags)"
  ```
- [ ] Display version with `git-issue --version`

**Build Optimization**
- [ ] Strip debug info for smaller binaries:
  ```bash
  go build -ldflags "-s -w"
  ```

### Distribution

**Binary Naming**
- Pattern: `git-issue-{os}-{arch}`
- Examples:
  - `git-issue-darwin-arm64`
  - `git-issue-darwin-amd64`
  - `git-issue-linux-amd64`

**Installation Methods**
- [ ] Direct download from GitHub releases
- [ ] `go install github.com/Allra-Fintech/git-issue@latest`
- [ ] Manual build from source

### Documentation

- [ ] Update README.md with release information
- [ ] Document installation methods
- [ ] Document build process for contributors

## Success Criteria

- [ ] Makefile with all targets working
- [ ] Cross-compilation working for all platforms
- [ ] Binaries build successfully
- [ ] Version information embedded in binary
- [ ] GitHub Actions CI configured (if applicable)
- [ ] Release process documented
- [ ] First release created and published

## Testing

- [ ] Test built binaries on macOS ARM64
- [ ] Test built binaries on macOS AMD64
- [ ] Test built binaries on Linux AMD64
- [ ] Verify binary size is reasonable
- [ ] Verify binaries run on fresh systems (no Go installed)

## Dependencies

- Requires #002 (Project Setup)
- Requires #007 (Testing) - tests should pass before release
