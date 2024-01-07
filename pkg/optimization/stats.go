package optimization

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type CustomStatsBuffer struct {
	charRawParticles    [][]float64
	charFlatEnergy      [][]float64
	WeightedER          [][]float64
	erPerParticleEvent  [][]float64
	rawPerParticleEvent [][]float64
}

func NewOptimizerStat(core *core.Core) (stats.CollectorCustomStats[CustomStatsBuffer], error) {
	if !core.Flags.IgnoreBurstEnergy {
		// This data doesn't mean much without the IgnoreBurstEnergy flag set
		// So the stat collector disables itself when this flag isn't set
		return &CustomStatsBuffer{}, nil
	}
	out := CustomStatsBuffer{
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
		isParticle := args[4].(bool)
		ind := character.Index

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

func (b CustomStatsBuffer) Flush(core *core.Core) CustomStatsBuffer {
	return b
}

type CustomAggBuffer struct {
}

func NewCustomAggBuffer() CustomAggBuffer {
	return CustomAggBuffer{}
}

func (agg *CustomAggBuffer) Add(b CustomStatsBuffer) {

}
