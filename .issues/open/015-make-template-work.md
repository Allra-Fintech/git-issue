---
id: "015"
assignee: ""
labels: []
created: 2025-11-24T11:30:27.643813+09:00
updated: 2025-11-24T11:30:27.643813+09:00
---

# Make template work

## Description

The `gi create` command is not using the `.issues/template.md` file when creating new issues. New issues are created with minimal content instead of the structured template that includes Description, Requirements, and Success Criteria sections.

## Current Behavior

When running `gi create`, new issues are created with only:
- YAML frontmatter (id, assignee, labels, created, updated)
- A simple title heading
- Empty content

## Expected Behavior

New issues should be created based on `.issues/template.md`, which includes:
- YAML frontmatter (with appropriate fields populated)
- Title heading
- Description section with placeholder text
- Requirements section with example bullet points
- Success Criteria section with checkboxes

## Steps to Reproduce

1. Run `gi create` and enter an issue title
2. Open the newly created issue file
3. Observe that the file does not contain the template sections

## Impact

Users must manually add common sections (Description, Requirements, Success Criteria) to each new issue, reducing productivity and consistency across issues.

## Success Criteria

- [ ] `gi create` command reads `.issues/template.md` when creating new issues
- [ ] Template sections (Description, Requirements, Success Criteria) appear in new issues
- [ ] YAML frontmatter fields are properly populated (id, timestamps, etc.)
- [ ] Title from user input replaces "Issue Title" placeholder in template
