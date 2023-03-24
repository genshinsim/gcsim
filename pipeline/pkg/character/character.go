package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/avatar"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	PackageName   string   `yaml:"package_name,omitempty"`
	GenshinID     int64    `yaml:"genshin_id,omitempty"`
	TravelerSubID int64    `yaml:"traveler_sub_id,omitempty"`
	Key           string   `yaml:"key,omitempty"`
	Shortcuts     []string `yaml:"shortcuts,omitempty"`

	//extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Generator struct {
	GeneratorConfig
	src   *avatar.DataSource
	chars []Config
	data  map[int64]*model.AvatarData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            make(map[int64]*model.AvatarData),
	}

	src, err := avatar.NewDataSource(g.Excels)
	if err != nil {
		return nil, err
	}
	g.src = src

	chars, err := ParseCharConfig(g.Root)
	if err != nil {
		return nil, err
	}
	g.chars = chars

	keyCheck := make(map[string]bool)

	for _, v := range chars {
		if _, ok := g.data[v.GenshinID]; ok {
			continue
		}
		if _, ok := keyCheck[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		char, err := src.GetAvatarData(v.GenshinID)
		if err != nil {
			return nil, err
		}
		char.Key = v.Key
		g.data[v.GenshinID] = char
	}

	return g, nil
}
