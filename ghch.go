package ghch

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/octokit/go-octokit/octokit"
)

type Ghch struct {
	RepoPath string
	Remote   string
	client   *octokit.Client
}

func New(repo string) *Ghch {
	return &Ghch{
		RepoPath: repo,
		// XXX authentication
		client: octokit.NewClient(nil),
	}
}

func (gh *Ghch) cmd(argv ...string) (string, error) {
	arg := []string{"-C", gh.RepoPath}
	arg = append(arg, argv...)
	cmd := exec.Command("git", arg...)
	cmd.Env = append(os.Environ(), "LANG=C")

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return b.String(), err
}

var verReg = regexp.MustCompile(`^v?[0-9]+(?:\.[0-9]+){0,2}$`)

func (gh *Ghch) versions() []string {
	out, _ := gh.cmd("tag")
	return parseVerions(out)
}

func parseVerions(out string) []string {
	rawTags := strings.Split(out, "\n")
	var versions []*semver.Version
	for _, tag := range rawTags {
		t := strings.TrimSpace(tag)
		if verReg.MatchString(t) {
			v, _ := semver.NewVersion(t)
			versions = append(versions, v)
		}
	}
	sort.Sort(sort.Reverse(semver.Collection(versions)))
	var vers = make([]string, len(versions))
	for i, v := range versions {
		vers[i] = v.Original()
	}
	return vers
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

func (gh *Ghch) MergedPRs(from, to string) (prs []*octokit.PullRequest) {
	owner, repo := gh.ownerAndRepo()
	nums := gh.mergedPRNums(from, to)
	for _, num := range nums {
		url, _ := octokit.PullRequestsURL.Expand(octokit.M{"owner": owner, "repo": repo, "number": num})
		pr, r := gh.client.PullRequests(url).One()
		if r.HasError() {
			log.Print(r.Err)
			continue
		}
		prs = append(prs, reducePR(pr))
	}
	return
}

var prMergeReg = regexp.MustCompile(`^[a-f0-9]{7} Merge pull request #([0-9]+) from`)

func (gh *Ghch) mergedPRNums(from, to string) (nums []int) {
	if from == "" {
		vers := gh.versions()
		if len(vers) < 1 {
			return
		}
		from = vers[0]
	}
	revisionRange := fmt.Sprintf("%s..%s", from, to)
	out, err := gh.cmd("log", revisionRange, "--merges", "--oneline")
	if err != nil {
		return
	}
	return parseMergedPRNums(out)
}

func parseMergedPRNums(out string) (nums []int) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if matches := prMergeReg.FindStringSubmatch(line); len(matches) > 1 {
			i, _ := strconv.Atoi(matches[1])
			nums = append(nums, i)
		}
	}
	return
}
