---
id: "003"
assignee: ""
labels: [feature, backend]
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T00:00:00Z
---

# Storage and File System Operations

**Parent Issue:** #001

## Description

Implement the storage layer for managing issue files on the file system, including the `.issues/` directory structure, counter management, and file operations.

## Tasks

### Implement Storage Layer (pkg/issue/storage.go)

- [ ] `InitializeRepo()` - Create `.issues/` directory structure:
  - `.issues/open/`
  - `.issues/closed/`
  - `.issues/.counter`
  - `.issues/template.md`
- [ ] `GetNextID()` - Read and atomically increment counter
- [ ] `SaveIssue(issue, dir)` - Write issue to specified directory (open/closed)
- [ ] `LoadIssue(id)` - Read issue from file system by ID (searches both open/ and closed/)
- [ ] `MoveIssue(id, fromDir, toDir)` - Move issue file between open/closed directories
- [ ] `ListIssues(dir)` - Get all issues from a directory (open or closed)
- [ ] `FindIssueFile(id)` - Search for issue file by ID pattern (`{id}-*.md`)
- [ ] `DeleteIssue(id)` - Remove issue file (for cleanup/testing)

### Implement Parser (pkg/issue/parser.go)

- [ ] `ParseMarkdown(content)` - Parse YAML frontmatter + Markdown body into Issue struct
- [ ] `SerializeIssue(issue)` - Convert Issue struct to Markdown with YAML frontmatter
- [ ] `GenerateSlug(title)` - Generate URL-safe slug from title (lowercase, hyphens, no special chars)
- [ ] Handle YAML frontmatter structure:
  ```yaml
  ---
  id: "001"
  assignee: username
  labels: [bug, backend]
  created: 2025-11-14T10:30:00Z
  updated: 2025-11-14T14:20:00Z
  ---
  ```
  **Note:** Status is NOT a field - it's determined by directory location (open/ or closed/)

## Technical Considerations

- **Atomic Operations:** File moves must be atomic (use `os.Rename`)
- **Concurrency:** `.counter` file increments must handle race conditions (use file locking)
- **File Naming:** Pattern `{id}-{slug}.md` where ID is zero-padded 3 digits
- **Error Handling:** Handle missing files, invalid YAML, permission errors

## Success Criteria

- [ ] All storage functions implemented and working
- [ ] YAML frontmatter parsing works correctly
- [ ] File operations are atomic and safe
- [ ] Counter increments correctly with concurrency protection
- [ ] Unit tests for all storage operations
- [ ] Unit tests for parser (YAML + Markdown)

## Dependencies

- Requires #002 (Project Setup) to be completed
