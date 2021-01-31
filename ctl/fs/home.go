package fs

import (
	"errors"
	"os"
	"path/filepath"
)

func FoldHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to locate home directory")
	}
	return filepath.Join(home, ".fold"), nil
}
