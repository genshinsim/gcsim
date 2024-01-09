package optimization

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

const FavCritRateBias = 8

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}
type SubstatOptimizerDetails struct {
	charRelevantSubstats   map[keys.Char][]attributes.Stat
	artifactSets4Star      []keys.Set
	substatValues          []float64
	mainstatValues         []float64
	charSubstatFinal       [][]int
	charSubstatLimits      [][]int
	charSubstatRarityMod   []float64
	charProfilesInitial    []info.CharacterProfile
	charWithFavonius       []bool
	charProfilesERBaseline []info.CharacterProfile
	charProfilesCopy       []info.CharacterProfile
	charMaxExtraERSubs     []float64
	simcfg                 *info.ActionList
	gcsl                   ast.Node
	simopt                 simulator.Options
	cfg                    string
	fixedSubstatCount      int
	indivSubstatLiquidCap  int
	totalLiquidSubstats    int
}

// Calculate per-character per-substat "gradients" at initial state using finite differences
// We need to use the ErCalc mode to remove noise from energy, and the ExpectedCritDmg mode
// to remove noise from random crit. This allows us to run a very small 25 iterations per gradient calculation
// TODO: Add setting which allows the user to increase the number of iterations (for cases
// with inherent randomness like Widsith or random delays)
func (stats *SubstatOptimizerDetails) optimizeNonERSubstats() []string {
	var (
		opDebug   []string
		charDebug []string
	)
	origIter := stats.simcfg.Settings.Iterations
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	stats.simcfg.Settings.Iterations = 50
	stats.simcfg.Characters = stats.charProfilesCopy

	for idxChar := range stats.charProfilesCopy {
		charDebug = stats.optimizeNonErSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar])
		opDebug = append(opDebug, charDebug...)
	}
	stats.simcfg.Settings.IgnoreBurstEnergy = false
	stats.simcfg.Settings.Iterations = origIter
	return opDebug
}

// Calculate per-character per-substat "gradients" at initial state using finite differences
func (stats *SubstatOptimizerDetails) optimizeERAndDMGSubstats() []string {
	var (
		opDebug   []string
		charDebug []string
	)

	stats.simcfg.Characters = stats.charProfilesCopy

	for idxChar := range stats.charProfilesCopy {
		charDebug = stats.optimizeERAndDMGSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar])
		opDebug = append(opDebug, charDebug...)
	}
	return opDebug
}

func (stats *SubstatOptimizerDetails) getCharSubstatTotal(idxChar int) int {
	sum := 0
	for _, count := range stats.charSubstatFinal[idxChar] {
		sum += count
	}
	return sum
}

// This function assumes that there are now all subs allocated. For every sub of ER we gain, we will lose one sub of damage
// We compare the damage loss of 1 DMG sub against the damage gain of 1 ER sub.
// We deallocate that DMG sub and allocate 1 ER sub if it would be an overall gain
// Repeat until we cannot allocate ER subs or the DMG loss would be greater than the gain
// The ER dmg gain is prone to noise, so we need to do more iterations
func (stats *SubstatOptimizerDetails) optimizeERAndDMGSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}
	totalSubs := stats.getCharSubstatTotal(idxChar)
	if totalSubs != stats.totalLiquidSubstats {
		opDebug = append(opDebug, fmt.Sprint("Character has", totalSubs, "total liquid subs allocated but expected", stats.totalLiquidSubstats))
	}
	// fmt.Println(char.Base.Key.Pretty(), "has", totalSubs, "total liquid substats")
	for stats.charMaxExtraERSubs[idxChar] > 0.0 && stats.charSubstatFinal[idxChar][attributes.ER] < stats.charSubstatLimits[idxChar][attributes.ER] {
		origIter := stats.simcfg.Settings.Iterations
		stats.simcfg.Settings.IgnoreBurstEnergy = true
		stats.simcfg.Settings.Iterations = 25
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, -1)
		stats.simcfg.Settings.IgnoreBurstEnergy = false
		stats.simcfg.Settings.Iterations = 350
		erGainGradient := stats.calculateSubstatGradientsForChar(idxChar, []attributes.Stat{attributes.ER}, 1)
		stats.simcfg.Settings.Iterations = origIter
		lowestLoss := 0.0
		for idxSubstat, gradient := range substatGradients {
			substat := relevantSubstats[idxSubstat]
			if stats.charSubstatFinal[idxChar][substat] > 0 && gradient < lowestLoss {
				lowestLoss = gradient
			}
		}
		if erGainGradient[0]+lowestLoss <= 0 {
			break
		}
		allocDebug := stats.allocateSomeSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats, -1)
		opDebug = append(opDebug, allocDebug...)
		stats.charSubstatFinal[idxChar][attributes.ER] += 1
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] += float64(1) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
		stats.charMaxExtraERSubs[idxChar] -= 1
	}
	opDebug = append(opDebug, "Final Liquid Substat Counts: "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	return opDebug
}

