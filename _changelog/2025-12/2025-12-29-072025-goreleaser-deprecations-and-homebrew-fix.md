# GoReleaser Deprecation Fixes and Homebrew Install Fix

**Date**: December 29, 2025

## Summary

Fixed GoReleaser v2 deprecation warnings and resolved a Homebrew installation failure caused by using binary format instead of tar.gz archives.

## Problem Statement

Two issues were discovered after the initial release automation:

### Pain Points

- **Deprecation warnings in CI**: GoReleaser v2.13 showed warnings about deprecated config options (`archives.format` and `brews`)
- **Homebrew install failure**: Running `brew install swarupdonepudi/tap/gitr` failed with `No such file or directory - gitr`

### Root Cause

The binary format downloads files like `gitr-darwin-arm64`, but the generated Homebrew formula tries `bin.install "gitr"` which doesn't exist. The file retains its original name when downloaded as a raw binary.

## Solution

### 1. Fix Deprecation Warnings

Updated `.goreleaser.yaml` to use GoReleaser v2 syntax:

| Deprecated | Replacement |
|------------|-------------|
| `archives.format: binary` | `archives.formats: [tar.gz]` |
| `brews:` | `homebrews:` |

### 2. Fix Homebrew Install

Changed from binary format to tar.gz archives:

```yaml
# Before (broken)
archives:
  - name_template: "{{ .ProjectName }}-{{ .Os }}-{{ .Arch }}"
    formats:
      - binary

# After (working)
archives:
  - name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats:
          - zip
```

**Why this works**: tar.gz archives extract to a directory containing a binary named `gitr`, which the formula can find and install correctly.

## Implementation Details

### Release Artifact Changes

**Before (v1.0.9)**:
```
gitr-darwin-arm64      (raw binary)
gitr-darwin-amd64      (raw binary)
gitr-linux-arm64       (raw binary)
gitr-linux-amd64       (raw binary)
```

**After (v1.0.10+)**:
```
gitr_1.0.10_darwin_arm64.tar.gz   (contains: gitr)
gitr_1.0.10_darwin_amd64.tar.gz   (contains: gitr)
gitr_1.0.10_linux_arm64.tar.gz    (contains: gitr)
gitr_1.0.10_linux_amd64.tar.gz    (contains: gitr)
gitr_1.0.10_windows_amd64.zip     (contains: gitr.exe)
gitr_1.0.10_windows_arm64.zip     (contains: gitr.exe)
```

## Benefits

- **Clean CI logs**: No more deprecation warnings in release workflow
- **Working Homebrew install**: `brew install swarupdonepudi/tap/gitr` succeeds
- **Standard artifact naming**: Follows conventional `project_version_os_arch.tar.gz` pattern
- **Future-proof**: Using current GoReleaser v2 best practices

## Impact

### For Users
- Homebrew installation now works correctly
- Can install with `brew install swarupdonepudi/tap/gitr`

### For CI/CD
- No deprecation warnings in release logs
- Cleaner release output

## Related Work

- Follows: `2025-12-29-071405-goreleaser-v2-and-auto-versioning.md`
- Part of the release automation initiative

---

**Status**: ✅ Production Ready
**Timeline**: ~15 minutes debugging and fix

