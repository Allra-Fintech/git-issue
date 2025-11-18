---
id: "013"
assignee: ""
labels:
    - bug
    - fixed
created: 2025-11-18T10:20:19.765626+09:00
updated: 2025-11-18T10:24:47.5259+09:00
---

# Fix filename preservation bug when closing/opening issues with modified titles

## Problem

When users modified issue titles after creation, the `close` and `open` commands created duplicate files with malformed filenames:

1. **Korean/Special Character Titles**: Generated malformed filenames like `001-.md`
2. **Modified English Titles**: Created new files with updated slugs, leaving originals behind

### Root Cause

Both `close` and `open` commands called `SaveIssue()` before `MoveIssue()`, which:
- Generated new slugs based on current `issue.Title` (which could be modified by users)
- Created new files instead of updating existing ones
- Left original files in place, resulting in duplicates


---


## Solution

1. **Modified `MoveIssue` in `pkg/storage.go`**:
   - Update `updated` timestamp after successful file relocation
   - Preserve original filename during move

2. **Removed redundant `SaveIssue` calls in `cmd/close.go` and `cmd/open.go`**:
   - Eliminated duplicate file creation
   - Simplified command logic

3. **Added regression tests**:
   - `TestRunClosePreservesFilenameWithKoreanTitle`
   - `TestRunClosePreservesFilenameWhenTitleModified`
   - `TestRunOpenPreservesFilenameWithKoreanTitle`
   - `TestRunOpenPreservesFilenameWhenTitleModified`
   - Updated `TestMoveIssue` to verify timestamp updates

## Test Results

- ✅ All existing tests pass
- ✅ New regression tests pass
- ✅ Code coverage maintained: 86.7% (cmd), 76.0% (pkg)

## Files Changed

- `pkg/storage.go`: Add timestamp update in `MoveIssue`
- `cmd/close.go`: Remove `SaveIssue` call
- `cmd/open.go`: Remove `SaveIssue` call
- `pkg/storage_test.go`: Add timestamp verification
- `cmd/close_test.go`: Add filename preservation tests
- `cmd/open_test.go`: Add filename preservation tests