// This calculation starts all the relevant substats at maximum allocated liquid.
// This reduces the chances of hitting a local maximum of stacking atk%
// It uses a gradient to determine which substat would cause the least damage loss
// when removing. This continues until we are within the total liquid limits
// Initially substats are removed 4 or 2 at a time to speed up the computations
// TODO: Allow the user to specify the removal rate?
// TODO: Multistart gradient descent/ascent from 0 allocated liquid and compare?
func (stats *SubstatOptimizerDetails) optimizeNonErSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	// Reset favonius char crit rate
	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[attributes.CR] -= FavCritRateBias * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
	}

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}

	// start from max liquid in all relevant substats
	for _, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat]) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
		stats.charSubstatFinal[idxChar][substat] = stats.charSubstatLimits[idxChar][substat]
	}
	totalSubs := stats.getCharSubstatTotal(idxChar)
	for totalSubs > stats.totalLiquidSubstats {
		amount := -1
		if totalSubs-stats.totalLiquidSubstats >= 8 {
			amount = -4
		} else if totalSubs-stats.totalLiquidSubstats >= 4 {
			amount = -2
		}
		substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, amount)
		allocDebug := stats.allocateSomeSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats, amount)
		opDebug = append(opDebug, allocDebug...)
		totalSubs = stats.getCharSubstatTotal(idxChar)
	}
	opDebug = append(opDebug, "Liquid Substat Counts: "+PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
	return opDebug
}

func (stats *SubstatOptimizerDetails) allocateSomeSubstatGradientsForChar(
	idxChar int,
	_ info.CharacterProfile,
	substatGradient []float64,
	relevantSubstats []attributes.Stat,
	amount int,
) []string {
	var opDebug []string
	sorted := newSlice(substatGradient...)
	sort.Sort(sort.Reverse(sorted))

	for _, idxSubstat := range sorted.idx {
		substat := relevantSubstats[idxSubstat]

		if amount > 0 {
			if stats.charSubstatFinal[idxChar][substat] < stats.charSubstatLimits[idxChar][substat] {
				stats.charSubstatFinal[idxChar][substat] += amount
				stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
				// fmt.Println("Current Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
				// opDebug = append(opDebug, "Current Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
				return opDebug
			}
		}

		if stats.charSubstatFinal[idxChar][substat] > 0 {
			amount = clamp[int](-stats.charSubstatFinal[idxChar][substat], amount, amount)
			stats.charSubstatFinal[idxChar][substat] += amount
			stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
			// fmt.Println("Current Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
			// opDebug = append(opDebug, "Current Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))
			return opDebug
		}
	}

	// TODO: No relevant substat can be allocated/deallocated, alloc/dealloc some random other substat??
	opDebug = append(opDebug, "Couldn't alloc/dealloc anything?????")
	return opDebug
}

func (stats *SubstatOptimizerDetails) calculateSubstatGradientsForChar(
	idxChar int,
	relevantSubstats []attributes.Stat,
	amount int,
) []float64 {
	stats.simcfg.Characters = stats.charProfilesCopy

	a := NewDamageAggBuffer(stats.simcfg)
	simulator.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now(), OptimizerDmgStat, a.Add)
	a.Flush()

	// TODO: Test if median or mean gives better results
	initialMedian := percentile(a.ExpectedDps, 0.50)

	substatGradients := make([]float64, len(relevantSubstats))
	// Build "gradient" by substat
	for idxSubstat, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

		stats.simcfg.Characters = stats.charProfilesCopy

		a := NewDamageAggBuffer(stats.simcfg)
		simulator.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now(), OptimizerDmgStat, a.Add)
		a.Flush()

		substatGradients[idxSubstat] = percentile(a.ExpectedDps, 0.50) - initialMedian
		// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
		// might give them "too much" cr (= max out liquid cr subs or overcap crit beyond 100%) but that's probably not a big deal
		if stats.charWithFavonius[idxChar] && substat == attributes.CR {
			substatGradients[idxSubstat] += 1000 * float64(amount)
		}
		stats.charProfilesCopy[idxChar].Stats[substat] -= float64(amount) * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	}
	return substatGradients
}

// TODO: Seems like this should be configurable
func (stats *SubstatOptimizerDetails) getNonErSubstatsToOptimizeForChar(char info.CharacterProfile) []attributes.Stat {
	// Get relevant substats, and add additional ones for special characters if needed
	relevantSubstats := []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD, attributes.EM}
	// RIP crystallize...
	if keys.CharKeyToEle[char.Base.Key] == attributes.Geo {
		relevantSubstats = []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD}
	}
	return relevantSubstats
}

