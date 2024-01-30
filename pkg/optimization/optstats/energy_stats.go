package optstats

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type CustomEnergyStatsBuffer struct {
	ErNeeded   [][]float64
	WeightedER [][]float64
}

func OptimizerERStat(core *core.Core) (CollectorCustomStats[CustomEnergyStatsBuffer], error) {
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

	for ind := range core.Player.Chars() {
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
			raw := amount / er
			charRawParticles[ind] += raw
			erPerParticleEvent[ind] = append(erPerParticleEvent[ind], er)
			rawPerParticleEvent[ind] = append(rawPerParticleEvent[ind], raw)
		} else {
			if amount < 0 {
				charFlatEnergy[ind] -= min(-amount, preEnergy)
			} else {
				// log.Println("Flat energy gained by", character.Base.Key, out.charFlatEnergy[ind])
				charFlatEnergy[ind] += amount
			}
		}
		return false
	}, "substat-opt-energy-log")

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
			out.WeightedER[ind] = append(out.WeightedER[ind], char.Stat(attributes.ER))
		} else {
			out.WeightedER[ind] = append(out.WeightedER[ind], wERsum/wsum)
		}

		erNeeded := 999999999999.9
		if charRawParticles[ind] > 0 {
			erNeeded = max((char.EnergyMax-charFlatEnergy[ind])/charRawParticles[ind], 1.0)
		}
		erPerParticleEvent[ind] = nil
		rawPerParticleEvent[ind] = nil

		out.ErNeeded[ind] = append(out.ErNeeded[ind], erNeeded)

		charRawParticles[ind] = 0
		charFlatEnergy[ind] = 0
		burstCount[ind]++
		// log.Println("After burst", char.Base.Key, out.charFlatEnergy[ind])
		return false
	}, "substat-opt-energy-burst-log")

	return &out, nil
}

func (b CustomEnergyStatsBuffer) Flush(core *core.Core) CustomEnergyStatsBuffer {
	for i := range core.Player.Chars() {
		if len(b.ErNeeded[i]) == 0 {
			b.ErNeeded[i] = append(b.ErNeeded[i], 1)
			b.WeightedER[i] = append(b.WeightedER[i], 1)
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
	charCount := len(cfg.Characters)
	return CustomEnergyAggBuffer{
		WeightedER:         make([][]float64, charCount),
		ErNeeded:           make([][]float64, charCount),
		AdditionalErNeeded: make([][]float64, charCount),
	}
}

func (agg *CustomEnergyAggBuffer) Add(b CustomEnergyStatsBuffer) {
	charCount := len(b.WeightedER)
	for i := 0; i < charCount; i++ {
		burstCount := len(b.WeightedER[i])
		if burstCount == 0 {
			agg.WeightedER[i] = append(agg.WeightedER[i], 1.0)
			agg.ErNeeded[i] = append(agg.ErNeeded[i], 1.0)
			agg.AdditionalErNeeded[i] = append(agg.AdditionalErNeeded[i], 0.0)
		}

		weightedEr := 99999999999.0 // some very large initial value
		erNeeded := 1.0
		additionalNeeded := erNeeded - weightedEr

		// j starts at 1 to ignore the first burst
		for j := 1; j < burstCount; j++ {
			weightedEr = min(weightedEr, b.WeightedER[i][j])
			erNeeded = max(erNeeded, b.ErNeeded[i][j])
			additionalNeeded = max(additionalNeeded, b.ErNeeded[i][j]-b.WeightedER[i][j])
		}

		agg.WeightedER[i] = append(agg.WeightedER[i], weightedEr)
		agg.ErNeeded[i] = append(agg.ErNeeded[i], erNeeded)
		agg.AdditionalErNeeded[i] = append(agg.AdditionalErNeeded[i], additionalNeeded)
	}
}

func (agg *CustomEnergyAggBuffer) Flush() {
	charCount := len(agg.WeightedER)

	for i := 0; i < charCount; i++ {
		slices.Sort(agg.WeightedER[i])
		slices.Sort(agg.ErNeeded[i])
		slices.Sort(agg.AdditionalErNeeded[i])
	}
}
