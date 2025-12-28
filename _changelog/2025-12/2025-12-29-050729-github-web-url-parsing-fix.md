# Fix GitHub Web URL Parsing for Clone Command

**Date**: December 29, 2025

## Summary

Fixed a bug where `gitr clone` failed when using GitHub web URLs containing path segments like `/tree/main`, `/blob/main/file.go`, or `/pull/123`. The URL parser now correctly extracts the repository path by stripping these GitHub-specific path patterns, enabling users to directly clone from any GitHub web URL they copy from their browser.

## Problem Statement

When users copied a GitHub URL directly from their browser while viewing a specific branch, file, or pull request, `gitr clone` would fail because it incorrectly interpreted the full path as the repository name.

### Pain Points

- **Failed clones from browser URLs**: Copying `https://github.com/owner/repo/tree/main` would attempt to clone `owner/repo/tree/main.git` instead of `owner/repo.git`
- **Confusing error messages**: Git would report "not a valid repository name" without indicating the URL parsing issue
- **Broken user workflow**: Users had to manually edit URLs to remove path segments before cloning
- **Inconsistent behavior**: `/blob/` URLs were handled, but `/tree/` and other patterns were not

### Example Failure

```bash
$ gitr clone https://github.com/sarwarbeing-ai/Agentic_Design_Patterns/tree/main

# Attempted to clone:
git@github.com:sarwarbeing-ai/Agentic_Design_Patterns/tree/main.git

# Error:
fatal: remote error:
  sarwarbeing-ai/Agentic_Design_Patterns/tree/main is not a valid repository name
```

## Solution

Updated the `GetRepoPath` function in `pkg/url/url.go` to recognize and strip all common GitHub web URL path patterns before extracting the repository path.

### Patterns Now Handled

| Pattern | Use Case |
|---------|----------|
| `/tree/` | Branch or directory browsing |
| `/blob/` | File viewing (already handled) |
| `/commits/` | Commit history |
| `/pull/` | Pull request pages |
| `/issues/` | Issue pages |
| `/compare/` | Branch comparison |

## Implementation Details

### Before

```go
case config.GitHub:
    if strings.Contains(url, "/blob/") {
        return url[strings.Index(url, host)+1+len(host) : strings.Index(url, "/blob/")], nil
    } else {
        return url[strings.Index(url, host)+1+len(host):], nil
    }
```

### After

```go
case config.GitHub:
    repoPath := url[strings.Index(url, host)+1+len(host):]
    // GitHub web URL patterns that should be stripped to extract owner/repo
    githubPatterns := []string{"/tree/", "/blob/", "/commits/", "/pull/", "/issues/", "/compare/"}
    for _, pattern := range githubPatterns {
        if idx := strings.Index(repoPath, pattern); idx != -1 {
            return repoPath[:idx], nil
        }
    }
    return repoPath, nil
```

### Key Design Decisions

1. **Pattern list approach**: Rather than multiple if-else statements, using a slice of patterns makes the code more maintainable and extensible
2. **First match wins**: Patterns are checked in order; the first match determines where to truncate
3. **Preserve existing behavior**: Plain repository URLs without path segments continue to work unchanged

## Files Changed

| File | Change |
|------|--------|
| `pkg/url/url.go` | Updated `GetRepoPath` to handle GitHub URL patterns |
| `pkg_test/url/url_test.go` | Added comprehensive test cases for all patterns |

## Testing

Added 13 new test cases covering:

- Basic repository URLs (with and without `.git` suffix)
- `/tree/` URLs with branches and nested paths
- `/blob/` URLs with file paths
- `/commits/`, `/pull/`, `/issues/`, `/compare/` URLs

```bash
$ go test ./pkg_test/url/... -v
=== RUN   TestGetRepoPath
=== RUN   TestGetRepoPath/GitHub_URL_patterns_should_extract_correct_repo_path
--- PASS: TestGetRepoPath (0.00s)
```

### Verified Fix

```bash
$ gitr clone --dry https://github.com/sarwarbeing-ai/Agentic_Design_Patterns/tree/main

+------------+---------------------------------------------------------------------+
| remote     | https://github.com/sarwarbeing-ai/Agentic_Design_Patterns/tree/main |
+------------+---------------------------------------------------------------------+
| repo-name  | Agentic_Design_Patterns                                             |
+------------+---------------------------------------------------------------------+
| ssh-url    | git@github.com:sarwarbeing-ai/Agentic_Design_Patterns.git           |
+------------+---------------------------------------------------------------------+
| clone-path | /Users/swarup/scm/github.com/sarwarbeing-ai/Agentic_Design_Patterns |
+------------+---------------------------------------------------------------------+
```

## Benefits

- **Seamless browser-to-terminal workflow**: Copy any GitHub URL directly and clone without modification
- **Reduced friction**: No need to manually strip path segments from URLs
- **Comprehensive coverage**: Handles all common GitHub web URL patterns
- **Better error prevention**: Users won't encounter confusing "invalid repository name" errors
- **Consistent behavior**: All GitHub web URL types now work uniformly

## Impact

**Users**: Can now clone repositories by copying any GitHub URL from their browser, including:
- Branch view URLs (`/tree/branch-name`)
- File view URLs (`/blob/main/path/to/file`)
- Pull request URLs (`/pull/123`)
- Issue URLs (`/issues/456`)
- Commit history URLs (`/commits/main`)
- Compare URLs (`/compare/main...feature`)

**Developers**: Clean, extensible pattern-matching approach makes adding new patterns trivial.

## Related Work

This fix complements the existing GitLab URL handling which uses `/-/` as the delimiter for web URL paths. Both providers now have robust URL parsing for browser-copied URLs.

---

**Status**: ✅ Production Ready  
**Scope**: Small (URL parsing fix with test coverage)

