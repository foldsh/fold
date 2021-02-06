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

	NotAService        = errors.New("not a valid service")
	ServicePathInvalid = errors.New("cannot build an absolute path to the service")
)

type Service struct {
	Name    string   `mapstructure:"name"`
	Path    string   `mapstructure:"path"`
	Mounts  []string `mapstructure:"mounts"`
	Project *Project
}

func (s *Service) Validate() bool {
	matched, _ := regexp.MatchString(ServiceNameRegex, s.Name)
	return matched
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

func (s *Service) Start(ctx context.Context, out io.Writer, net *container.Network) error {
	img, err := s.Build(ctx, out)
	if err != nil {
		return err
	}
	container := s.Project.backend.NewContainer(s.containerName(), *img)
	err = s.Project.backend.RunContainer(container)
	if err != nil {
		return err
	}
	err = s.Project.backend.AddToNetwork(net, container)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Stop() error {
	container, err := s.Project.backend.GetContainer(s.containerName())
	if err != nil {
		return err
	}
	if container == nil {
		// There is no container for this service, no need do anything.
		return nil
	}
	err = s.Project.backend.StopContainer(container)
	if err != nil {
		s.Project.logger.Debugf("failed bo stop container %s: %v", container.Name, err)
		return err
	}
	err = s.Project.backend.RemoveContainer(container)
	if err != nil {
		s.Project.logger.Debugf("failed bo remove container %s: %v", container.Name, err)
		return err
	}
	return nil
}

func (s *Service) Build(ctx context.Context, out io.Writer) (*container.Image, error) {
	img, err := s.img()
	if err != nil {
		return nil, err
	}
	s.Project.logger.Debugf("preparing to build service %s with tag %s", s.Name, img.Name)
	ib, err := container.NewImageBuilder(ctx, s.Project.logger, out)
	if err != nil {
		s.Project.logger.Debugf("failed bo construct image builder %v", err)
		return nil, err
	}
	err = ib.Build(*img)
	if err != nil {
		s.Project.logger.Debugf("failed bo build image %v", err)
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
