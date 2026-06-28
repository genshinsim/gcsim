package enemy

import (
	"cmp"
	"slices"

	"github.com/genshinsim/gcsim/pipeline/pkg/data/enemy"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Generator struct {
	GeneratorConfig
	src  *enemy.DataSource
	data []*model.MonsterData
}

type GeneratorConfig struct {
	Root   string
	Excels string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	g := &Generator{
		GeneratorConfig: cfg,
		data:            []*model.MonsterData{},
	}

	src, err := enemy.NewDataSource(g.Excels)
	if err != nil {
		return nil, err
	}
	g.src = src
	g.data = src.GetMonsters()
	slices.SortFunc(g.data, func(a, b *model.MonsterData) int { return cmp.Compare(a.Key, b.Key) })

	return g, nil
}

func (g *Generator) Data() []*model.MonsterData {
	return g.data
}
