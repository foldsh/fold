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

func FoldTemplates(foldHome string) string {
	return filepath.Join(foldHome, "templates")
}

func FoldBin(foldHome string) string {
	return filepath.Join(foldHome, "bin")
}
