package project

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/logging"
)

var (
	NotAFoldProject = errors.New("fold.yaml not found")
	InvalidConfig   = errors.New("invalid config file")
	CantWriteConfig = errors.New("can't write fold.yaml")
)

type Project struct {
	Name       string     `mapstructure:"name"`
	Maintainer string     `mapstructure:"maintainer"`
	Email      string     `mapstructure:"email"`
	Repository string     `mapstructure:"repository"`
	Services   []*Service `mapstructure:"services"`

	logger  logging.Logger
	backend Backend
}

func Load(logger logging.Logger, searchPaths ...string) (*Project, error) {
	if len(searchPaths) == 0 {
		searchPaths = []string{"."}
	}
	return load(logger, searchPaths)
}

func IsAFoldProject(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "fold.yaml")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (p *Project) ConfigureBackend(b Backend) {
	p.backend = b
}

func (p *Project) ConfigureLogger(l logging.Logger) {
	p.logger = l
}

func (p *Project) SaveConfig(to ...string) error {
	if len(to) == 0 {
		to = []string{"."}
	}
	return saveConfig(p, to)
}

func (p *Project) AddService(svc Service) {
	svc.Project = p
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

func (p *Project) GetServices(paths ...string) []*Service {
	// TODO this just ignore invalid services
	var services []*Service
	for _, path := range paths {
		service, err := p.GetService(path)
		if err == nil {
			services = append(services, service)
		}
	}
	return services
}

func (p *Project) Up(ctx context.Context, out io.Writer, services ...*Service) error {
	net := p.network()
	p.logger.Debugf("creating network %v", net)
	err := p.backend.CreateNetworkIfNotExists(net)
	if err != nil {
		p.logger.Debugf("failed to create network for project %s: %v", p.Name, err)
		return err
	}

	for _, service := range services {
		err = service.Start(ctx, out, net)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) Down() error {
	for _, service := range p.Services {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	net := p.network()
	err := p.backend.RemoveNetworkIfExists(net)
	if err != nil {
		p.logger.Debugf("failed to remove network for project %s: %v", p.Name, err)
		return err
	}
	return nil
}

func (p *Project) network() *container.Network {
	name := fmt.Sprintf("foldnet-%s", p.Name)
	return p.backend.NewNetwork(name)
}
