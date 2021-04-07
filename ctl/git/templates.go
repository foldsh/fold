package git

import (
	"fmt"
	"io"
	"os"

	"github.com/foldsh/fold/version"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func UpdateTemplates(out io.Writer, templatesPath string) error {
	fmt.Printf("%s", templatesPath)
	_, err := os.Stat(templatesPath)
	if !os.IsNotExist(err) {
		fmt.Printf("removing %s", templatesPath)
		// If the directory exists then we need to remove it first.
		if err := os.RemoveAll(templatesPath); err != nil {
			return err
		}
	}
	return cloneTemplates(out, templatesPath)
}

func cloneTemplates(out io.Writer, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           "https://github.com/foldsh/templates",
		Progress:      out,
		ReferenceName: plumbing.NewTagReferenceName(version.FoldVersion.String()),
		Depth:         1,
		SingleBranch:  true,
	})
	return err
}
