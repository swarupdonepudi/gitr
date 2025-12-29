package web

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/swarupdonepudi/gitr/pkg/config"
	"os"
)

func PrintGitrWebInfo(p config.ScmProvider, host, remoteUrl, webUrl, repoPath, repoName, branch string) {
	println("")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{"remote", remoteUrl})
	t.AppendSeparator()
	t.AppendRow(table.Row{"provider", p})
	t.AppendSeparator()
	t.AppendRow(table.Row{"host", host})
	t.AppendSeparator()
	t.AppendRow(table.Row{"repo-path", repoPath})
	t.AppendSeparator()
	t.AppendRow(table.Row{"repo-name", repoName})
	t.AppendSeparator()
	t.AppendRow(table.Row{"branch", branch})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-web", webUrl})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-remote", GetRemUrl(p, webUrl, branch)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-commits", GetCommitsUrl(p, repoPath, branch)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-branches", GetBranchesUrl(p, webUrl)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-tags", GetTagsUrl(p, webUrl)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-releases", GetReleasesUrl(p, webUrl)})
	t.AppendSeparator()
	t.AppendRow(table.Row{"url-pipelines", GetPipelinesUrl(p, webUrl)})
	t.AppendSeparator()
	t.Render()
	println("")
}

func GetWebUrl(p config.ScmProvider, scheme config.HttpScheme, host, repoPath string) string {
	switch p {
	default:
		return fmt.Sprintf("%s://%s/%s", scheme, host, repoPath)
	}
}

func GetRemUrl(p config.ScmProvider, webUrl, repoBranch string) string {
	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/tree/%s", webUrl, repoBranch)
	case config.BitBucketCloud:
		return fmt.Sprintf("%s/branch/%s", webUrl, repoBranch)
	default:
		return fmt.Sprintf("%s/tree/%s", webUrl, repoBranch)
	}
}

func GetPrsUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.GitHub:
		return fmt.Sprintf("%s/pulls", webUrl)
	case config.GitLab:
		return fmt.Sprintf("%s/-/merge_requests", webUrl)
	case config.BitBucketDatacenter, config.BitBucketCloud:
		return fmt.Sprintf("%s/pull-requests", webUrl)
	default:
		return ""
	}
}

func GetBranchesUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/branches", webUrl)
	default:
		return fmt.Sprintf("%s/branches", webUrl)
	}
}

func GetCommitsUrl(p config.ScmProvider, webUrl, repoBranch string) string {
	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/commits/%s", webUrl, repoBranch)
	default:
		return fmt.Sprintf("%s/commits/%s", webUrl, repoBranch)
	}
}

func GetTagsUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/tags", webUrl)
	default:
		return fmt.Sprintf("%s/tags", webUrl)
	}
}

func GetIssuesUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.BitBucketDatacenter, config.BitBucketCloud:
		return ""
	case config.GitLab:
		return fmt.Sprintf("%s/-/issues", webUrl)
	default:
		return fmt.Sprintf("%s/issues", webUrl)
	}
}

func GetReleasesUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.GitHub:
		return fmt.Sprintf("%s/releases", webUrl)
	case config.GitLab:
		return fmt.Sprintf("%s/-/releases", webUrl)
	default:
		return ""
	}
}

func GetPipelinesUrl(p config.ScmProvider, webUrl string) string {
	switch p {
	case config.GitLab:
		return fmt.Sprintf("%s/-/pipelines", webUrl)
	case config.GitHub:
		return fmt.Sprintf("%s/actions", webUrl)
	case config.BitBucketCloud:
		return fmt.Sprintf("%s/addon/pipelines/home", webUrl)
	default:
		return ""
	}
}
