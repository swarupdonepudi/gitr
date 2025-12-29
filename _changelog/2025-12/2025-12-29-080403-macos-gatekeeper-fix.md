# macOS Gatekeeper Quarantine Fix

**Date**: December 29, 2025

## Summary

Added a post-install hook to the Homebrew Cask that automatically removes the macOS quarantine attribute, eliminating the Gatekeeper "cannot verify developer" warning on first run.

## Problem Statement

When users installed gitr via Homebrew Cask, macOS Gatekeeper displayed a warning dialog:

> "gitr" Not Opened - Apple could not verify "gitr" is free of malware that may harm your Mac or compromise your privacy.

This happened because the binary isn't code-signed with an Apple Developer certificate.

## Solution

Added a `post_install` hook to the `homebrew_casks` configuration that removes the quarantine extended attribute after installation:

```yaml
hooks:
  post_install: |
    system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/gitr"]
```

## Implementation Details

The quarantine attribute (`com.apple.quarantine`) is an extended file attribute that macOS sets on files downloaded from the internet. The post-install hook removes this attribute only from the gitr binary—it does not affect system-wide Gatekeeper settings.

## Benefits

- **Seamless installation**: Users no longer see the Gatekeeper warning
- **No manual steps**: Quarantine removal happens automatically during install
- **Security preserved**: Only affects the gitr binary, not system-wide settings

## Impact

### For Users
- Clean installation experience without security popups
- No need to manually run `xattr` commands or bypass Gatekeeper

---

**Status**: ✅ Production Ready
**Timeline**: ~5 minutes

