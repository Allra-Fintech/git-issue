---
id: "009"
assignee: ""
labels: []
created: 2025-11-14T17:23:12.796868+09:00
updated: 2025-11-14T17:23:12.796868+09:00
---

# Build in CI/CD

## Description

Automate the build pipeline so that every change merged into `main` produces fresh binaries as GitHub Actions artifacts. This ensures contributors and AI agents can always download the latest `git-issue` binary without waiting for a tagged release.

## Tasks

- [ ] Add a GitHub Actions workflow (e.g., `.github/workflows/build.yml`) that triggers on `push` to `main`.
  - [ ] Run `make build` using the repo's version injection (ensure `VERSION` is derived from the commit SHA when tags are absent).
  - [ ] Execute `make test` and `make lint` before building to guarantee quality.
- [ ] Matrix the build step for supported GOOS/GOARCH combinations (macOS arm64/amd64, Linux amd64) or call `make build-all`.
- [ ] Upload the resulting binaries (`git-issue-<os>-<arch>`) as workflow artifacts with a retention policy (at least 7 days).
- [ ] Gate the workflow so it only runs on the canonical repository (skip forks).
- [ ] Document in `DEVELOPMENT.md` how to download CI-built artifacts for quick testing.

## Acceptance Criteria

- Workflow runs automatically when `main` updates and finishes within 10 minutes.
- Each run publishes binaries for macOS (arm64/amd64) and Linux (amd64) as downloadable artifacts.
- Failing tests or lint steps block the upload stage.
- Documentation references the workflow and explains how to grab the artifacts.
