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

// Find optimal ER cutoffs for each character
// We use the ignore_burst_energy mode to determine how much ER is needed for each character to successfully do the
// multiple rotations 75% of the time.
func (stats *SubstatOptimizerDetails) optimizeERSubstats() []string {
	var opDebug []string
	stats.simcfg.Settings.Iterations = 350

	// For now going to ignore Raiden, since typically she won't be running maximum ER subs just to battery. The scaling isn't that strong
	// From minimum subs (0.1102 ER) to maximum subs (0.6612 ER) she restores 4 more flat energy per rotation.
	// She starts at 4 liquid so it's +/- 2 flat energy
	stats.findOptimalERforChars()

	// Fix ER at previously found values then optimize all other substats
	opDebug = append(opDebug, "Initial Calculated ER Liquid Substats by character:")
	printVal := ""
	for i := range stats.charProfilesInitial {
		printVal += fmt.Sprintf(
			"%v: %.4g, ",
			stats.charProfilesInitial[i].Base.Key.String(),
			float64(stats.charSubstatFinal[i][attributes.ER])*stats.substatValues[attributes.ER],
		)
	}
	opDebug = append(opDebug, printVal)

	return opDebug
}

// TODO: Allow the user to specify the initial ER bias? Setting the bias to positive values will mean that the ER vs DMG step runs longer
// But the ER vs DMG step should be more accurate than this function
func (stats *SubstatOptimizerDetails) findOptimalERforChars() {
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	// characters start at maximum ER
	stats.simcfg.Characters = stats.charProfilesERBaseline

	a := optstats.NewEnergyAggBuffer(stats.simcfg)
	optstats.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now(), optstats.OptimizerERStat, a.Add)
	a.Flush()
	for idxChar := range stats.charProfilesERBaseline {
		// erDiff is the amount of ER we need
		erLen := len(a.AdditionalErNeeded[idxChar])
		erDiff := percentile(a.AdditionalErNeeded[idxChar], 0.75)

		// find the closest whole count of ER subs
		erSubs := int(math.Round(erDiff / stats.substatValues[attributes.ER]))
		erSubs = clamp[int](0, erSubs, stats.charSubstatLimits[idxChar][attributes.ER]-stats.charSubstatFinal[idxChar][attributes.ER])
		stats.charMaxExtraERSubs[idxChar] = math.Ceil(a.AdditionalErNeeded[idxChar][erLen-1]/stats.substatValues[attributes.ER]) - float64(stats.charSubstatFinal[idxChar][attributes.ER])
		stats.charProfilesCopy[idxChar] = stats.charProfilesERBaseline[idxChar].Clone()
		stats.charSubstatFinal[idxChar][attributes.ER] += erSubs
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] += float64(erSubs) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
	}
	stats.simcfg.Settings.IgnoreBurstEnergy = false
}

// Add some points into CR/CD to reduce crit variance and have reasonable baseline stats
// Also helps to slightly better evaluate the impact of favonius
// Current concern is that optimization on 2nd stage doesn't perform very well due to messed up rotation
func (stats *SubstatOptimizerDetails) calculateERBaseline() {
	for i := range stats.charProfilesInitial {
		stats.charProfilesERBaseline[i] = stats.charProfilesInitial[i].Clone()
		// Need special exception to Raiden due to her burst mechanics
		// TODO: Don't think there's a better solution without an expensive recursive solution to check across all Raiden ER states
		// Practically high ER substat Raiden is always currently unoptimal, so we just set her initial stacks low
		erSubs := 0
		if stats.charProfilesInitial[i].Base.Key == keys.Raiden {
			erSubs = 4
		}
		stats.charSubstatFinal[i][attributes.ER] = erSubs

		stats.charProfilesERBaseline[i].Stats[attributes.ER] += float64(erSubs) * stats.substatValues[attributes.ER]

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
