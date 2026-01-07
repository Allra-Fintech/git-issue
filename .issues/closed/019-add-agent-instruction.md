---
id: "019"
assignee: ""
labels: []
created: 2026-01-07T14:16:07.924716+09:00
updated: 2026-01-07T14:28:50.454501+09:00
---

# Add AI Agent Instruction Guide to gi init Output

## Description

Enhance the `gi init` command output to display sample AI agent instructions after successful initialization. This will guide users to set up their coding agents (Claude Code, Cursor, etc.) to work seamlessly with git-issue.

## Requirements

- After displaying the `.issues/` directory structure, add a new section that shows:
  - A brief explanation about AI agent integration
  - Sample instruction content that users can copy to their CLAUDE.md, AGENTS.md, or .cursorrules files
  - The sample should teach AI agents how to:
    - Find issues by ID pattern (`.issues/open/{id}-*.md` or `.issues/closed/{id}-*.md`)
    - Understand issue file structure (YAML frontmatter + Markdown body)
    - Determine status by directory location (open/ vs closed/, not a YAML field)
    - Parse and work with issue references (e.g., "#001", "issue 001")
- Suggest common agent instruction file names:
  - CLAUDE.md (for Claude Code)
  - AGENTS.md (for general AI agents)
  - .cursorrules (for Cursor IDE)
- Keep the output concise but informative
- Make it easy for users to understand the benefit of setting this up

## Success Criteria

- [ ] `gi init` command prints AI agent instruction guide after directory structure
- [ ] Sample instruction text is clear and copy-paste ready
- [ ] Output mentions multiple agent instruction file types
- [ ] Instructions cover key concepts: file naming, finding by ID, YAML frontmatter, status determination
- [ ] Users understand how to set up their AI agents to work with git-issue
