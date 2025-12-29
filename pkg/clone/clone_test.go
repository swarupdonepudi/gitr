package clone

import (
	"github.com/swarupdonepudi/gitr/pkg/config"
	"testing"
)

func TestScmHome(t *testing.T) {
	var tests = []struct {
		scmHostHomeDir string
		scmHomeDir     string
		expectedVal    string
		expectErr      error
	}{
		{scmHostHomeDir: "/home/john/scm/host", scmHomeDir: "", expectedVal: "/home/john/scm/host", expectErr: nil},
		{scmHostHomeDir: "", scmHomeDir: "/home/john/scm", expectedVal: "/home/john/scm", expectErr: nil},
		{scmHostHomeDir: "/home/john/scm/host/custom/path", scmHomeDir: "/home/john/scm", expectedVal: "/home/john/scm/host/custom/path", expectErr: nil},
	}
	t.Run("validate scm home evaluation", func(t *testing.T) {
		for _, test := range tests {
			returnedVal, err := getScmHome(test.scmHostHomeDir, test.scmHomeDir)
			if err != test.expectErr {
				t.Errorf("expecting %v err but got %v err", test.expectErr, err)
			}
			if returnedVal != test.expectedVal {
				t.Errorf("expecting %s but got %s", test.expectedVal, returnedVal)
			}
		}
	})
}

func TestGetClonePath(t *testing.T) {
	type getClonePathInput struct {
		cfg      *config.GitrConfig
		inputUrl string
		creDir   bool
	}
	var tests = []struct {
		testName    string
		input       *getClonePathInput
		expectPath  string
		expectedErr error
	}{
		{
			testName: "creDir=false, alwaysCreDir=false, includeHost=false, the repo should be cloned to scm home dir",
			input: &getClonePathInput{cfg: &config.GitrConfig{Scm: &config.Scm{HomeDir: "/Users/joe/scm", Hosts: []*config.ScmHost{{Hostname: "github.com", Provider: config.GitHub, Clone: &config.CloneConfig{
				HomeDir:              "",
				AlwaysCreDir:         false,
				IncludeHostForCreDir: true,
			}}}}}, inputUrl: "https://github.com/kubernetes-sigs/kind", creDir: false},
			expectPath:  "/Users/joe/scm/kind",
			expectedErr: nil,
		},
		{
			testName: "creDir=false, alwaysCreDir=true, includeHost=false then clone loc should follow repo path on scm",
			input: &getClonePathInput{cfg: &config.GitrConfig{Scm: &config.Scm{HomeDir: "/Users/joe/scm", Hosts: []*config.ScmHost{{Hostname: "github.com", Provider: config.GitHub, Clone: &config.CloneConfig{
				HomeDir:              "",
				AlwaysCreDir:         true,
				IncludeHostForCreDir: false,
			}}}}}, inputUrl: "https://github.com/kubernetes-sigs/kind", creDir: false},
			expectPath:  "/Users/joe/scm/kubernetes-sigs/kind",
			expectedErr: nil,
		},
		{
			testName: "creDir=true, alwaysCreDir=false, includeHost=false then clone loc should follow repo path on scm",
			input: &getClonePathInput{cfg: &config.GitrConfig{Scm: &config.Scm{HomeDir: "/Users/joe/scm", Hosts: []*config.ScmHost{{Hostname: "github.com", Provider: config.GitHub, Clone: &config.CloneConfig{
				HomeDir:              "",
				AlwaysCreDir:         false,
				IncludeHostForCreDir: false,
			}}}}}, inputUrl: "https://github.com/kubernetes-sigs/kind", creDir: true},
			expectPath:  "/Users/joe/scm/kubernetes-sigs/kind",
			expectedErr: nil,
		},
		{
			testName: "creDir=true, alwaysCreDir=false, includeHost=true then clone loc should mimic scm",
			input: &getClonePathInput{cfg: &config.GitrConfig{Scm: &config.Scm{HomeDir: "/Users/joe/scm", Hosts: []*config.ScmHost{{Hostname: "github.com", Provider: config.GitHub, Clone: &config.CloneConfig{
				HomeDir:              "",
				AlwaysCreDir:         false,
				IncludeHostForCreDir: true,
			}}}}}, inputUrl: "https://github.com/kubernetes-sigs/kind", creDir: true},
			expectPath:  "/Users/joe/scm/github.com/kubernetes-sigs/kind",
			expectedErr: nil,
		},
		{
			testName: "creDir=true, alwaysCreDir=true, includeHost=true then clone loc should mimic scm",
			input: &getClonePathInput{cfg: &config.GitrConfig{Scm: &config.Scm{HomeDir: "/Users/joe/scm", Hosts: []*config.ScmHost{{Hostname: "github.com", Provider: config.GitHub, Clone: &config.CloneConfig{
				HomeDir:              "",
				AlwaysCreDir:         false,
				IncludeHostForCreDir: true,
			}}}}}, inputUrl: "https://github.com/kubernetes-sigs/kind", creDir: true},
			expectPath:  "/Users/joe/scm/github.com/kubernetes-sigs/kind",
			expectedErr: nil,
		},
	}

	t.Run("test get clone path", func(t *testing.T) {
		for _, tc := range tests {
			t.Run(tc.testName, func(t *testing.T) {
				returnedVal, err := GetClonePath(tc.input.cfg, tc.input.inputUrl, tc.input.creDir)
				if err != tc.expectedErr {
					t.Errorf("expecting %v err but got %v err", tc.expectedErr, err)
				}
				if returnedVal != tc.expectPath {
					t.Errorf("expecting %s but got %s", tc.expectPath, returnedVal)
				}
			})
		}
	})
}
