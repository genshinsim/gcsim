package weapon

import (
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/weapon"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	PackageName string   `yaml:"package_name,omitempty"`
	StructName  string   `yaml:"struct_name,omitempty"`
	GenshinID   int32    `yaml:"genshin_id,omitempty"`
	Key         string   `yaml:"key,omitempty"`
	Shortcuts   []string `yaml:"shortcuts,omitempty"`

	SkipDataFunc bool `yaml:"skip_data_func,omitempty"`

	// extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Generator struct {
	GeneratorConfig
	src   *weapon.DataSource
	weaps []Config
	data  map[string]*model.WeaponData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            make(map[string]*model.WeaponData),
	}

	src, err := weapon.NewDataSource(g.Excels)
	if err != nil {
		return nil, err
	}
	g.src = src

	weaps, err := ParseWeaponConfig(g.Root)
	if err != nil {
		return nil, err
	}
	g.weaps = weaps

	for _, v := range weaps {
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		w, err := src.GetWeaponData(v.GenshinID)
		if err != nil {
			log.Printf("Error loading %v data: %v; skipping\n", v.Key, err)
			continue
		}
		w.Key = v.Key
		g.data[v.Key] = w
	}

	return g, nil
}

func (g *Generator) Data() []*model.WeaponData {
	keys := make([]string, 0, len(g.data))
	for k := range g.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var res []*model.WeaponData
	for _, k := range keys {
		res = append(res, g.data[k])
	}
	return res
}
