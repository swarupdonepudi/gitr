package git

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetGitRepo returns a git repository by walking the file system upwards from the provided directory
// and returns an error when no git repository is found before reaching the "/"(root) directory
func GetGitRepo(folder string) (*git.Repository, error) {
	for true {
		gitRepo, err := git.PlainOpen(folder)
		if err != nil {
			if folder == "/" {
				break
			}
			folder = filepath.Dir(folder)
			continue
		}
		return gitRepo, nil
	}
	return nil, errors.New("git repository not found in the folder tree")
}

// GetGitRemoteUrl returns the first url in the first remote found in the git repository object
// and returns an errors either if there is no remotes or if there is no urls for the first remote.
func GetGitRemoteUrl(r *git.Repository) (string, error) {
	remotes, err := r.Remotes()
	if err != nil {
		return "", errors.Wrap(err, "failed to get remotes from git repo")
	}
	if len(remotes) == 0 {
		return "", errors.New("no remotes found")
	}
	if len(remotes[0].Config().URLs) == 0 {
		return "", errors.Errorf("urls not found for %s remote", remotes[0].Config().Name)
	}
	return remotes[0].Config().URLs[0], nil
}

// GetGitBranch returns the name of the current branch
func GetGitBranch(r *git.Repository) (string, error) {
	head, err := r.Head()
	if err != nil {
		return "", errors.Wrap(err, "failed to get head from git repo")
	}
	return strings.ReplaceAll(head.Name().String(), "refs/heads/", ""), nil
}

// DoesBranchExistOnRemote checks if a branch exists on the remote repository
// by checking local remote-tracking branches (e.g., refs/remotes/origin/branch-name)
// This method uses local information and doesn't require network access or authentication
func DoesBranchExistOnRemote(r *git.Repository, branchName string) bool {
	remotes, err := r.Remotes()
	if err != nil || len(remotes) == 0 {
		log.Debugf("no remotes found")
		return false
	}

	remoteName := remotes[0].Config().Name

	// Check for remote-tracking branch (e.g., refs/remotes/origin/branch-name)
	remoteTrackingRef := "refs/remotes/" + remoteName + "/" + branchName

	refs, err := r.References()
	if err != nil {
		log.Debugf("failed to get references: %v", err)
		return false
	}

	exists := false
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().String() == remoteTrackingRef {
			exists = true
			return errors.New("found") // stop iteration
		}
		return nil
	})

	return exists
}

// GetDefaultBranch returns the default branch of the remote repository
// by checking the local remote HEAD reference (e.g., refs/remotes/origin/HEAD)
// This method uses local information and doesn't require network access or authentication
func GetDefaultBranch(r *git.Repository) (string, error) {
	remotes, err := r.Remotes()
	if err != nil || len(remotes) == 0 {
		return "", errors.New("no remotes found")
	}

	remoteName := remotes[0].Config().Name

	// Try to get the remote HEAD reference (e.g., refs/remotes/origin/HEAD)
	remoteHeadRef := "refs/remotes/" + remoteName + "/HEAD"

	ref, err := r.Reference(plumbing.ReferenceName(remoteHeadRef), true)
	if err == nil {
		// The reference was found, extract the branch name
		// It will be in format refs/remotes/origin/main
		targetRef := ref.Name().String()
		if ref.Type() == plumbing.SymbolicReference {
			targetRef = ref.Target().String()
		}

		// Extract branch name from refs/remotes/origin/main
		defaultBranch := strings.TrimPrefix(targetRef, "refs/remotes/"+remoteName+"/")
		if defaultBranch != "" && defaultBranch != targetRef {
			log.Debugf("found default branch from remote HEAD: %s", defaultBranch)
			return defaultBranch, nil
		}
	}

	// Fallback: try common default branch names by checking remote-tracking branches
	log.Debugf("remote HEAD not found, trying common defaults")
	commonDefaults := []string{"main", "master"}
	for _, defaultBranch := range commonDefaults {
		if DoesBranchExistOnRemote(r, defaultBranch) {
			log.Debugf("using common default branch: %s", defaultBranch)
			return defaultBranch, nil
		}
	}

	return "", errors.New("unable to determine default branch")
}
