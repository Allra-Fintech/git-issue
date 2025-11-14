---
id: "011"
assignee: ""
labels: []
created: 2025-11-14T23:12:56.327214+09:00
updated: 2025-11-14T23:17:46.213815+09:00
---

# github binary path is wrong in the README

## Description

Following the release installation instructions on a clean macOS machine installs the binary into `/usr/local/bin`, but the README never reminds users to add that directory to their `PATH`. On systems where `/usr/local/bin` is missing from the default `PATH`, running `gi` fails with `/usr/local/bin/gi: line 1: Not: command not found`.

## Steps to Reproduce

1. Follow the "From Release" instructions to move the downloaded binary into `/usr/local/bin`.
2. Keep `/usr/local/bin` out of `PATH` (default on some setups, e.g. bare zsh config).
3. Run `gi` from the terminal.

## Expected Behavior

The README documents the need to ensure `/usr/local/bin` is part of `PATH`, so the binary runs successfully after installation.

## Actual Behavior

Users copy the commands verbatim, then `gi` fails to execute because the shell cannot find the binary.

## Additional Context

- The README points users to release assets with filenames such as `gi-darwin-arm64`, but `.github/workflows/build.yml` currently only uploads those binaries as CI artifacts (scoped to the workflow run SHA) and never attaches them to GitHub Releases with the expected names/paths.
- Someone following the README links cannot download the binaries because they do not exist at the advertised release URLs, even if `/usr/local/bin` is set up correctly.

## Acceptance Criteria

- README explicitly calls out that `/usr/local/bin` must be in the `PATH` (or offers an alternative install location).
- Copy-pasting the platform-specific snippet plus the follow-up instructions lets users immediately run `gi`.
- Build workflow publishes binaries to the release assets (matching the names referenced in the README), not just as temporary CI artifacts.
