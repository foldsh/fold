package ctl

import "testing"

func TestVersionString(t *testing.T) {
	v := SemVer{
		major: 1,
		minor: 2,
		patch: 3,
	}

	if v.String() != "1.2.3" {
		t.Errorf("Expection version to be '1.2.3' but found %s", v.String())
	}
}
