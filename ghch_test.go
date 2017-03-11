package ghch

import (
	"reflect"
	"testing"
)

func TestParsePRLogs(t *testing.T) {

	input := `6191693 Merge pull request #225 from mackerelio/fix-test-for-invaild-toml
dbb1d50 Merge pull request #224 from mackerelio/retry-retire
922eb55 Merge pull request #223 from mackerelio/remove_vet
aa18e36 Merge pull request #222 from mackerelio/fix-comments
71da053 Merge pull request #221 from yukiyan/fix-typo
0925081 Merge pull request #217 from mackerelio/remove-usr-local-bin-again
b4bc51b Merge pull request #216 from mackerelio/bump-version-0.30.2
a8ea16b Merge pull request #215 from mackerelio/revert-9e0c8ab1
98b28f1 Merge pull request #214 from mackerelio/bump-version-0.30.1
19a0010 Merge pull request #213 from mackerelio/workaround-amd64
9e0c8ab Merge pull request #211 from mackerelio/usr-bin
7d278aa Merge pull request #210 from mackerelio/bump-version-0.30.0
ce37096 Merge pull request #208 from mackerelio/refactor-net-interface
8a07070 Merge pull request #207 from mackerelio/subcommand-init
a39ca5e Merge pull request #209 from mackerelio/remove-cpu-flags
c0ed1f1 Merge pull request #205 from mackerelio/interface-ips
8cc281c Merge pull request #202 from mackerelio/remove-deprecated-sensu
afeb5e5 Merge pull request #161 from mackerelio/remove-uptime
a75f8b2 Merge pull request #206 from mackerelio/bump-version-0.29.2
fd40654 Merge pull request #174 from mackerelio/travis-docker
a8665f5 Merge branch 'master' into travis-docker
bdb2271 Merge pull request #203 from mackerelio/alternative-build
2ac5301 Merge pull request #199 from mackerelio/fix-deb
7c79f92 Merge pull request #201 from mackerelio/bump-version-0.29.1
32a3e1f Merge pull request #200 from mackerelio/bump-version-0.29.0
b4b8c2c Merge pull request #197 from hanazuki/check-timeouts
a30e851 Merge pull request #198 from mackerelio/dont-ignore-logging-level_string
2ec717e Merge pull request #196 from mackerelio/refactor-around-start
ca345ea Merge pull request #195 from mackerelio/introduce-motemen-go-cli
843b32e Merge pull request #194 from mackerelio/remove-deprecated
87375ec Merge pull request #193 from mackerelio/bump-version-0.28.1
82ccaa3 Merge branch 'master' of github.com:mackerelio/mackerel-agent
4a6d83c Merge pull request #192 from mackerelio/deb_init_d_stop_retval
5b0a536 Merge pull request #191 from mackerelio/gofmt-on-travis
`
	expect := []*mergedPRLog{
		{num: 225, branch: "mackerelio/fix-test-for-invaild-toml"},
		{num: 224, branch: "mackerelio/retry-retire"},
		{num: 223, branch: "mackerelio/remove_vet"},
		{num: 222, branch: "mackerelio/fix-comments"},
		{num: 221, branch: "yukiyan/fix-typo"},
		{num: 217, branch: "mackerelio/remove-usr-local-bin-again"},
		{num: 216, branch: "mackerelio/bump-version-0.30.2"},
		{num: 215, branch: "mackerelio/revert-9e0c8ab1"},
		{num: 214, branch: "mackerelio/bump-version-0.30.1"},
		{num: 213, branch: "mackerelio/workaround-amd64"},
		{num: 211, branch: "mackerelio/usr-bin"},
		{num: 210, branch: "mackerelio/bump-version-0.30.0"},
		{num: 208, branch: "mackerelio/refactor-net-interface"},
		{num: 207, branch: "mackerelio/subcommand-init"},
		{num: 209, branch: "mackerelio/remove-cpu-flags"},
		{num: 205, branch: "mackerelio/interface-ips"},
		{num: 202, branch: "mackerelio/remove-deprecated-sensu"},
		{num: 161, branch: "mackerelio/remove-uptime"},
		{num: 206, branch: "mackerelio/bump-version-0.29.2"},
		{num: 174, branch: "mackerelio/travis-docker"},
		{num: 203, branch: "mackerelio/alternative-build"},
		{num: 199, branch: "mackerelio/fix-deb"},
		{num: 201, branch: "mackerelio/bump-version-0.29.1"},
		{num: 200, branch: "mackerelio/bump-version-0.29.0"},
		{num: 197, branch: "hanazuki/check-timeouts"},
		{num: 198, branch: "mackerelio/dont-ignore-logging-level_string"},
		{num: 196, branch: "mackerelio/refactor-around-start"},
		{num: 195, branch: "mackerelio/introduce-motemen-go-cli"},
		{num: 194, branch: "mackerelio/remove-deprecated"},
		{num: 193, branch: "mackerelio/bump-version-0.28.1"},
		{num: 192, branch: "mackerelio/deb_init_d_stop_retval"},
		{num: 191, branch: "mackerelio/gofmt-on-travis"},
	}
	if !reflect.DeepEqual(parseMergedPRLogs(input), expect) {
		t.Errorf("somthing went wrong")
	}
}
