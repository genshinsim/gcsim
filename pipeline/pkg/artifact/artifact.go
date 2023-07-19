package artifact

import (
	"fmt"
	"log"

	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	Key       string `yaml:"key,omitempty"`
	TextMapId string `yaml:"text_map_id,omitempty"`

	//extra fields to be populate but not read from yaml
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

	textIdCheck := make(map[string]bool)

	for _, v := range a {
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		if _, ok := textIdCheck[v.TextMapId]; ok {
			return nil, fmt.Errorf("duplicated text map id %v found; second instance at %v", v.TextMapId, v.RelativePath)
		}
		g.data[v.Key] = &model.ArtifactData{
			TextMapId: v.TextMapId,
		}
		log.Printf("%v loaded ok\n", v.Key)
	}

	return g, nil
}
