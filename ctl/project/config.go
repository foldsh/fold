package project

import (
	"errors"
	"path/filepath"

	"github.com/foldsh/fold/logging"
	"github.com/spf13/viper"
)

// Looks for fold.yaml in the current directory and loads it.
func load(logger logging.Logger, location string) (*Project, error) {
	v := newViper(location)
	var fileNotFound viper.ConfigFileNotFoundError
	err := v.ReadInConfig()
	if err != nil {
		if errors.As(err, &fileNotFound) {
			logger.Debugf("config file not found %v", err)
			return nil, NotAFoldProject
		} else {
			logger.Debugf("config file invalid %v", err)
			return nil, InvalidConfig
		}
	}
	if !validateConfig(v) {
		logger.Debugf("invalid config: must set name, maintainer, email and repository")
		return nil, InvalidConfig
	}
	var p *Project
	err = v.Unmarshal(&p)
	if err != nil {
		logger.Debugf("failed to unmarshal config %v", err)
		return nil, InvalidConfig
	}
	p.logger = logger
	for _, s := range p.Services {
		s.project = p
	}
	return p, nil
}

func saveConfig(p *Project, to string) error {
	v := newViper(to)
	v.Set("name", p.Name)
	v.Set("maintainer", p.Maintainer)
	v.Set("email", p.Email)
	v.Set("repository", p.Repository)
	v.Set("services", p.Services)
	if err := v.WriteConfigAs(filepath.Join(to, "fold.yaml")); err != nil {
		p.logger.Debugf("Failed to write config %+v", err)
		return CantWriteConfig
	}
	return nil
}

func newViper(configPath string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("fold")
	v.SetConfigType("yaml")
	return v
}

func validateConfig(v *viper.Viper) bool {
	n := v.IsSet("name")
	m := v.IsSet("maintainer")
	e := v.IsSet("email")
	r := v.IsSet("repository")
	return n && m && e && r
}
