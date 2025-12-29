package web

import (
	"fmt"
	"github.com/pkg/errors"
	gitrgit "github.com/swarupdonepudi/gitr/pkg/git"
	"github.com/swarupdonepudi/gitr/pkg/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/swarupdonepudi/gitr/pkg/config"
)

// GetFileURL returns the browser URL for a single file in the repo.
//
//	base   – repo web URL, e.g. https://github.com/org/repo
//	ref    – branch name or commit SHA
//	rel    – path inside repo, **forward-slash** format
//
// Provider-specific rules:
//
//	GitHub             : <base>/blob/<ref>/<rel>
//	GitLab             : <base>/-/blob/<ref>/<rel>
//	Bitbucket Cloud/DC : <base>/src/<ref>/<rel>
func GetFileURL(p config.ScmProvider, base, ref, rel string) string {
	rel = strings.TrimPrefix(rel, "/") // safety

	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/blob/%s/%s", base, ref, rel)
	case config.BitBucketCloud, config.BitBucketDatacenter:
		return fmt.Sprintf("%s/src/%s/%s", base, ref, rel)
	default: // GitHub and similar
		return fmt.Sprintf("%s/blob/%s/%s", base, ref, rel)
	}
}

// FileURLFromPwd returns the provider-specific web URL for fileName,
// where fileName is given **relative to the current working directory**.
func FileURLFromPwd(fileName string) (string, error) {
	wd, _ := os.Getwd()

	// repo, remote, ref
	repo, err := gitrgit.GetGitRepo(wd)
	if err != nil {
		return "", errors.Wrap(err, "git repo not found")
	}
	remote, err := gitrgit.GetGitRemoteUrl(repo)
	if err != nil {
		return "", errors.Wrap(err, "remote URL not found")
	}
	ref, err := gitrgit.GetGitBranch(repo)
	if err != nil {
		return "", errors.Wrap(err, "branch not found")
	}

	// provider & base URL
	cfg, err := config.NewGitrConfig()
	if err != nil {
		return "", err
	}
	host := url.GetHostname(remote)
	hostCfg, err := config.GetScmHost(cfg, host)
	if err != nil {
		return "", err
	}
	repoPath, err := url.GetRepoPath(remote, host, hostCfg.Provider)
	if err != nil {
		return "", err
	}
	base := GetWebUrl(hostCfg.Provider, hostCfg.Scheme, host, repoPath)

	// path inside repo
	wt, _ := repo.Worktree()
	abs := filepath.Join(wd, fileName)
	rel, err := filepath.Rel(wt.Filesystem.Root(), abs)
	if err != nil {
		return "", errors.Wrap(err, "relative path calc failed")
	}

	// final link
	return GetFileURL(hostCfg.Provider, base, ref, filepath.ToSlash(rel)), nil
}
