---
id: "010"
assignee: ""
labels: []
created: 2025-11-14T21:44:51.946687+09:00
updated: 2025-11-14T22:23:23.349942+09:00
---

# Change binary file name from git-issue to gi

## Description

The CLI binary is currently published as `git-issue`, which is a mouthful to type frequently. We want to ship future builds under the shorter `gi` name for faster invocation while still keeping the project name (and repo) as-is. The change needs to cover the build pipeline, release assets, documentation, and any developer tooling that references the executable.

Key considerations:

- Update Makefile defaults, release assets, and completion scripts to emit/install `gi`.
- Ensure packaging instructions and docs explain the rename.

## Requirements

- Update build scripts (Makefile, GoReleaser/GHA workflow if applicable) to produce binaries named `gi-*`.
- Remove references to the old binary name in README, DEVELOPMENT, RELEASE docs, and completion instructions.
- Ensure CI/CD release assets adopt the new naming convention without breaking the existing release workflow.
- Verify the CLI still reports the correct version metadata after the rename.

## Success Criteria

- [ ] `make build`/`make build-all` output `gi` binaries and the install step places the executable under the new name.
- [ ] Release artifacts (local and CI) are named `gi-<os>-<arch>` and documented accordingly.
- [ ] Docs clearly state the new name.
- [ ] Manual smoke test shows `gi` command executes the CLI with expected behavior.
