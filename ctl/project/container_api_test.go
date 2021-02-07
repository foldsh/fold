package project_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/foldsh/fold/ctl/container"
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
		api := NewMockContainerAPI(ctrl)

		proj := tc.project
		proj.ConfigureContainerAPI(api)
		proj.ConfigureLogger(logging.NewTestLogger())
		t.Run(tc.project.Name, func(t *testing.T) {
			out := new(bytes.Buffer)
			netName := fmt.Sprintf("foldnet-%s", proj.Name)
			net := &container.Network{Name: netName}
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
			for i, svc := range proj.Services {
				// TODO mock out image builder too
				containerName := fmt.Sprintf("%s.%s", svc.Id(), svc.Name)
				container := &container.Container{ID: fmt.Sprintf("%d", i), Name: containerName}
				api.
					EXPECT().
					RunContainer(container)
				api.
					EXPECT().
					AddToNetwork(net, container)
			}
			proj.Up(context.Background(), out, proj.Services...)

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
				api.
					EXPECT().
					RemoveContainer(container)
			}
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

func makeProject(name string, services ...*project.Service) *project.Project {
	p := &project.Project{
		Name: name,
	}
	for _, svc := range services {
		svc.Project = p
		p.Services = append(p.Services, svc)
	}
	return p
}