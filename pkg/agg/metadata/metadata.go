package metadata

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	minRun run
	maxRun run
	count  int
}

type run struct {
	seed uint64
	dps  float64
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		minRun: run{
			dps: math.MaxFloat64,
		},
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	if result.DPS < b.minRun.dps {
		b.minRun = run{seed: result.Seed, dps: result.DPS}
	}
	if result.DPS > b.maxRun.dps {
		b.maxRun = run{seed: result.Seed, dps: result.DPS}
	}
	b.count += 1
}

func (b buffer) Flush(result *agg.Result) {
	result.MinSeed = b.minRun.seed
	result.MaxSeed = b.maxRun.seed
	result.Iterations = b.count
}
