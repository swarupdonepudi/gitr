package clone

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/leftbin/go-util/pkg/file"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	intssh "github.com/swarupdonepudi/gitr/internal/ssh"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"github.com/swarupdonepudi/gitr/pkg/ui"
	"github.com/swarupdonepudi/gitr/pkg/url"
	"golang.org/x/crypto/ssh"
)

func Clone(cfg *config.GitrConfig, inputUrl string, token string, creDir, dry bool) (repoLocation string, err error) {
	// Strip query parameters and fragments from URLs (handles browser URLs with tracking params like ?utm_source=...)
	inputUrl = url.StripQueryParams(inputUrl)

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
			ui.RepoAlreadyExists(repoLocation)
			return repoLocation, nil
		}
		if err := os.RemoveAll(repoLocation); err != nil {
			return "", errors.Wrapf(err, "failed to remove %s dir", repoLocation)
		}
	}
	if url.IsGitUrl(inputUrl) {
		if url.IsGitSshUrl(inputUrl) {
			ui.Cloning(inputUrl)
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
		ui.Warn("Unsupported URL Format", "gitr does not support clone using browser URLs for BitBucket. Please use SSH or HTTPS clone URLs instead.")
		return "", nil
	}
	sshCloneUrl := GetSshCloneUrl(s.Hostname, repoPath)
	ui.Cloning(sshCloneUrl)
	if err := sshClone(sshCloneUrl, repoLocation); err != nil {
		// Check if the error indicates the repository doesn't exist
		// In this case, HTTP fallback won't help and would show a confusing auth error
		if isRepoNotFoundError(err) {
			return "", errors.New("repository not found. Please verify the URL exists and you have access")
		}
		// Clean up the directory from failed SSH clone before trying HTTP
		if err := os.RemoveAll(repoLocation); err != nil {
			log.Debugf("failed to clean up directory after SSH clone failure: %v", err)
		}
		log.Debugf("SSH clone failed, trying HTTP fallback: %v", err)
		httpCloneUrl := GetHttpCloneUrl(s.Hostname, repoPath, s.Scheme)
		if err := httpClone(httpCloneUrl, repoLocation); err != nil {
			return "", errors.Wrap(err, "error cloning the repo using http")
		}
	}
	return repoLocation, nil
}

// isRepoNotFoundError checks if the error indicates the repository doesn't exist
func isRepoNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	// These patterns are specific to repository-not-found errors from git hosts.
	// Avoid overly broad patterns like "not found" which match unrelated errors
	// (e.g., "host not found", "key not found", network errors).
	notFoundPatterns := []string{
		"repository not found",
		"repo not found",
		"remote: repository not found",
		"project not found",
		"the project you were looking for could not be found", // GitLab
		"error: repository '", // Start of git's "repository 'X' not found" message
	}
	for _, pattern := range notFoundPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}

func GetClonePath(cfg *config.GitrConfig, inputUrl string, creDir bool) (string, error) {
	// Strip query parameters and fragments from URLs (handles browser URLs with tracking params like ?utm_source=...)
	inputUrl = url.StripQueryParams(inputUrl)

	s, err := config.GetScmHost(cfg, url.GetHostname(inputUrl))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get scm host for %s", url.GetHostname(inputUrl))
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

	// Create fancy progress display
	display := ui.NewCloneProgressDisplay()
	display.Start()

	// go-git will output progress
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      url,
		Progress: display.Writer(),
	})

	// Stop the progress display
	display.Stop()

	return err
}

func httpsGitClone(repoUrl, token, clonePath string) error {
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to created dir %s", clonePath)
	}

	// Create fancy progress display
	display := ui.NewCloneProgressDisplay()
	display.Start()

	// go-git will output progress
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: display.Writer(),
		Auth: &http.BasicAuth{
			Username: "abc123", // this can be anything except an empty string
			Password: token,
		},
	})

	// Stop the progress display
	display.Stop()

	if err != nil {
		return errors.Wrapf(err, "failed to clone repo using personal access token %s", token)
	}
	return err
}

func sshClone(repoUrl, clonePath string) error {
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create dir %s", clonePath)
	}

	// Create fancy progress display
	display := ui.NewCloneProgressDisplay()
	display.Start()

	// Use --progress to force git to show progress even when stderr is not a TTY
	cmd := exec.Command("git", "clone", "--progress", repoUrl, clonePath)

	// Capture stderr for both progress parsing and error detection
	var stderrBuf strings.Builder
	cmd.Stderr = io.MultiWriter(display.Writer(), &stderrBuf)
	cmd.Stdout = io.Discard // Suppress stdout since we're showing fancy progress

	err := cmd.Run()

	// Stop the progress display
	display.Stop()

	if err != nil {
		stderrStr := stderrBuf.String()
		// Include stderr in the error so we can detect specific failure reasons
		return errors.Errorf("clone failed: %s", stderrStr)
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
