package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/foldsh/fold/ctl/fs"
	"github.com/foldsh/fold/version"
	"github.com/spf13/viper"
)

var (
	CreateConfigError = errors.New("failed to create the default foldctl config file")
	ReadConfigError   = errors.New("failed to read the foldctl config file")
)

type Config struct {
	AccessToken string `mapstructure:"access-token"`
	Version     string `mapstructure:"version"`

	FoldHome      string
	FoldTemplates string
}

func Load(path string) (*Config, error) {
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
				return nil, CreateConfigError
			}
			// We successfully wrote the default config, so lets try to load it.
			continue
		} else {
			// There was some other error, likely the config file was malformed.
			// Bail and return the error.
			return nil, ReadConfigError
		}
	}
	return &Config{
		AccessToken:   viper.GetString("access-token"),
		Version:       viper.GetString("version"),
		FoldHome:      path,
		FoldTemplates: fs.FoldTemplates(path),
	}, nil
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
	writer.Set("access-token", "")
	writer.Set("version", version.FoldVersion.String())
	err := writer.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}
