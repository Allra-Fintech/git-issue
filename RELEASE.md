# Release Guide

This document describes how to cut and publish a `git-issue` release.

## Versioning

`git-issue` follows [Semantic Versioning](https://semver.org/):

- `MAJOR` for backwards-incompatible changes
- `MINOR` for new functionality that is backwards compatible
- `PATCH` for backwards-compatible bug fixes

Release tags use the format `vMAJOR.MINOR.PATCH` (for example `v1.2.0`).

## Release Checklist

1. Ensure `main` is green and contains the commits you want to release.
2. Update any docs if needed (README, CHANGELOG, issue templates, etc.).
3. Run the full test and lint suite:
   ```bash
   make fmt
   make lint
   make test-coverage
   ```
4. Verify binaries build cleanly (with version + ldflags):
   ```bash
   make clean
   make build
   make build-all
   ```
5. Manually smoke-test binaries on supported platforms if possible.
6. Update `.issues/` records if the release closes outstanding items.

## Creating a Release

1. Decide the new version number (e.g., `v1.0.0`).
2. Tag the commit and push the tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. In GitHub, create a new release from the tag:
   - Title: `v1.0.0`
   - Use the tag description as release notes (include highlights, contributors, breaking changes).
4. Upload the cross-compiled binaries produced by `make build-all`:
   - `git-issue-darwin-arm64`
   - `git-issue-darwin-amd64`
   - `git-issue-linux-amd64`
5. (Optional) Attach SHA256 checksums for each binary.

## Testing Releases

For each platform:

```bash
chmod +x git-issue-<os>-<arch>
./git-issue-<os>-<arch> --version
./git-issue-<os>-<arch> --help
```

Perform a quick workflow smoke-test:

1. `git-issue init`
2. `git-issue create "Smoke test issue"`
3. `git-issue list`
4. `git-issue close 001`

## Post-Release

1. Announce the release (internal Slack, email, etc.).
2. Update downstream consumers (bots, scripts) if new features require it.
3. Open follow-up issues for anything deferred.
