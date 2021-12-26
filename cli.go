package ghch

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/google/go-github/v41/github"
	"github.com/jessevdk/go-flags"
)

// Ghch is main application struct
type Ghch struct {
	RepoPath    string `short:"r" long:"repo" default:"." description:"git repository path"`
	BaseURL     string
	GitPath     string `short:"g" long:"git" default:"git" description:"git path"`
	From        string `short:"f" long:"from" description:"git commit revision range start from"`
	To          string `short:"t" long:"to" description:"git commit revision range end to"`
	Latest      bool   `          long:"latest" description:"output changes between latest two semantic versioned tags"`
	Token       string `          long:"token" description:"github token"`
	Verbose     bool   `short:"v" long:"verbose"`
	Remote      string `          long:"remote" default:"origin" description:"default remote name"`
	Format      string `short:"F" long:"format" description:"json or markdown"`
	All         bool   `short:"A" long:"all" description:"output all changes"`
	NextVersion string `short:"N" long:"next-version"`
	Write       bool   `short:"w" description:"write result to file"`
	ChangelogMd string
	// Tmpl string
	OutStream io.Writer

	client *github.Client
}

// Run the ghch
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	return (&cli{OutStream: outStream, ErrStream: errStream}).Run(ctx, argv)
}

type cli struct {
	OutStream, ErrStream io.Writer
}

func (cl *cli) Run(ctx context.Context, argv []string) error {
	log.SetOutput(cl.ErrStream)
	gh, err := cl.parseArgs(argv)
	if err != nil {
		if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrHelp {
				fmt.Fprint(cl.OutStream, err)
				return nil
			}
			return ferr
		}
		return err
	}
	if err := gh.Run(); err != nil {
		return err
	}
	return nil
}

func (cl *cli) parseArgs(args []string) (*Ghch, error) {
	gh := &Ghch{
		OutStream: cl.OutStream,
	}
	p := flags.NewParser(gh, flags.HelpFlag|flags.PassDoubleDash)
	p.Usage = fmt.Sprintf("[OPTIONS]\n\nVersion: %s (rev: %s)", version, revision)
	rest, err := p.ParseArgs(args)
	if gh.Write {
		gh.Format = "markdown"
		gh.ChangelogMd = "CHANGELOG.md"
		if len(rest) > 0 {
			gh.ChangelogMd = rest[0]
		}
	}
	return gh, err
}
