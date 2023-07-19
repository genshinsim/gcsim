package artifact

import (
	"fmt"
	"os"

	"github.com/genshinsim/gcsim/pipeline/pkg/pipeline"
	"gopkg.in/yaml.v2"
)

func ParseArtifactConfig(root string) ([]Config, error) {
	c, err := pipeline.WalkConfigYml(root)
	if err != nil {
		return nil, err
	}
	return read(c)
}

func read(c []string) ([]Config, error) {
	var res []Config
	for _, p := range c {
		cfg, err := readArtifact(p)
		if err != nil {
			return nil, err
		}
		res = append(res, cfg)
	}
	return res, nil
}

func readArtifact(path string) (Config, error) {
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
