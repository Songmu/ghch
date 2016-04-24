package ghch

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

type Ghch struct {
	RepoPath string
}

func (gh *Ghch) Cmd(argv ...string) (string, error) {
	arg := []string{"-C", gh.RepoPath}
	arg = append(arg, argv...)
	cmd := exec.Command("git", arg...)
	cmd.Env = append(os.Environ(), "LANG=C")

	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	return b.String(), err
}

var verReg = regexp.MustCompile(`^v?[0-9]+(?:\.[0-9]+){0,2}$`)

func (gh *Ghch) Versions() []string {
	out, _ := gh.Cmd("tag")
	rawTags := strings.Split(out, "\n")

	var versions []*semver.Version
	for _, tag := range rawTags {
		if verReg.MatchString(tag) {
			v, _ := semver.NewVersion(tag)
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
