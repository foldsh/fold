package version

import "testing"

func TestVersionString(t *testing.T) {
	v := SemVer{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}

	if v.String() != "1.2.3" {
		t.Errorf("Expection version to be '1.2.3' but found %s", v.String())
	}
}
