package url

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"regexp"
	"strings"
)

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
