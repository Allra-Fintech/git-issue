---
id: "018"
assignee: ""
labels: []
created: 2025-11-27T10:01:37.532453+09:00
updated: 2025-11-27T10:01:37.532453+09:00
---

# increment counter

## Description

When creating a new issue, if the counter points to an ID that already exists in the closed directory, the `gi create` command fails with an error instead of automatically finding the next available ID.

**Current behavior:**
```
$ gi create aa
Error: failed to save issue: issue 003 exists in closed directory, cannot save to open
```

This happens because `GetNextID()` in `pkg/storage.go` simply reads the counter value and increments it, without checking if that ID is already occupied by a closed issue.

**Expected behavior:**
The system should automatically skip occupied IDs and use the next available one, updating the counter accordingly.

## Requirements

- Modify `GetNextID()` function in `pkg/storage.go:98` to check if the current counter ID exists
- If ID exists in either open or closed directory, increment and check again until an available ID is found
- Update the counter file to reflect the next available ID
- Maintain backward compatibility with existing counter behavior
- Handle edge cases (e.g., large gaps in ID sequence)

## Success Criteria

- [x] `gi create` successfully creates issues even when counter points to closed issue ID
- [x] Counter automatically increments to next available ID
- [x] System handles multiple sequential occupied IDs correctly
- [x] No breaking changes to existing functionality
- [x] Unit tests cover the auto-increment logic
- [x] Works correctly in both scenarios:
  - Counter = 3, issue 003 in closed → creates issue 004
  - Counter = 5, issues 005-007 in closed → creates issue 008

## Implementation Notes

Fixed in `pkg/storage.go:97-132`. The `GetNextID()` function now:
1. Checks if the current counter value points to an existing issue
2. Loops through IDs until it finds an available one
3. Updates the counter to the next available ID

Added comprehensive unit tests in `pkg/storage_test.go`:
- `TestGetNextIDSkipsOccupiedIDs` - verifies skipping sequential occupied IDs
- `TestGetNextIDWithGapsInSequence` - verifies finding gaps in ID sequence

All tests passing with 77.5% code coverage.
