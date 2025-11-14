---
id: "006"
assignee: ""
labels:
    - feature
    - cli
created: 2025-11-14T00:00:00Z
updated: 2025-11-14T16:37:07.481391+09:00
---

# Implement Search Command

**Parent Issue:** #001

## Description

Implement full-text search functionality to find issues by content.

## Tasks

### Search Command (cmd/search.go)

```bash
git-issue search "Redis"
git-issue search "authentication" --status open --label bug
```

- [ ] Accept search query as argument
- [ ] Flags:
  - `--status <status>` - Filter by status (open/closed)
  - `--assignee <name>` - Filter by assignee
  - `--label <label>` - Filter by label
- [ ] Search in:
  - Issue title
  - Issue body/description
- [ ] Load all issues (open and/or closed based on status filter)
- [ ] Perform case-insensitive text search
- [ ] Apply additional filters (assignee, label)
- [ ] Display matching issues in table format:
  - ID
  - Title (with search term highlighted if possible)
  - Assignee
  - Labels
  - Status
- [ ] Show count of matches
- [ ] Handle no results gracefully

## Implementation Details

### Search Algorithm

```go
func searchIssues(query string, filters Filters) []Issue {
    // 1. Load issues based on status filter
    // 2. For each issue:
    //    - Check if query appears in title or body (case-insensitive)
    //    - Apply assignee filter
    //    - Apply label filter
    // 3. Return matching issues
}
```

### Optional Enhancements (Future)

- Regular expression support
- Search in specific fields only (`--title-only`, `--body-only`)
- Fuzzy matching
- Search result ranking by relevance

## Success Criteria

- [ ] Search command implemented
- [ ] Full-text search works in title and body
- [ ] Case-insensitive search
- [ ] All filters work correctly
- [ ] Results displayed in clear table format
- [ ] Performance is acceptable even with many issues
- [ ] Empty results handled gracefully

## Dependencies

- Requires #002 (Project Setup)
- Requires #003 (Storage and File System Operations)
- Requires #004 (Core CLI Commands - for table formatting patterns)
