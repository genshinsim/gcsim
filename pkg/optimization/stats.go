package optimization

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type CustomEnergyStatsBuffer struct {
	ErNeeded   [][]float64
	WeightedER [][]float64
}

func OptimizerERStat(core *core.Core) (stats.CollectorCustomStats[CustomEnergyStatsBuffer], error) {
	if !core.Flags.IgnoreBurstEnergy {
		// This data doesn't mean much without the IgnoreBurstEnergy flag set
		// So the stat collector disables itself when this flag isn't set
		return &CustomEnergyStatsBuffer{}, nil
	}

	out := CustomEnergyStatsBuffer{
		ErNeeded:   make([][]float64, len(core.Player.Chars())),
		WeightedER: make([][]float64, len(core.Player.Chars())),
	}
	burstCount := make([]int, len(core.Player.Chars()))
	erPerParticleEvent := make([][]float64, len(core.Player.Chars()))
	rawPerParticleEvent := make([][]float64, len(core.Player.Chars()))
	charRawParticles := make([]float64, len(core.Player.Chars()))
	charFlatEnergy := make([]float64, len(core.Player.Chars()))

	for ind, _ := range core.Player.Chars() {
		erPerParticleEvent[ind] = make([]float64, 0)
		rawPerParticleEvent[ind] = make([]float64, 0)
		erPerParticleEvent[ind] = append(erPerParticleEvent[ind], 0)
		rawPerParticleEvent[ind] = append(rawPerParticleEvent[ind], 0)
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
			charRawParticles[ind] += raw
			erPerParticleEvent[ind] = append(erPerParticleEvent[ind], 1+er)
			rawPerParticleEvent[ind] = append(rawPerParticleEvent[ind], raw)
		} else {
			if amount < 0 {
				charFlatEnergy[ind] -= max(-amount, preEnergy)
			} else {
				// log.Println("Flat energy gained by", character.Base.Key, out.charFlatEnergy[ind])
				charFlatEnergy[ind] += amount
			}
		}
		return false
	}, "stats-energy-log")

	core.Events.Subscribe(event.OnBurst, func(_ ...interface{}) bool {
		char := core.Player.ActiveChar()
		ind := char.Index

		wERsum := 0.0
		wsum := 0.0
		for i, raw := range rawPerParticleEvent[ind] {
			wERsum += erPerParticleEvent[ind][i] * raw
			wsum += raw
		}
		if wsum == 0 {
			out.WeightedER[ind] = append(out.WeightedER[ind], char.Stat(attributes.ER+1))
		} else {
			out.WeightedER[ind] = append(out.WeightedER[ind], wERsum/wsum)
		}

		erNeeded := 999999999999.9
		if charRawParticles[ind] > 0 {
			erNeeded = max((char.EnergyMax-charFlatEnergy[ind])/charRawParticles[ind], 1.0)
		}
		erPerParticleEvent[ind] = erPerParticleEvent[ind][:0]
		rawPerParticleEvent[ind] = rawPerParticleEvent[ind][:0]

		out.ErNeeded[ind] = append(out.ErNeeded[ind], erNeeded)

		charRawParticles[ind] = 0
		charFlatEnergy[ind] = 0
		burstCount[ind]++
		// log.Println("After burst", char.Base.Key, out.charFlatEnergy[ind])
		return false
	}, "stats-energy-burst-log")

	return &out, nil
}

func (b CustomEnergyStatsBuffer) Flush(core *core.Core) CustomEnergyStatsBuffer {
	for i, _ := range core.Player.Chars() {
		if len(b.ErNeeded[i]) == 0 {
			b.ErNeeded[i] = append(b.ErNeeded[i], 1)
		}
	}

	return b
}

type CustomEnergyAggBuffer struct {
	WeightedER         [][]float64
	ErNeeded           [][]float64
	AdditionalErNeeded [][]float64
}

func NewEnergyAggBuffer(cfg *info.ActionList) CustomEnergyAggBuffer {

	character_count := len(cfg.Characters)
	return CustomEnergyAggBuffer{
		WeightedER:         make([][]float64, character_count),
		ErNeeded:           make([][]float64, character_count),
		AdditionalErNeeded: make([][]float64, character_count),
	}
}

func (agg *CustomEnergyAggBuffer) Add(b CustomEnergyStatsBuffer) {
	char_count := len(b.WeightedER)
	for i := 0; i < char_count; i++ {

		burst_count := len(b.WeightedER[i])
		if burst_count == 0 {
			agg.WeightedER[i] = append(agg.WeightedER[i], 1.0)
			agg.ErNeeded[i] = append(agg.ErNeeded[i], 1.0)
			agg.AdditionalErNeeded[i] = append(agg.AdditionalErNeeded[i], 0.0)
		}

		weighted_er := 99999999999.0 // some very large initial value
		er_needed := 1.0
		additional_needed := 0.0
		for j := 0; j < burst_count; j++ {
			weighted_er = min(weighted_er, b.WeightedER[i][j])
			er_needed = max(er_needed, b.ErNeeded[i][j])
			additional_needed = max(additional_needed, b.ErNeeded[i][j]-b.WeightedER[i][j])
		}

		agg.WeightedER[i] = append(agg.WeightedER[i], weighted_er)
		agg.ErNeeded[i] = append(agg.ErNeeded[i], er_needed)
		agg.AdditionalErNeeded[i] = append(agg.AdditionalErNeeded[i], additional_needed)
	}
}

func (agg *CustomEnergyAggBuffer) Flush() {
	char_count := len(agg.WeightedER)

	for i := 0; i < char_count; i++ {
		slices.Sort(agg.WeightedER[i])
		slices.Sort(agg.ErNeeded[i])
		slices.Sort(agg.AdditionalErNeeded[i])
	}

}
