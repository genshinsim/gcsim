package simulation

import (
	"errors"
	"log"
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var subDist [attributes.DelimBaseStat]float64
var subUpgrade [attributes.DelimBaseStat][4]float64

// mainstat at lvl 20
var mainStat = map[attributes.Stat]float64{
	attributes.HP:       4780,
	attributes.ATK:      311,
	attributes.HPP:      0.466,
	attributes.ATKP:     0.466,
	attributes.DEFP:     0.583,
	attributes.PyroP:    0.466,
	attributes.HydroP:   0.466,
	attributes.CryoP:    0.466,
	attributes.ElectroP: 0.466,
	attributes.AnemoP:   0.466,
	attributes.GeoP:     0.466,
	attributes.DendroP:  0.466,
	attributes.PhyP:     0.466,
	attributes.EM:       186.5,
	attributes.ER:       0.518,
	attributes.CD:       0.622,
	attributes.CR:       0.311,
	attributes.Heal:     0.359,
}

func init() {
	subDist[attributes.HP] = 6
	subDist[attributes.ATK] = 6
	subDist[attributes.DEF] = 6
	subDist[attributes.HPP] = 4
	subDist[attributes.ATKP] = 4
	subDist[attributes.DEFP] = 4
	subDist[attributes.ER] = 4
	subDist[attributes.EM] = 4
	subDist[attributes.CR] = 3
	subDist[attributes.CD] = 3

	subUpgrade[attributes.HP] = [4]float64{209, 239, 269, 299}
	subUpgrade[attributes.DEF] = [4]float64{16, 19, 21, 23}
	subUpgrade[attributes.ATK] = [4]float64{14, 16, 18, 19}
	subUpgrade[attributes.HPP] = [4]float64{0.041, 0.047, 0.053, 0.058}
	subUpgrade[attributes.DEFP] = [4]float64{0.051, 0.058, 0.066, 0.073}
	subUpgrade[attributes.ATKP] = [4]float64{0.041, 0.047, 0.053, 0.058}
	subUpgrade[attributes.EM] = [4]float64{16, 19, 21, 23}
	subUpgrade[attributes.ER] = [4]float64{0.045, 0.052, 0.058, 0.065}
	subUpgrade[attributes.CR] = [4]float64{0.027, 0.031, 0.035, 0.039}
	subUpgrade[attributes.CD] = [4]float64{0.054, 0.062, 0.07, 0.078}
}

func generateRandSubs(r *info.RandomSubstats, rng *rand.Rand) ([]float64, error) {
	if r.Rarity != 5 {
		return nil, errors.New("sorry only 5 star artifacts supported currently")
	}
	stats := make([]float64, attributes.EndStatType)
	// main stats first
	stats[attributes.ATK] = mainStat[attributes.ATK]
	stats[attributes.HP] = mainStat[attributes.HP]
	stats[r.Sand] += mainStat[r.Sand]
	stats[r.Goblet] += mainStat[r.Goblet]
	stats[r.Circlet] += mainStat[r.Circlet]

	mains := [5]attributes.Stat{attributes.ATK, attributes.HP, r.Sand, r.Goblet, r.Circlet}

	for _, m := range mains {
		// weights
		var weight [attributes.DelimBaseStat]float64
		var picked [4]attributes.Stat
		copy(weight[:], subDist[:])
		weight[m] = 0

		//TODO: option to use boss
		upgrades := 4
		if rng.Float64() <= 0.2 {
			upgrades = 5
		}
		for i := 0; i < 4; i++ {
			// pick stat from weight
			s := randSub(weight, rng)
			if s == attributes.NoStat {
				log.Println("weights no good?")
				log.Println(subDist)
				log.Println(weight)
				return nil, errors.New("unexpected error picking random sub; none found")
			}
			weight[s] = 0
			stats[s] += subUpgrade[s][rng.Intn(4)]
			picked[i] = s
		}
		for i := 0; i < upgrades; i++ {
			// pick one out of 4 stats
			s := picked[rng.Intn(4)]
			stats[s] += subUpgrade[s][rng.Intn(4)]
		}
	}
	return stats, nil
}

func randSub(weights [attributes.DelimBaseStat]float64, rng *rand.Rand) attributes.Stat {
	var sumWeights float64
	for _, v := range weights {
		sumWeights += v
	}
	pick := rng.Float64() * sumWeights
	var cum float64
	for i, v := range weights {
		cum += v
		if pick <= cum {
			return attributes.Stat(i)
		}
	}
	return attributes.NoStat
}
