package ghch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jessevdk/go-flags"
)

type ghOpts struct {
	RepoPath string `short:"r" long:"repo" default:"." description:"git repository path"`
	GitPath  string `short:"g" long:"git" default:"git" description:"git path"`
	From     string `short:"f" long:"from" description:"git commit revision range start from"`
	To       string `short:"t" long:"to" description:"git commit revision range end to"`
	Token    string `          long:"token" description:"github token"`
	//Format   string `short:"F" long:"format" default:"json" description:"json or markdown"`
	Verbose bool   `short:"v" long:"verbose"`
	Remote  string `          long:"remote" default:"origin"`
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
	r := gh.mergedPRs(opts.From, opts.To)
	jsn, _ := json.MarshalIndent(r, "", "  ")
	fmt.Fprintln(cli.OutStream, string(jsn))
	return exitCodeOK
}

func parseArgs(args []string) (*ghOpts, error) {
	opts := &ghOpts{}
	_, err := flags.ParseArgs(opts, args)
	return opts, err
}
