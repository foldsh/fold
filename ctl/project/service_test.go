package project_test

import (
	"testing"

	"github.com/foldsh/fold/ctl/project"
	"github.com/stretchr/testify/assert"
)

func TestServiceNameValidation(t *testing.T) {
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
			var wasValid bool
			if result == nil {
				wasValid = true
			}
			assert.Equalf(
				t,
				tc.expectation,
				wasValid,
				"For case %s expected valid to be %t but found %v.",
				tc.service.Name,
				tc.expectation,
				result,
			)
		})
	}
}
