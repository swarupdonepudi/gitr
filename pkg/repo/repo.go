package repo

import (
	"github.com/pkg/errors"
	"github.com/swarupdonepudi/gitr/pkg/config"
	gitrgit "github.com/swarupdonepudi/gitr/pkg/git"
	"github.com/swarupdonepudi/gitr/pkg/url"
	"os"
)

func GetRepoPathOnHost(remoteUrl string, gitrCfg *config.GitrConfig) (string, error) {
	scmHost := url.GetHostname(remoteUrl)
	scmHostCfg, err := config.GetScmHost(gitrCfg, scmHost)
	if err != nil {
		return "", errors.Wrap(err, "failed to get scm host config")
	}
	repoPath, err := url.GetRepoPath(remoteUrl, scmHost, scmHostCfg.Provider)
	if err != nil {
		return "", errors.Wrap(err, "failed to get repo path")
	}
	return repoPath, nil
}

func GetPathOnScmFromPwd() (string, error) {
	gitrCfg, err := config.NewGitrConfig()
	if err != nil {
		return "", errors.Wrap(err, "failed to read gitr config")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to read current dir")
	}
	gitRepo, err := gitrgit.GetGitRepo(pwd)
	if err != nil {
		return "", errors.Wrap(err, "failed to get git repo")
	}
	remoteUrl, err := gitrgit.GetGitRemoteUrl(gitRepo)
	if err != nil {
		return "", errors.Wrap(err, "failed to get remote url")
	}
	projectPath, err := GetRepoPathOnHost(remoteUrl, gitrCfg)
	if err != nil {
		return "", errors.Wrap(err, "failed to get repo path on host")
	}
	return projectPath, nil
}
