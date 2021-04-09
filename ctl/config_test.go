package ctl_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/foldsh/fold/version"
)

func TestLoadCtlConfig(t *testing.T) {
	viper.Reset()
	err := ctl.Load("./testdata/")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Fatal("Failed to load config")
	}
	if viper.Get("version") != "1.2.3" {
		t.Fatalf("Expected '1.2.3' but found %s", viper.Get("version"))
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

	if viper.Get("version") != version.FoldVersion.String() {
		t.Fatalf("Expected '%s' but found %s", version.FoldVersion.String(), viper.Get("version"))
	}
	if viper.Get("access-token") != "" {
		t.Fatalf("Expected '' but found %v", viper.Get("access-token"))
	}
}
