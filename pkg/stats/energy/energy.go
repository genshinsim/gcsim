package energy

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
	events                    [][]stats.EnergyEvent
	energySpent               []float64
	flatEnergyPerBurst        [][]float64
	charRawParticles          []float64
	charFlatEnergy            []float64
	rawParticleEnergyPerBurst [][]float64
	erPerParticleEvent        [][]float64
	rawPerParticleEvent       [][]float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		events:                    make([][]stats.EnergyEvent, len(core.Player.Chars())),
		energySpent:               make([]float64, len(core.Player.Chars())),
		charRawParticles:          make([]float64, len(core.Player.Chars())),
		charFlatEnergy:            make([]float64, len(core.Player.Chars())),
		flatEnergyPerBurst:        make([][]float64, len(core.Player.Chars())),
		rawParticleEnergyPerBurst: make([][]float64, len(core.Player.Chars())),
		erPerParticleEvent:        make([][]float64, len(core.Player.Chars())),
		rawPerParticleEvent:       make([][]float64, len(core.Player.Chars())),
	}

	for i := 0; i < len(out.flatEnergyPerBurst); i++ {
		out.flatEnergyPerBurst[i] = make([]float64, 0)
	}

	for i := 0; i < len(out.rawParticleEnergyPerBurst); i++ {
		out.rawParticleEnergyPerBurst[i] = make([]float64, 0)
	}

	for i := 0; i < len(out.erPerParticleEvent); i++ {
		out.erPerParticleEvent[i] = make([]float64, 0)
	}

	for i := 0; i < len(out.rawPerParticleEvent); i++ {
		out.rawPerParticleEvent[i] = make([]float64, 0)
	}

	core.Events.Subscribe(event.OnBurst, func(_ ...interface{}) bool {
		char := core.Player.ActiveChar()
		ind := char.Index
		out.energySpent[ind] += char.EnergyMax

		out.flatEnergyPerBurst[ind] = append(out.flatEnergyPerBurst[ind], out.charFlatEnergy[ind])
		out.rawParticleEnergyPerBurst[ind] = append(out.rawParticleEnergyPerBurst[ind], out.charRawParticles[ind])
		out.charFlatEnergy[ind] = 0
		out.charRawParticles[ind] = 0

		return false
	}, "stats-burst-energy-log")

	core.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		character := args[0].(*character.CharWrapper)
		ind := character.Index
		preEnergy := args[1].(float64)
		amount := args[2].(float64)
		source := args[3].(string)
		isParticle := args[4].(bool)

		event := stats.EnergyEvent{
			Frame:   core.F,
			Source:  source,
			Gained:  character.Energy - preEnergy,
			Wasted:  preEnergy + amount - character.Energy,
			Current: character.Energy,
		}

		if core.Player.Active() == ind {
			event.FieldStatus = stats.OnField
		} else {
			event.FieldStatus = stats.OffField
		}

		out.events[ind] = append(out.events[ind], event)

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
	}, "stats-energy-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.events); c++ {
		result.Characters[c].EnergyEvents = b.events[c]
	}

	for c := 0; c < len(core.Player.Chars()); c++ {
		result.Characters[c].EnergySpent = b.energySpent[c]
		result.Characters[c].ErNeeded = 1.0

		// ignore the first burst
		for i := 1; i < len(b.flatEnergyPerBurst[c]); i++ {
			erNeeded := (core.Player.Chars()[c].EnergyMax - b.flatEnergyPerBurst[c][i]) / b.rawParticleEnergyPerBurst[c][i]
			if erNeeded > result.Characters[c].ErNeeded {
				result.Characters[c].ErNeeded = erNeeded
			}
		}
		wERsum := 0.0
		wsum := 0.0
		for i, raw := range b.rawParticleEnergyPerBurst[c] {
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
