---
id: "021"
assignee: ""
labels:
    - feature
created: 2026-02-06T14:27:24.912527+09:00
updated: 2026-02-06T14:40:03.891202+09:00
---

# Add view command to open issue markdown in default program

## Description

Add a new `gi view` command that opens an issue's markdown file using the system's default program associated with the `.md` extension. This allows users to quickly view or edit issues in their preferred markdown editor or viewer (e.g., Typora, VS Code, Obsidian).

Usage: `gi view 001`

## Requirements

- Accept an issue ID as argument (e.g., `gi view 001`)
- Look up the issue file in `.issues/open/` and `.issues/closed/`
- Open the file using the OS default program for `.md` files:
  - macOS: use `open` command
  - Linux: use `xdg-open` command
- Show an error if the issue ID is not found

## Success Criteria

- [ ] `gi view <id>` opens the issue markdown in the default associated program
- [ ] Works on macOS (`open`) and Linux (`xdg-open`)
- [ ] Displays clear error when issue ID does not exist
- [ ] Unit tests for file lookup logic
- [ ] Update README file
