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
	Name    string
	Path    string
	Mounts  []string
	Port    int
	project *Project
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
	h.Write([]byte(s.project.Name))
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
	s.project.logger.Debugf("%v %v", s, img, net)
	con := s.project.api.NewContainer(s.containerName(), *img)
	con.NetworkAlias = s.Name
	if s.Port != 0 {
		con.Ports = []int{s.Port}
	}
	serviceHome, err := s.AbsPath()
	if err != nil {
		return err
	}
	var mounts []container.Mount
	for _, m := range s.Mounts {
		src := filepath.Join(serviceHome, m)
		dst := filepath.Join(img.WorkDir, m)
		mounts = append(mounts, container.Mount{Src: src, Dst: dst})
	}
	con.Mounts = mounts
	con.Environment = map[string]string{"FOLD_SERVICE_NAME": s.Name}

	err = s.project.api.RunContainer(net, con)
	if err != nil {
		return err
	}
	s.project.logger.Infof("Service %s is up in container %s", s.Name, con.Name)
	return nil
}

func (s *Service) Stop() error {
	container, err := s.project.api.GetContainer(s.containerName())
	if err != nil {
		return err
	}
	if container == nil {
		// There is no container for this service, no need do anything.
		return nil
	}
	s.project.logger.Infof("Stopping container %s", container.Name)
	err = s.project.api.StopContainer(container)
	if err != nil {
		s.project.logger.Debugf("Failed to stop container %s: %v", container.Name, err)
		return err
	}
	return nil
}

func (s *Service) Build(ctx context.Context, out io.Writer) (*container.Image, error) {
	img, err := s.img()
	if err != nil {
		return nil, err
	}
	s.project.logger.Debugf("Preparing to build service %s with tag %s", s.Name, img.Name)
	err = s.project.api.BuildImage(img)
	if err != nil {
		s.project.logger.Debugf("Failed bo build image %v", err)
		return nil, err
	}
	return img, nil
}

func (s *Service) img() (*container.Image, error) {
	path, err := s.AbsPath()
	if err != nil {
		return nil, err
	}
	tag := fmt.Sprintf("foldlocal/%s/%s:latest", s.Id(), s.Name)
	img := &container.Image{
		Src:  path,
		Name: tag,
	}
	return img, nil
}

func (s *Service) containerName() string {
	return fmt.Sprintf("%s.%s", s.Id(), s.Name)
}
