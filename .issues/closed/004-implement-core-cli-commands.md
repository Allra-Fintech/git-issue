---
id: "004"
assignee: ""
labels: [feature, cli]
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T00:00:00Z
---

# Implement Core CLI Commands

**Parent Issue:** #001

## Description

Implement the core CLI commands using Cobra framework: init, create, list, and show.

## Tasks

### Root Command (cmd/root.go)

- [ ] Set up Cobra root command
- [ ] Global flags: `-h, --help`
- [ ] Version information
- [ ] Register all subcommands

### Init Command (cmd/init.go)

```bash
git-issue init
```

- [ ] Check if `.issues/` already exists (warn user)
- [ ] Call `storage.InitializeRepo()`
- [ ] Create directory structure
- [ ] Initialize `.counter` file with value "1"
- [ ] Create default `template.md`
- [ ] Display success message with directory structure

### Create Command (cmd/create.go)

```bash
git-issue create [title]
git-issue create "Fix bug" --assignee jonghun --label bug --label backend
```

- [ ] Accept title as required argument
- [ ] Flags:
  - `--assignee <name>` - Set assignee
  - `--label <label>` - Add labels (repeatable)
- [ ] Generate next ID using `storage.GetNextID()`
- [ ] Create issue file: `.issues/open/{id}-{slug}.md`
- [ ] Set timestamps (created, updated)
- [ ] Display created issue summary
- [ ] User can edit the created file manually to add description

### List Command (cmd/list.go)

```bash
git-issue list
git-issue list --all
git-issue list --assignee jonghun --label bug
```

- [ ] Flags:
  - `--all, -a` - Include closed issues
  - `--assignee <name>` - Filter by assignee
  - `--label <label>` - Filter by label
  - `--status <status>` - Filter by status (open/closed)
- [ ] Load issues using `storage.ListIssues()`
- [ ] Apply filters
- [ ] Display issues in table format using `tablewriter`:
  - ID
  - Title
  - Assignee
  - Labels
  - Status
- [ ] Color coding: green for open, red for closed

### Show Command (cmd/show.go)

```bash
git-issue show 001
```

- [ ] Accept issue ID as argument
- [ ] Find issue using `storage.FindIssueFile(id)`
- [ ] Load and parse issue
- [ ] Display full issue content:
  - Metadata (ID, status, assignee, labels, timestamps)
  - Title
  - Full description/body
- [ ] Handle missing issue error

## Success Criteria

- [ ] All 4 commands implemented and working
- [ ] Filters work correctly on list command
- [ ] Table output is properly formatted
- [ ] Color output works on supported terminals
- [ ] Error handling for missing files, invalid input
- [ ] Help messages are clear and helpful

## Dependencies

- Requires #002 (Project Setup)
- Requires #003 (Storage and File System Operations)
