package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/version"
)

func TestLoadCtlConfig(t *testing.T) {
	viper.Reset()
	cfg, err := config.Load("./testdata/")
	require.Nil(t, err)

	assert.Equal(t, "1.2.3", cfg.Version, "version should equal 1.2.3")
	assert.Equal(t, "ABCDEF123456", cfg.AccessToken, "access token should equal ABCDEF123456")
}

func TestConfigCreatedIfNotPresent(t *testing.T) {
	viper.Reset()
	dir, err := ioutil.TempDir("", "fold.ctl.config.test")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	cfg, err := config.Load(dir)
	require.Nil(t, err)

	v := version.FoldVersion.String()
	assert.Equal(t, v, cfg.Version, fmt.Sprintf("version should equal %s", v))
	assert.Equal(t, "", cfg.AccessToken, "access token should be empty string")
}
