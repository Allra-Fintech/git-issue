---
id: "020"
assignee: ""
labels: []
created: 2026-01-26T11:37:46.674936+09:00
updated: 2026-01-26T11:37:46.674936+09:00
---

# Update AI instructions to use gi create command instead of creating files directly

## Description

Currently, the AI instructions in `CLAUDE.md` focus on how to **find and read** existing issues, but do not instruct AI agents on how to **create** new issues. As a result, AI agents may create issue files directly by writing to the `.issues/open/` directory instead of using the proper `gi create` command.

This is problematic because:
1. The `gi create` command handles ID generation from `.counter` file
2. It generates the proper slug from the title
3. It creates the correct YAML frontmatter with timestamps
4. It ensures consistent file naming format

## Requirements

- Update `CLAUDE.md` to add instructions for creating new issues
- Update the `gi init` command output in `cmd/init.go` (`printAIAgentInstructions()` function) to include instructions for creating issues
- Instruct AI agents to use `gi create "issue title"` command
- Document available flags: `--assignee`, `--label`

## Success Criteria

- [ ] CLAUDE.md includes a "Creating Issues" section
- [ ] `gi init` output includes instructions for creating issues using `gi create`
- [ ] AI agents use `gi create` command instead of manually creating files
- [ ] Instructions cover common flags for issue creation
