package git

import (
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func UpdateTemplates(out io.Writer, templatesPath, version string) error {
	_, err := os.Stat(templatesPath)
	if !os.IsNotExist(err) {
		// If the directory exists then we need to remove it first.
		if err := os.RemoveAll(templatesPath); err != nil {
			return err
		}
	}
	return cloneTemplates(out, templatesPath, version)
}

func cloneTemplates(out io.Writer, path, version string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           "https://github.com/foldsh/templates",
		Progress:      out,
		ReferenceName: plumbing.NewTagReferenceName(version),
		Depth:         1,
		SingleBranch:  true,
	})
	return err
}
