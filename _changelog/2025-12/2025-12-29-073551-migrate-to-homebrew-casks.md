# Migrate from brews to homebrew_casks

**Date**: December 29, 2025

## Summary

Migrated GoReleaser configuration from the deprecated `brews` section to `homebrew_casks`, following the official deprecation notice since GoReleaser v2.10. This aligns with Homebrew's recommendation to use Casks for pre-compiled binaries.

## Problem Statement

GoReleaser v2.10+ deprecates the `brews` configuration in favor of `homebrew_casks`:

### Pain Points

- Deprecation warning in CI: `brews is being phased out in favor of homebrew_casks`
- The `brews` approach generates "hackyish" formulas for pre-compiled binaries
- Casks are the proper Homebrew mechanism for distributing pre-compiled binaries

## Solution

Replaced `brews` with `homebrew_casks` in `.goreleaser.yaml` and created a tap migration file for seamless user transition.

### Key Changes

**GoReleaser config (`.goreleaser.yaml`)**:
- `brews:` → `homebrew_casks:`
- Removed `directory: Formula` (Casks default to `Casks/`)
- Changed `install: |` to `binaries: [gitr]`

**Homebrew tap (`homebrew-tap`)**:
- Created `tap_migrations.json` for automatic user migration

## Implementation Details

### Before

```yaml
brews:
  - name: gitr
    directory: Formula
    install: |
      bin.install "gitr"
```

### After

```yaml
homebrew_casks:
  - name: gitr
    binaries:
      - gitr
```

### Tap Migration

Created `tap_migrations.json` in homebrew-tap:

```json
{
  "gitr": "swarupdonepudi/tap/gitr"
}
```

This ensures existing users running `brew upgrade` automatically migrate to the new Cask.

## Benefits

- **No more deprecation warnings**: Uses current GoReleaser v2 best practices
- **Proper Homebrew pattern**: Casks are designed for pre-compiled binaries
- **Seamless migration**: Existing users auto-upgrade via tap_migrations.json
- **Future-proof**: Aligned with GoReleaser's direction

## Impact

### For Users
- **Existing users**: `brew upgrade gitr` handles migration automatically
- **New users**: Install via `brew install --cask swarupdonepudi/tap/gitr`

### For CI/CD
- Clean release logs without deprecation warnings

## Migration Steps

1. Push `tap_migrations.json` to homebrew-tap
2. Release new version with `homebrew_casks` config
3. Delete old `Formula/gitr.rb` after successful cask release

## Related Work

- Follows: `2025-12-29-072025-goreleaser-deprecations-and-homebrew-fix.md`
- Part of the release automation initiative

---

**Status**: ✅ Production Ready
**Timeline**: ~20 minutes

