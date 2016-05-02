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
	"sync"

	"github.com/Masterminds/semver"
	"github.com/octokit/go-octokit/octokit"
	"github.com/tcnksm/go-gitconfig"
)

type ghch struct {
	repoPath string
	remote   string
	verbose  bool
	token    string
	client   *octokit.Client
}

func (gh *ghch) initialize() *ghch {
	var auth octokit.AuthMethod
	gh.setToken()
	if gh.token != "" {
		auth = octokit.TokenAuth{AccessToken: gh.token}
	}
	gh.client = octokit.NewClient(auth)
	return gh
}

func (gh *ghch) setToken() {
	if gh.token != "" {
		return
	}
	if gh.token = os.Getenv("GITHUB_TOKEN"); gh.token != "" {
		return
	}
	gh.token, _ = gitconfig.GithubToken()
	return
}

func (gh *ghch) cmd(argv ...string) (string, error) {
	arg := []string{"-C", gh.repoPath}
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

func (gh *ghch) versions() []string {
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

func (gh *ghch) getRemote() string {
	if gh.remote != "" {
		return gh.remote
	}
	return "origin"
}

var repoURLReg = regexp.MustCompile(`([^/:]+)/([^/]+?)(?:\.git)?$`)

func (gh *ghch) ownerAndRepo() (owner, repo string) {
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

func (gh *ghch) mergedPRs(from, to string) (prs []*octokit.PullRequest) {
	owner, repo := gh.ownerAndRepo()
	nums := gh.mergedPRNums(from, to)

	var wg sync.WaitGroup
	prCh := make(chan *octokit.PullRequest)

	go func() {
		for pr := range prCh {
			prs = append(prs, pr)
		}
	}()

	for _, num := range nums {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			url, _ := octokit.PullRequestsURL.Expand(octokit.M{"owner": owner, "repo": repo, "number": num})
			pr, r := gh.client.PullRequests(url).One()
			if r.HasError() {
				log.Print(r.Err)
				return
			}
			if !gh.verbose {
				pr = reducePR(pr)
			}
			prCh <- pr
		}(num)
	}
	wg.Wait()
	close(prCh)

	return
}

func (gh *ghch) getLatestSemverTag() string {
	vers := gh.versions()
	if len(vers) < 1 {
		return ""
	}
	return vers[0]
}

var prMergeReg = regexp.MustCompile(`^[a-f0-9]{7} Merge pull request #([0-9]+) from`)

func (gh *ghch) mergedPRNums(from, to string) (nums []int) {
	if from == "" {
		if to == "" {
			return
		}
		from, _ = gh.cmd("rev-list", "--max-parents=0", "HEAD")
		from = strings.TrimSpace(from)
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
