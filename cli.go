package ghch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/octokit/go-octokit/octokit"
)

type ghOpts struct {
	RepoPath string `short:"r" long:"repo" default:"." description:"git repository path"`
	GitPath  string `short:"g" long:"git" default:"git" description:"git path"`
	From     string `short:"f" long:"from" description:"git commit revision range start from"`
	To       string `short:"t" long:"to" description:"git commit revision range end to"`
	Token    string `          long:"token" description:"github token"`
	Verbose  bool   `short:"v" long:"verbose"`
	Remote   string `          long:"remote" default:"origin"`
	// All     bool   `short:"A" long:"all" `
	// Format   string `short:"F" long:"format" default:"json" description:"json or markdown"`
	// Tmpl string
}

const (
	exitCodeOK = iota
	exitCodeParseFlagError
	exitCodeErr
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

func (cli *CLI) Run(argv []string) int {
	opts, err := parseArgs(argv)
	if err != nil {
		return exitCodeParseFlagError
	}

	gh := (&ghch{
		remote:   opts.Remote,
		repoPath: opts.RepoPath,
		verbose:  opts.Verbose,
		token:    opts.Token,
	}).initialize()

	if opts.From == "" && opts.To == "" {
		opts.From = gh.getLatestSemverTag()
	}
	r := gh.getResult(opts.From, opts.To)
	jsn, _ := json.MarshalIndent(r, "", "  ")
	fmt.Fprintln(cli.OutStream, string(jsn))
	return exitCodeOK
}

func parseArgs(args []string) (*ghOpts, error) {
	opts := &ghOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}

func (gh *ghch) getResult(from, to string) result {
	r := gh.mergedPRs(from, to)
	t, err := gh.getChangedAt(to)
	if err != nil {
		log.Print(err)
	}
	return result{
		PullRequests: r,
		FromRevision: from,
		ToRevision:   to,
		ChangedAt:    t,
	}
}

type result struct {
	PullRequests []*octokit.PullRequest `json:"pull_requests"`
	FromRevision string                 `json:"from_revision"`
	ToRevision   string                 `json:"to_revision"`
	ChangedAt    time.Time              `json:"changed_at"`
}
