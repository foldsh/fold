package project_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/foldsh/fold/ctl"
	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/gateway"
	"github.com/foldsh/fold/ctl/mocks"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
)

var workflowTests = []struct {
	project *project.Project
}{
	{project: makeProject("no-services")},
	{project: makeProject("one-service", &project.Service{Name: "one", Path: "./one"})},
	{
		project: makeProject(
			"two-services",
			&project.Service{Name: "one", Path: "./one"},
			&project.Service{Name: "two", Path: "./two"},
		),
	},
}

// The goal with this is to run a load of example projects through a
// project 'lifecycle'. I.e. we bring up all the services, one by one
// and then bring it down. This should result in a consistent pattern
// of calls to the container api, which we will make assertions about.
func TestProjectUp(t *testing.T) {
	for _, tc := range workflowTests {
		api := &mocks.ContainerAPI{}
		proj := tc.project
		proj.ConfigureContainerAPI(api)
		t.Run(tc.project.Name, func(t *testing.T) {
			out := new(bytes.Buffer)
			netName := fmt.Sprintf("foldnet-%s", proj.Name)
			net := &container.Network{Name: netName}
			// Should set up network
			api.
				On("NewNetwork", netName).
				Return(net)
			api.
				On("NetworkExists", net).
				Return(false, nil)
			api.
				On("CreateNetwork", net).
				Return(nil)

			// Should start the gateway
			gwSvc := proj.NewService("foldgw")
			gwImgName := (&gateway.Gateway{}).ImageName()
			gwImg := &container.Image{Name: gwImgName}
			gwContainerName := fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)
			gwContainer := &container.Container{ID: fmt.Sprintf("%d", 0), Name: gwContainerName}
			api.
				On("GetContainer", gwContainerName).
				Return(nil, nil)
			api.
				On("GetImage", gwImgName).
				Return(nil, nil)
			api.
				On("PullImage", gwImgName).
				Return(gwImg, nil)
			api.
				On("NewContainer", gwContainerName, *gwImg).
				Return(gwContainer)
			api.
				On("RunContainer",
					net,
					&container.Container{
						ID:           fmt.Sprintf("%d", 0),
						Name:         gwContainerName,
						NetworkAlias: gwSvc.Name,
						Environment:  map[string]string{"FOLD_SERVICE_NAME": gwSvc.Name},
					},
				).
				Return(nil)

			// Should set up the services
			for i, svc := range proj.Services {
				path, err := svc.AbsPath()
				if err != nil {
					t.Errorf("failed to get service abs path %v", err)
				}
				img := &container.Image{
					Name: fmt.Sprintf("foldlocal/%s/%s:latest", svc.Id(), svc.Name),
					Src:  path,
				}
				containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
				api.
					On("GetContainer", containerName).
					Return(nil, nil)
				api.
					On("BuildImage", img).
					Return(nil)
				api.
					On("NewContainer", containerName, *img).
					Return(&container.Container{ID: fmt.Sprintf("%d", i), Name: containerName})
				modifiedCon := &container.Container{
					ID:           fmt.Sprintf("%d", i),
					Name:         containerName,
					NetworkAlias: svc.Name,
					Environment:  map[string]string{"FOLD_SERVICE_NAME": svc.Name},
				}
				api.
					On("RunContainer", net, modifiedCon).
					Return(nil)
			}
			proj.Up(out, proj.Services...)

			// Should get the logs
			for _, svc := range proj.Services {
				api.
					On("ContainerLogs", mock.Anything).
					Return(&container.LogStream{}, nil)
				svc.Logs()
			}
		})
		api.AssertExpectations(t)
	}
}
func TestProjectDown(t *testing.T) {
	for _, tc := range workflowTests {
		api := &mocks.ContainerAPI{}
		proj := tc.project
		proj.ConfigureContainerAPI(api)
		t.Run(tc.project.Name, func(t *testing.T) {
			netName := fmt.Sprintf("foldnet-%s", proj.Name)
			net := &container.Network{Name: netName}

			gwSvc := proj.NewService("foldgw")
			gwContainerName := fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)
			gwContainer := &container.Container{ID: fmt.Sprintf("%d", 0), Name: gwContainerName}

			// Should take down the services
			for i, svc := range proj.Services {
				containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
				container := &container.Container{ID: fmt.Sprintf("%d", i), Name: containerName}
				api.
					On("GetContainer", containerName).
					Return(container, nil)
				api.
					On("StopContainer", container).
					Return(nil)
			}
			// Should take down the gateway
			api.
				On("GetContainer", gwContainerName).
				Return(gwContainer, nil)
			api.
				On("StopContainer", gwContainer).
				Return(nil)
			// Should take down the network
			api.
				On("NewNetwork", netName).
				Return(&container.Network{Name: netName})
			api.
				On("NetworkExists", net).
				Return(true, nil)
			api.
				On("RemoveNetwork", net).
				Return(nil)
			proj.Down()
		})
		api.AssertExpectations(t)
	}
}

func TestUpDoesntDuplicateResources(t *testing.T) {
	// For this test we will set the mocked calls to the container API
	// to return resources. This should result in a short circuit,
	// and no attempt should be made to create the resources again.
	api := &mocks.ContainerAPI{}
	proj := makeProject("one-service", &project.Service{Name: "one", Path: "./one"})
	proj.ConfigureContainerAPI(api)
	svc, err := proj.GetService("./one")
	require.Nil(t, err)

	netName := fmt.Sprintf("foldnet-%s", proj.Name)
	net := &container.Network{Name: netName}
	// Should reuse network
	api.
		On("NewNetwork", netName).
		Return(net)
	api.
		On("NetworkExists", net).
		Return(true, nil)
	// Should reuse gateway
	gwSvc := proj.NewService("foldgw")
	gwContainerName := fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)
	gwContainer := &container.Container{ID: fmt.Sprintf("%d", 0), Name: gwContainerName}
	api.
		On("GetContainer", fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)).
		Return(gwContainer, nil)

	containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
	container := &container.Container{ID: fmt.Sprintf("%d", 0), Name: containerName}
	// Should reuse service container
	api.
		On("GetContainer", containerName).
		Return(container, nil)
	err = proj.Up(&bytes.Buffer{}, svc)
	assert.Nil(t, err)
	api.AssertExpectations(t)
}

func makeProject(name string, services ...*project.Service) *project.Project {
	p := &project.Project{
		Name: name,
	}
	p.ConfigureCmdCtx(
		ctl.NewCmdCtx(
			context.Background(),
			logging.NewTestLogger(),
			&config.Config{},
			output.NewColorOutput(),
		),
	)
	for _, s := range services {
		svc := p.NewService(s.Name)
		svc.Port = s.Port
		p.Services = append(p.Services, svc)
	}
	return p
}
