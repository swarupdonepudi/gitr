# gitr(git rapid): the missing link between git cli and web browsers

**tl;dr** in short, gitr does this:

You are in your terminal inside a subdirectory of a git repo, and you just pushed a commit. Now you want to see what
your readme looks like on the web page of the repo, or you want to see if the pipeline has been triggered for this
commit.

_before gitr:_

1. open web browser
2. go to gitlab.com
3. search for the repo or visually locate the repo on the home page
4. click on the repo
5. after the home page is loaded then locate the icon on the page to go to the pipelines section and click on it

_with gitr:_

from any subdirectory under any git repo in the terminal

1. `gitr web` to check readme or `gitr pipe` to check pipeline

> same for prs, tags, releases, issues, branches, commits etc...

- [why should i use it?](#why-should-i-use-it)
- [install](#install)
- [supported providers](#scm-providers)
- [features](#features)
  - [gitr web](#gitr-web)
  - [gitr clone](#gitr-clone)
- [gitr config](#config-file)
- [on-prem scm deployments](#on-prem-scm-deployments)
- [dry-run](#dry-run)
- [aliases](#aliases)
- [cleanup](#cleanup)
- [contribute](#contribute)

# why should i use it?

`gitr` reduces a ton of clicking and waiting for the git repo web pages to load by **taking you directly** to the page
that you would like to see right from the command line. You may think that it's not a lot of waiting but trust me, and
it adds up, and you will notice the difference once you start using `gitr`.

`gitr` relies on the contents of `.git` folder and combines it with the provider(think of gitlab, github and bitbucket)
knowledge that is built into it to smartly navigate you to the right web page of your repo right from the command line.

### install

`gitr` can be easily installed on mac using brew. While it can also be installed on linux, windows using
the [binary](https://github.com/swarupdonepudi/gitr/releases), it has not been tested on those platforms.

```
brew tap swarupdonepudi/homebrew-tap
brew install gitr
```

### scm providers

`gitr` currently, supports the following scm providers. We can always add support for more if there is demand.

- github
- gitlab
- bitbucket

### features

| feature     | description                                                                                                                                                                                                                                                     |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| gitr config | display gitr config                                                                                                                                                                                                                                             |
| gitr web    | open a repo and different parts of a repo in web browser from command line                                                                                                                                                                                      |
| gitr clone  | organize git repos cloned from different scm providers and also retain their hierarchy on the scm provider on laptops which is not possible with the default `git clone <clone-url>`. This is particularly useful for gitlab as it supports a nested hierarchy. |

#### gitr web

`gitr` opens the following parts of your git repo in the web browser right from the command line

> note: the below commands will work only when executed from root or subdirectory of any git repo

| command       | description                                             |
| ------------- | ------------------------------------------------------- |
| gitr web      | open home page of the repo in the browser               |
| gitr rem      | open local checkout branch of the repo in the browser   |
| gitr prs      | open prs/mrs of the repo in the browser                 |
| gitr pipe     | open pipelines/actions of the repo in the browser       |
| gitr issues   | open issues of the repo in the browser                  |
| gitr releases | open releases of the repo in the browser                |
| gitr tags     | open tags of the repo in the browser                    |
| gitr commits  | open commits of the local branch of repo in the browser |
| gitr branches | open branches of the repo in the browser                |

#### gitr clone

This may not seem like a very useful feature at first. If you would like every git repo that you clone to land in a
deterministic location then use `gitr clone <clone-url>` instead of `git clone <clone-url>`.

By default `gitr clone` does the same exact thing as `git clone`.

To mimic the folder structure that your repo is on the scm then use `-c` flag.

Running `gitr clone git@gitlab.mycompany.net:parent/subgroup1/subgroup2/repo.git -c` command from `~/scm` folder, the
repo is cloned to `~/scm/subgroup1/subgroup2/repo` location.

Adding the below configuration to `~/.gitr.yaml`, every repo cloned using `gitr` will be cloned to the location that
mimics scm provider.

```yaml
scm:
  homeDir: /Users/swarup/scm
  hosts:
    - scheme: https
      hostname: gitlab.mycompany.net
      provider: gitlab
      clone:
        homeDir: ""
        includeHostForCreDir: true
        alwaysCreDir: true
## more config
```

With above config `gitr` will clone all repos to `scm.homeDir` location, regardless of where you
run `gitr clone <clone-url>` command.

Because `scm.hosts.[0].includeHostForCreDir` is set to `true`, `gitr` will clone the repo to a folder with the name of
the hostname of the scm under `scm.homeDir`.

Because `scm.hosts.[0].alwaysCreDir` is set to `true`, `gitr` will clone the repo to the same path as that is in
the `<clone-url>`.

example:

```shell
gitr clone git@gitlab.mycompany.net:parent/subgroup1/subgroup2/repo.git
```

The repo gets cloned to `{scmHome}/{scmHostname}/{repoPath}/{repoName}`
i.e `/Users/swarup/scm/gitlab.mycompany.net/parent/subgroup1/subgroup2/repo` location.

Below is a snapshot of a nicely organized directory structure for all repos from different scm systems.

```shell
> tree ~/scm

gitlab.com
└── swarup
    └── group1
        └── app-1
            ├── Dockerfile
            ├── Makefile
            ├── README.md
            ├── main.go
    └── group2
        └── sub-group
            └──app-2
               ├── Dockerfile
               ├── Makefile
               ├── README.md
               ├── main.go
github.com
└── swarupdonepudi
    ├── gitr
    │ ├── go.mod
    │ ├── lib
    │ ├── lib_test
    │ └── main.go
    └── homebrew-gitr
        ├── Formula
        └── README.md
    ahmetb
    └── kubectx
        ├── kubens
        ├── kubectx
        └── README.md
```

### config file

The first time you run any `gitr` comand, `gitr` will automatically create a config file and stores it at `${HOME}/.gitr.yaml` location with the below config if the config file does not already exist.

```yaml
scm:
  copyCloneLocationCdCmdToClipboard: false
  homeDir: ""
  hosts:
    - hostname: github.com
      provider: github
      defaultBranch: master
      clone:
        homeDir: ""
        alwaysCreDir: false
        includeHostForCreDir: false
      scheme: https
    - hostname: gitlab.com
      provider: gitlab
      defaultBranch: main
      clone:
        homeDir: ""
        alwaysCreDir: false
        includeHostForCreDir: false
      scheme: https
    - hostname: bitbucket.org
      provider: bitbucket-cloud
      defaultBranch: master
      clone:
        homeDir: ""
        alwaysCreDir: false
        includeHostForCreDir: false
      scheme: https
```

You can customize the config per your requirements. Below is the config options supported in `~/.gitr.yaml`

| config                                  | default | description                                                                                                                               |
| --------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| copyCloneLocationCdCmdToClipboard       | false   | if this value is set, then gitr will add "cd <cloned-location>" text to your os clipboard. this is useful for quick navigation            |
| scm.homeDir                             | ""      | if this value is set, then gitr clone will always clone the repos to this path, regardless of where you run `gitr clone` command from     |
| scm.hosts.[].scheme                     | ""      | http scheme of scm system allowed: http or https                                                                                          |
| scm.hosts.[].hostname                   | ""      | hostname of scm system                                                                                                                    |
| scm.hosts.[].provider                   | ""      | provider of scm system. allowed values are github, gitlab and bitbucket                                                                   |
| scm.hosts.[].defaultBranch              | ""      | this is the value of the default branch configured on the scm                                                                             |
| scm.hosts.[].clone.homeDir              | ""      | if this is to non-empty string, then gitr clone will consider the value as the home directory while cloning the repos from this host      |
| scm.hosts.[].clone.alwaysCreDir         | false   | if this is set to true, then gitr clone will always create the directories present in the clone urlwhile cloning the repos from this host |
| scm.hosts.[].clone.includeHostForCreDir | false   | if this is set to true, then gitr clone will always prefix the hostname to the clone path                                                 |

#### example config file

```yaml
scm:
  copyCloneLocationCdCmdToClipboard: false
  homeDir: /Users/swarupd/scm
  hosts:
    - hostname: github.com
      provider: github
      defaultBranch: master
      clone:
        homeDir: ""
        alwaysCreDir: true
        includeHostForCreDir: true
      scheme: https
    - hostname: gitlab.com
      provider: gitlab
      defaultBranch: main
      clone:
        homeDir: ""
        alwaysCreDir: true
        includeHostForCreDir: true
      scheme: https
    - hostname: bitbucket.org
      provider: bitbucket-cloud
      defaultBranch: master
      clone:
        homeDir: ""
        alwaysCreDir: true
        includeHostForCreDir: true
      scheme: https
```

### on-prem scm deployments

`gitr` can work with on-prem deployments of the supported scm providers i.e github, gitlab and bitbucket.

add the below shown config to `~/.gitr.yaml` file

example:

```yaml
scm:
  hosts:
    - hostname: gitlab.mycompany.net
      scheme: https
      provider: gitlab
      clone:
        homeDir: ""
        includeHostForCreDir: true
        alwaysCreDir: true
```

multiple on-prem deployments can be added to `~/.gitr.yaml` file

```yaml
scm:
  hosts:
    - hostname: gitlab.mycompany.net
      scheme: https
      provider: gitlab
      clone:
        homeDir: ""
        includeHostForCreDir: true
        alwaysCreDir: true
    - hostname: bitbucket.mycompany.net
      scheme: https
      provider: bitbucket-datacenter
      clone:
        homeDir: ""
        includeHostForCreDir: true
        alwaysCreDir: true
```

### dry run

if `gitr` does not work as expected, it is possible to see what urls `gitr` uses by using `--dry` option. This option is
available for both web and clone features.

> note: when `--dry` or `-d` flags are used, gitr will not open up the repo or clone the repo. it simply displays the info to the console that it would use in non-dry mode

- `--dry` flag passed to `gitr web` command

```shell
> gitr web -d

+---------------+----------------------------------------------------+
| remote        | git@github.com:swarupdonepudi/gitr.git             |
+---------------+----------------------------------------------------+
| provider      | github                                             |
+---------------+----------------------------------------------------+
| host          | github.com                                         |
+---------------+----------------------------------------------------+
| repo-path     | swarupdonepudi/gitr                                |
+---------------+----------------------------------------------------+
| repo-name     | gitr                                               |
+---------------+----------------------------------------------------+
| branch        | master                                             |
+---------------+----------------------------------------------------+
| url-web       | https://github.com/swarupdonepudi/gitr             |
+---------------+----------------------------------------------------+
| url-remote    | https://github.com/swarupdonepudi/gitr/tree/master |
+---------------+----------------------------------------------------+
| url-commits   | swarupdonepudi/gitr/commits/master                 |
+---------------+----------------------------------------------------+
| url-branches  | https://github.com/swarupdonepudi/gitr/branches    |
+---------------+----------------------------------------------------+
| url-tags      | https://github.com/swarupdonepudi/gitr/tags        |
+---------------+----------------------------------------------------+
| url-releases  | https://github.com/swarupdonepudi/gitr/releases    |
+---------------+----------------------------------------------------+
| url-pipelines | https://github.com/swarupdonepudi/gitr/actions     |
+---------------+----------------------------------------------------+

```

- `--dry` flag passed to `gitr clone` command

```shell
> gitr clone git@github.com:swarupdonepudi/gitr.git --dry

+------------+---------------------------------------------------+
| remote     | git@github.com:swarupdonepudi/gitr.git            |
+------------+---------------------------------------------------+
| provider   | github                                            |
+------------+---------------------------------------------------+
| host       | github.com                                        |
+------------+---------------------------------------------------+
| repo-name  | gitr                                              |
+------------+---------------------------------------------------+
| ssh-url    | git@github.com:swarupdonepudi/gitr.git            |
+------------+---------------------------------------------------+
| http-url   | https://github.com/swarupdonepudi/gitr.git        |
+------------+---------------------------------------------------+
| create-dir | true                                              |
+------------+---------------------------------------------------+
| scm-home   | /Users/swarupd/scm                                |
+------------+---------------------------------------------------+
| clone-path | /Users/swarupd/scm/github.com/swarupdonepudi/gitr |
+------------+---------------------------------------------------+

```

### aliases

if you are interested in typing fewer keystrokes, add the below aliases to your `.zshrc` or `.bashrc` file.

this will obviate the need to type `gitr` everytime you want to use it, just type the sub command ex: pipe.

```shell
alias clone="gitr clone "
alias web="gitr web "
alias pipe="gitr pipe "
alias rem="gitr rem "
alias commits="gitr commits "
alias prs="gitr prs "
alias issues="gitr issues "
alias releases="gitr releases "
alias tags="gitr tags "
alias branches="gitr branches "
```

### cleanup

```
brew uninstall gitr
brew untap swarupdonepudi/homebrew-tap
```

### contribute

`gitr` was built to share my passion for extreme productivity with other productivity geeks. Life is too short, and I
ain't wasting time clicking and typing around that does not return any value. For those of you who share the same
passion, I hope you find this project both useful and interesting.

I am also pretty sure that productivity geeks are never content and will always look for more. So, if you see
opportunities to improve this, I will be a happy man to see new issues and pull-requests.
