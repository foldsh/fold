package project_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
)

func TestProjectConfigValidation(t *testing.T) {
	logger := logging.NewTestLogger()
	dir, err := ioutil.TempDir("", "testProjectConfig")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to create temporary test directory")
	}
	defer os.RemoveAll(dir)
	if project.IsAFoldProject(dir) {
		t.Errorf("Empty directory should not be a valid fold project")
	}
	cfgPath := filepath.Join(dir, "fold.yaml")

	ctx := newCmdCtx(logger, dir)
	_, err = project.Load(ctx, dir)
	if !errors.Is(err, project.NotAFoldProject) {
		t.Errorf("Empty directory should not be a valid fold project")
	}

	err = ioutil.WriteFile(cfgPath, []byte("not valid yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}
	if !project.IsAFoldProject(dir) {
		t.Errorf("After creating a fold.yaml file it should be a valid fold project")
	}
	_, err = project.Load(ctx, dir)
	if !errors.Is(err, project.InvalidConfig) {
		t.Errorf("Non yaml file should lead to an InvalidConfig error.")
	}

	err = ioutil.WriteFile(cfgPath, []byte("valid: yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}

	cfg, err := project.Load(ctx, dir)
	if !errors.Is(err, project.InvalidConfig) {
		logger.Debugf("Loaded config: %+v", cfg)
		t.Errorf("Expected InvalidConfig but got %v\n", err)
	}

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
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}
	cfg, err = project.Load(ctx, dir)
	if err != nil {
		t.Errorf("Expected to load config but got error %v", err)
	}
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
		if s != nil {
			t.Errorf("Service should be nil when asking for an invalid service")
		}
		var notAService project.NotAService
		if !errors.As(err, &notAService) {
			t.Errorf("Expected NotAService but got %v", notAService)
		}
	}
	validServicePaths := []string{"a", "./a", "a/", "./a/"}
	for _, path := range validServicePaths {
		s, err := p.GetService(path)
		if err != nil {
			t.Errorf("Expected error to be nil but found %v", err)
		}
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
			if wasValid != tc.expectation {
				t.Errorf(
					"For case %s expected valid to be %t but found %v.",
					tc.project.Name,
					tc.expectation,
					result,
				)
			}
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
