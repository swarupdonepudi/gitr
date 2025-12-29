package config_test

import (
	"github.com/swarupdonepudi/gitr/pkg/config"
	"testing"
)

func TestGitrConfig(t *testing.T) {
	defaultCloneConfig := &config.CloneConfig{
		HomeDir:              "",
		AlwaysCreDir:         false,
		IncludeHostForCreDir: false,
	}

	gc1 := &config.GitrConfig{
		Scm: &config.Scm{
			Hosts: []*config.ScmHost{
				{Scheme: config.Https, Hostname: "github.com", Provider: config.GitHub, DefaultBranch: "master", Clone: defaultCloneConfig},
				{Scheme: config.Https, Hostname: "gitlab.com", Provider: config.GitLab, DefaultBranch: "main", Clone: defaultCloneConfig},
				{Scheme: config.Https, Hostname: "bitbucket.org", Provider: config.BitBucketCloud, DefaultBranch: "master", Clone: defaultCloneConfig},
				{Scheme: config.Https, Hostname: "gitlab.mycompany.com", Provider: config.GitLab, DefaultBranch: "main", Clone: defaultCloneConfig},
			},
		},
	}

	t.Run("scm provider lookup using gitlab.com", func(t *testing.T) {
		scm, _ := config.GetScmHost(gc1, "gitlab.com")
		if scm.Provider != config.GitLab {
			t.Errorf("expecting %s provider and got %s", config.GitLab, scm.Provider)
		}
	})
	t.Run("scm provider lookup using github.com", func(t *testing.T) {
		scm, _ := config.GetScmHost(gc1, "github.com")
		if scm.Provider != config.GitHub {
			t.Errorf("expecting %s provider and got %s", config.GitHub, scm.Provider)
		}
	})
	t.Run("scm provider lookup using bitbucket.org", func(t *testing.T) {
		scm, _ := config.GetScmHost(gc1, "bitbucket.org")
		if scm.Provider != config.BitBucketCloud {
			t.Errorf("expecting %s provider and got %s", config.BitBucketCloud, scm.Provider)
		}
	})
	t.Run("scm provider lookup using for custom internal host", func(t *testing.T) {
		scm, _ := config.GetScmHost(gc1, "gitlab.mycompany.com")
		if scm.Provider != config.GitLab {
			t.Errorf("expecting %s provider and got %s", config.GitLab, scm.Provider)
		}
	})
	t.Run("default branch for gitlab.com should be main", func(t *testing.T) {
		scm, _ := config.GetScmHost(gc1, "gitlab.com")
		if scm.DefaultBranch != "main" {
			t.Errorf("expecting default branch to be %s and got %s", "main", scm.DefaultBranch)
		}
	})
}
