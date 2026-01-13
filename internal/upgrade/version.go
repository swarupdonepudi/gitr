package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	// githubReleasesURL is the GitHub API endpoint for releases
	githubReleasesURL = "https://api.github.com/repos/swarupdonepudi/gitr/releases"
	// httpTimeout is the timeout for HTTP requests
	httpTimeout = 10 * time.Second
)

// GitHubRelease represents a GitHub release from the API
type GitHubRelease struct {
	TagName    string `json:"tag_name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

// GetLatestVersion fetches the latest version from GitHub releases
func GetLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	req, err := http.NewRequest("GET", githubReleasesURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch releases: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var releases []GitHubRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return "", fmt.Errorf("failed to parse releases: %w", err)
	}

	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found")
	}

	// Find the highest semantic version
	var highest *semver.Version
	var highestTag string

	for _, release := range releases {
		if release.Draft || release.Prerelease {
			continue
		}

		v, err := semver.NewVersion(release.TagName)
		if err != nil {
			continue
		}

		if highest == nil || v.GreaterThan(highest) {
			highest = v
			highestTag = release.TagName
		}
	}

	if highestTag == "" {
		return "", fmt.Errorf("no valid releases found")
	}

	return highestTag, nil
}

// CompareVersions returns true if latestVersion is newer than currentVersion
func CompareVersions(currentVersion, latestVersion string) bool {
	if currentVersion == latestVersion {
		return false
	}

	// Normalize by removing 'v' prefix
	current := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	if current == "dev" || current == "" {
		return true
	}

	currentSem, errCurrent := semver.NewVersion(current)
	latestSem, errLatest := semver.NewVersion(latest)

	if errCurrent != nil || errLatest != nil {
		return current != latest
	}

	return latestSem.GreaterThan(currentSem)
}
