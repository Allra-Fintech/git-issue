---
id: "005"
assignee: ""
labels:
    - feature
    - cli
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T06:20:58.289536858Z
---

# Implement Issue Management Commands

**Parent Issue:** #001

## Description

Implement commands for managing issue lifecycle: close, open, and edit.

## Tasks

### Close Command (cmd/close.go)

```bash
git-issue close 001
git-issue close 001 --commit
```

- [ ] Accept issue ID as argument
- [ ] Flags:
  - `--commit, -c` - Auto-commit the change
- [ ] Find issue in `.issues/open/`
- [ ] Update issue:
  - Set status to "closed"
  - Update timestamp
- [ ] Move file to `.issues/closed/` using `storage.MoveIssue()`
- [ ] If `--commit` flag:
  - Check if directory is a git repository
  - Stage changes: `git add .issues/`
  - Commit with message: `"Close issue #001"`
- [ ] Display success message

### Open Command (cmd/open.go)

```bash
git-issue open 001
git-issue open 001 --commit
```

- [ ] Accept issue ID as argument
- [ ] Flags:
  - `--commit, -c` - Auto-commit the change
- [ ] Find issue in `.issues/closed/`
- [ ] Update issue:
  - Set status to "open"
  - Update timestamp
- [ ] Move file to `.issues/open/` using `storage.MoveIssue()`
- [ ] If `--commit` flag:
  - Check if directory is a git repository
  - Stage changes: `git add .issues/`
  - Commit with message: `"Reopen issue #001"`
- [ ] Display success message

### Edit Command (cmd/edit.go)

```bash
git-issue edit 001
```

- [ ] Accept issue ID as argument
- [ ] Find issue file using `storage.FindIssueFile(id)`
- [ ] Get editor from `$EDITOR` environment variable (fallback to vim)
- [ ] Open issue file in editor
- [ ] After editing:
  - Parse and validate YAML frontmatter
  - Update timestamp
  - Save changes
- [ ] Display success message
- [ ] Handle errors:
  - Missing issue
  - Invalid YAML after edit
  - Editor not found

## Technical Considerations

- **Git Integration:** Check if current directory is a git repository before git operations
- **Git Safety:** Never use `--force` or destructive operations
- **Editor Support:** Support common editors (vim, nano, emacs, vscode, etc.)
- **Validation:** Ensure YAML frontmatter is valid after editing

## Success Criteria

- [ ] Close command works and moves files correctly
- [ ] Open command works and moves files correctly
- [ ] Edit command opens correct editor
- [ ] `--commit` flag works for both close and open
- [ ] Timestamps are updated correctly
- [ ] Git operations are safe and validated
- [ ] Error handling for non-git directories
- [ ] Error handling for missing issues

## Dependencies

- Requires #002 (Project Setup)
- Requires #003 (Storage and File System Operations)
