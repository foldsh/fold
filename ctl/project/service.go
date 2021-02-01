package project

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"regexp"
)

var ServiceNameRegex = `^[a-z][a-z-_]+$`

type Service struct {
	Name   string   `mapstructure:"name"`
	Path   string   `mapstructure:"path"`
	Mounts []string `mapstructure:"mounts"`

	project *Project
}

func (s *Service) Validate() bool {
	matched, _ := regexp.MatchString(ServiceNameRegex, s.Name)
	return matched
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
	return filepath.Abs(s.Path)
}
