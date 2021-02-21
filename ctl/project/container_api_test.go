package project_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/gateway"
	"github.com/foldsh/fold/ctl/project"
	"github.com/foldsh/fold/logging"
	"github.com/golang/mock/gomock"
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
func TestProjectWorkflow(t *testing.T) {
	for _, tc := range workflowTests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		api := NewMockContainerAPI(ctrl)

		proj := tc.project
		proj.ConfigureContainerAPI(api)
		proj.ConfigureLogger(logging.NewTestLogger())
		t.Run(tc.project.Name, func(t *testing.T) {
			out := new(bytes.Buffer)
			netName := fmt.Sprintf("foldnet-%s", proj.Name)
			net := &container.Network{Name: netName}
			// Should set up network
			api.
				EXPECT().
				NewNetwork(netName).
				Return(net)
			api.
				EXPECT().
				NetworkExists(net).
				Return(false, nil)
			api.
				EXPECT().
				CreateNetwork(net)

			// Should start the gateway
			gwSvc := proj.NewService("foldgw")
			gwImgName := (&gateway.Gateway{}).ImageName()
			gwImg := &container.Image{Name: gwImgName}
			gwContainerName := fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)
			gwContainer := &container.Container{ID: fmt.Sprintf("%d", 0), Name: gwContainerName}
			api.
				EXPECT().
				GetContainer(gwContainerName).
				Return(nil, nil)
			api.
				EXPECT().
				GetImage(gwImgName).
				Return(nil, nil)
			api.
				EXPECT().
				PullImage(gwImgName).
				Return(gwImg, nil)
			api.
				EXPECT().
				NewContainer(gwContainerName, *gwImg).
				Return(gwContainer)
			api.
				EXPECT().
				RunContainer(
					net,
					&container.Container{
						ID:           fmt.Sprintf("%d", 0),
						Name:         gwContainerName,
						NetworkAlias: gwSvc.Name,
						Environment:  map[string]string{"FOLD_SERVICE_NAME": gwSvc.Name},
					},
				)

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
					EXPECT().
					GetContainer(containerName).
					Return(nil, nil)
				api.
					EXPECT().
					BuildImage(img)
				api.
					EXPECT().
					NewContainer(containerName, *img).
					Return(&container.Container{ID: fmt.Sprintf("%d", i), Name: containerName})
				modifiedCon := &container.Container{
					ID:           fmt.Sprintf("%d", i),
					Name:         containerName,
					NetworkAlias: svc.Name,
					Environment:  map[string]string{"FOLD_SERVICE_NAME": svc.Name},
				}
				api.
					EXPECT().
					RunContainer(gomock.Eq(net), gomock.Eq(modifiedCon))
			}
			proj.Up(context.Background(), out, proj.Services...)

			// Should take down the services
			for i, svc := range proj.Services {
				containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
				container := &container.Container{ID: fmt.Sprintf("%d", i), Name: containerName}
				api.
					EXPECT().
					GetContainer(containerName).
					Return(container, nil)
				api.
					EXPECT().
					StopContainer(container)
			}
			// Should take down the gateway
			api.
				EXPECT().
				GetContainer(gwContainerName).
				Return(gwContainer, nil)
			api.
				EXPECT().
				StopContainer(gwContainer)
			// Should take down the network
			api.
				EXPECT().
				NewNetwork(netName).
				Return(&container.Network{Name: netName})
			api.
				EXPECT().
				NetworkExists(net).
				Return(true, nil)
			api.
				EXPECT().
				RemoveNetwork(net)
			proj.Down()
		})
	}
}

func TestUpDoesntDuplicateResources(t *testing.T) {
	// For this test we will set the mocked calls to the container API
	// to return resources. This should result in a short circuit,
	// and no attempt should be made to create the resources again.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	api := NewMockContainerAPI(ctrl)

	proj := makeProject("one-service", &project.Service{Name: "one", Path: "./one"})
	svc, err := proj.GetService("./one")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	proj.ConfigureContainerAPI(api)
	proj.ConfigureLogger(logging.NewTestLogger())

	netName := fmt.Sprintf("foldnet-%s", proj.Name)
	net := &container.Network{Name: netName}
	// Should reuse network
	api.
		EXPECT().
		NewNetwork(netName).
		Return(net)
	api.
		EXPECT().
		NetworkExists(net).
		Return(true, nil)
	// Should reuse gateway
	gwSvc := proj.NewService("foldgw")
	gwContainerName := fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)
	gwContainer := &container.Container{ID: fmt.Sprintf("%d", 0), Name: gwContainerName}
	api.
		EXPECT().
		GetContainer(fmt.Sprintf("%s.%s", gwSvc.Id(), gwSvc.Name)).
		Return(gwContainer, nil)

	containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
	container := &container.Container{ID: fmt.Sprintf("%d", 0), Name: containerName}
	// Should reuse service container
	api.
		EXPECT().
		GetContainer(containerName).
		Return(container, nil)
	err = proj.Up(context.Background(), &bytes.Buffer{}, svc)
	if err != nil {
		t.Errorf("expected no error but found %v", err)
	}
}

func makeProject(name string, services ...*project.Service) *project.Project {
	p := &project.Project{
		Name: name,
	}
	p.ConfigureLogger(logging.NewTestLogger())
	for _, s := range services {
		svc := p.NewService(s.Name)
		svc.Port = s.Port
		p.Services = append(p.Services, svc)
	}
	return p
}
