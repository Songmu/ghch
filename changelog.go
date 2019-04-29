package ghch

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/github"
)

// Changelog contains Sectionst
type Changelog struct {
	Sections []Section `json:"Sections"`
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

// Section contains changes between two revisions
type Section struct {
	PullRequests []*github.PullRequest `json:"pull_requests"`
	FromRevision string                `json:"from_revision"`
	ToRevision   string                `json:"to_revision"`
	ChangedAt    time.Time             `json:"changed_at"`
	Owner        string                `json:"owner"`
	Repo         string                `json:"repo"`
	HTMLURL      string                `json:"html_url"`
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
