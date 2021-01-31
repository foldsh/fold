package ctl

import "fmt"

var Version version

func init() {
	Version = version{
		major: 0,
		minor: 0,
		patch: 1,
	}
}

type version struct {
	major int
	minor int
	patch int
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}
