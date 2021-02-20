// Utilities for loading foldctl config with viper.
// This is distinct from the fold project config!
package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/version"
)

var (
	foldHome      string
	foldBin       string
	foldTemplates string

	couldNotCreateDefaultConfig = errors.New("failed to create the default foldctl config file")
	couldNotReadConfigFile      = errors.New("failed to read the foldctl config file")
)

func init() {
	home, err := fs.FoldHome()
	exitIfErr(err, "Failed to locate fold home directory at ~/.fold.")
	foldHome = home
	foldBin = fs.FoldBin(foldHome)
	foldTemplates = fs.FoldTemplates(foldHome)
}

func loadConfigAtPath(path string) error {
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	var (
		err          error
		fileNotFound viper.ConfigFileNotFoundError
	)
	for {
		err = viper.ReadInConfig()
		if err == nil {
			// The config was successfully read in, so break and try to Ummarshal it.
			break
		} else if errors.As(err, &fileNotFound) {
			// The config file didn't exist, so create it.
			err = writeDefaultConfig(path)
			if err != nil {
				// Creating the config failed, so we bail and return an error.
				return couldNotCreateDefaultConfig
			}
			// We successfully wrote the default config, so lets try to load it.
			continue
		} else {
			// There was some other error, likely the config file was malformed.
			// Bail and return the error.
			return couldNotReadConfigFile
		}
	}
	return nil
}

func makeFoldConfig() (ctl.Config, error) {
	var cfg ctl.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, errors.New("failed to unmarshal the foldctl config file ")
	}
	return cfg, nil
}

func writeDefaultConfig(path string) error {
	if err := os.MkdirAll(path, fs.DIR_PERMISSIONS); err != nil {
		fmt.Println(err)
		return err
	}
	writer := viper.New()
	writer.AddConfigPath(path)
	writer.SetConfigName("config")
	writer.SetConfigType("yaml")
	writer.Set("version", version.FoldVersion.String())
	writer.Set("name", "")
	writer.Set("email", "")
	writer.Set("access-token", "")
	err := writer.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}
