package ghch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Songmu/gitsemvers"
	"github.com/octokit/go-octokit/octokit"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (gh *Ghch) Run() error {
	gh.initialize()
	if gh.All {
		return gh.runAll()
	}
	return gh.run()
}

func (gh *Ghch) runAll() error {
	chlog := Changelog{}
	vers := append(gh.versions(), "")
	prevRev := ""
	for _, rev := range vers {
		r, err := gh.getSection(rev, prevRev)
		if err != nil {
			return err
		}
		if prevRev == "" && gh.NextVersion != "" {
			r.ToRevision = gh.NextVersion
		}
		chlog.Sections = append(chlog.Sections, r)
		prevRev = rev
	}

	if gh.Format != "markdown" { // json
		encoder := json.NewEncoder(gh.OutStream)
		encoder.SetIndent("", "  ")
		return encoder.Encode(chlog)
	}
	results := make([]string, len(chlog.Sections))
	for i, v := range chlog.Sections {
		results[i], _ = v.toMkdn()
	}
	if gh.Write {
		content := "# Changelog\n\n" + strings.Join(results, "\n\n")
		if err := ioutil.WriteFile(gh.ChangelogMd, []byte(content), 0644); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(gh.OutStream, strings.Join(results, "\n\n"))
	}
	return nil
}

func (gh *Ghch) run() error {
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
		return err
	}
	if r.ToRevision == "" && gh.NextVersion != "" {
		r.ToRevision = gh.NextVersion
	}

	if gh.Format != "markdown" { // json
		encoder := json.NewEncoder(gh.OutStream)
		encoder.SetIndent("", "  ")
		return encoder.Encode(r)
	}
	str, err := r.toMkdn()
	if err != nil {
		return err
	}
	if gh.Write {
		content := ""
		if exists(gh.ChangelogMd) {
			byt, err := ioutil.ReadFile(gh.ChangelogMd)
			if err != nil {
				return err
			}
			content = insertNewChangelog(byt, str)
		} else {
			content = "# Changelog\n\n" + str + "\n"
		}
		if err := ioutil.WriteFile(gh.ChangelogMd, []byte(content), 0644); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(gh.OutStream, str)
	}
	return nil
}

func (gh *Ghch) initialize() {
	if gh.Write {
		gh.Format = "markdown"
		if gh.ChangelogMd != "" {
			gh.ChangelogMd = "CHANGELOG.md"
		}
	}
	if gh.OutStream == nil {
		gh.OutStream = os.Stdout
	}

	var auth octokit.AuthMethod
	gh.setToken()
	if gh.Token != "" {
		auth = octokit.TokenAuth{AccessToken: gh.Token}
	}

	gh.setBaseURL()
	if gh.BaseURL != "" {
		gh.client = octokit.NewClientWith(gh.BaseURL, "Octokit Go", auth, nil)
	} else {
		gh.client = octokit.NewClient(auth)
	}
}

func (gh *Ghch) setToken() {
	if gh.Token != "" {
		return
	}
	if gh.Token = os.Getenv("GITHUB_TOKEN"); gh.Token != "" {
		return
	}
	gh.Token, _ = gitconfig.GithubToken()
	return
}

func (gh *Ghch) setBaseURL() {
	if gh.BaseURL != "" {
		return
	}
	if gh.BaseURL = os.Getenv("GITHUB_API"); gh.BaseURL != "" {
		return
	}

	return
}

func (gh *Ghch) gitProg() string {
	if gh.GitPath != "" {
		return gh.GitPath
	}
	return "git"
}

func (gh *Ghch) cmd(argv ...string) (string, error) {
	arg := []string{"-C", gh.RepoPath}
	arg = append(arg, argv...)
	cmd := exec.Command(gh.gitProg(), arg...)
	cmd.Env = append(os.Environ(), "LANG=C")

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return b.String(), err
}

var verReg = regexp.MustCompile(`^v?[0-9]+(?:\.[0-9]+){0,2}$`)

func (gh *Ghch) versions() []string {
	sv := gitsemvers.Semvers{
		RepoPath: gh.RepoPath,
		GitPath:  gh.GitPath,
	}
	return sv.VersionStrings()
}

func (gh *Ghch) getRemote() string {
	if gh.Remote != "" {
		return gh.Remote
	}
	return "origin"
}

