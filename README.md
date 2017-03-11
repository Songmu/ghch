ghch
=======

[![Build Status](https://travis-ci.org/Songmu/ghch.png?branch=master)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/ghch?status.svg)](godoc)

[travis]: https://travis-ci.org/Songmu/ghch
[coveralls]: https://coveralls.io/r/Songmu/ghch?branch=master
[license]: https://github.com/Songmu/ghch/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/ghch

## Description

Generate changelog from git history, tags and merged pull requests

## Installation

    % go get github.com/Songmu/ghch/cmd/ghch

## Synopsis

    % ghch -r /path/to/repo [--format markdown]

## Options

```
-r, --repo=         git repository path (default: .)
-f, --from=         git commit revision range start from
-t, --to=           git commit revision range end to
-v, --verbose
-F, --format=       json or markdown (default: json)
-A, --all           output all changes
-N, --next-version=
-g, --git=          git path (default: git)
    --token=        github token
    --remote=       default remote name (default: origin)
```

## GITHUB Token

When github's api token is required in private repository etc., it is used in the following order of priority.

- command line option `--token`
- enviroment variable `GITHUB_TOKEN`
- `git config github.token`

## GitHub Enterprise

You can use `ghch` for GitHub Enterprise. Change API endpoint via the enviromental variable.

    $ export GITHUB_API=http://github.company.com/api/v3

## Requirements

git 1.8.5 or newer is required.

## Examples

### display changes from last versioned tag

    % ghch
    {
      "pull_requests": [
        {
          "html_url": "https://github.com/mackerelio/mackerel-agent/pull/221",
          "title": "Fix typo",
          "number": 221,
          "state": "closed",
          "user": {
            "login": "yukiyan",
            "avatar_url": "https://avatars.githubusercontent.com/u/7304122?v=3",
            "type": "User"
          },
          "body": "Just fixing a typo 😄 ",
          "created_at": "2016-04-19T08:27:30Z",
          "updated_at": "2016-04-25T01:51:15Z",
          "merged_at": "2016-04-25T01:51:11Z",
          ...
          "merged_by": {
            "login": "stefafafan",
            "avatar_url": "https://avatars.githubusercontent.com/u/3520520?v=3",
            "type": "User"
          }
        },
        ...
      ],
      "from_revision": "v0.30.2",
      "to_revision": "",
      "changed_at": "2016-04-27T19:05:49+09:00",
      "owner": "mackerelio",
      "repo": "mackerel-agent"
    }

### display changes from last versioned tag in markdown

    % ghch --format=markdown --next-version=v0.30.3
    ## [v0.30.3](https://github.com/mackerelio/mackerel-agent/releases/tag/v0.30.3) (2016-04-27)

    * retry retirement when api request failed [#224](https://github.com/mackerelio/mackerel-agent/pull/224) ([Songmu](https://github.com/Songmu))
    * Fix comments [#222](https://github.com/mackerelio/mackerel-agent/pull/222) ([stefafafan](https://github.com/stefafafan))
    * Remove go get cmd/vet [#223](https://github.com/mackerelio/mackerel-agent/pull/223) ([itchyny](https://github.com/itchyny))
    * [nit] [plugin.checks.foo ] is valid toml now [#225](https://github.com/mackerelio/mackerel-agent/pull/225) ([Songmu](https://github.com/Songmu))
    * Remove usr local bin again [#217](https://github.com/mackerelio/mackerel-agent/pull/217) ([Songmu](https://github.com/Songmu))
    * Fix typo [#221](https://github.com/mackerelio/mackerel-agent/pull/221) ([yukiyan](https://github.com/yukiyan))

### display all changes

    % ghch --format=markdown --next-version=v0.30.3 --all
    ...

### display changes between specified two revisions

    % ghch --from v0.9.0 --to v0.9.1
    ...

## Author

[Songmu](https://github.com/Songmu)
