# Clone URL Validation and Graceful Error Handling

**Date**: January 15, 2026

## Summary

Added comprehensive URL validation to the `gitr clone` command that catches invalid inputs before they cause panics. Most notably, the implementation detects when clipboard contents (from a previous gitr clone's `cd` command) are accidentally pasted instead of a repository URL, providing a self-aware error message that explains exactly what happened.

## Problem Statement

When running `gitr clone` with invalid input like `cd /path/to/repo` (accidentally pasted from clipboard after another gitr clone finished), the CLI would panic with an unhelpful "index out of range" error deep in the URL parsing code.

### Pain Points

- **Cryptic panic**: `panic: runtime error: index out of range [1] with length 1` gave no indication of what went wrong
- **No input validation**: The `GetHostname` function assumed valid URL format without checking
- **Clipboard collision**: gitr's own feature (copying `cd` command to clipboard) could cause user errors when working with multiple terminals
- **No guidance**: Users had no idea what they did wrong or how to fix it

### Real-World Scenario

```bash
# Terminal 1: Clone finishes, copies "cd ~/scm/github.com/owner/repo" to clipboard
✓  Repository cloned successfully
   A cd command has been copied to your clipboard.

# Terminal 2: User had copied a URL, but clipboard was overwritten
# Pasting now gives the cd command instead of the URL
$ gitr clone cd /Users/swarup/scm/github.com/openbao/openbao-helm
panic: runtime error: index out of range [1] with length 1
```

## Solution

Implemented early URL validation with pattern detection that provides context-aware error messages:

1. **URL Validation Function**: New `ValidateCloneURL()` in `pkg/url/url.go` that checks input before any parsing
2. **Pattern Detection**: Identifies specific invalid patterns (cd commands, local paths, malformed URLs)
3. **Smart Error Messages**: Context-aware messages that explain what happened and how to fix it
4. **Shell-Split Detection**: Handles the case where pasted `cd /path` gets split into separate arguments by the shell

### Key Components

**`pkg/url/url.go`**:
- `InvalidURLReason` enum for categorizing validation failures
- `URLValidationResult` struct with details about why validation failed
- `ValidateCloneURL()` function with pattern matching for valid/invalid URLs

**`pkg/ui/errors.go`**:
- `InvalidCloneURLCdCommand()` - self-aware message about clipboard being overwritten
- `InvalidCloneURLLocalPath()` - guidance for local path mistakes
- `InvalidCloneURLFormat()` - generic fallback with examples

**`cmd/gitr/root/clone.go`**:
- Early validation before clone logic
- Special handling for shell-split `cd` + `path` arguments

## Implementation Details

### URL Validation Patterns

The validator checks for:

```go
// Valid patterns
git@host:path           // SSH short form
ssh://user@host/path    // SSH URL form
https://host/path       // HTTPS URLs
http://host/path        // HTTP URLs

// Invalid patterns detected
cd /path/to/repo        // Clipboard cd command
/absolute/path          // Local filesystem path
~/relative/path         // Home-relative path
./current/relative      // Current-relative path
```

### Shell-Split Detection

When user pastes `cd /path/to/repo`, the shell splits it:
- `args[0]` = `"cd"`
- `args[1]` = `"/path/to/repo"`

The clone handler detects this specific pattern and reconstructs the original command for the error message.

### Error Message Example

```
✗  Invalid Clone URL

   Received a cd command instead of a repository URL:
   cd /Users/swarup/scm/github.com/openbao/openbao-helm

   Hint: Your clipboard was likely overwritten by another gitr clone that just finished
   Hint: Re-copy the repository URL and try again
   Hint: Example: gitr clone https://github.com/owner/repo
```

## Benefits

### User Experience
- **Self-aware messaging**: The error acknowledges that gitr's own clipboard feature likely caused the issue
- **Actionable guidance**: Clear hints on how to resolve the problem
- **No more panics**: Graceful failures with helpful output instead of stack traces

### Developer Experience
- **Reusable validation**: `ValidateCloneURL()` can be used elsewhere if needed
- **Typed reasons**: `InvalidURLReason` enum enables programmatic handling of different failure types
- **Extensible patterns**: Easy to add new invalid URL patterns as discovered

### Code Quality
- **Fail fast**: Validation happens before any complex parsing logic
- **Clear separation**: Validation logic separated from URL parsing logic
- **Testable**: Validation function is pure and easily testable

## Impact

- **All gitr users** benefit from clearer error messages when clone fails
- **Multi-terminal workflows** are now more resilient to clipboard collisions
- **Future debugging** is easier with categorized validation failures

## Files Changed

| File | Changes |
|------|---------|
| `pkg/url/url.go` | +107 lines: URL validation types and function |
| `pkg/ui/errors.go` | +37 lines: Invalid URL error messages |
| `cmd/gitr/root/clone.go` | +29 lines: Validation integration |

## Related Work

- Builds on the Terminal UX improvements from December 2025 (`2025-12-29-062623`)
- Uses the same `ui.Error()` pattern with hints established in that work
- Complements the Modern Clone Progress Display (`2026-01-03-140029`)

---

**Status**: ✅ Production Ready
**Timeline**: Single session
