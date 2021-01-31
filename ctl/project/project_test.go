package project

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

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
	if IsAFoldProject(dir) {
		t.Errorf("Empty directory should not be a valid fold project")
	}
	cfgPath := filepath.Join(dir, "fold.yaml")

	_, err = load(logger, dir)
	if !errors.Is(err, NotAFoldProject) {
		t.Errorf("Empty directory should not be a valid fold project")
	}

	err = ioutil.WriteFile(cfgPath, []byte("not valid yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}
	if !IsAFoldProject(dir) {
		t.Errorf("After creating a fold.yaml file it should be a valid fold project")
	}
	_, err = load(logger, dir)
	if !errors.Is(err, InvalidConfig) {
		t.Errorf("Non yaml file should lead to an InvalidConfig error.")
	}

	err = ioutil.WriteFile(cfgPath, []byte("valid: yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to config file.")
	}

	cfg, err := load(logger, dir)
	if !errors.Is(err, InvalidConfig) {
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
	cfg, err = load(logger, dir)
	if err != nil {
		t.Errorf("Expected to load config but got error %v", err)
	}
	expectation := &Project{
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

	project := &Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services: []*Service{
			&Service{Name: "a", Path: "a"},
		},
	}
	saveConfig(project, dir)

	loaded, err := load(logger, dir)
	diffConfig(t, project, loaded)
}

func TestProjectServiceUtils(t *testing.T) {
	project := &Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services:   []*Service{},
		logger:     logging.NewTestLogger(),
	}
	project.AddService(Service{Name: "a", Path: "a"})
	expectation := &Project{
		Name:       "foo",
		Maintainer: "tom",
		Email:      "tom@tom.com",
		Repository: "github.com/tom",
		Services: []*Service{
			&Service{Name: "a", Path: "a"},
		},
	}
	diffConfig(t, expectation, project)

	invalidServicePaths := []string{"./b", "b", "./b/", "./foo/a", ".a/a/"}
	for _, p := range invalidServicePaths {
		s, err := project.GetService(p)
		if s != nil {
			t.Errorf("Service should be nil when asking for an invalid service")
		}
		if !errors.Is(err, NotAService) {
			t.Errorf("Expected NotAService but got %v", err)
		}
	}
	validServicePaths := []string{"a", "./a", "a/", "./a/"}
	for _, p := range validServicePaths {
		s, err := project.GetService(p)
		if err != nil {
			t.Errorf("Expected error to be nil but found %v", err)
		}
		diffConfig(t, &Service{Name: "a", Path: "a"}, s)
	}
}

func diffConfig(t *testing.T, expectation, actual interface{}) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(Service{}, Project{}),
	); diff != "" {
		t.Errorf("Expected loaded config to equal input (-want +got):\n%s", diff)
	}
}
