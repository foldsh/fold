package version

import "fmt"

var FoldVersion = SemVer{Major: 0, Minor: 0, Patch: 3}

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func (v SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}
