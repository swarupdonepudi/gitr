package clone

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/leftbin/go-util/pkg/file"
	"github.com/leftbin/go-util/pkg/shell"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	intssh "github.com/swarupdonepudi/gitr/internal/ssh"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/url"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func Clone(cfg *config.GitrConfig, inputUrl string, token string, creDir, dry bool) (repoLocation string, err error) {
	s, err := config.GetScmHost(cfg, url.GetHostname(inputUrl))
	if err != nil {
		return "", errors.Wrapf(err, "failed to clone git repo with %s url", inputUrl)
	}
	repoPath, err := url.GetRepoPath(inputUrl, s.Hostname, s.Provider)
	if err != nil {
		return "", errors.Wrap(err, "failed to get repo path")
	}
	repoLocation, err = GetClonePath(cfg, inputUrl, creDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get clone path")
	}
	if dry {
		err := printGitrCloneInfo(cfg, inputUrl, creDir || s.Clone.AlwaysCreDir)
		if err != nil {
			return "", errors.Wrap(err, "failed to print gitr clone info")
		}
		return repoLocation, nil
	}
	if file.IsDirExists(repoLocation) {
		if file.IsDirExists(filepath.Join(repoLocation, ".git")) {
			log.Info("repo already exists. skipping cloning...")
			return repoLocation, nil
		}
		if err := os.RemoveAll(repoLocation); err != nil {
			return "", errors.Wrapf(err, "failed to remove %s dir", repoLocation)
		}
	}
	if url.IsGitUrl(inputUrl) {
		if url.IsGitSshUrl(inputUrl) {
			log.Debugf("cloning using ssh url %s", inputUrl)
			if err := sshClone(inputUrl, repoLocation); err != nil {
				return "", errors.Wrap(err, "error cloning the repo")
			}
			return repoLocation, nil
		}
		if token == "" {
			token, err = getHttpsCloneToken(s.Hostname)
			if err != nil {
				return "", errors.Wrap(err, "failed to check if https clone token is configured")
			}
		}
		if token != "" {
			if err := httpsGitClone(inputUrl, token, repoLocation); err != nil {
				return "", errors.Wrap(err, "error cloning the repo")
			}
			return repoLocation, nil
		}

	}
	if s.Provider == config.BitBucketDatacenter || s.Provider == config.BitBucketCloud {
		log.Warn("gitr does not support clone using browser urls for bitbucket-datacenter & bitbucket.org")
		return "", nil
	}
	sshCloneUrl := GetSshCloneUrl(s.Hostname, repoPath)
	log.Debugf("cloning using ssh url %s", sshCloneUrl)
	if err := sshClone(sshCloneUrl, repoLocation); err != nil {
		log.Warn("failed to clone repo using ssh. trying http clone...")
		httpCloneUrl := GetHttpCloneUrl(s.Hostname, repoPath, s.Scheme)
		if err := httpClone(httpCloneUrl, repoLocation); err != nil {
			return "", errors.Wrap(err, "error cloning the repo using http")
		}
	}
	return repoLocation, nil
}

func GetClonePath(cfg *config.GitrConfig, inputUrl string, creDir bool) (string, error) {
	s, err := config.GetScmHost(cfg, url.GetHostname(inputUrl))
	if err != nil {
		log.Fatalf("failed to get scm host. err: %v", err)
	}
	repoPath, err := url.GetRepoPath(inputUrl, s.Hostname, s.Provider)
	if err != nil {
		return "", errors.Wrap(err, "failed to get repo path")
	}
	repoName := url.GetRepoName(repoPath)
	scmHome, err := getScmHome(s.Clone.HomeDir, cfg.Scm.HomeDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get scm home dir")
	}
	clonePath := ""
	if creDir || s.Clone.AlwaysCreDir {
		if s.Clone.IncludeHostForCreDir {
			clonePath = fmt.Sprintf("%s/%s", s.Hostname, repoPath)
		} else {
			clonePath = repoPath
		}
	} else {
		clonePath = repoName
	}
	if scmHome != "" {
		clonePath = fmt.Sprintf("%s/%s", scmHome, clonePath)
	}
	return clonePath, nil
}

