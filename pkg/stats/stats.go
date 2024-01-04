package stats

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Config struct {
	Name string
	New  NewStatsFunc
}

type Collector interface {
	Flush(core *core.Core, result *Result)
}

type NewStatsFunc func(core *core.Core) (Collector, error)

var (
	mu         sync.Mutex
	collectors = map[string]Config{}
)

func Register(cfg Config) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := collectors[cfg.Name]; ok {
		panic("duplicate stats collector registered: " + cfg.Name)
	}
	collectors[cfg.Name] = cfg
}

func Collectors() map[string]Config {
	return collectors
}
