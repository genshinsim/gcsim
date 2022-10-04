package metadata

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/agg/util"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	minRun   run
	maxRun   run
	duration *util.FloatBuffer
}

type run struct {
	seed uint64
	dps  float64
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		duration: util.NewFloatBuffer(cfg.Settings.Iterations),
		minRun: run{
			dps: math.MaxFloat64,
		},
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result, i int) {
	duration := float64(result.Duration) / 60
	b.duration.Add(duration, i)

	if result.DPS < b.minRun.dps {
		b.minRun = run{seed: result.Seed, dps: result.DPS}
	}
	if result.DPS > b.maxRun.dps {
		b.maxRun = run{seed: result.Seed, dps: result.DPS}
	}
}

func (b *buffer) Flush(result *agg.Result) {
	result.Duration = b.duration.Flush()
	result.MinSeed = b.minRun.seed
	result.MaxSeed = b.maxRun.seed
}
