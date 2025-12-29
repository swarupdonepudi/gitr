# macOS Gatekeeper Quarantine Fix

**Date**: December 29, 2025

## Summary

Added a post-install hook to the Homebrew Cask that automatically removes the macOS quarantine attribute, eliminating the Gatekeeper "cannot verify developer" warning on first run.

## Problem Statement

When users installed gitr via Homebrew Cask, macOS Gatekeeper displayed a warning dialog:

> "gitr" Not Opened - Apple could not verify "gitr" is free of malware that may harm your Mac or compromise your privacy.

This happened because the binary isn't code-signed with an Apple Developer certificate.

## Failed Attempts

Finding the correct GoReleaser syntax required multiple attempts:

### Attempt 1: `hooks.post_install` (single key)
```yaml
hooks:
  post_install: |
    system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/gitr"]
```
**Error**: `field post_install not found in type config.HomebrewCaskHooks`

### Attempt 2: `postflight` (top-level field)
```yaml
postflight: |
  system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/gitr"]
```
**Error**: `field postflight not found in type config.HomebrewCask`

### Why these failed
- GoReleaser's config expects **nested keys**: `hooks:` → `post:` → `install:`
- There is no single `post_install` key or top-level `postflight` field
- The Homebrew Cask DSL concept of "postflight" is exposed via the nested `hooks.post.install` structure in GoReleaser

## Solution (Found via Deep Research)

The correct syntax uses **nested keys** under `hooks`:

```yaml
hooks:
  post:
    install: |
      if OS.mac?
        system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/gitr"]
      end
```

This was confirmed by:
1. [GoReleaser official documentation](https://goreleaser.com/customization/homebrew_casks/)
2. Real-world examples from [FileBrowser](https://github.com/filebrowser/filebrowser) and [SLIM](https://github.com/agntcy/slim) projects
3. GoReleaser v2.13 release notes confirming hooks support was added in that version

## Implementation Details

The quarantine attribute (`com.apple.quarantine`) is an extended file attribute that macOS sets on files downloaded from the internet. The post-install hook:

1. Checks if running on macOS (`if OS.mac?`)
2. Removes the quarantine attribute from the staged binary
3. Only affects the gitr binary—does not change system-wide Gatekeeper settings

GoReleaser translates this into a `postflight do ... end` block in the generated Ruby Cask file.

## Benefits

- **Seamless installation**: Users no longer see the Gatekeeper warning
- **No manual steps**: Quarantine removal happens automatically during install
- **Security preserved**: Only affects the gitr binary, not system-wide settings
- **Cross-platform safe**: The `if OS.mac?` check ensures it only runs on macOS

## Key Learnings

1. GoReleaser v2.13+ is required for `homebrew_casks` hooks support
2. The syntax is `hooks.post.install`, not `hooks.post_install` or `postflight`
3. Always wrap macOS-specific commands in `if OS.mac?` for portability
4. Real-world GitHub examples are invaluable when documentation is unclear

## Impact

### For Users
- Clean installation experience without security popups
- No need to manually run `xattr` commands or bypass Gatekeeper

---

**Status**: ✅ Production Ready
**Timeline**: ~1 hour (including research and debugging)
