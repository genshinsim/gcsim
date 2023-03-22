package character

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseCharConfig(root string) ([]Config, error) {
	c, err := walk(root)
	if err != nil {
		return nil, err
	}
	return read(c)
}

func read(c []string) ([]Config, error) {
	var res []Config
	for _, path := range c {
		p := path + "/config.yml"
		if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
			continue
		}
		cfg, err := readChar(p)
		if err != nil {
			return nil, err
		}
		res = append(res, cfg)
	}
	return res, nil
}

func readChar(path string) (Config, error) {
	c := Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("error reading %v: %v", path, err)
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf("error parsing config %v: %v", path, err)
	}

	return c, nil
}

func walk(root string) ([]string, error) {
	var c []string

	err := filepath.Walk(root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unexpected error walking: %v", err)
			}

			//we're only interested in finding directories
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
