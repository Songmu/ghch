package ghch

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/octokit/go-octokit/octokit"
)

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
	client *octokit.Client
}

const (
	exitCodeOK = iota
	exitCodeParseFlagError
	exitCodeErr
)

// CLI is struct for command line tool
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Run the ghch
func (cli *CLI) Run(argv []string) int {
	log.SetOutput(cli.ErrStream)
	p, gh, err := parseArgs(argv)
	if err != nil {
		if ferr, ok := err.(*flags.Error); !ok || ferr.Type != flags.ErrHelp {
			p.WriteHelp(cli.ErrStream)
		}
		return exitCodeParseFlagError
	}
	gh = gh.Initialize()

	if gh.All {
		chlog := Changelog{}
		vers := append(gh.versions(), "")
		prevRev := ""
		for _, rev := range vers {
			r, err := gh.getSection(rev, prevRev)
			if err != nil {
				log.Print(err)
				return exitCodeErr
			}
			if prevRev == "" && gh.NextVersion != "" {
				r.ToRevision = gh.NextVersion
			}
			chlog.Sections = append(chlog.Sections, r)
			prevRev = rev
		}

		if gh.Format == "markdown" {
			results := make([]string, len(chlog.Sections))
			for i, v := range chlog.Sections {
				results[i], _ = v.toMkdn()
			}

			if gh.Write {
				content := "# Changelog\n\n" + strings.Join(results, "\n\n")
				err := ioutil.WriteFile(gh.ChangelogMd, []byte(content), 0644)
				if err != nil {
					log.Print(err)
					return exitCodeErr
				}
			} else {
				fmt.Fprintln(cli.OutStream, strings.Join(results, "\n\n"))
			}
		} else {
			jsn, _ := json.MarshalIndent(chlog, "", "  ")
			fmt.Fprintln(cli.OutStream, string(jsn))
		}
	} else {
		if gh.Latest {
			vers := gh.versions()
			if len(vers) > 0 {
				gh.To = vers[0]
			}
			if gh.From == "" && len(vers) > 1 {
				gh.From = vers[1]
			}
		} else if gh.From == "" && gh.To == "" {
			gh.From = gh.getLatestSemverTag()
		}
		r, err := gh.getSection(gh.From, gh.To)
		if err != nil {
			log.Print(err)
			return exitCodeErr
		}
		if r.ToRevision == "" && gh.NextVersion != "" {
			r.ToRevision = gh.NextVersion
		}
		if gh.Format == "markdown" {
			str, err := r.toMkdn()
			if err != nil {
				log.Print(err)
				return exitCodeErr
			}
			if gh.Write {
				content := ""
				if exists(gh.ChangelogMd) {
					byt, err := ioutil.ReadFile(gh.ChangelogMd)
					if err != nil {
						log.Print(err)
						return exitCodeErr
					}
					content = insertNewChangelog(byt, str)
				} else {
					content = "# Changelog\n\n" + str + "\n"
				}
				err = ioutil.WriteFile(gh.ChangelogMd, []byte(content), 0644)
				if err != nil {
					log.Print(err)
					return exitCodeErr
				}
			} else {
				fmt.Fprintln(cli.OutStream, str)
			}
		} else {
			jsn, _ := json.MarshalIndent(r, "", "  ")
			fmt.Fprintln(cli.OutStream, string(jsn))
		}
	}
	return exitCodeOK
}

func insertNewChangelog(orig []byte, section string) string {
	var bf bytes.Buffer
	lineSnr := bufio.NewScanner(bytes.NewReader(orig))
	inserted := false
	for lineSnr.Scan() {
		line := lineSnr.Text()
		if !inserted && strings.HasPrefix(line, "## ") {
			bf.WriteString(section)
			bf.WriteString("\n\n")
			inserted = true
		}
		bf.WriteString(line)
		bf.WriteString("\n")
	}
	if !inserted {
		bf.WriteString(section)
	}
	return bf.String()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func parseArgs(args []string) (*flags.Parser, *Ghch, error) {
	gh := &Ghch{}
	p := flags.NewParser(gh, flags.Default)
	p.Usage = fmt.Sprintf("[OPTIONS]\n\nVersion: %s (rev: %s)", version, revision)
	rest, err := p.ParseArgs(args)
	if gh.Write {
		gh.Format = "markdown"
		gh.ChangelogMd = "CHANGELOG.md"
		if len(rest) > 0 {
			gh.ChangelogMd = rest[0]
		}
	}
	return p, gh, err
}

func (gh *Ghch) getSection(from, to string) (Section, error) {
	if from == "" {
		from, _ = gh.cmd("rev-list", "--max-parents=0", "HEAD")
		from = strings.TrimSpace(from)
		if len(from) > 12 {
			from = from[:12]
		}
	}
	r, err := gh.mergedPRs(from, to)
	if err != nil {
		return Section{}, err
	}
	t, err := gh.getChangedAt(to)
	if err != nil {
		return Section{}, err
	}
	owner, repo := gh.ownerAndRepo()
	htmlURL, err := gh.htmlURL(owner, repo)
	if err != nil {
		return Section{}, err
	}
	return Section{
		PullRequests: r,
		FromRevision: from,
		ToRevision:   to,
		ChangedAt:    t,
		Owner:        owner,
		Repo:         repo,
		HTMLURL:      htmlURL,
	}, nil
}

// Changelog contains Sectionst
type Changelog struct {
	Sections []Section `json:"Sections"`
}

// Section contains changes between two revisions
type Section struct {
	PullRequests []*octokit.PullRequest `json:"pull_requests"`
	FromRevision string                 `json:"from_revision"`
	ToRevision   string                 `json:"to_revision"`
	ChangedAt    time.Time              `json:"changed_at"`
	Owner        string                 `json:"owner"`
	Repo         string                 `json:"repo"`
	HTMLURL      string                 `json:"html_url"`
}

var tmplStr = `{{$ret := . -}}
## [{{.ToRevision}}]({{.HTMLURL}}/compare/{{.FromRevision}}...{{.ToRevision}}) ({{.ChangedAt.Format "2006-01-02"}})
{{range .PullRequests}}
* {{.Title}} [#{{.Number}}]({{.HTMLURL}}) ([{{.User.Login}}]({{.User.HTMLURL}}))
{{- end}}`

var mdTmpl *template.Template

func init() {
	var err error
	mdTmpl, err = template.New("md-changelog").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
}

func (rs Section) toMkdn() (string, error) {
	var b bytes.Buffer
	err := mdTmpl.Execute(&b, rs)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
