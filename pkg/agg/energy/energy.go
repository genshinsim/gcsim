package energy

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	sourceEnergy []map[string]*calc.StreamStats

	rawParticles [][][]float64
	flatEnergy   [][][]float64
	weightedER   [][][]float64
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		sourceEnergy: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		rawParticles: make([][][]float64, len(cfg.Characters)),
		flatEnergy:   make([][][]float64, len(cfg.Characters)),
		weightedER:   make([][][]float64, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.sourceEnergy[i] = make(map[string]*calc.StreamStats)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	for i := range result.Characters {
		sourceEnergy := make(map[string]float64)
		for _, ev := range result.Characters[i].EnergyEvents {
			sourceEnergy[ev.Source] += ev.Gained + ev.Wasted
		}
		for k, v := range sourceEnergy {
			if _, ok := b.sourceEnergy[i][k]; !ok {
				b.sourceEnergy[i][k] = &calc.StreamStats{}
			}
			b.sourceEnergy[i][k].Add(v)
		}
		b.weightedER[i] = append(b.weightedER[i], result.Characters[i].EnergyInfo.WeightedER)
		b.flatEnergy[i] = append(b.flatEnergy[i], result.Characters[i].EnergyInfo.FlatEnergyPerBurst)
		b.rawParticles[i] = append(b.rawParticles[i], result.Characters[i].EnergyInfo.RawParticlesPerBurst)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.TotalSourceEnergy = make([]*model.SourceStats, len(b.sourceEnergy))
	for i, c := range b.sourceEnergy {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.TotalSourceEnergy[i] = &model.SourceStats{
			Sources: source,
		}
	}
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
