package project_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

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

	_, err = project.Load(logger, dir)
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
	_, err = project.Load(logger, dir)
	if !errors.Is(err, project.InvalidConfig) {
		t.Errorf("Non yaml file should lead to an InvalidConfig error.")
	}

	err = ioutil.WriteFile(cfgPath, []byte("valid: yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}

	cfg, err := project.Load(logger, dir)
	if !errors.Is(err, project.InvalidConfig) {
		logger.Debugf("Loaded config: %+v", cfg)
		t.Errorf("Expected InvalidConfig but got %v\n", err)
	}

	s := `name: blah
maintainer: blah
email: blah
repository: blah`
	err = ioutil.WriteFile(cfgPath, []byte(s), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}
	cfg, err = project.Load(logger, dir)
	if err != nil {
		t.Errorf("Expected to load config but got error %v", err)
	}
	expectation := &project.Project{
		Name:       "blah",
		Maintainer: "blah",
		Email:      "blah",
		Repository: "blah",
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
			&project.Service{Name: "a", Path: "a"},
		},
	}
	p.SaveConfig(dir)

	loaded, err := project.Load(logger, dir)
	diffConfig(t, p, loaded)
}

func TestProjectServiceUtils(t *testing.T) {
	p := &project.Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services:   []*project.Service{},
	}
	p.ConfigureLogger(logging.NewTestLogger())
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
		if !errors.Is(err, project.NotAService) {
			t.Errorf("Expected NotAService but got %v", err)
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

func diffConfig(t *testing.T, expectation, actual interface{}) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(project.Service{}, project.Project{}),
		cmpopts.IgnoreFields(project.Service{}, "Project"),
	); diff != "" {
		t.Errorf("Expected loaded config to equal input (-want +got):\n%s", diff)
	}
}
