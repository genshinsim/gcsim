package pipeline

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Walk searches from the provided root folder and provides a list of paths
// of all sub folders
func Walk(root string) ([]string, error) {
	var c []string

	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unexpected error walking: %w", err)
			}

			// we're only interested in finding directories
			switch {
			case !info.IsDir():
				return nil
			case strings.HasSuffix(path, root):
				return nil
			}

			c = append(c, path)
			return nil
		})

	if err != nil {
		return nil, err
	}

	return c, nil
}

// WalkConfigYml searches from the provided root folder (ignoring root), and provide
// a slice of paths to "config.yml" or "config.yaml"
func WalkConfigYml(root string) ([]string, error) {
	c, err := Walk(root)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, path := range c {
		p := path + "/config.yml"
		if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
			continue
		}
		res = append(res, p)
	}
	return res, err
}
