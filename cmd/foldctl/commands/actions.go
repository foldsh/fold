// This file contains a number of 'actions' which are common to many commands.
// It avoids duplicating things like common bits of resource acquisition across
// commands, and ensures that error handling is consistent for all of them.
// This also serves to make the commands much less verbose. When there are errors
// the only course of action is to stop the command with an appropriate help message.
// We can therefore capture a lot of the errors that library code throws in here.
package commands

import (
	"errors"
	"fmt"

	"github.com/foldsh/fold/ctl/container"
	"github.com/foldsh/fold/ctl/project"
)

func loadProject() *project.Project {
	p, err := project.Load(logger)
	if err != nil {
		if errors.Is(err, project.NotAFoldProject) {
			exitWithMessage(
				"This is not a fold project root.",
				"Please either initialise a project or cd to a project root.",
			)
		} else if errors.Is(err, project.InvalidConfig) {
			exitWithMessage(
				"Fold config is invalid.",
				"Please check that the yaml is valid and that you have spelled all the keys correctly.",
			)
		} else {
			exitWithMessage("Failed to load fold config. Please ensure you're in a fold project root.")
		}
	}
	return p
}

func saveProjectConfig(p *project.Project) {
	err := p.SaveConfig()
	exitIfError(
		err,
		"Failed to save fold config.",
		"Please check you have permission to write files in this directory.",
	)
}

func getService(p *project.Project, path string) *project.Service {
	service, err := p.GetService(path)
	exitIfError(
		err,
		fmt.Sprintf("The path %s is not a registered service.", path),
		"Please check the path you typed or, if this is a mistake, make sure that the service",
		"is registered in your fold.yaml file.",
	)
	return service
}

func getDockerClient() container.DockerClient {
	dc, err := container.NewDockerClient(logger)
	exitIfError(err, "Failed to create DockerClient. Ensure that the docker daemon is running.")
	return dc
}

func getContainerRuntime(outPrefix string) container.ContainerRuntime {
	dc := getDockerClient()
	rt := container.NewRuntime(
		commandCtx,
		logger,
		newStreamLinePrefixer(serr, blue(outPrefix)),
		dc,
	)
	return rt
}

func buildService(service *project.Service) *container.ImageSpec {
	absPath, err := service.AbsPath()
	exitIfError(err, servicePathInvalid)
	logger.Debugf("absolute path to service inferred as %s", absPath)
	tag := fmt.Sprintf("foldlocal/%s/%s", service.Id(), service.Name)
	print("Preparing to build service %s with tag %s", service.Name, tag)
	dc, err := container.NewDockerClient(logger)
	exitIfError(err, "Failed to create DockerClient")
	img := &container.ImageSpec{
		Src:          absPath,
		Name:         tag,
		Logger:       logger,
		Out:          newStreamLinePrefixer(serr, blue("docker: ")),
		DockerClient: dc,
	}
	err = img.Build(commandCtx)
	exitIfError(
		err,
		"Failed to build the service.",
		"Check the build logs above for more information on why this happened.",
	)
	return img
}

func getFoldLocalNet(rt container.ContainerRuntime) *container.Network {
	return rt.NewNetwork("foldlocalnet")
}

func getOrCreateFoldNet(rt container.ContainerRuntime) *container.Network {
	net := getFoldLocalNet(rt)
	err := net.CreateIfNotExists()
	exitIfError(err, "Failed to start up the network")
	return net
}

func removeFoldLocalNet(rt container.ContainerRuntime) {
	net := getFoldLocalNet(rt)
	if exists, err := net.Exists(); err != nil {
		exitIfError(err, cantReachDocker)
	} else if exists {
		err = net.Remove()
		if err != nil {
			exitIfError(err, "Failed to stop the network")
		}
	} else {
		print("Network is not up.")
	}
}

func getOrCreateContainer(
	rt container.ContainerRuntime, net *container.Network, service *project.Service,
) *container.Container {
	c, err := rt.GetContainer(service.Name)
	exitIfError(err, cantReachDocker)
	if c != nil {
		return c
	}

	img := buildService(service)
	c = rt.NewContainer(service.Name, img.Name)
	err = c.Start()
	exitIfError(err, "Failed to start container")
	err = c.JoinNetwork(net)
	exitIfError(err, "Failed to join network")
	return c
}

func getAllContainers(rt container.ContainerRuntime) []*container.Container {
	containers, err := rt.AllContainers()
	exitIfError(err, cantReachDocker)
	return containers
}

func stopAndRemoveContainer(c *container.Container) {
	err := c.Stop()
	exitIfError(err, "Failed to stop container %s", c.Name)
	err = c.Remove()
	exitIfError(err, "Failed to remove container %s", c.Name)
}
