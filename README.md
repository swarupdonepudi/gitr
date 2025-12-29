# gitr

<div align="center">

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev)
[![GitHub](https://img.shields.io/badge/GitHub-supported-181717?logo=github)](https://github.com)
[![GitLab](https://img.shields.io/badge/GitLab-supported-FC6D26?logo=gitlab)](https://gitlab.com)
[![Bitbucket](https://img.shields.io/badge/Bitbucket-supported-0052CC?logo=bitbucket)](https://bitbucket.org)

**Clone to organized paths. Open PRs, pipelines, branches instantly. One CLI, zero browser tabs.**

[Quick Start](#quick-start) • [Features](#features) • [CLI Reference](#cli-reference) • [Website](https://swarupdonepudi.github.io/gitr) • [Contributing](#contributing)

</div>

---

## What is gitr?

`gitr` solves two daily frustrations: **1)** repos scattered everywhere, **2)** clicking through GitHub/GitLab/Bitbucket to find PRs, pipelines, issues.

```bash
$ gitr clone git@github.com:owner/repo.git  # → ~/scm/github.com/owner/repo
$ gitr prs        # Opens PRs/MRs in browser
$ gitr pipe       # Opens pipelines/actions
```

**[📚 See it in action on the website →](https://swarupdonepudi.github.io/gitr)**

---

## Why gitr?

| Before gitr 😫 | After gitr ✨ |
|---|---|
| "Where should I clone this repo?" | `gitr clone <url>` → deterministic path |
| Create folders manually | Auto-creates `host/owner/repo` structure |
| Click through tabs to find PRs | `gitr prs` → instant navigation |
| Hunt for that repo you cloned | Always at `~/scm/{provider}/{owner}/{repo}` |

---

## Quick Start

### Install
**macOS:** `brew install swarupdonepudi/tap/gitr`  
**Go:** `go install github.com/swarupdonepudi/gitr@latest`  
**Binary:** Download from [releases](https://github.com/swarupdonepudi/gitr/releases)

### Usage
```bash
gitr clone https://github.com/kubernetes/kubernetes
# → Clones to: ~/scm/github.com/kubernetes/kubernetes

cd ~/scm/github.com/kubernetes/kubernetes
gitr web          # Opens repo homepage
gitr prs          # Opens pull requests
```

**[🎯 Full tutorial →](https://swarupdonepudi.github.io/gitr#quickstart)**

---

## Features

| Feature | Description |
|---------|-------------|
| 🗂️ **Organized Cloning** | Clone repos to `~/scm/{host}/{owner}/{repo}` structure |
| 🌐 **Instant Web Nav** | Open PRs, pipelines, issues, branches from terminal |
| 🏢 **Enterprise Ready** | Works with on-prem GitHub/GitLab/Bitbucket |
| 🔐 **Multi Auth** | SSH keys + HTTPS tokens support |
| 👀 **Dry Run** | Preview paths/URLs with `gitr --dry <command>` |

**[📖 Full feature documentation →](https://swarupdonepudi.github.io/gitr)**

---

## CLI Reference

### Clone Commands
```bash
gitr clone <url>              # Clone to deterministic path
gitr clone <url> -c           # Create full directory hierarchy
gitr clone <url> --dry        # Preview without cloning
gitr clone <url> --token=xxx  # Clone with HTTPS token
```

### Web Navigation Commands
**Run inside any git repository:**

| Command | Opens |
|---------|-------|
| `gitr web` | Repository homepage |
| `gitr rem` | Current branch in web UI |
| `gitr prs` | Pull Requests / Merge Requests |
| `gitr pipe` | Pipelines / Actions |
| `gitr issues` | Issues |
| `gitr commits` | Commits for current branch |
| `gitr branches` | All branches |
| `gitr tags` | All tags |
| `gitr releases` | Releases page |

### Utility Commands
```bash
gitr config show    # Show current configuration
gitr config edit    # Edit ~/.gitr.yaml in $EDITOR
gitr path <url>     # Show deterministic path for URL
gitr --dry <cmd>    # Preview mode (no changes)
```

**[📖 Complete CLI docs →](https://swarupdonepudi.github.io/gitr#cli)**

---

## Configuration

`gitr` auto-creates `~/.gitr.yaml` on first run. Quick example:

```yaml
scm:
  homeDir: /Users/you/scm
  hosts:
    - hostname: github.com
      provider: github
      clone:
        alwaysCreDir: true
        includeHostForCreDir: true
    - hostname: gitlab.mycompany.net  # On-prem support
      provider: gitlab
      scheme: https
```

**Supports:** On-prem instances • Per-host clone rules • SSH config (`~/.ssh/config`) • HTTPS tokens (`~/.personal_access_tokens/{hostname}`)

**[⚙️ Full configuration guide →](https://swarupdonepudi.github.io/gitr#cli)**

---

## Supported Providers

✅ **GitHub** (github.com + Enterprise) • ✅ **GitLab** (gitlab.com + Self-hosted) • ✅ **Bitbucket** (bitbucket.org + Datacenter)

---

## Example: Organized Workspace

After using `gitr clone`:

```
~/scm/
├── github.com/kubernetes/kubernetes/
├── github.com/swarupdonepudi/gitr/
├── gitlab.com/team/project/backend/
└── gitlab.mycompany.net/org/infra/terraform/
```

**Power user aliases** (add to `.zshrc`):
```bash
alias clone="gitr clone"
alias prs="gitr prs"
alias pipe="gitr pipe"
```

---

## Links

📚 **[Documentation](https://swarupdonepudi.github.io/gitr)** • 🎯 **[Tutorial](https://swarupdonepudi.github.io/gitr#quickstart)** • ❓ **[FAQ](https://swarupdonepudi.github.io/gitr#faq)** • 📦 **[Releases](https://github.com/swarupdonepudi/gitr/releases)**

---

## Contributing

`gitr` was built to share extreme productivity with other productivity geeks. Issues and pull requests are welcome!

```bash
gitr clone https://github.com/swarupdonepudi/gitr
make build && make test
```

---

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

---

<div align="center">

**Built with ❤️ for developers who value their time**

[⭐ Star on GitHub](https://github.com/swarupdonepudi/gitr) • [🌐 Visit Website](https://swarupdonepudi.github.io/gitr)

</div>
