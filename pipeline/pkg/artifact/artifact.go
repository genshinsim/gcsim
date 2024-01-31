package artifact

import (
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	Key       string `yaml:"key,omitempty"`
	TextMapId int64  `yaml:"text_map_id,omitempty"`

	// extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Generator struct {
	GeneratorConfig
	artifacts []Config
	data      map[string]*model.ArtifactData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            make(map[string]*model.ArtifactData),
	}

	a, err := ParseArtifactConfig(g.Root)
	if err != nil {
		return nil, err
	}
	g.artifacts = a

	textIDCheck := make(map[int64]bool)

	for _, v := range a {
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		if _, ok := textIDCheck[v.TextMapId]; ok {
			return nil, fmt.Errorf("duplicated text map id %v found; second instance at %v", v.TextMapId, v.RelativePath)
		}
		textIDCheck[v.TextMapId] = true
		g.data[v.Key] = &model.ArtifactData{
			TextMapId: v.TextMapId,
			Key:       v.Key,
		}
		log.Printf("%v loaded ok\n", v.Key)
	}

	return g, nil
}

func (g *Generator) Data() []*model.ArtifactData {
	keys := make([]string, 0, len(g.data))
	for k := range g.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var res []*model.ArtifactData
	for _, k := range keys {
		res = append(res, g.data[k])
	}
	return res
}
