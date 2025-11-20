---
id: "014"
assignee: ""
labels:
    - bug
created: 2025-11-20T14:57:43.364376+09:00
updated: 2025-11-20T15:14:10.062556+09:00
---

# Create new issue when close

## Problem

Closing an issue with `gi close 004 -c` is creating a brand new issue file in `./.issues/open/` instead of only moving the existing one to `./.issues/closed/`.

Observed console output from `~/work/cargo-note-backend`:

```
gi close 004 -c
✓ Closed issue #004
[feat/balance-check 9901acb] Close issue #004
 8 files changed, 630 insertions(+), 1 deletion(-)
 rename .issues/{open => closed}/004-add-transaction-history-endpoint.md (100%)
 create mode 100644 .issues/open/004-save-transaction-history-to-database.md
 ...
```

## Expected Behavior

- The close command should move the existing issue file from `./.issues/open/` to `./.issues/closed/` and stop there.
- No new issue files should be created during close, regardless of the issue title.

## Actual Behavior

- The existing file was moved to `./.issues/closed/004-add-transaction-history-endpoint.md`.
- A new file `./.issues/open/004-save-transaction-history-to-database.md` was created alongside the move, leaving the issue appearing open again with a new slug.

## Steps to Reproduce

1. In a repo with `.issues/open/004-add-transaction-history-endpoint.md`, run `gi close 004 -c`.
2. Inspect `.issues/open/` and `.issues/closed/`.

## Requirements

- Ensure `gi close` only relocates the targeted issue file; it must not create any new `.issues/open/*.md` file as part of the operation.
- Preserve the original filename/slug when closing, even if the issue title has changed.
- Add a regression test that covers closing an issue after its title or slug has been modified.

## Success Criteria

- [ ] Running `gi close 004 -c` results in exactly one file in `./.issues/closed/` for ID 004 and zero files in `./.issues/open/` with ID 004.
- [ ] Closing and reopening flows handle title changes without creating duplicate files.
- [ ] New test(s) fail on current main and pass after the fix.
