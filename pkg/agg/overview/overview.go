package overview

import (
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/agg/util"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	dps         *util.FloatBuffer
	totalDamage *util.FloatBuffer
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		dps:         util.NewFloatBuffer(cfg.Settings.Iterations),
		totalDamage: util.NewFloatBuffer(cfg.Settings.Iterations),
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result, i int) {
	b.dps.Add(result.DPS, i)
	b.totalDamage.Add(result.TotalDamage, i)
}

func (b *buffer) Flush(result *agg.Result) {
	result.DPS = b.dps.Flush()
	result.TotalDamage = b.totalDamage.Flush()
}
