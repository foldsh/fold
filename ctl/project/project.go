package project

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/gateway"
	"github.com/foldsh/fold/ctl/output"
)

var (
	ProjectNameRegex = `^[a-zA-Z][a-zA-Z-_]+$`

	NotAFoldProject = errors.New("fold.yaml not found")
	InvalidConfig   = errors.New("invalid config file")
	CantWriteConfig = errors.New("can't write fold.yaml")
)

type Project struct {
	Name       string
	Maintainer string
	Email      string
	Repository string
	Services   []*Service

	gatewayPort int
	ctx         *ctl.CmdCtx
	api         ContainerAPI
}

func Load(ctx *ctl.CmdCtx, projectPath string) (*Project, error) {
	return load(ctx, projectPath)
}

func Home() (string, error) {
	if abs, err := filepath.Abs("."); err != nil {
		return "", errors.New("can't locate fold home directory")
	} else {
		return abs, nil
	}
}

func IsAFoldProject(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "fold.yaml")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (p *Project) SaveConfig(projectPath string) error {
	return saveConfig(p, projectPath)
}

func (p *Project) ConfigureGatewayPort(port int) {
	p.gatewayPort = port
}

func (p *Project) ConfigureContainerAPI(b ContainerAPI) {
	p.api = b
}

func (p *Project) ConfigureCmdCtx(ctx *ctl.CmdCtx) {
	p.ctx = ctx
}

func (p *Project) NewService(name string) *Service {
	return &Service{Name: name, project: p, ctx: p.ctx}
}

func (p *Project) Validate() error {
	matched, _ := regexp.MatchString(ProjectNameRegex, p.Name)
	if !matched {
		return fmt.Errorf(
			"%s is not a valid project name, it must match the regex %s",
			p.Name,
			ProjectNameRegex,
		)
	}
	for _, svc := range p.Services {
		err := svc.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AddService(svc Service) {
	svc.project = p
	p.Services = append(p.Services, &svc)
}

// This works with either a path to a service or just a service name.
func (p *Project) GetService(path string) (*Service, error) {
	name := filepath.Clean(path)
	p.ctx.Logger.Debugf("fetching service for path %s", name)
	for _, svc := range p.Services {
		if svc.Name == name {
			return svc, nil
		}
	}
	p.ctx.Logger.Debugf("no service found for path %s", name)
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

func (p *Project) Up(out io.Writer, services ...*Service) error {
	p.ctx.InformHeader("Bringing up the fold development server for project %s...", p.Name)

	// Ensure network
	net := p.network()
	exists, err := p.api.NetworkExists(net)
	if err != nil {
		return err
	} else if !exists {
		p.ctx.InformBody("Creating the local network for project %s...", p.Name)
		if err := p.api.CreateNetwork(net); err != nil {
			p.ctx.Logger.Debugf("Failed to create network for project %s: %v", p.Name, err)
			return err
		}
	}
	if err := p.startGateway(net); err != nil {
		return err
	}

	// Bring up services
	for _, service := range services {
		p.ctx.Informf("Starting container for service %s...", service.Name)
		// Check if the service is already up.
		container, err := p.api.GetContainer(service.containerName())
		if err != nil {
			p.ctx.Logger.Debugf(
				"Failed to check if container for service %s already exists",
				service.Name,
			)
			return err
		}
		if container != nil {
			p.ctx.Informf("Service %s is already up, no need to do anything", service.Name)
			service.container = container
			continue
		}
		// Build the service
		p.ctx.Informf("Building image for service %s...", service.Name)
		img, err := service.Build(p.ctx.Context, out)
		if err != nil {
			return err
		}
		// Start the service
		p.ctx.Informf("Starting service %s...", service.Name)
		if err := service.Start(img, net); err != nil {
			return err
		}
	}
	p.ctx.Informf("The fold development server is now ready")
	return nil
}

func (p *Project) Down() error {
	p.ctx.Informf("Taking down the fold development server for project %s...", p.Name)

	// Take down services - doing this first ensures that we remove fold containers even
	// if the user has done something like delete the network manually.
	for _, service := range p.Services {
		if err := service.Stop(); err != nil {
			return err
		}
	}

	// Stop the gateway
	if err := p.stopGateway(); err != nil {
		return err
	}

	// Determine if we need to take down the network.
	net := p.network()
	exists, err := p.api.NetworkExists(net)
	p.ctx.Informf("checking network exists %v %v", exists, err)
	if err != nil {
		return err
	} else if !exists {
		p.ctx.Informf("Local network for project %s is not up, nothing to do.", p.Name)
		return nil
	}
	// It exists, so remove it.
	p.ctx.Informf("Taking down the local network for project %s...", p.Name)
	if err = p.api.RemoveNetwork(net); err != nil {
		p.ctx.Logger.Debugf("failed to remove network for project %s: %v", p.Name, err)
		return err
	}

	p.ctx.Inform(output.Success("The fold development server has been taken down successfully"))
	return nil
}

func (p *Project) network() *container.Network {
	name := fmt.Sprintf("foldnet-%s", p.Name)
	return p.api.NewNetwork(name)
}

func (p *Project) gateway() *gateway.Gateway {
	return &gateway.Gateway{Port: p.gatewayPort}
}

func (p *Project) startGateway(net *container.Network) error {
	gw := p.gateway()
	p.ctx.Informf("Starting fold local gateway on port %d...", gw.Port)
	if con, err := p.isGatewayUp(gw); err != nil {
		return err
	} else if con != nil {
		p.ctx.Informf("Gateway is already up, nothing to do.")
		return nil
	}
	imgName := gw.ImageName()
	img, err := p.pullImageIfRequired(imgName)
	if err != nil {
		p.ctx.Logger.Debugf("failed to pull the image for the gateway: %v", err)
		return fmt.Errorf("failed to pull image %s", imgName)
	}
	gwService := p.gatewayService(gw)
	err = gwService.Start(img, net)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) stopGateway() error {
	p.ctx.Informf("Stopping fold local gateway...")
	gw := p.gateway()
	if con, err := p.isGatewayUp(gw); err != nil {
		p.ctx.Logger.Debugf("Failed to confirm if gateway is running: %v", err)
		return err
	} else if con == nil {
		p.ctx.Informf("Gateway is not up, nothing to do.")
		return nil
	} else {
		if err := p.api.StopContainer(con); err != nil {
			p.ctx.Logger.Debugf("Failed to stop the gateway: %v", err)
			return err
		}
	}
	return nil
}

func (p *Project) isGatewayUp(gw *gateway.Gateway) (*container.Container, error) {
	p.ctx.Logger.Debugf("Checking if gateway is running")
	svc := p.gatewayService(gw)
	con, err := p.api.GetContainer(svc.containerName())
	if err != nil {
		p.ctx.Logger.Debugf("Failed to check if gateway is already up: %v", err)
		return nil, err
	}
	if con == nil {
		return nil, nil
	}
	return con, nil
}

func (p *Project) gatewayService(gw *gateway.Gateway) *Service {
	svc := p.NewService("foldgw")
	svc.Port = gw.Port
	return svc
}

func (p *Project) pullImageIfRequired(image string) (*container.Image, error) {
	img, err := p.api.GetImage(image)
	if err != nil {
		return nil, err
	}
	if img != nil {
		p.ctx.Logger.Debugf("Image %s already available locally, nothing to do.", image)
		return img, nil
	}
	p.ctx.Logger.Debugf("Pulling image %s", image)
	img, err = p.api.PullImage(image)
	return img, nil
}