var repoURLReg = regexp.MustCompile(`([^/:]+)/([^/]+?)(?:\.git)?$`)

func (gh *Ghch) ownerAndRepo() (owner, repo string) {
	out, _ := gh.cmd("remote", "-v")
	remotes := strings.Split(out, "\n")
	for _, r := range remotes {
		fields := strings.Fields(r)
		if len(fields) > 1 && fields[0] == gh.getRemote() {
			if matches := repoURLReg.FindStringSubmatch(fields[1]); len(matches) > 2 {
				return matches[1], matches[2]
			}
		}
	}
	return
}

func (gh *Ghch) htmlURL(owner, repo string) (htmlURL string, err error) {
	re, r := gh.client.Repositories().One(nil, octokit.M{"owner": owner, "repo": repo})
	if r.Err != nil {
		if rerr, ok := r.Err.(*octokit.ResponseError); ok {
			if rerr.Response != nil && rerr.Response.StatusCode == http.StatusNotFound {
				return
			}
		}
		err = r.Err
		log.Print(r.Err)
		return
	}

	htmlURL = re.HTMLURL
	return
}

func (gh *Ghch) mergedPRs(from, to string) (prs []*octokit.PullRequest, err error) {
	owner, repo := gh.ownerAndRepo()
	prlogs, err := gh.mergedPRLogs(from, to)
	if err != nil {
		return
	}
	prs = make([]*octokit.PullRequest, 0, len(prlogs))
	prsWithNil := make([]*octokit.PullRequest, len(prlogs))
	errsWithNil := make([]error, len(prlogs))

	var wg sync.WaitGroup

	for i, prlog := range prlogs {
		wg.Add(1)
		go func(i int, prlog *mergedPRLog) {
			defer wg.Done()
			url, _ := octokit.PullRequestsURL.Expand(octokit.M{"owner": owner, "repo": repo, "number": prlog.num})
			pr, r := gh.client.PullRequests(url).One()
			if r.Err != nil {
				if rerr, ok := r.Err.(*octokit.ResponseError); ok {
					if rerr.Response != nil && rerr.Response.StatusCode == http.StatusNotFound {
						return
					}
				}
				errsWithNil[i] = r.Err
				log.Print(r.Err)
				return
			}
			// replace repoowner:branch-name to repo-owner/branch-name
			if strings.Replace(pr.Head.Label, ":", "/", 1) != prlog.branch {
				return
			}
			if !gh.Verbose {
				pr = reducePR(pr)
			}
			prsWithNil[i] = pr
		}(i, prlog)
	}
	wg.Wait()
	for _, pr := range prsWithNil {
		if pr != nil {
			prs = append(prs, pr)
		}
	}
	for _, e := range errsWithNil {
		if e != nil {
			err = e
		}
	}
	return
}

func (gh *Ghch) getLatestSemverTag() string {
	vers := gh.versions()
	if len(vers) < 1 {
		return ""
	}
	return vers[0]
}

type mergedPRLog struct {
	num    int
	branch string
}

func (gh *Ghch) mergedPRLogs(from, to string) (nums []*mergedPRLog, err error) {
	revisionRange := fmt.Sprintf("%s..%s", from, to)
	out, err := gh.cmd("log", revisionRange, "--merges", "--oneline")
	if err != nil {
		return []*mergedPRLog{}, err
	}
	return parseMergedPRLogs(out), nil
}

var prMergeReg = regexp.MustCompile(`^[a-f0-9]+ Merge pull request #([0-9]+) from (\S+)`)

func parseMergedPRLogs(out string) (prs []*mergedPRLog) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if matches := prMergeReg.FindStringSubmatch(line); len(matches) > 2 {
			i, _ := strconv.Atoi(matches[1])
			prs = append(prs, &mergedPRLog{
				num:    i,
				branch: matches[2],
			})
		}
	}
	return
}

func (gh *Ghch) getChangedAt(rev string) (time.Time, error) {
	if rev == "" {
		rev = "HEAD"
	}
	out, err := gh.cmd("show", "-s", rev+"^{commit}", `--format=%ct`)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get changed at from git revision. `git show` failed")
	}
	out = strings.TrimSpace(out)
	i, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get changed at from git revision. ParseInt failed")
	}
	return time.Unix(i, 0), nil
}
