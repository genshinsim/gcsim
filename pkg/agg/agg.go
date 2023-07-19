package agg

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type Aggregator interface {
	Add(result stats.Result)
	// TODO: Merge(other Aggregator) Aggregator for multi-threaded aggregations (optional optimization)
	Flush(result *model.SimulationStatistics)
}

type NewAggFunc func(cfg *ast.ActionList) (Aggregator, error)

var (
	mu          sync.Mutex
	aggregators []NewAggFunc
)

func Register(f NewAggFunc) {
	mu.Lock()
	defer mu.Unlock()
	aggregators = append(aggregators, f)
}

func Aggregators() []NewAggFunc {
	return aggregators
}
