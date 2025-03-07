package optimization

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/optimization/optstats"
)

const FavCritRateBias = 8

// Start at minimum ER
// This is to ensure that the recommended code for preventing ER substat allocation works:
//
//	if .char.burst.ready && .char.energy == .char.energymax {
//	  char burst;
//	}
func (stats *SubstatOptimizerDetails) calculateERBaseline() {
	for i := range stats.charProfilesInitial {
		stats.charProfilesERBaseline[i] = stats.charProfilesInitial[i].Clone()
		// Need special exception to Raiden due to her burst mechanics
		// TODO: Don't think there's a better solution without an expensive recursive solution to check across all Raiden ER states
		// Practically high ER substat Raiden is always currently unoptimal, so we just set her initial stacks lowish
		erSubs := 0
		if stats.charProfilesInitial[i].Base.Key == keys.Raiden {
			erSubs = 4
		}
		stats.charSubstatFinal[i][attributes.ER] = erSubs

		stats.charProfilesERBaseline[i].Stats[attributes.ER] += float64(erSubs) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[i]

		if strings.Contains(stats.charProfilesInitial[i].Weapon.Name, "favonius") {
			stats.calculateERBaselineHandleFav(i)
		}
	}
}

// Current strategy for favonius is to just boost this character's crit values a bit extra for optimal ER calculation purposes
// Then at next step of substat optimization, should naturally see relatively big DPS increases for that character if higher crit matters a lot
// TODO: Do we need a better special case for favonius?
func (stats *SubstatOptimizerDetails) calculateERBaselineHandleFav(i int) {
	stats.charProfilesERBaseline[i].Stats[attributes.CR] += FavCritRateBias * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[i]
	stats.charWithFavonius[i] = true
}

// Find optimal ER cutoffs for each character
// We use the ignore_burst_energy mode to determine how much ER is needed for each character to successfully do the
// multiple rotations 75% of the time.
// TODO: Add option for user to set the percentile used for optimization
func (stats *SubstatOptimizerDetails) optimizeERSubstats() {
	stats.simcfg.Settings.Iterations = 350

	// For now going to ignore Raiden, since typically she won't be running maximum ER subs just to battery. The scaling isn't that strong
	// From minimum subs (0.1102 ER) to maximum subs (0.6612 ER) she restores 4 more flat energy per rotation.
	// She starts at 4 liquid so it's +/- 2 flat energy
	stats.findOptimalERforChars()

	// Fix ER at previously found values then optimize all other substats
	stats.optimizer.logger.Info("Initial Calculated ER Liquid Substats by character:")
	output := ""
	for i := range stats.charProfilesInitial {
		output +=
			fmt.Sprintf("%v: %.4g, ",
				stats.charProfilesInitial[i].Base.Key.String(),
				float64(stats.charSubstatFinal[i][attributes.ER])*stats.substatValues[attributes.ER]*stats.charSubstatRarityMod[i],
			)
	}
	stats.optimizer.logger.Info(output)
}

func (stats *SubstatOptimizerDetails) findOptimalERforChars() {
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	// characters start at minimum ER
	stats.simcfg.Characters = stats.charProfilesERBaseline

	seed := time.Now().UnixNano()
	a := optstats.NewEnergyAggBuffer(stats.simcfg)
	_, err := optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, seed, optstats.OptimizerERStat, a.Add)
	if err != nil {
		stats.optimizer.logger.Fatal(err.Error())
	}
	a.Flush()
	for idxChar := range stats.charProfilesERBaseline {
		// erDiff is the amount of ER we need
		erLen := len(a.AdditionalErNeeded[idxChar])
		if stats.optimizer.verbose {
			hist := fmtHist(a.ErNeeded[idxChar], float64(int(a.ErNeeded[idxChar][0]*10))/10.0, 0.05)
			stats.optimizer.logger.Infof("%v: ER Needed Distribution", stats.charProfilesInitial[idxChar].Base.Key.Pretty())
			for _, val := range hist {
				stats.optimizer.logger.Infoln(val)
			}
		}
		erDiff := percentile(a.AdditionalErNeeded[idxChar], 0.8)

		erSubVal := stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]

		// find the closest whole count of ER subs
		// TODO: is ceil better than round? Maybe round with some kind of bias?
		erSubs := int(math.Round(erDiff / erSubVal))

		// Raiden doesn't start at 0 ER subs so need to subtract that out
		erSubs = clamp[int](0, erSubs, stats.charSubstatLimits[idxChar][attributes.ER]-stats.charSubstatFinal[idxChar][attributes.ER])
		stats.charMaxExtraERSubs[idxChar] = math.Ceil(a.AdditionalErNeeded[idxChar][erLen-1]/erSubVal) - float64(stats.charSubstatFinal[idxChar][attributes.ER])
		stats.charProfilesCopy[idxChar] = stats.charProfilesERBaseline[idxChar].Clone()
		stats.charSubstatFinal[idxChar][attributes.ER] += erSubs
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] += float64(erSubs) * erSubVal
	}
	stats.simcfg.Settings.IgnoreBurstEnergy = false
}
