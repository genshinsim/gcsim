package artifact

import (
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/artifact"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	Key   string `yaml:"key,omitempty"`
	SetID int64  `yaml:"set_id,omitempty"`

	// extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Generator struct {
	GeneratorConfig
	src       *artifact.DataSource
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

	src, err := artifact.NewDataSource(g.Excels)
	if err != nil {
		return nil, err
	}
	g.src = src

	a, err := ParseArtifactConfig(g.Root)
	if err != nil {
		return nil, err
	}
	g.artifacts = a

	setIDCheck := make(map[int64]bool)

	for _, v := range a {
		if v.SetID == 0 {
			fmt.Printf("[SKIP] invalid set with set id 0: %v\n", v.RelativePath)
			continue
		}
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		if _, ok := setIDCheck[v.SetID]; ok {
			return nil, fmt.Errorf("duplicated set id %v found; second instance at %v", v.SetID, v.RelativePath)
		}
		setIDCheck[v.SetID] = true
		ad, err := g.src.GetSetData(v.SetID)
		if err != nil {
			return nil, err
		}
		ad.Key = v.Key
		g.data[v.Key] = ad
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
