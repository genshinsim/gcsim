package energy_burst

import (
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(agg.Config{
		Name: "energy_burst",
		New:  NewAgg,
	})
}

type buffer struct {
	rawParticles [][][]float64
	flatEnergy   [][][]float64
	weightedER   [][][]float64
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		rawParticles: make([][][]float64, len(cfg.Characters)),
		flatEnergy:   make([][][]float64, len(cfg.Characters)),
		weightedER:   make([][][]float64, len(cfg.Characters)),
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		b.weightedER[i] = append(b.weightedER[i], result.Characters[i].EnergyInfo.WeightedER)
		b.flatEnergy[i] = append(b.flatEnergy[i], result.Characters[i].EnergyInfo.FlatEnergyPerBurst)
		b.rawParticles[i] = append(b.rawParticles[i], result.Characters[i].EnergyInfo.RawParticlesPerBurst)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.CharacterEnergyInfo = make([]*model.CharacterEnergyInfo, len(b.weightedER))
	for char := range b.weightedER {
		energyInfo := make([]*model.EnergyPerBurstInfo, len(b.weightedER[char]))
		for iter := range b.weightedER[char] {
			energyInfo[iter] = &model.EnergyPerBurstInfo{
				WeightedEr:     b.weightedER[char][iter],
				FlatEnergy:     b.flatEnergy[char][iter],
				ParticleEnergy: b.rawParticles[char][iter],
			}
		}
		result.CharacterEnergyInfo[char] = &model.CharacterEnergyInfo{BurstEnergyInfo: energyInfo}
	}
}
