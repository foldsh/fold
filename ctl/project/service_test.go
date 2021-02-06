package project_test

import (
	"testing"

	"github.com/foldsh/fold/ctl/project"
)

func TestServiceValidate(t *testing.T) {
	cases := []struct {
		service     project.Service
		expectation bool
	}{
		{project.Service{Name: "foo"}, true},
		{project.Service{Name: "bar"}, true},
		{project.Service{Name: "foo-bar"}, true},
		{project.Service{Name: "foo_bar"}, true},
		{project.Service{Name: "foo/bar"}, false},
		{project.Service{Name: "1foo/bar"}, false},
		{project.Service{Name: "1"}, false},
		{project.Service{Name: "1foo"}, false},
		{project.Service{Name: ""}, false},
		{project.Service{Name: "-------"}, false},
		{project.Service{Name: "a------"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.service.Name, func(t *testing.T) {
			result := tc.service.Validate()
			if result != tc.expectation {
				t.Errorf(
					"For case %s expected %t but found %t.",
					tc.service.Name,
					tc.expectation,
					result,
				)
			}
		})
	}
}
