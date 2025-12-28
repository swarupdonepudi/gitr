package url_test

import (
	"github.com/plantoncloud/gitr/pkg/config"
	"github.com/plantoncloud/gitr/pkg/url"
	"testing"
)

func TestIsGitUrl(t *testing.T) {
	var positiveUrlTests = []struct {
		url      string
		isGitUrl bool
	}{
		{"git@github.com:swarupdonepudi/gitr.git", true},
		{"https://github.com/plantoncloud/gitr.git", true},
		{"https://github.com/plantoncloud/gitr", false},
		{"git@github.com:swarupdonepudi/gitr", false},
	}
	var negativeUrlTests = []struct {
		url      string
		isGitUrl bool
	}{
		{"https://github.com/plantoncloud/gitr", false},
		{"git@github.com:swarupdonepudi/gitr", false},
	}
	t.Run("urls ending with .git should be git urls", func(t *testing.T) {
		for _, u := range positiveUrlTests {
			if url.IsGitUrl(u.url) != u.isGitUrl {
				t.Errorf("expected url %s as git url", u.url)
			}
		}
	})
	t.Run("urls not ending with .git should not be git urls", func(t *testing.T) {
		for _, u := range negativeUrlTests {
			if url.IsGitUrl(u.url) != u.isGitUrl {
				t.Errorf("expected url %s as not git url", u.url)
			}
		}
	})
}

func TestIsGitSshUrl(t *testing.T) {
	var positiveUrlTests = []struct {
		url         string
		isGitSshUrl bool
	}{
		{"git@github.com:swarupdonepudi/gitr.git", true},
		{"ssh://github.com/plantoncloud/gitr.git", true},
	}
	var negativeUrlTests = []struct {
		url         string
		isGitSshUrl bool
	}{
		{"https://github.com/plantoncloud/gitr", false},
		{"github.com:swarupdonepudi/gitr.git", false},
	}
	t.Run("urls prefixed with ssh or git should be git ssh urls", func(t *testing.T) {
		for _, u := range positiveUrlTests {
			if url.IsGitSshUrl(u.url) != u.isGitSshUrl {
				t.Errorf("expected url %s as git url", u.url)
			}
		}
	})
	t.Run("urls not prefixed with ssh or git should not be git ssh urls", func(t *testing.T) {
		for _, u := range negativeUrlTests {
			if url.IsGitSshUrl(u.url) != u.isGitSshUrl {
				t.Errorf("expected url %s as not git url", u.url)
			}
		}
	})
}

func TestIsGitHttpUrlHasUsername(t *testing.T) {
	var usernameTests = []struct {
		url         string
		hasUsername bool
	}{
		{"https://swarup@github.com:swarupdonepudi/gitr.git", true},
		{"https://swarupd@github.com:swarupdonepudi/gitr", true},
		{"https://github.com/plantoncloud/gitr", false},
		{"github.com:swarupdonepudi/gitr.git", false},
	}

	t.Run("username in http url", func(t *testing.T) {
		for _, u := range usernameTests {
			if url.IsGitHttpUrlHasUsername(u.url) != u.hasUsername {
				t.Errorf("expected %v but received %v for %s ", u.hasUsername, url.IsGitHttpUrlHasUsername(u.url), u.url)
			}
		}
	})
}

func TestIsGitRepoName(t *testing.T) {
	var repoNameTests = []struct {
		repoPath string
		repoName string
	}{
		{"swarupdonepudi/gitr.git", "gitr.git"},
		{"parent/sub/sub2/project-name.git", "project-name.git"},
		{"parent/sub/sub2/sub/project-name.git", "project-name.git"},
		{"no-path.git", "no-path.git"},
		{"parent/sub/git-repo", "git-repo"},
		{"parent/git-repo", "git-repo"},
		{"git-repo", "git-repo"},
	}

	t.Run("repo name from repo path", func(t *testing.T) {
		for _, u := range repoNameTests {
			if url.GetRepoName(u.repoPath) != u.repoName {
				t.Errorf("expected %s but got %s for %s path", u.repoName, url.GetRepoName(u.repoPath), u.repoPath)
			}
		}
	})
}

func TestGetRepoPath(t *testing.T) {
	var githubTests = []struct {
		url      string
		host     string
		expected string
	}{
		// Basic repository URLs
		{"https://github.com/owner/repo", "github.com", "owner/repo"},
		{"https://github.com/owner/repo.git", "github.com", "owner/repo"},
		// Tree URLs (branch/directory browsing)
		{"https://github.com/sarwarbeing-ai/Agentic_Design_Patterns/tree/main", "github.com", "sarwarbeing-ai/Agentic_Design_Patterns"},
		{"https://github.com/owner/repo/tree/feature-branch", "github.com", "owner/repo"},
		{"https://github.com/owner/repo/tree/main/src/pkg", "github.com", "owner/repo"},
		// Blob URLs (file viewing)
		{"https://github.com/owner/repo/blob/main/README.md", "github.com", "owner/repo"},
		{"https://github.com/owner/repo/blob/v1.0.0/file.go", "github.com", "owner/repo"},
		// Commits URLs
		{"https://github.com/owner/repo/commits/main", "github.com", "owner/repo"},
		// Pull request URLs
		{"https://github.com/owner/repo/pull/123", "github.com", "owner/repo"},
		// Issues URLs
		{"https://github.com/owner/repo/issues/456", "github.com", "owner/repo"},
		// Compare URLs
		{"https://github.com/owner/repo/compare/main...feature", "github.com", "owner/repo"},
	}

	t.Run("GitHub URL patterns should extract correct repo path", func(t *testing.T) {
		for _, tc := range githubTests {
			result, err := url.GetRepoPath(tc.url, tc.host, config.GitHub)
			if err != nil {
				t.Errorf("unexpected error for url %s: %v", tc.url, err)
				continue
			}
			if result != tc.expected {
				t.Errorf("expected %s but got %s for url %s", tc.expected, result, tc.url)
			}
		}
	})
}
