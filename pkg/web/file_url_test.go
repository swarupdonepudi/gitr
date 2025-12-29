package web

import (
	"github.com/swarupdonepudi/gitr/pkg/config"
	"testing"
)

func TestGetFileURL(t *testing.T) {
	cases := []struct {
		p       config.ScmProvider
		base    string
		ref     string
		rel     string
		wantURL string
	}{
		{config.GitHub, "https://github.com/acme/repo", "main", "docs/readme.md",
			"https://github.com/acme/repo/blob/main/docs/readme.md"},
		{config.GitLab, "https://gitlab.com/acme/repo", "main", "docs/readme.md",
			"https://gitlab.com/acme/repo/-/blob/main/docs/readme.md"},
		{config.BitBucketCloud, "https://bitbucket.org/acme/repo", "main", "docs/readme.md",
			"https://bitbucket.org/acme/repo/src/main/docs/readme.md"},
	}

	for _, c := range cases {
		if got := GetFileURL(c.p, c.base, c.ref, c.rel); got != c.wantURL {
			t.Errorf("GetFileURL(%s) = %s, want %s", c.p, got, c.wantURL)
		}
	}
}
