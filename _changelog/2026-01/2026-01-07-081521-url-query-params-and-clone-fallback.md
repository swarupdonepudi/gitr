# Fix URL Query Parameter Handling and Clone Fallback Logic

**Date**: January 7, 2026

## Summary

Fixed multiple issues preventing successful clones when URLs contain query parameters (common when copying from browsers) and improved the SSH-to-HTTP fallback behavior. URLs like `https://github.com/owner/repo?utm_source=chatgpt.com` now work correctly, and failed SSH clones properly fall back to HTTP without noise or false error messages.

## Problem Statement

Users copying repository URLs from browsers often include tracking parameters like `?utm_source=chatgpt.com`. These parameters were being passed through to git operations, causing clone failures.

### Pain Points

- **Clone failures**: URLs with query params caused "invalid pkt-len" errors
- **Path pollution**: The `path` command returned paths containing query strings (e.g., `~/scm/github.com/owner/repo?utm_source=chatgpt.com`)
- **False "Repository Not Found" errors**: The `isRepoNotFoundError` function used an overly broad pattern (`"not found"`) that matched unrelated errors
- **HTTP fallback failures**: When SSH clone failed, the directory wasn't cleaned up, causing "repository already exists" errors when falling back to HTTP
- **Noisy output**: Users saw "SSH Clone Failed / Trying HTTP clone instead..." even when clones ultimately succeeded

## Solution

Implemented URL sanitization at the entry points of clone operations, tightened error detection patterns, added cleanup between fallback attempts, and removed unnecessary user-facing messages.

### Key Components

1. **URL Sanitization** (`pkg/url/url.go`)
   - New `StripQueryParams()` function removes query parameters and fragments
   - Handles both HTTP(S) and SSH URL formats
   - Uses proper URL parsing with string fallback

2. **Clone Entry Points** (`pkg/clone/clone.go`)
   - `Clone()` and `GetClonePath()` now sanitize URLs before processing
   - Added directory cleanup before HTTP fallback
   - Changed SSH fallback message to debug-only logging

3. **Error Detection** (`pkg/clone/clone.go`)
   - Removed overly broad `"not found"` pattern from `isRepoNotFoundError()`
   - Added specific patterns for actual repo-not-found errors

4. **UI Polish** (`pkg/ui/success.go`, `pkg/clone/clone.go`)
   - Reduced excessive blank lines between clone messages
   - Removed trailing `\n\n` from `Cloning()` function
   - Removed redundant `fmt.Println()` calls after progress display

## Implementation Details

### URL Sanitization

```go
// StripQueryParams removes query parameters and fragments from URLs
func StripQueryParams(inputUrl string) string {
    // Handle SSH URLs (git@host:path)
    if strings.HasPrefix(inputUrl, "git@") && !strings.Contains(inputUrl, "://") {
        if idx := strings.Index(inputUrl, "?"); idx != -1 {
            return inputUrl[:idx]
        }
        return inputUrl
    }
    // For HTTP(S) URLs, use proper URL parsing
    parsed, err := url.Parse(inputUrl)
    if err != nil {
        // Fallback to string manipulation
        if idx := strings.Index(inputUrl, "?"); idx != -1 {
            return inputUrl[:idx]
        }
        return inputUrl
    }
    parsed.RawQuery = ""
    parsed.Fragment = ""
    return parsed.String()
}
```

### Tighter Error Pattern Matching

```go
// Before: Too broad - matched "host not found", "key not found", etc.
notFoundPatterns := []string{
    "repository not found",
    "not found",  // ❌ Too broad!
}

// After: Specific to actual repo-not-found errors
notFoundPatterns := []string{
    "repository not found",
    "repo not found",
    "remote: repository not found",
    "project not found",
    "the project you were looking for could not be found",
}
```

### Clean Fallback Logic

```go
if err := sshClone(sshCloneUrl, repoLocation); err != nil {
    if isRepoNotFoundError(err) {
        return "", errors.New("repository not found...")
    }
    // Clean up directory from failed SSH clone
    os.RemoveAll(repoLocation)
    // Debug-only log, not user-facing
    log.Debugf("SSH clone failed, trying HTTP fallback: %v", err)
    // Try HTTP
    if err := httpClone(httpCloneUrl, repoLocation); err != nil {
        return "", errors.Wrap(err, "error cloning...")
    }
}
```

## Benefits

- **Browser URL support**: Copy-paste URLs from any browser now work without manual cleanup
- **Reliable fallback**: SSH-to-HTTP fallback works correctly without "already exists" errors
- **Accurate errors**: No more false "Repository Not Found" when SSH auth fails
- **Clean output**: Users only see relevant information, no debug noise
- **Compact display**: Reduced vertical whitespace in clone output

## Impact

### User Experience

- **Before**: Copying a GitHub link from ChatGPT resulted in clone failure
- **After**: Same URL clones successfully with no user intervention

### Output Quality

```bash
# Before (noisy)
↓  Cloning git@github.com:owner/repo.git


!  SSH Clone Failed

   Trying HTTP clone instead...


✓  Repository cloned successfully

# After (clean)
↓  Cloning git@github.com:owner/repo.git

✓  Repository cloned successfully
```

## Testing

Added comprehensive tests for `StripQueryParams()`:

```go
func TestStripQueryParams(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        // URLs with tracking params
        {"https://github.com/owner/repo?utm_source=chatgpt.com", "https://github.com/owner/repo"},
        // URLs with fragments
        {"https://github.com/owner/repo#readme", "https://github.com/owner/repo"},
        // Clean URLs unchanged
        {"https://github.com/owner/repo.git", "https://github.com/owner/repo.git"},
        // SSH URLs
        {"git@github.com:owner/repo.git?foo=bar", "git@github.com:owner/repo.git"},
    }
    // ...
}
```

## Files Changed

| File | Changes |
|------|---------|
| `pkg/url/url.go` | Added `StripQueryParams()` function |
| `pkg/url/url_test.go` | Added tests for query param stripping |
| `pkg/clone/clone.go` | URL sanitization, error patterns, cleanup, debug logging |
| `pkg/ui/success.go` | Reduced trailing newlines in `Cloning()` |

## Related Work

This fix improves URL handling robustness, complementing existing browser URL pattern support for GitHub tree/blob/pull URLs.

---

**Status**: ✅ Production Ready
**Timeline**: ~30 minutes

