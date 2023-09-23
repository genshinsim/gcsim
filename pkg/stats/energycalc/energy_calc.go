package energycalc

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	stats.Register(NewStat)
}

type buffer struct {
	charRawParticles    []float64
	charFlatEnergy      []float64
	erNeeded            []float64
	erPerParticleEvent  [][]float64
	rawPerParticleEvent [][]float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		charRawParticles:    make([]float64, len(core.Player.Chars())),
		charFlatEnergy:      make([]float64, len(core.Player.Chars())),
		erNeeded:            make([]float64, len(core.Player.Chars())),
		erPerParticleEvent:  make([][]float64, len(core.Player.Chars())),
		rawPerParticleEvent: make([][]float64, len(core.Player.Chars())),
	}

	for i := range out.erNeeded {
		out.erNeeded[i] = 1.0
	}

	for i := 0; i < len(out.erPerParticleEvent); i++ {
		out.erPerParticleEvent[i] = make([]float64, 0)
	}

	for i := 0; i < len(out.rawPerParticleEvent); i++ {
		out.rawPerParticleEvent[i] = make([]float64, 0)
	}
	burstCount := make([]int, len(core.Player.Chars()))
	core.Events.Subscribe(event.OnBurst, func(_ ...interface{}) bool {
		char := core.Player.ActiveChar()
		ind := char.Index

		if burstCount[ind] > 0 {
			erNeeded := (char.EnergyMax - out.charFlatEnergy[ind]) / out.charRawParticles[ind]
			if erNeeded > out.erNeeded[ind] {
				out.erNeeded[ind] = erNeeded
			}
		}
		out.charFlatEnergy[ind] = 0
		out.charRawParticles[ind] = 0
		burstCount[ind] += 1
		return false
	}, "stats-burst-energy-calc-log")

	core.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		character := args[0].(*character.CharWrapper)
		ind := character.Index
		preEnergy := args[1].(float64)
		amount := args[2].(float64)
		isParticle := args[4].(bool)

		er := character.Stat(attributes.ER)

		if isParticle {
			raw := amount / (1.0 + er)
			out.charRawParticles[ind] += raw
			out.erPerParticleEvent[ind] = append(out.erPerParticleEvent[ind], 1+er)
			out.rawPerParticleEvent[ind] = append(out.rawPerParticleEvent[ind], raw)
		} else {
			if amount < 0 && preEnergy < -amount {
				out.charFlatEnergy[ind] -= preEnergy
			} else {
				out.charFlatEnergy[ind] += amount
			}
		}

		return false
	}, "stats-energy-calc-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(core.Player.Chars()); c++ {
		result.Characters[c].ErNeeded = b.erNeeded[c]

		wERsum := 0.0
		wsum := 0.0
		for i, raw := range b.rawPerParticleEvent[c] {
			wERsum += b.erPerParticleEvent[c][i] * raw
			wsum += raw
		}
		if wsum == 0 {
			result.Characters[c].WeightedER = core.Player.ActiveChar().Stat(attributes.ER + 1)
		} else {
			result.Characters[c].WeightedER = wERsum / wsum
		}
	}
}
