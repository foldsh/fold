package project

import "testing"

func TestServiceValidate(t *testing.T) {
	cases := []struct {
		service     Service
		expectation bool
	}{
		{Service{Name: "foo"}, true},
		{Service{Name: "bar"}, true},
		{Service{Name: "foo-bar"}, true},
		{Service{Name: "foo_bar"}, true},
		{Service{Name: "foo/bar"}, false},
		{Service{Name: "1foo/bar"}, false},
		{Service{Name: "1"}, false},
		{Service{Name: "1foo"}, false},
		{Service{Name: ""}, false},
		{Service{Name: "-------"}, false},
		{Service{Name: "a------"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.service.Name, func(t *testing.T) {
			result := tc.service.Validate()
			if result != tc.expectation {
				t.Errorf("For case %s expected %t but found %t.", tc.service.Name, tc.expectation, result)
			}
		})
	}
}
