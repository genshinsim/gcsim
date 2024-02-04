package metadata

import (
	"math"
	"slices"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "metadata",
		New:  NewAgg,
	})
}

type Runs []run

type buffer struct {
	runs Runs
}

type run struct {
	seed uint64
	dps  float64
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		runs: make(Runs, 0, cfg.Settings.Iterations),
	}
	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	b.runs = append(b.runs, run{seed: result.Seed, dps: result.DPS})
}

func (b buffer) Flush(result *model.SimulationStatistics) {
	iterations := len(b.runs)
	result.Iterations = uint32(iterations)

	slices.SortStableFunc(b.runs, func(k, j run) int {
		if math.Abs(k.dps-j.dps) < agg.FloatEqDelta {
			return 0
		}
		if k.dps < j.dps {
			return -1
		}
		return 1
	})
	c1, c2 := agg.GetPercentileIndexes(b.runs)

	result.MinSeed = strconv.FormatUint(b.runs[0].seed, 10)
	result.MaxSeed = strconv.FormatUint(b.runs[iterations-1].seed, 10)
	result.P25Seed = strconv.FormatUint(agg.Median(b.runs[:c1]).seed, 10)
	result.P50Seed = strconv.FormatUint(agg.Median(b.runs).seed, 10)
	result.P75Seed = strconv.FormatUint(agg.Median(b.runs[c2:]).seed, 10)
}
