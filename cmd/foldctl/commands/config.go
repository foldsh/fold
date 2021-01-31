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
)

var foldHome string

func init() {
	home, err := fs.FoldHome()
	if err != nil {
		panic(errors.New("failed to locate home directory"))
	}
	foldHome = home
}

func loadConfig() error {
	return loadConfigAtPath(foldHome)
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
				return errors.New("failed to create the default foldctl config file")
			}
			// We successfully wrote the default config, so lets try to load it.
			continue
		} else {
			// There was some other error, likely the config file was malformed.
			// Bail and return the error.
			return errors.New("failed to read the foldctl config file")
		}
	}
	return nil
}

func loadFoldConfig() (ctl.Config, error) {
	var cfg ctl.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, errors.New("failed to unmarshal the foldctl config file ")
	}
	return cfg, nil
}

func writeDefaultConfig(path string) error {
	if err := os.MkdirAll(path, DIR_PERMISSIONS); err != nil {
		fmt.Println(err)
		return err
	}
	writer := viper.New()
	writer.AddConfigPath(path)
	writer.SetConfigName("config")
	writer.SetConfigType("yaml")
	writer.Set("version", ctl.Version.String())
	writer.Set("name", "")
	writer.Set("email", "")
	writer.Set("access-token", "")
	err := writer.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}
