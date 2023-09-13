package metadata

import (
	"sort"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
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
	result.Iterations = uint32(b.runs.Len())

	sort.Sort(b.runs)
	result.MinSeed = strconv.FormatUint(b.runs[0].seed, 10)
	result.MaxSeed = strconv.FormatUint(b.runs[b.runs.Len()-1].seed, 10)

	l := b.runs.Len()
	var c1 int
	var c2 int
	if l%2 == 0 {
		c1 = l / 2
		c2 = l / 2
	} else {
		c1 = (l - 1) / 2
		c2 = c1 + 1
	}

	result.P25Seed = strconv.FormatUint(b.runs[:c1].median().seed, 10)
	result.P50Seed = strconv.FormatUint(b.runs.median().seed, 10)
	result.P75Seed = strconv.FormatUint(b.runs[c2:].median().seed, 10)
}

func (r Runs) Len() int           { return len(r) }
func (r Runs) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Runs) Less(i, j int) bool { return r[i].dps < r[j].dps }

// assumes already sorted
func (r Runs) median() run {
	l := r.Len()

	if l == 0 {
		return run{seed: 0, dps: -1}
	}
	// if length of array is even, median is between r[l/2] and r[l/2+1]
	// since need a seed that was used, r[l/2] is close enough
	return r[l/2]
}
