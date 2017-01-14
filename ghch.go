package ghch

import (
	"bytes"
	"fmt"
	"log"
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

const version = "0.0.4"

type ghch struct {
	repoPath string
	gitPath  string
	remote   string
	verbose  bool
	token    string
	baseURL  string
	client   *octokit.Client
}

func (gh *ghch) initialize() *ghch {
	var auth octokit.AuthMethod
	gh.setToken()
	if gh.token != "" {
		auth = octokit.TokenAuth{AccessToken: gh.token}
	}

	gh.setBaseURL()

	if gh.baseURL != "" {
		gh.client = octokit.NewClientWith(gh.baseURL, "Octokit Go", auth, nil)
		return gh
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

func (gh *ghch) setBaseURL() {
	if gh.baseURL != "" {
		return
	}
	if gh.baseURL = os.Getenv("GITHUB_API"); gh.baseURL != "" {
		return
	}

	return
}

func (gh *ghch) gitProg() string {
	if gh.gitPath != "" {
		return gh.gitPath
	}
	return "git"
}

func (gh *ghch) cmd(argv ...string) (string, error) {
	arg := []string{"-C", gh.repoPath}
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

func (gh *ghch) versions() []string {
	sv := gitsemvers.Semvers{
		RepoPath: gh.repoPath,
		GitPath:  gh.gitPath,
	}
	return sv.VersionStrings()
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
	prs = make([]*octokit.PullRequest, 0, len(nums))
	prsWithNil := make([]*octokit.PullRequest, len(nums))

	var wg sync.WaitGroup

	for i, num := range nums {
		wg.Add(1)
		go func(i, num int) {
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
			prsWithNil[i] = pr
		}(i, num)
	}
	wg.Wait()
	for _, pr := range prsWithNil {
		if pr != nil {
			prs = append(prs, pr)
		}
	}
	return
}

func (gh *ghch) getLatestSemverTag() string {
	vers := gh.versions()
	if len(vers) < 1 {
		return ""
	}
	return vers[0]
}

var prMergeReg = regexp.MustCompile(`#([0-9]+)`)

func (gh *ghch) mergedPRNums(from, to string) (nums []int) {
	if from == "" {
		from, _ = gh.cmd("rev-list", "--max-parents=0", "HEAD")
		from = strings.TrimSpace(from)
	}

	revisionRange := fmt.Sprintf("%s..%s", from, to)
	out, err := gh.cmd("log", revisionRange, "--pretty=%s")
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

func (gh *ghch) getChangedAt(rev string) (time.Time, error) {
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