func setUpSshAuth(hostname string) (*ssh2.PublicKeys, error) {
	sshKeyPath, err := intssh.GetKeyPath(hostname)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ssh config path")
	}
	if !file.IsFileExists(sshKeyPath) {
		log.Debugf("%s file not found", sshKeyPath)
		return nil, errors.Errorf("ssh auth not found")
	}
	pem, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s file", sshKeyPath)
	}
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse private key %s", sshKeyPath)
	}
	return &ssh2.PublicKeys{User: "git", Signer: signer}, nil
}

func getHttpsCloneToken(hostname string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home dir")
	}
	pAccessTokenFilePath := filepath.Join(homeDir, ".personal_access_tokens", hostname)
	pAccessTokenFileAbsPath, err := file.GetAbsPath(pAccessTokenFilePath)
	if !file.IsFileExists(pAccessTokenFileAbsPath) {
		return "", nil
	}
	pem, err := ioutil.ReadFile(pAccessTokenFileAbsPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %s file", pAccessTokenFileAbsPath)
	}
	return string(pem), nil
}

func httpClone(url, clonePath string) error {
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to created dir %s", clonePath)
	}
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func httpsGitClone(repoUrl, token, clonePath string) error {
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to created dir %s", clonePath)
	}
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: "abc123", // this can be anything except an empty string
			Password: token,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to clone repo using personal access token %s", token)
	}
	return err
}

func sshClone(repoUrl, clonePath string) error {
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create dir %s", clonePath)
	}
	if err := shell.RunCmd(exec.Command("git", "clone", repoUrl, clonePath)); err != nil {
		return errors.Wrapf(err, "failed to clone")
	}
	return nil
}

func GetSshCloneUrl(hostname, repoPath string) string {
	return fmt.Sprintf("git@%s:%s.git", hostname, repoPath)
}

func GetHttpCloneUrl(hostname, repoPath string, scheme config.HttpScheme) string {
	return fmt.Sprintf("%s://%s/%s.git", scheme, hostname, repoPath)
}

func printGitrCloneInfo(cfg *config.GitrConfig, inputUrl string, creDir bool) error {
	s, err := config.GetScmHost(cfg, url.GetHostname(inputUrl))
	repoPath, err := url.GetRepoPath(inputUrl, s.Hostname, s.Provider)
	if err != nil {
		return errors.Wrap(err, "failed to get repo path")
	}
	repoName := url.GetRepoName(repoPath)
	scmHome, err := getScmHome(s.Clone.HomeDir, cfg.Scm.HomeDir)
	if err != nil {
		return errors.Wrap(err, "failed to get scm home dir")
	}
	clonePath, err := GetClonePath(cfg, inputUrl, creDir)
	if err != nil {
		return errors.Wrap(err, "failed to get clone path")
	}
	println("")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{"remote", inputUrl})
	t.AppendSeparator()
	t.AppendRow(table.Row{"provider", s.Provider})
	t.AppendSeparator()
	t.AppendRow(table.Row{"host", s.Hostname})
	t.AppendSeparator()
	t.AppendRow(table.Row{"repo-name", repoName})
	t.AppendSeparator()
	t.AppendRow(table.Row{"ssh-url", GetSshCloneUrl(s.Hostname, repoPath)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"http-url", GetHttpCloneUrl(s.Hostname, repoPath, s.Scheme)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"create-dir", s.Clone.AlwaysCreDir || creDir})
	t.AppendSeparator()
	t.AppendRow(table.Row{"scm-home", scmHome})
	t.AppendSeparator()
	t.AppendRow(table.Row{"clone-path", clonePath})
	t.AppendSeparator()
	t.Render()
	println("")
	return nil
}

func getScmHome(scmHostHomeDir, scmHomeDir string) (string, error) {
	if scmHostHomeDir != "" {
		return scmHostHomeDir, nil
	}
	if scmHomeDir != "" {
		return scmHomeDir, nil
	}
	getwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current dir")
	}
	return getwd, nil
}
