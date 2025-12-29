# Automated Release Pipeline with GoReleaser and Homebrew

**Date**: December 29, 2025

## Summary

Updated the GoReleaser configuration to enable fully automated releases that build cross-platform binaries and automatically update the Homebrew formula. When a tag is pushed, GitHub Actions now builds for darwin/linux/windows, creates a GitHub release, and pushes the updated formula to the homebrew-tap repository.

## Problem Statement

The release process for gitr required manual intervention at multiple stages:

### Pain Points

- The `.goreleaser.yaml` configuration referenced the old repository owner (`plantoncloud`) instead of `swarupdonepudi`
- The ldflags path was incorrect, preventing version information from being embedded in builds
- Linux builds were missing, limiting the audience for the CLI tool
- The Homebrew formula had to be manually updated after each release with new versions and checksums
- No automation existed to keep the homebrew-tap in sync with new releases

## Solution

Updated the GoReleaser configuration to:
1. Reference the correct repository owner (`swarupdonepudi`)
2. Use the correct ldflags path for version embedding
3. Build for all major platforms (darwin, linux, windows)
4. Automatically update the Homebrew tap after each release

### Key Components

**Updated `.goreleaser.yaml`**:
- Fixed `ldflags` path: `github.com/swarupdonepudi/gitr/cmd/gitr/root.VersionLabel`
- Added `linux` to `goos` targets
- Updated all owner references to `swarupdonepudi`
- Added explicit token configuration for homebrew-tap access
- Set `format: binary` for raw binary downloads
- Preserved caveats from existing Homebrew formula

## Implementation Details

### GoReleaser Configuration Changes

```yaml
builds:
  - ldflags:
      - -s -w -X github.com/swarupdonepudi/gitr/cmd/gitr/root.VersionLabel={{.Version}}
    goos:
      - darwin
      - linux
      - windows
    goarch:
      - amd64
      - arm64

brews:
  - repository:
      owner: swarupdonepudi
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
```

### GitHub Actions Integration

The existing `.github/workflows/release.yml` already:
- Triggers on tag pushes matching `v*`
- Runs GoReleaser with the `HOMEBREW_TAP_GITHUB_TOKEN` secret
- Has proper permissions for creating releases

### Required Setup

A GitHub Personal Access Token with `repo` scope is required as the `HOMEBREW_TAP_GITHUB_TOKEN` repository secret to allow GoReleaser to push formula updates to the homebrew-tap repository.

## Benefits

- **Zero-touch releases**: Run `make release version=vX.Y.Z` and everything is automated
- **Cross-platform support**: Binaries now built for darwin, linux, and windows on both amd64 and arm64
- **Correct versioning**: Version information properly embedded via ldflags
- **Homebrew auto-update**: Formula automatically updated with new version and checksums
- **Preserved UX**: Caveats and installation instructions maintained in the formula

## Impact

### For Users
- Linux users can now install and use gitr
- Homebrew users get updates immediately after a release
- Version command shows correct version information

### For Maintainers
- Release process reduced from multiple manual steps to a single command
- No need to manually edit the Homebrew formula
- Consistent, reproducible release artifacts

## Related Work

- Builds on the existing GitHub Actions workflow in `.github/workflows/release.yml`
- Updates the Homebrew formula at `homebrew-tap/Formula/gitr.rb`

---

**Status**: ✅ Production Ready
**Timeline**: Single session implementation

