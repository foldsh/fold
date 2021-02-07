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

	logger logging.Logger
	api    ContainerAPI
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

func (p *Project) ConfigureContainerAPI(b ContainerAPI) {
	p.api = b
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
	p.logger.Debugf("no service found for path %s", name)
	return nil, NotAService{path}
}

func (p *Project) GetServices(paths ...string) ([]*Service, error) {
	// TODO this just ignore invalid services
	var services []*Service
	for _, path := range paths {
		service, err := p.GetService(path)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (p *Project) Up(ctx context.Context, out io.Writer, services ...*Service) error {
	p.logger.Infof("Bringing up the fold development server for project %s...", p.Name)

	// Ensure network
	net := p.network()
	exists, err := p.api.NetworkExists(net)
	if err != nil {
		return err
	} else if !exists {
		p.logger.Infof("Creating the local network for project %s...", p.Name)
		err := p.api.CreateNetwork(net)
		if err != nil {
			p.logger.Debugf("Failed to create network for project %s: %v", p.Name, err)
			return err
		}
	}

	// Bring up services
	for _, service := range services {
		err = service.Start(ctx, out, net)
		if err != nil {
			return err
		}
	}
	p.logger.Infof("The fold development server is now ready")
	return nil
}

func (p *Project) Down() error {
	p.logger.Infof("Taking down the fold development server for project %s...", p.Name)

	// Take down services - doing this first ensures that we remove fold containers even
	// if the user has done something like delete the network manually.
	for _, service := range p.Services {
		err := service.Stop()
		if err != nil {
			return err
		}
	}

	// Determine if we need to take down the network.
	net := p.network()
	exists, err := p.api.NetworkExists(net)
	if err != nil {
		return err
	} else if !exists {
		p.logger.Infof("Local network for project %s is not up, nothing to do.", p.Name)
		return nil
	}
	// It exists, so remove it.
	p.logger.Infof("Taking down the local network for project %s...", p.Name)
	err = p.api.RemoveNetwork(net)
	if err != nil {
		p.logger.Debugf("failed to remove network for project %s: %v", p.Name, err)
		return err
	}

	p.logger.Infof("The fold development server has been taken down successfully")
	return nil
}

func (p *Project) network() *container.Network {
	name := fmt.Sprintf("foldnet-%s", p.Name)
	return p.api.NewNetwork(name)
}
