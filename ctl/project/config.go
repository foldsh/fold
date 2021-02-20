package project

import (
	"errors"
	"path/filepath"

	"github.com/foldsh/fold/logging"
	"github.com/spf13/viper"
)

type onDiskProject struct {
	Name       string          `mapstructure:"name"`
	Maintainer string          `mapstructure:"maintainer"`
	Email      string          `mapstructure:"email"`
	Repository string          `mapstructure:"repository"`
	Services   []onDiskService `mapstructure:"services"`
}

type onDiskService struct {
	Name   string   `mapstructure:"name"`
	Path   string   `mapstructure:"path"`
	Mounts []string `mapstructure:"mounts"`
}

func (odp *onDiskProject) unmarshall(logger logging.Logger) *Project {
	p := &Project{
		Name:       odp.Name,
		Maintainer: odp.Maintainer,
		Email:      odp.Email,
		Repository: odp.Repository,
	}
	p.ConfigureLogger(logger)
	for _, ods := range odp.Services {
		svc := &Service{Name: ods.Name, Path: ods.Path, Mounts: ods.Mounts}
		svc.project = p
		p.Services = append(p.Services, svc)
	}
	return p
}

func marshallProject(p *Project) *onDiskProject {
	odp := &onDiskProject{
		Name:       p.Name,
		Maintainer: p.Maintainer,
		Email:      p.Email,
		Repository: p.Repository,
	}
	for _, s := range p.Services {
		ods := onDiskService{Name: s.Name, Path: s.Path, Mounts: s.Mounts}
		odp.Services = append(odp.Services, ods)
	}
	return odp
}

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
	var odp *onDiskProject
	err = v.Unmarshal(&odp)
	if err != nil {
		logger.Debugf("failed to unmarshal config %v", err)
		return nil, InvalidConfig
	}
	return odp.unmarshall(logger), nil
}

func saveConfig(p *Project, to string) error {
	odp := marshallProject(p)
	v := newViper(to)
	v.Set("name", odp.Name)
	v.Set("maintainer", odp.Maintainer)
	v.Set("email", odp.Email)
	v.Set("repository", odp.Repository)
	v.Set("services", odp.Services)
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
