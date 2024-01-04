package energy_burst

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
	charRawParticles    [][]float64
	charFlatEnergy      [][]float64
	WeightedER          [][]float64
	erPerParticleEvent  [][]float64
	rawPerParticleEvent [][]float64
}

func NewStat(core *core.Core) (stats.Collector, error) {
	out := buffer{
		charRawParticles:    make([][]float64, len(core.Player.Chars())),
		charFlatEnergy:      make([][]float64, len(core.Player.Chars())),
		WeightedER:          make([][]float64, len(core.Player.Chars())),
		erPerParticleEvent:  make([][]float64, len(core.Player.Chars())),
		rawPerParticleEvent: make([][]float64, len(core.Player.Chars())),
	}

	burstCount := make([]int, len(core.Player.Chars()))
	for ind := 0; ind < len(core.Player.Chars()); ind++ {
		out.charRawParticles[ind] = append(out.charRawParticles[ind], 0)
		out.charFlatEnergy[ind] = append(out.charFlatEnergy[ind], 0)
		out.erPerParticleEvent[ind] = make([]float64, 0)
		out.rawPerParticleEvent[ind] = make([]float64, 0)
		out.erPerParticleEvent[ind] = append(out.erPerParticleEvent[ind], 0)
		out.rawPerParticleEvent[ind] = append(out.rawPerParticleEvent[ind], 0)
	}

	core.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		character := args[0].(*character.CharWrapper)
		preEnergy := args[1].(float64)
		amount := args[2].(float64)
		source := args[3].(string)
		isParticle := args[4].(bool)
		ind := character.Index

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

		er := character.Stat(attributes.ER)

		if isParticle {
			raw := amount / (1.0 + er)
			out.charRawParticles[ind][burstCount[ind]] += raw
			out.erPerParticleEvent[ind] = append(out.erPerParticleEvent[ind], 1+er)
			out.rawPerParticleEvent[ind] = append(out.rawPerParticleEvent[ind], raw)
		} else {
			if amount < 0 {
				out.charFlatEnergy[ind][burstCount[ind]] -= max(-amount, preEnergy)
			} else {
				// log.Println("Flat energy gained by", character.Base.Key, out.charFlatEnergy[ind])
				out.charFlatEnergy[ind][burstCount[ind]] += amount
			}
		}
		return false
	}, "stats-energy-log")

	core.Events.Subscribe(event.OnBurst, func(_ ...interface{}) bool {
		char := core.Player.ActiveChar()
		ind := char.Index

		wERsum := 0.0
		wsum := 0.0
		for i, raw := range out.rawPerParticleEvent[ind] {
			wERsum += out.erPerParticleEvent[ind][i] * raw
			wsum += raw
		}
		if wsum == 0 {
			out.WeightedER[ind] = append(out.WeightedER[ind], char.Stat(attributes.ER+1))
		} else {
			out.WeightedER[ind] = append(out.WeightedER[ind], wERsum/wsum)
		}
		out.erPerParticleEvent[ind] = make([]float64, 0)
		out.rawPerParticleEvent[ind] = make([]float64, 0)

		burstCount[ind]++
		out.charRawParticles[ind] = append(out.charRawParticles[ind], 0)
		out.charFlatEnergy[ind] = append(out.charFlatEnergy[ind], 0)
		// log.Println("After burst", char.Base.Key, out.charFlatEnergy[ind])
		return false
	}, "stats-energy-burst-log")

	return &out, nil
}

func (b buffer) Flush(core *core.Core, result *stats.Result) {
	for c := 0; c < len(b.charRawParticles); c++ {
		data := stats.EnergyInfo{
			RawParticlesPerBurst: b.charRawParticles[c],
			FlatEnergyPerBurst:   b.charFlatEnergy[c],
			WeightedER:           b.WeightedER[c],
		}
		result.Characters[c].EnergyInfo = data
	}
}
