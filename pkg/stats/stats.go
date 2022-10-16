package stats

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core"
)

type StatsCollector interface {
	Flush(core *core.Core, result *Result)
}

type NewStatsFunc func(core *core.Core) (StatsCollector, error)

var (
	mu         sync.Mutex
	collectors []NewStatsFunc
)

func Register(f NewStatsFunc) {
	mu.Lock()
	defer mu.Unlock()
	collectors = append(collectors, f)
}

func Collectors() []NewStatsFunc {
	return collectors
}