// Find optimal ER cutoffs for each character
// We use the ignore_burst_energy mode to determine how much ER is needed for each character to successfully do the
// multiple rotations 75% of the time.
// TODO: Overall it seems like 75% of the time is still bad for the DPS, maybe go through the results
// aggregator and add a p90 to reduce the number of times the 350 iterations of the ER vs DMG step needs
// to be run
func (stats *SubstatOptimizerDetails) optimizeERSubstats() []string {
	var opDebug []string

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

func clamp[T Ordered](min, val, max T) T {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// TODO: Allow the user to specify the initial ER bias? Setting the bias to positive values will mean that the ER vs DMG step runs longer
// But the ER vs DMG step should be more accurate than this function
func (stats *SubstatOptimizerDetails) findOptimalERforChars() {
	stats.simcfg.Settings.IgnoreBurstEnergy = true
	// characters start at maximum ER
	stats.simcfg.Characters = stats.charProfilesERBaseline

	a := NewEnergyAggBuffer(stats.simcfg)
	simulator.RunWithConfigCustomStats(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now(), OptimizerERStat, a.Add)
	a.Flush()
	for idxChar := range stats.charProfilesERBaseline {
		// fmt.Printf("Found character %s has ER len of %d\n", stats.charProfilesERBaseline[idxChar].Base.Key.String(), len(a.AdditionalErNeeded[idxChar]))

		// erDiff is the amount of excess ER we have
		erLen := len(a.AdditionalErNeeded[idxChar])
		erDiff := -percentile(a.AdditionalErNeeded[idxChar], 0.75)

		// find the closest whole count of ER subs
		erStack := int(math.Round(erDiff / stats.substatValues[attributes.ER]))
		erStack = clamp[int](0, erStack, stats.charSubstatFinal[idxChar][attributes.ER])
		stats.charMaxExtraERSubs[idxChar] = float64(erStack) + a.AdditionalErNeeded[idxChar][erLen-1]/stats.substatValues[attributes.ER]
		stats.charProfilesCopy[idxChar] = stats.charProfilesERBaseline[idxChar].Clone()
		stats.charSubstatFinal[idxChar][attributes.ER] -= erStack
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] -= float64(erStack) * stats.substatValues[attributes.ER] * stats.charSubstatRarityMod[idxChar]
	}
	stats.simcfg.Settings.IgnoreBurstEnergy = false
}

func (stats *SubstatOptimizerDetails) setInitialSubstats(fixedSubstatCount int) {
	stats.cloneStatsWithFixedAllocations(fixedSubstatCount)
	stats.calculateERBaseline()
}

// Copy to save initial character state with fixed allocations (2 of each substat)
func (stats *SubstatOptimizerDetails) cloneStatsWithFixedAllocations(fixedSubstatCount int) {
	for i := range stats.simcfg.Characters {
		stats.charProfilesInitial[i] = stats.simcfg.Characters[i].Clone()
		for idxStat, stat := range stats.substatValues {
			if stat == 0 {
				continue
			}
			if attributes.Stat(idxStat) == attributes.ER {
				stats.charProfilesInitial[i].Stats[idxStat] += float64(fixedSubstatCount) * stat
			} else {
				stats.charProfilesInitial[i].Stats[idxStat] += float64(fixedSubstatCount) * stat * stats.charSubstatRarityMod[i]
			}
		}
	}
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
		erStack := stats.charSubstatLimits[i][attributes.ER]
		if stats.charProfilesInitial[i].Base.Key == keys.Raiden {
			erStack = 4
		}
		stats.charSubstatFinal[i][attributes.ER] = erStack

		stats.charProfilesERBaseline[i].Stats[attributes.ER] += float64(erStack) * stats.substatValues[attributes.ER]

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

func NewSubstatOptimizerDetails(
	cfg string,
	simopt simulator.Options,
	simcfg *info.ActionList,
	gcsl ast.Node,
	indivLiquidCap int,
	totalLiquidSubstats int,
	fixedSubstatCount int,
) *SubstatOptimizerDetails {
	s := SubstatOptimizerDetails{}
	s.cfg = cfg
	s.simopt = simopt
	s.simcfg = simcfg
	s.fixedSubstatCount = fixedSubstatCount
	s.indivSubstatLiquidCap = indivLiquidCap
	s.totalLiquidSubstats = totalLiquidSubstats

	s.artifactSets4Star = []keys.Set{
		keys.ResolutionOfSojourner,
		keys.TinyMiracle,
		keys.Berserker,
		keys.Instructor,
		keys.TheExile,
		keys.DefendersWill,
		keys.BraveHeart,
		keys.MartialArtist,
		keys.Gambler,
		keys.Scholar,
		keys.PrayersForWisdom,
		keys.PrayersForDestiny,
		keys.PrayersForIllumination,
		keys.PrayersToSpringtime,
	}

	s.substatValues = make([]float64, attributes.EndStatType)
	s.mainstatValues = make([]float64, attributes.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	s.substatValues[attributes.ATKP] = 0.0496
	s.substatValues[attributes.CR] = 0.0331
	s.substatValues[attributes.CD] = 0.0662
	s.substatValues[attributes.EM] = 19.82
	s.substatValues[attributes.ER] = 0.0551
	s.substatValues[attributes.HPP] = 0.0496
	s.substatValues[attributes.DEFP] = 0.062
	s.substatValues[attributes.ATK] = 16.54
	s.substatValues[attributes.DEF] = 19.68
	s.substatValues[attributes.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	s.mainstatValues[attributes.ATKP] = 0.466
	s.mainstatValues[attributes.CR] = 0.311
	s.mainstatValues[attributes.CD] = 0.622
	s.mainstatValues[attributes.EM] = 186.5
	s.mainstatValues[attributes.ER] = 0.518
	s.mainstatValues[attributes.HPP] = 0.466
	s.mainstatValues[attributes.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	s.charRelevantSubstats = map[keys.Char][]attributes.Stat{
		keys.Albedo:      {attributes.DEFP},
		keys.Hutao:       {attributes.HPP},
		keys.Kokomi:      {attributes.HPP},
		keys.Zhongli:     {attributes.HPP},
		keys.Itto:        {attributes.DEFP},
		keys.Yunjin:      {attributes.DEFP},
		keys.Noelle:      {attributes.DEFP},
		keys.Gorou:       {attributes.DEFP},
		keys.Yelan:       {attributes.HPP},
		keys.Candace:     {attributes.HPP},
		keys.Nilou:       {attributes.HPP},
		keys.Layla:       {attributes.HPP},
		keys.Neuvillette: {attributes.HPP},
		keys.Furina:      {attributes.HPP},
	}

	// Final output array that holds [character][substat_count]
	s.charSubstatFinal = make([][]int, len(simcfg.Characters))
	for i := range simcfg.Characters {
		s.charSubstatFinal[i] = make([]int, attributes.EndStatType)
	}
	s.charMaxExtraERSubs = make([]float64, len(simcfg.Characters))
	s.charSubstatLimits = make([][]int, len(simcfg.Characters))
	s.charSubstatRarityMod = make([]float64, len(simcfg.Characters))
	s.charProfilesInitial = make([]info.CharacterProfile, len(simcfg.Characters))

	// Need to make an exception in energy calcs for these characters for optimization purposes
	s.charWithFavonius = make([]bool, len(simcfg.Characters))
	// Give all characters max ER to set initial state
	s.charProfilesERBaseline = make([]info.CharacterProfile, len(simcfg.Characters))
	s.charProfilesCopy = make([]info.CharacterProfile, len(simcfg.Characters))
	s.gcsl = gcsl

	return &s
}

// Obtain substat count limits based on main stats and also determine 4* set status
// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
// Most people will have 1 5* artifact which messes things up
// TODO: Check whether taking like an average of the two stat values is good enough?
func (stats *SubstatOptimizerDetails) setStatLimits() bool {
	profileIncludesFourStar := false

	for i := range stats.simcfg.Characters {
		stats.charSubstatLimits[i] = make([]int, attributes.EndStatType)
		for idxStat, stat := range stats.mainstatValues {
			if stat == 0 {
				continue
			}
			if stats.simcfg.Characters[i].Stats[idxStat] == 0 {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap
			} else {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap - (stats.fixedSubstatCount * int(math.Round(stats.simcfg.Characters[i].Stats[idxStat]/stats.mainstatValues[idxStat])))
			}
		}

		// Display warning message for 4* sets
		stats.charSubstatRarityMod[i] = 1
		for set := range stats.simcfg.Characters[i].Sets {
			for _, fourStar := range stats.artifactSets4Star {
				if set == fourStar {
					profileIncludesFourStar = true
					stats.charSubstatRarityMod[i] = 0.8
				}
			}
		}
	}

	return profileIncludesFourStar
}

// Helper function to pretty print substat counts. Stolen from similar function that takes in the float array
func PrettyPrintStatsCounts(statsCounts []int) string {
	var sb strings.Builder
	for i, v := range statsCounts {
		if v > 0 {
			sb.WriteString(attributes.StatTypeString[i])
			sb.WriteString(": ")
			sb.WriteString(fmt.Sprintf("%v", v))
			sb.WriteString(" ")
		}
	}
	return strings.Trim(sb.String(), " ")
}
