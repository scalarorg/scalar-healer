package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func RootPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err // Failed to get current working directory
	}

	// Find project root by looking for go.mod file
	root := cwd
	for {
		_, err = os.Stat(filepath.Join(root, "go.mod"))
		if err == nil {
			break // Found project root
		}
		parent := filepath.Dir(root)
		if parent == root {
			return "", fmt.Errorf("project root not found") // Reached root without finding go.mod
		}
		root = parent
	}
	return root, nil
}
