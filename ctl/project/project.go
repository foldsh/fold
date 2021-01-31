package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/foldsh/fold/logging"
)

type Project struct {
	Name       string     `mapstructure:"name"`
	Maintainer string     `mapstructure:"maintainer"`
	Email      string     `mapstructure:"email"`
	Repository string     `mapstructure:"repository"`
	Services   []*Service `mapstructure:"services"`

	logger logging.Logger
}

var (
	NotAFoldProject = errors.New("fold.yaml not found")
	InvalidConfig   = errors.New("invalid config file")
	CantWriteConfig = errors.New("can't write fold.yaml")

	NotAService = errors.New("not a valid service")
)

func IsAFoldProject(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "fold.yaml")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Load(logger logging.Logger) (*Project, error) {
	return load(logger, ".")
}

func (p *Project) SaveConfig() error {
	return saveConfig(p, ".")
}

func (p *Project) AddService(svc Service) {
	svc.project = p
	p.Services = append(p.Services, &svc)
}

// This works with either a path to a service or just a service name.
func (p *Project) GetService(path string) (*Service, error) {
	name := filepath.Clean(path)
	p.logger.Debugf("fetching service for path %s", name)
	for _, svc := range p.Services {
		if svc.Name == name {
			return svc, nil
		}
	}
	return nil, NotAService
}

// Looks for fold.yaml in the current directory and loads it.
func load(logger logging.Logger, path string) (*Project, error) {
	v := newViper(path)
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

func saveConfig(p *Project, path string) error {
	v := newViper(path)
	v.Set("name", p.Name)
	v.Set("maintainer", p.Maintainer)
	v.Set("email", p.Email)
	v.Set("repository", p.Repository)
	v.Set("services", p.Services)
	err := v.WriteConfigAs(filepath.Join(path, "fold.yaml"))
	if err != nil {
		fmt.Printf("%v", err)
		return CantWriteConfig
	}
	return nil
}

func newViper(path string) *viper.Viper {
	v := viper.New()
	v.AddConfigPath(path)
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
