package character

import (
	"fmt"
	"os"
	"strings"

	"github.com/genshinsim/gcsim/pipeline/pkg/pipeline"
	"gopkg.in/yaml.v3"
)

func ParseCharConfig(root string) ([]Config, error) {
	c, err := pipeline.WalkConfigYml(root)
	if err != nil {
		return nil, err
	}
	return read(c)
}

func read(c []string) ([]Config, error) {
	var res []Config
	for _, p := range c {
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
		return c, fmt.Errorf("error reading %v: %w", path, err)
	}
	c.RelativePath = strings.TrimSuffix(path, "/config.yml")
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf("error parsing config %v: %w", path, err)
	}

	return c, nil
}
