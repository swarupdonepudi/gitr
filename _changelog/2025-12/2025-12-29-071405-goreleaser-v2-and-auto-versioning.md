# GoReleaser v2 Fix and Auto-Version Bumping

**Date**: December 29, 2025

## Summary

Fixed the GoReleaser version mismatch that was causing release failures and added automatic semantic version bumping to the Makefile. Releases now require zero manual version specification—just run `make release` and the next patch version is automatically determined and tagged.

## Problem Statement

The initial release automation setup had two issues:

### Pain Points

- **GoReleaser version mismatch**: The GitHub Action was downloading GoReleaser v1.26.2, but the config file uses `version: 2` format which requires GoReleaser v2.x
- **Manual version specification**: Every release required explicitly passing `version=vX.Y.Z`, which was error-prone and tedious
- **No version preview**: No way to see what the next version would be without actually releasing

### CI Error

```
error=only configurations files on version: 1 are supported, yours is version: 2
```

## Solution

Two-part fix addressing both the CI failure and developer experience:

1. **Update GitHub Actions workflow** to use GoReleaser v2
2. **Add auto-version bumping** with configurable bump type (default: patch)

### Key Components

**`.github/workflows/release.yml`**:
- Updated `goreleaser-action` from v5 to v6
- Set explicit version constraint `~> v2` for GoReleaser v2.x

**`Makefile`**:
- Added `bump` variable with default value `patch`
- Added `next-version` target for previewing without releasing
- Rewrote `release` target with auto-version logic

## Implementation Details

### GitHub Actions Update

```yaml
- uses: goreleaser/goreleaser-action@v6
  with:
    distribution: goreleaser
    version: "~> v2"
    args: release --clean
```

### Makefile Auto-Version Logic

```makefile
# bump: major, minor, or patch (default)
bump ?= patch

release: test build-check
	@latest=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
	# Parse semver components...
	# Increment based on bump type...
	git tag -a $$version -m "$$version"; \
	git push origin $$version
```

### Version Bump Behavior

| Command | Current | Result |
|---------|---------|--------|
| `make release` | v1.0.8 | v1.0.9 |
| `make release bump=minor` | v1.0.8 | v1.1.0 |
| `make release bump=major` | v1.0.8 | v2.0.0 |

## Benefits

- **Zero-friction releases**: Just `make release` with no arguments needed
- **Consistent versioning**: No typos or version confusion
- **Flexibility**: Override with `bump=minor` or `bump=major` when needed
- **Preview capability**: `make next-version` shows what would be released
- **CI compatibility**: GoReleaser v2 features now fully supported

## Impact

### For Maintainers
- Release process reduced to single command with no required arguments
- Eliminates version string typos and inconsistencies
- Can preview releases before committing to them

### For CI/CD
- Releases now complete successfully with GoReleaser v2
- Homebrew formula updates work as expected

## Related Work

- Continues from: `2025-12-29-070331-automated-release-pipeline.md`
- Part of the overall release automation initiative

---

**Status**: ✅ Production Ready
**Timeline**: ~30 minutes

