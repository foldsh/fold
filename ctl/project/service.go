package project

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"

	"github.com/foldsh/fold/ctl/container"
)

var (
	ServiceNameRegex = `^[a-z][a-z-_]+$`

	ServicePathInvalid = errors.New("cannot build an absolute path to the service")
)

type Service struct {
	Name    string  `mapstructure:"name"`
	Path    string  `mapstructure:"path"`
	Mounts  []Mount `mapstructure:"mounts"`
	Port    int
	Project *Project
}

type Mount struct {
	Src string `mapstructure:"src"`
	Dst string `mapstructure:"dst"`
}

type NotAService struct {
	Service string
}

func (nas NotAService) Error() string {
	return fmt.Sprintf("%s is not a registered service", nas.Service)
}

func (s *Service) Validate() error {
	matched, _ := regexp.MatchString(ServiceNameRegex, s.Name)
	if matched {
		return nil
	}
	return fmt.Errorf(
		"%s is not a valid service name, it must match the regex %s",
		s.Name,
		ServiceNameRegex,
	)
}

func (s *Service) Id() string {
	h := sha1.New()
	h.Write([]byte(s.Project.Name))
	h.Write([]byte(s.Name))
	hashString := fmt.Sprintf("%x", h.Sum(nil))
	// Just the first 7 characters will be fine
	return hashString[0:7]
}

func (s *Service) AbsPath() (string, error) {
	p, err := filepath.Abs(s.Path)
	if err != nil {
		return "", ServicePathInvalid
	}
	return p, nil
}

func (s *Service) Start(img *container.Image, net *container.Network) error {
	s.Project.logger.Debugf("%v %v", s, img, net)
	con := s.Project.api.NewContainer(s.containerName(), *img)
	con.NetworkAlias = s.Name
	con.Ports = []int{s.Port}
	var mounts []container.Mount
	for _, m := range s.Mounts {
		mounts = append(mounts, container.Mount{Src: m.Src, Dst: m.Dst})
	}
	con.Mounts = mounts

	err := s.Project.api.RunContainer(net, con)
	if err != nil {
		return err
	}
	s.Project.logger.Infof("Service %s is up in container %s", s.Name, con.Name)
	return nil
}

func (s *Service) Stop() error {
	container, err := s.Project.api.GetContainer(s.containerName())
	if err != nil {
		return err
	}
	if container == nil {
		// There is no container for this service, no need do anything.
		return nil
	}
	s.Project.logger.Infof("Stopping container %s", container.Name)
	err = s.Project.api.StopContainer(container)
	if err != nil {
		s.Project.logger.Debugf("Failed to stop container %s: %v", container.Name, err)
		return err
	}
	return nil
}

func (s *Service) Build(ctx context.Context, out io.Writer) (*container.Image, error) {
	img, err := s.img()
	if err != nil {
		return nil, err
	}
	s.Project.logger.Debugf("Preparing to build service %s with tag %s", s.Name, img.Name)
	err = s.Project.api.BuildImage(img)
	if err != nil {
		s.Project.logger.Debugf("Failed bo build image %v", err)
		return nil, err
	}
	return img, nil
}

func (s *Service) img() (*container.Image, error) {
	path, err := s.AbsPath()
	if err != nil {
		return nil, err
	}
	tag := fmt.Sprintf("foldlocal/%s/%s", s.Id(), s.Name)
	img := &container.Image{
		Src:  path,
		Name: tag,
	}
	return img, nil
}

func (s *Service) containerName() string {
	return fmt.Sprintf("%s.%s", s.Id(), s.Name)
}
