package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/internal/testutils"
)

func TestLoadCtlConfig(t *testing.T) {
	viper.Reset()
	err := loadConfigAtPath("./testdata/")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to load config")
	}
	if viper.Get("version") != "1.2.3" {
		t.Fatalf("Expected '1.2.3' but found %s", viper.Get("version"))
	}
	if viper.Get("name") != "John Smith" {
		t.Fatalf("Expected 'John Smith' but found %s", viper.Get("name"))
	}
	if viper.Get("email") != "test@test.com" {
		t.Fatalf("Expected test@test.com but found %s", viper.Get("email"))
	}
	if viper.Get("access-token") != "ABCDEF123456" {
		t.Fatalf("Expected 'ABCDEF123456' but found %s", viper.Get("access-token"))
	}
}

func TestConfigCreatedIfNotPresent(t *testing.T) {
	viper.Reset()
	dir, err := ioutil.TempDir("", "ctl-test")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to create temporary test directory")
	}
	defer os.RemoveAll(dir)
	err = loadConfigAtPath(dir)
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to create default config")
	}

	if viper.Get("version") != ctl.Version.String() {
		t.Fatalf("Expected '%s' but found %s", ctl.Version.String(), viper.Get("version"))
	}
	if viper.Get("name") != "" {
		t.Fatalf("Expected '' but found %s", viper.Get("name"))
	}
	if viper.Get("email") != "" {
		t.Fatalf("Expected '' but found %s", viper.Get("email"))
	}
	if viper.Get("access-token") != "" {
		t.Fatalf("Expected '' but found %v", viper.Get("access-token"))
	}
}

func TestMakeFoldConfig(t *testing.T) {
	viper.Reset()
	err := loadConfigAtPath("./testdata/")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to load config")
	}
	cfg, err := makeFoldConfig()
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed unmarshal config")
	}
	expectation := ctl.Config{
		Version:     "1.2.3",
		Name:        "John Smith",
		Email:       "test@test.com",
		AccessToken: "ABCDEF123456",
	}
	testutils.Diff(t, expectation, cfg, "Parsed config did not match expectation")
}
