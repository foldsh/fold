package project_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
)

func TestProjectConfigValidation(t *testing.T) {
	logger := logging.NewTestLogger()
	dir, err := ioutil.TempDir("", "testProjectConfig")
	require.Nil(t, err)
	defer os.RemoveAll(dir)
	assert.False(
		t,
		project.IsAFoldProject(dir),
		"Empty directory should not be a valid fold project",
	)

	cfgPath := filepath.Join(dir, "fold.yaml")

	ctx := newCmdCtx(logger, dir)
	_, err = project.Load(ctx, dir)
	assert.ErrorIs(t, err, project.NotAFoldProject)

	err = ioutil.WriteFile(cfgPath, []byte("not valid yaml"), 0644)
	require.Nil(t, err, "Failed to write to config file.")

	assert.True(
		t,
		project.IsAFoldProject(dir),
		"After creating a fold.yaml file it should be a valid fold project",
	)
	_, err = project.Load(ctx, dir)
	assert.ErrorIs(t, err, project.InvalidConfig)

	err = ioutil.WriteFile(cfgPath, []byte("valid: yaml"), 0644)
	require.Nil(t, err, "Failed to write to config file.")

	cfg, err := project.Load(ctx, dir)
	assert.ErrorIsf(t, err, project.InvalidConfig, "Expected InvalidConfig but got %v\n", err)

	s := `name: blah
maintainer: blah
email: blah
repository: blah
services:
- name: blah
  path: ./blah
  mounts:
  - ./foo
  - ./bar
`
	err = ioutil.WriteFile(cfgPath, []byte(s), 0644)
	require.Nil(t, err, "Failed to write to config file.")

	cfg, err = project.Load(ctx, dir)
	assert.Nilf(t, err, "Expected to load config but got error %v", err)

	expectation := &project.Project{
		Name:       "blah",
		Maintainer: "blah",
		Email:      "blah",
		Repository: "blah",
		Services: []*project.Service{
			{
				Name:   "blah",
				Path:   "./blah",
				Mounts: []string{"./foo", "./bar"},
			},
		},
	}
	diffConfig(t, expectation, cfg)
}

func TestProjectLoadingAndSaving(t *testing.T) {
	logger := logging.NewTestLogger()
	dir, err := ioutil.TempDir("", "testProjectLoadAndSave")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to create temporary test directory")
	}
	defer os.RemoveAll(dir)

	p := &project.Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services: []*project.Service{
			&project.Service{
				Name:   "a",
				Path:   "a",
				Mounts: []string{"./foo", "./baz"},
			},
		},
	}
	p.SaveConfig(dir)
	ctx := newCmdCtx(logger, dir)
	loaded, err := project.Load(ctx, dir)
	diffConfig(t, p, loaded)
}

func TestProjectAddService(t *testing.T) {
	logger := logging.NewTestLogger()
	dir, err := ioutil.TempDir("", "testProjectAddService")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to create temporary test directory")
	}
	defer os.RemoveAll(dir)

	p := &project.Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services: []*project.Service{
			&project.Service{
				Name:   "a",
				Path:   "a",
				Mounts: []string{"./a", "./b"},
			},
		},
	}
	p.ConfigureCmdCtx(newCmdCtx(logging.NewTestLogger(), ""))
	if err := p.SaveConfig(dir); err != nil {
		t.Fatalf("%+v", err)
	}

	p.AddService(project.Service{
		Name:   "foo",
		Path:   "bar",
		Mounts: []string{"./foo", "./bar"},
	})

	if err := p.SaveConfig(dir); err != nil {
		t.Fatalf("%+v", err)
	}

	ctx := newCmdCtx(logger, dir)
	if loaded, err := project.Load(ctx, dir); err != nil {
		t.Fatalf("%+v", err)
	} else {
		diffConfig(t, p, loaded)
	}
}

func TestProjectServiceUtils(t *testing.T) {
	p := &project.Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services:   []*project.Service{},
	}
	p.ConfigureCmdCtx(newCmdCtx(logging.NewTestLogger(), ""))
	p.AddService(project.Service{Name: "a", Path: "a"})
	expectation := &project.Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services: []*project.Service{
			&project.Service{Name: "a", Path: "a"},
		},
	}
	diffConfig(t, expectation, p)

	invalidServicePaths := []string{"./b", "b", "./b/", "./foo/a", ".a/a/"}
	for _, path := range invalidServicePaths {
		s, err := p.GetService(path)
		assert.Nil(t, s, "Service should be nil when asking for an invalid service")

		var notAService project.NotAService
		assert.ErrorAsf(t, err, &notAService, "Expected NotAService but got %v", notAService)
	}
	validServicePaths := []string{"a", "./a", "a/", "./a/"}
	for _, path := range validServicePaths {
		s, err := p.GetService(path)
		assert.Nilf(t, err, "Expected error to be nil but found %v", err)
		diffConfig(t, &project.Service{Name: "a", Path: "a"}, s)
	}
}

func TestProjectValidation(t *testing.T) {
	cases := []struct {
		project     project.Project
		expectation bool
	}{
		{project.Project{Name: "Foo"}, true},
		{project.Project{Name: "foo"}, true},
		{project.Project{Name: "Bar"}, true},
		{project.Project{Name: "bar"}, true},
		{project.Project{Name: "foo-bar"}, true},
		{project.Project{Name: "Foo-Bar"}, true},
		{project.Project{Name: "foo_bar"}, true},
		{project.Project{Name: "FooBar"}, true},
		{project.Project{Name: "foo/bar"}, false},
		{project.Project{Name: "1foo/bar"}, false},
		{project.Project{Name: "1"}, false},
		{project.Project{Name: "1foo"}, false},
		{project.Project{Name: ""}, false},
		{project.Project{Name: "-------"}, false},
		{project.Project{Name: "--FooBar--"}, false},
		{project.Project{Name: "a------"}, true},
		{project.Project{Name: "A------"}, true},
		{project.Project{Name: "FooBar", Services: []*project.Service{{Name: "foo"}}}, true},
		{project.Project{Name: "FooBar", Services: []*project.Service{{Name: "-foo"}}}, false},
		{
			project.Project{
				Name:     "FooBar",
				Services: []*project.Service{{Name: "foo"}, {Name: "-foo"}},
			},
			false,
		},
		{
			project.Project{
				Name:     "FooBar",
				Services: []*project.Service{{Name: "foo"}, {Name: "bar"}},
			},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.project.Name, func(t *testing.T) {
			result := tc.project.Validate()
			var wasValid bool
			if result == nil {
				wasValid = true
			}
			assert.Equalf(
				t,
				tc.expectation,
				wasValid,
				"For case %s expected valid to be %t but found %v.",
				tc.project.Name,
				tc.expectation,
				result,
			)
		})
	}
}

func diffConfig(t *testing.T, expectation, actual interface{}) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(project.Service{}, project.Project{}),
	); diff != "" {
		t.Errorf("Expected loaded config to equal input (-want +got):\n%s", diff)
	}
}

func newCmdCtx(logger logging.Logger, dir string) *ctl.CmdCtx {
	return ctl.NewCmdCtx(
		context.Background(),
		logger,
		&config.Config{FoldHome: dir},
		output.NewColorOutput(),
	)
}
