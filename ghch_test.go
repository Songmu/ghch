package ghch

import (
	"reflect"
	"testing"
)

func TestParsePRNums(t *testing.T) {

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
	expect := []int{
		225,
		224,
		223,
		222,
		221,
		217,
		216,
		215,
		214,
		213,
		211,
		210,
		208,
		207,
		209,
		205,
		202,
		161,
		206,
		174,
		203,
		199,
		201,
		200,
		197,
		198,
		196,
		195,
		194,
		193,
		192,
		191,
	}
	if !reflect.DeepEqual(parseMergedPRNums(input), expect) {
		t.Errorf("somthing went wrong")
	}
}
