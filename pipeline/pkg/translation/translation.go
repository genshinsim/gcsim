package translation

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/textmap"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Generator struct {
	GeneratorConfig
}

type GeneratorConfig struct {
	Characters []*model.AvatarData
	Weapons    []*model.WeaponData
	Artifacts  []*model.ArtifactData
	Languages  map[string]string // map of languages and their corresponding textmap
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	return &Generator{
		GeneratorConfig: cfg,
	}, nil
}

func (g *Generator) DumpUIJSON(path string) error {
	// delete existing
	err := g.writeTranslationJSON(path + "/names.generated.json")
	if err != nil {
		return err
	}
	return nil
}

type outData struct {
	CharacterNames map[string]string
	WeaponNames    map[string]string
	ArtifactNames  map[string]string
}

func (g *Generator) writeTranslationJSON(path string) error {
	// sort keys first
	var keys []string
	for k := range g.Languages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make(map[string]outData)

	for _, k := range keys {
		data := outData{
			CharacterNames: make(map[string]string),
			WeaponNames:    make(map[string]string),
			ArtifactNames:  make(map[string]string),
		}
		// load generator for this language
		tp := g.Languages[k]
		src, err := textmap.NewTextMapSource(tp)
		if err != nil {
			return fmt.Errorf("error creating text map src for %v: %w", k, err)
		}
		// go through all char/weap/art and get names
		for _, v := range g.Characters {
			s, err := src.Get(v.NameTextHashMap)
			if err != nil {
				fmt.Printf("error getting string for char %v id %v\n", v.Key, v.NameTextHashMap)
				continue
			}
			data.CharacterNames[v.Key] = s
		}
		for _, v := range g.Weapons {
			s, err := src.Get(v.NameTextHashMap)
			if err != nil {
				fmt.Printf("error getting string for weapon %v id %v\n", v.Key, v.NameTextHashMap)
				continue
			}
			data.WeaponNames[v.Key] = s
		}
		for _, v := range g.Artifacts {
			s, err := src.Get(v.TextMapId)
			if err != nil {
				fmt.Printf("error getting string for set %v id %v\n", v.Key, v.TextMapId)
				continue
			}
			data.ArtifactNames[v.Key] = s
		}

		out[k] = data
	}

	data, err := json.MarshalIndent(out, "", "   ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
