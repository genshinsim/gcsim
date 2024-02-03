package character

import (
	"fmt"
	"log"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/avatar"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Config struct {
	PackageName    string   `yaml:"package_name,omitempty"`
	CharStructName string   `yaml:"char_struct_name,omitempty"`
	GenshinID      int32    `yaml:"genshin_id,omitempty"`
	SubID          int32    `yaml:"sub_id,omitempty"`
	Key            string   `yaml:"key,omitempty"`
	Shortcuts      []string `yaml:"shortcuts,omitempty"`

	// skill data generation
	GenerateSkillData bool                      `yaml:"generate_skill_data"`
	SkillDataMapping  map[string]map[int]string `yaml:"skill_data_mapping"`

	// extra fields to be populate but not read from yaml
	RelativePath string `yaml:"-"`
}

type Generator struct {
	GeneratorConfig
	src   *avatar.DataSource
	chars []Config
	data  map[string]*model.AvatarData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            make(map[string]*model.AvatarData),
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

	for _, v := range chars {
		// validate char config
		err := validateConfig(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %v config: %w", v.Key, err)
		}
		if _, ok := g.data[v.Key]; ok {
			return nil, fmt.Errorf("duplicated key %v found; second instance at %v", v.Key, v.RelativePath)
		}
		char, err := src.GetAvatarData(v.GenshinID, v.SubID)
		if err != nil {
			log.Printf("Error loading %v data: %v; skipping\n", v.Key, err)
			continue
		}
		char.Key = v.Key
		g.data[v.Key] = char
		log.Printf("%v loaded ok\n", v.Key)
	}

	return g, nil
}

func validateConfig(cfg Config) error {
	// make sure no duplicated variable names
	names := make(map[string]bool)

	for _, v := range cfg.SkillDataMapping {
		for _, varname := range v {
			if _, ok := names[varname]; ok {
				return fmt.Errorf("duplicate var name %v", varname)
			}
			names[varname] = true
		}
	}

	return nil
}

func (g *Generator) Data() []*model.AvatarData {
	keys := make([]string, 0, len(g.data))
	for k := range g.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var res []*model.AvatarData
	for _, k := range keys {
		res = append(res, g.data[k])
	}
	return res
}
