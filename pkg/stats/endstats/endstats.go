// package endstats collects a snapshot of relevant information right when a simulation
// ends, such as ending energy etc..
package endstats

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(stats.Config{
		Name: "ending-stats",
		New:  NewStat,
	})
}

type buffer struct {
	endingEnergy []float64 // ending energy for each character
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		endingEnergy: make([]float64, len(core.Player.Chars())),
	}

	core.Events.Subscribe(event.OnSimEndedSuccessfully, func(args ...interface{}) bool {
		for i, c := range core.Player.Chars() {
			out.endingEnergy[i] = c.Energy
		}
		return true
	}, "stats-ending-energy")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.endingEnergy); c++ {
		result.EndStats[c].EndingEnergy = b.endingEnergy[c]
	}
}
