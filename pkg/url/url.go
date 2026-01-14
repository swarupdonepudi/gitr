package url

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/swarupdonepudi/gitr/pkg/config"
)

// InvalidURLReason describes why a URL is invalid
type InvalidURLReason int

const (
	// InvalidURLReasonUnknown is a generic invalid URL reason
	InvalidURLReasonUnknown InvalidURLReason = iota
	// InvalidURLReasonCdCommand indicates the input looks like a cd command (clipboard overwrite scenario)
	InvalidURLReasonCdCommand
	// InvalidURLReasonLocalPath indicates the input looks like a local filesystem path
	InvalidURLReasonLocalPath
	// InvalidURLReasonEmpty indicates the input is empty
	InvalidURLReasonEmpty
	// InvalidURLReasonNoHost indicates the URL has no recognizable host
	InvalidURLReasonNoHost
)

// URLValidationResult contains the result of URL validation
type URLValidationResult struct {
	Valid        bool
	Reason       InvalidURLReason
	Input        string
	DetectedCmd  string // For cd command detection, stores the full cd command
	DetectedPath string // For path detection, stores the detected path
}

// ValidateCloneURL checks if the input looks like a valid git clone URL
// Returns a validation result with details about why validation failed
func ValidateCloneURL(input string) URLValidationResult {
	input = strings.TrimSpace(input)

	// Check for empty input
	if input == "" {
		return URLValidationResult{
			Valid:  false,
			Reason: InvalidURLReasonEmpty,
			Input:  input,
		}
	}

	// Check if it looks like a cd command (clipboard overwrite scenario)
	// This happens when another gitr clone finishes and overwrites clipboard
	if strings.HasPrefix(input, "cd ") {
		path := strings.TrimPrefix(input, "cd ")
		return URLValidationResult{
			Valid:        false,
			Reason:       InvalidURLReasonCdCommand,
			Input:        input,
			DetectedCmd:  input,
			DetectedPath: path,
		}
	}

	// Check if it's a local filesystem path (starts with / or ~/ or ./ or contains backslash for Windows)
	if strings.HasPrefix(input, "/") || strings.HasPrefix(input, "~/") ||
		strings.HasPrefix(input, "./") || strings.HasPrefix(input, "../") ||
		(len(input) >= 2 && input[1] == ':') { // Windows path like C:\
		return URLValidationResult{
			Valid:        false,
			Reason:       InvalidURLReasonLocalPath,
			Input:        input,
			DetectedPath: input,
		}
	}

	// Valid URL patterns:
	// 1. SSH: git@host:path or ssh://user@host/path
	// 2. HTTPS: https://host/path or http://host/path
	// 3. Browser URLs: https://github.com/owner/repo (no .git suffix)

	// Check for SSH URL patterns
	if strings.HasPrefix(input, "git@") || strings.HasPrefix(input, "ssh://") {
		// git@host:path pattern
		if strings.HasPrefix(input, "git@") {
			if strings.Contains(input, ":") && !strings.Contains(input, "://") {
				// Looks like git@host:path - valid SSH URL
				return URLValidationResult{Valid: true, Input: input}
			}
		}
		// ssh://user@host/path pattern
		if strings.HasPrefix(input, "ssh://") {
			return URLValidationResult{Valid: true, Input: input}
		}
	}

	// Check for HTTP(S) URLs
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// Basic validation: should have a host after the scheme
		if strings.Contains(input, "://") {
			parts := strings.SplitN(input, "://", 2)
			if len(parts) == 2 && len(parts[1]) > 0 && strings.Contains(parts[1], "/") {
				return URLValidationResult{Valid: true, Input: input}
			}
			// URL like https://github.com (no path) - still valid for some operations
			if len(parts) == 2 && len(parts[1]) > 0 {
				return URLValidationResult{Valid: true, Input: input}
			}
		}
	}

	// If we get here, it doesn't match any known valid pattern
	return URLValidationResult{
		Valid:  false,
		Reason: InvalidURLReasonNoHost,
		Input:  input,
	}
}

// StripQueryParams removes query parameters and fragments from URLs.
// This handles URLs copied from browsers that include tracking params like ?utm_source=...
func StripQueryParams(inputUrl string) string {
	// Handle SSH URLs (git@host:path) - they don't have query params
	if strings.HasPrefix(inputUrl, "git@") && !strings.Contains(inputUrl, "://") {
		// For SSH format git@host:path, query params would be after the path
		if idx := strings.Index(inputUrl, "?"); idx != -1 {
			return inputUrl[:idx]
		}
		if idx := strings.Index(inputUrl, "#"); idx != -1 {
			return inputUrl[:idx]
		}
		return inputUrl
	}

	// For HTTP(S) and ssh:// URLs, use proper URL parsing
	parsed, err := url.Parse(inputUrl)
	if err != nil {
		// If parsing fails, fall back to simple string manipulation
		if idx := strings.Index(inputUrl, "?"); idx != -1 {
			return inputUrl[:idx]
		}
		if idx := strings.Index(inputUrl, "#"); idx != -1 {
			return inputUrl[:idx]
		}
		return inputUrl
	}

	// Reconstruct URL without query params or fragment
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func IsGitUrl(repoUrl string) bool {
	return strings.HasSuffix(repoUrl, ".git")
}

func IsGitSshUrl(repoUrl string) bool {
	return strings.HasPrefix(repoUrl, "ssh://") || strings.HasPrefix(repoUrl, "git@")
}

func IsGitHttpUrlHasUsername(repoUrl string) bool {
	matched, err := regexp.MatchString("https*:\\/\\/.*@+.*", repoUrl)
	if err != nil {
		log.Warnf("error matching regex in %s url. err: %v", repoUrl, err)
		return false
	}
	return matched

}

func GetRepoName(repoPath string) string {
	return strings.Split(repoPath, "/")[strings.Count(repoPath, "/")]
}

func GetHostname(url string) string {
	if url == "" {
		return ""
	}
	if IsGitSshUrl(url) {
		if strings.HasPrefix(url, "ssh://") {
			return strings.Split(strings.Split(url, "@")[1], "/")[0]
		}
		return strings.Split(strings.Split(url, "@")[1], ":")[0]
	}
	if IsGitHttpUrlHasUsername(url) {
		return strings.Split(strings.Split(url, "@")[1], "/")[0]
	}
	return strings.Split(strings.Split(url, "://")[1], "/")[0]
}

func GetRepoPath(url, host string, p config.ScmProvider) (string, error) {
	if IsGitUrl(url) {
		return url[strings.Index(url, host)+1+len(host) : strings.Index(url, ".git")], nil
	}
	switch p {
	case config.GitLab:
		if strings.Contains(url, "/-/") {
			return url[strings.Index(url, host)+1+len(host) : strings.Index(url, "/-/")], nil
		} else {
			return url[strings.Index(url, host)+1+len(host):], nil
		}
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
	default:
		return "", errors.Errorf("provider %s not supported for browser urls", p)
	}
}

func OpenInBrowser(url string) {
	if url != "" {
		_ = open.Run(url)
	}
}
