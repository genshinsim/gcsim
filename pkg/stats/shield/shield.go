package shield

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/stats"
)

var elements = [...]attributes.Element{
	attributes.Anemo,
	attributes.Cryo,
	attributes.Electro,
	attributes.Geo,
	attributes.Hydro,
	attributes.Pyro,
	attributes.Dendro,
	attributes.Physical,
}

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	shields []stats.ShieldInterval
}

func NewStat(core *core.Core) (stats.StatsCollector, error) {
	out := buffer{}

	// TODO: Make cases of overwrite easy to detect in the data
	// want to keep all shield instances for maximal shield absorption info
	core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		shield := args[0].(shield.Shield)

		interval := stats.ShieldInterval{
			Start: core.F,
			End:   shield.Expiry(),
			Name:  shield.Desc(),
		}

		// normalized absoprtion across all damage element types
		bonus := core.Player.Shields.ShieldBonus()
		for _, e := range elements {
			interval.Absorption += shield.ShieldStrength(e, bonus)
		}
		interval.Absorption /= float64(len(elements))

		out.shields = append(out.shields, interval)
		return false
	}, "stats-shield-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	result.Shields = b.shields
}
