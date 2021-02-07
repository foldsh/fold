package ctl

import "fmt"

var FoldctlVersion SemVer

func init() {
	FoldctlVersion = SemVer{
		major: 0,
		minor: 0,
		patch: 1,
	}
}

type SemVer struct {
	major int
	minor int
	patch int
}

func (v SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}
