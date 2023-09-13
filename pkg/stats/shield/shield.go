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

const normalized = "normalized"

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	shields map[string][]stats.ShieldInterval
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		shields: make(map[string][]stats.ShieldInterval),
	}

	core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		shield := args[0].(shield.Shield)
		name := shield.Desc()
		bonus := core.Player.Shields.ShieldBonus()

		interval := stats.ShieldInterval{
			Start: core.F,
			End:   shield.Expiry(),
			HP:    make(map[string]float64),
		}

		var normalizedHP float64
		for _, e := range elements {
			hp := shield.ShieldStrength(e, bonus)
			interval.HP[e.String()] = hp
			normalizedHP += hp
		}
		interval.HP[normalized] = normalizedHP / float64(len(elements))

		// first instance of this shield type
		if _, ok := out.shields[name]; !ok {
			out.shields[name] = make([]stats.ShieldInterval, 0)
			out.shields[name] = append(out.shields[name], interval)
			return false
		}

		prevIndex := len(out.shields[name]) - 1
		prevInterval := out.shields[name][prevIndex]

		// shield refreshed before previous expired
		if prevInterval.End >= interval.Start {
			// same hp stats, merge intervals
			if same(prevInterval.HP, interval.HP) {
				// TODO: max not necessary here?
				prevInterval.End = max(prevInterval.End, interval.End)
				out.shields[name][prevIndex] = prevInterval
				return false
			}

			// diff, end prev interval early
			prevInterval.End = interval.Start
			out.shields[name][prevIndex] = prevInterval
		}

		out.shields[name] = append(out.shields[name], interval)
		return false
	}, "stats-shield-log")

	// TODO: Should be replaced with targeted events (IE on shield stats changes + char swap)
	core.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
		bonus := core.Player.Shields.ShieldBonus()

		for _, shield := range core.Player.Shields.List() {
			interval := stats.ShieldInterval{
				Start: core.F,
				End:   shield.Expiry(),
				HP:    make(map[string]float64),
			}

			var normalizedHP float64
			for _, e := range elements {
				hp := shield.ShieldStrength(e, bonus)
				interval.HP[e.String()] = hp
				normalizedHP += hp
			}
			interval.HP[normalized] = normalizedHP / float64(len(elements))

			prevIndex := len(out.shields[shield.Desc()]) - 1
			prevInterval := out.shields[shield.Desc()][prevIndex]
			if !same(prevInterval.HP, interval.HP) {
				if prevInterval.Start == interval.Start {
					// special case where shield gets recomputed on first frame
					out.shields[shield.Desc()][prevIndex] = interval
				} else {
					prevInterval.End = interval.Start
					out.shields[shield.Desc()][prevIndex] = prevInterval
					out.shields[shield.Desc()] = append(out.shields[shield.Desc()], interval)
				}
			}
		}
		return false
	}, "stats-shield-tick-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	shields := make([]stats.ShieldStats, 0, len(b.shields))
	for name, sb := range b.shields {
		shield := stats.ShieldStats{
			Name:      name,
			Intervals: sb,
		}
		shields = append(shields, shield)
	}

	result.ShieldResults = stats.ShieldResult{
		Shields:         shields,
		EffectiveShield: computeEffective(b.shields),
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func same(i, j map[string]float64) bool {
	for _, e := range elements {
		if i[e.String()] != j[e.String()] {
			return false
		}
	}
	return true
}
