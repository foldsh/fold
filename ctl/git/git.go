package git

import (
	"errors"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

func UpdateTemplates(out io.Writer, templatesPath string) error {
	if _, err := os.Stat(templatesPath); err != nil {
		if os.IsNotExist(err) {
			// direcotry does not exist
			cloneTemplates(out, templatesPath)
		} else {
			// other error
			return err
		}
	} else {
		err := pullTemplates(templatesPath)
		if err == nil || errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func cloneTemplates(out io.Writer, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      "https://github.com/foldsh/templates",
		Progress: out,
	})
	return err
}

func pullTemplates(path string) error {
	w, err := openWorkingTree(path)
	if err != nil {
		return err
	}
	return w.Pull(&git.PullOptions{RemoteName: "origin", Depth: 1})
}

func openWorkingTree(path string) (*git.Worktree, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.New("failed to open templates repository")
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, errors.New("failed to inspect working tree")
	}
	return w, nil
}
