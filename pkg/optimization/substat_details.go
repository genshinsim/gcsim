package optimization

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
	"math"
	"sort"
	"strings"
)

type SubstatOptimizerDetails struct {
	charRelevantSubstats   map[core.CharKey][]core.StatType
	artifactSets4Star      []string
	substatValues          []float64
	mainstatValues         []float64
	charSubstatFinal       [][]int
	charSubstatLimits      [][]int
	charSubstatRarityMod   []float64
	charProfilesInitial    []core.CharacterProfile
	charWithFavonius       []bool
	charProfilesERBaseline []core.CharacterProfile
	charProfilesCopy       []core.CharacterProfile
	simcfg                 core.SimulationConfig
	indivSubstatLiquidCap  int
	totalLiquidSubstats    int
}

// Calculate per-character per-substat "gradients" at initial state using finite differences
// In practical evaluations, adding small numbers of substats (<10) can be VERY noisy
// Therefore, "gradient" evaluations are done in groups of 10 substats
// Allocation strategy is to just max substats according to highest gradient to lowest
// TODO: Probably want to refactor to potentially run gradient step at least twice:
// once initially then another at 10 assigned liquid substats
// Fine grained evaluations are too expensive time wise, but can perhaps add in an option for people who want to sit around for a while
func (stats *SubstatOptimizerDetails) optimizeNonERSubstats(runner simRunner) []string {
	var (
		opDebug   []string
		charDebug []string
	)

	stats.simcfg.Characters.Profile = stats.charProfilesCopy

	// Get initial DPS value
	initialResult := runner(stats.simcfg)
	initialMean := initialResult.DPS.Mean

	opDebug = append(opDebug, "Calculating optimal substat distribution...")
	opDebug = append(opDebug, fmt.Sprintf("%v", initialMean))

	for idxChar, char := range stats.charProfilesCopy {
		charDebug = stats.optimizeNonErSubstatsForChar(idxChar, char, initialMean, runner)
		opDebug = append(opDebug, charDebug...)
	}

	return opDebug
}

func (stats *SubstatOptimizerDetails) optimizeNonErSubstatsForChar(idxChar int, char core.CharacterProfile, initialMean float64, runner simRunner) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	// Reset favonius char crit rate
	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[core.CR] -= 8 * stats.substatValues[core.CR] * stats.charSubstatRarityMod[idxChar]
	}

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}

	substatGradients, gradDebug := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, initialMean, runner)
	opDebug = append(opDebug, gradDebug...)

	allocDebug := stats.allocateSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats)
	opDebug = append(opDebug, allocDebug...)

	return opDebug
}

func (stats *SubstatOptimizerDetails) allocateSubstatGradientsForChar(idxChar int, char core.CharacterProfile, substatGradient []float64, relevantSubstats []core.StatType) []string {
	var opDebug []string

	sorted := NewSlice(substatGradient...)
	sort.Sort(sort.Reverse(sorted))

	printVal := ""
	for i, idxSorted := range sorted.idx {
		printVal += fmt.Sprintf("%v: %5.5g, ", relevantSubstats[idxSorted], sorted.slice[i])
	}
	opDebug = append(opDebug, printVal)

	for idxGrad, idxSubstat := range sorted.idx {
		allocDebug := stats.allocateSubstatGradientForChar(idxChar, char, sorted, idxGrad, idxSubstat, relevantSubstats)
		opDebug = append(opDebug, allocDebug...)
	}

	opDebug = append(opDebug, "Final Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))

	stats.resetFavoniusCritRateForChar(idxChar)

	return opDebug
}

func (stats *SubstatOptimizerDetails) resetFavoniusCritRateForChar(idxChar int) {
	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[core.CR] += 8 * stats.substatValues[core.CR] * stats.charSubstatRarityMod[idxChar]
	}
}

func (stats *SubstatOptimizerDetails) allocateSubstatGradientForChar(idxChar int, char core.CharacterProfile, sorted *Slice, idxGrad int, idxSubstat int, relevantSubstats []core.StatType) []string {
	var opDebug []string

	substatToMax := relevantSubstats[idxSubstat]

	// TODO: Improve this by adding a mix of CR/CD substats based on the ratio of gradient increase from CR/CD
	// If CR/CD is one of the selected substats, then adding them in a mix is generally most optimal
	// Use the ratio between gradient values to determine mix %
	// Need manual override here since gradient method from init does not always find this result
	var crCDSubstatRatio float64
	var gradStat float64
	switch substatToMax {
	case core.CR:
		gradCR := sorted.slice[idxGrad]
		gradCD := 0.0
		for i, idxSubstatTemp := range sorted.idx {
			if relevantSubstats[idxSubstatTemp] == core.CD {
				gradCD = sorted.slice[i]
			}
		}
		crCDSubstatRatio = gradCR / gradCD
	case core.CD:
		gradCD := sorted.slice[idxGrad]
		gradCR := 0.0
		for i, idxSubstatTemp := range sorted.idx {
			if relevantSubstats[idxSubstatTemp] == core.CR {
				gradCR = sorted.slice[i]
			}
		}
		crCDSubstatRatio = gradCR / gradCD
	default:
		gradStat = sorted.slice[idxGrad]
	}

	// If DPS change is really low, then it's usually better to just toss a few extra points into ER for stability
	if gradStat < 50 && crCDSubstatRatio == 0 {
		stats.assignSubstatsForChar(idxChar, char, core.ER, 4)
		opDebug = append(opDebug, "Low damage contribution from substats - adding some points to ER instead")
	}

	var globalLimit int
	var crLimit int
	var cdLimit int
	if crCDSubstatRatio > 0 {
		globalLimit, crLimit = stats.assignSubstatsForChar(idxChar, char, core.CR, 0)
		_, cdLimit = stats.assignSubstatsForChar(idxChar, char, core.CD, 0)

		// Continually add CR/CD to try to align CR/CD ratio to ratio of gradients until we hit a limit
		var currentRatio float64
		var amtCR int
		var amtCD int
		currentStat := core.CR
		// Debug to avoid runaway loops...
		var iteration int
		// Want this to continue until either global cap is reached, or we can neither add CR/CD
		for globalLimit > 0 && (crLimit > 0 || cdLimit > 0) && iteration < 100 {
			if stats.charSubstatFinal[idxChar][core.CD] == 0 {
				currentRatio = float64(stats.charSubstatFinal[idxChar][core.CR])
			} else {
				currentRatio = float64(stats.charSubstatFinal[idxChar][core.CR]) / float64(stats.charSubstatFinal[idxChar][core.CD])
			}

			if currentRatio > crCDSubstatRatio {
				amtCR = 0
				amtCD = 1
			} else if currentRatio <= crCDSubstatRatio {
				amtCR = 1
				amtCD = 0
			}

			// When we hit the limit on one stat, just try to fill the other up to max
			if crLimit == 0 {
				amtCD = 10
			}
			if cdLimit == 0 {
				amtCR = 10
			}

			if currentStat == core.CR {
				globalLimit, crLimit = stats.assignSubstatsForChar(idxChar, char, core.CR, amtCR)
				currentStat = core.CD
			} else if currentStat == core.CD {
				globalLimit, cdLimit = stats.assignSubstatsForChar(idxChar, char, core.CD, amtCD)
				currentStat = core.CR
			}
			iteration += 1
		}
	} else {
		globalLimit, _ = stats.assignSubstatsForChar(idxChar, char, substatToMax, 12)
	}
	if globalLimit == 0 {
		return opDebug
	}

	return opDebug
}

// Assigns substats and returns the remaining global limit and individual substat limit
func (stats *SubstatOptimizerDetails) assignSubstatsForChar(idxChar int, char core.CharacterProfile, substat core.StatType, amt int) (int, int) {
	totalSubstatCount := 0
	for _, val := range stats.charSubstatFinal[idxChar] {
		totalSubstatCount += val
	}

	baseLiquidSubstats := stats.totalLiquidSubstats
	for set, count := range char.Sets {
		for _, setfourstar := range stats.artifactSets4Star {
			if set == setfourstar {
				baseLiquidSubstats -= 2 * count
			}
		}
	}

	remainingLiquidSubstats := baseLiquidSubstats - totalSubstatCount
	// Minimum of individual limit, global limit, desired amount
	amtToAdd := MinInt(stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat], remainingLiquidSubstats, amt)
	stats.charSubstatFinal[idxChar][substat] += amtToAdd

	return remainingLiquidSubstats - amtToAdd, stats.charSubstatLimits[idxChar][substat] - stats.charSubstatFinal[idxChar][substat]
}

func (stats *SubstatOptimizerDetails) calculateSubstatGradientsForChar(idxChar int, relevantSubstats []core.StatType, initialMean float64, runner simRunner) ([]float64, []string) {
	var opDebug []string
	substatGradients := make([]float64, len(relevantSubstats))

	// Build "gradient" by substat
	for idxSubstat, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

		stats.simcfg.Characters.Profile = stats.charProfilesCopy
		substatEvalResult := runner(stats.simcfg)
		// opDebug = append(opDebug, fmt.Sprintf("%v: %v (%v)", substat.String(), substatEvalResult.DPS.Mean, substatEvalResult.DPS.SD))

		substatGradients[idxSubstat] = substatEvalResult.DPS.Mean - initialMean

		// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
		// might give them "too much" cr (= max out liquid cr subs) but that's probably not a big deal
		if stats.charWithFavonius[idxChar] && substat == core.CR {
			substatGradients[idxSubstat] += 1000
		}
		stats.charProfilesCopy[idxChar].Stats[substat] -= 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
	}

	return substatGradients, opDebug
}

// TODO: Seems like this should be configurable
func (stats *SubstatOptimizerDetails) getNonErSubstatsToOptimizeForChar(char core.CharacterProfile) []core.StatType {
	// Get relevant substats, and add additional ones for special characters if needed
	relevantSubstats := []core.StatType{core.ATKP, core.CR, core.CD, core.EM}
	// RIP crystallize...
	if core.CharKeyToEle[char.Base.Key] == core.Geo {
		relevantSubstats = []core.StatType{core.ATKP, core.CR, core.CD}
	}

	return relevantSubstats
}

// Find optimal ER cutoffs for each character
// For each character, do grid search to find optimal ER values
// TODO: Can maybe replace with some kind of gradient descent for speed improvements/allow for 1 ER substat moves?
// When I tried before, it was hard to define a good step size and penalty on high ER substats that generally worked well
// At least this version works semi-reliably...
func (stats *SubstatOptimizerDetails) optimizeERSubstats(tolMean float64, tolSD float64, runner simRunner) []string {
	var (
		charDebug []string
		opDebug   []string
	)

	for idxChar, char := range stats.charProfilesERBaseline {
		charDebug = stats.findOptimalERforChar(idxChar, char, tolMean, tolSD, runner)
		opDebug = append(opDebug, charDebug...)
	}

	// Need a separate optimization routine for strong battery characters (currently Raiden only, maybe EMC?)
	// Need to set all other character's ER substats at final value, then see added benefit from ER for global battery chars
	for i, char := range stats.charProfilesERBaseline {
		stats.charProfilesERBaseline[i].Stats[core.ER] = stats.charProfilesInitial[i].Stats[core.ER]

		if char.Base.Key == core.Raiden {
			stats.charSubstatFinal[i][core.ER] = 10
		}

		stats.charProfilesERBaseline[i].Stats[core.ER] += float64(stats.charSubstatFinal[i][core.ER]) * stats.substatValues[core.ER]
	}

	for i, char := range stats.charProfilesERBaseline {
		if char.Base.Key != core.Raiden {
			continue
		}
		opDebug = append(opDebug, "Raiden found in team comp - running secondary optimization routine...")
		charDebug = stats.findOptimalERforChar(i, char, tolMean, tolSD, runner)
		opDebug = append(opDebug, charDebug...)
	}

	// Fix ER at previously found values then optimize all other substats
	opDebug = append(opDebug, "Optimized ER Liquid Substats by character:")
	printVal := ""
	for i, char := range stats.charProfilesInitial {
		printVal += fmt.Sprintf("%v: %.4g, ", char.Base.Key.String(), float64(stats.charSubstatFinal[i][core.ER])*stats.substatValues[core.ER])
	}
	opDebug = append(opDebug, printVal)

	return opDebug
}

func (stats *SubstatOptimizerDetails) setCharProfilesCopy(charsToCopy []core.CharacterProfile) {
	for i, char := range charsToCopy {
		stats.charProfilesCopy[i] = char.Clone()
	}
}

func (stats *SubstatOptimizerDetails) findOptimalERforChar(idxChar int, char core.CharacterProfile, tolMean float64, tolSD float64, runner simRunner) []string {
	var debug []string
	var initialMean float64
	var initialSD float64

	debug = append(debug, fmt.Sprintf("%v", char.Base.Key))

	for erStack := 0; erStack <= 10; erStack += 2 {
		stats.charProfilesCopy[idxChar] = char.Clone()
		stats.charProfilesCopy[idxChar].Stats[core.ER] -= float64(erStack) * stats.substatValues[core.ER]

		stats.simcfg.Characters.Profile = stats.charProfilesCopy

		result := runner(stats.simcfg)
		debug = append(debug, fmt.Sprintf("%v: %v (%v)", stats.charSubstatFinal[idxChar][core.ER]-erStack, result.DPS.Mean, result.DPS.SD))

		if erStack == 0 {
			initialMean = result.DPS.Mean
			initialSD = result.DPS.SD
		}

		condition := result.DPS.Mean/initialMean-1 < -tolMean || result.DPS.SD/initialSD-1 > tolSD
		// For Raiden, we can't use DPS directly as a measure since she scales off of her own ER
		// Instead we ONLY use the SD tolerance as big jumps indicate the rotation is becoming more unstable
		if char.Base.Key == core.Raiden {
			condition = result.DPS.SD/initialSD-1 > tolSD
		}

		// If differences exceed tolerances, then immediately break
		if condition {
			// Reset character stats
			stats.charProfilesCopy[idxChar] = char.Clone()
			// Save ER value - optimal value is the value immediately prior, so we subtract 2
			stats.charSubstatFinal[idxChar][core.ER] -= erStack - 2
			break
		}

		// Reached minimum possible ER stacks, so optimal is the minimum amount of ER stacks
		if stats.charSubstatFinal[idxChar][core.ER]-erStack == 0 {
			// Reset character stats
			stats.charProfilesCopy[idxChar] = char.Clone()
			stats.charSubstatFinal[idxChar][core.ER] -= erStack
			break
		}
	}

	return debug
}

func (stats *SubstatOptimizerDetails) setInitialSubstats(fixedSubstatCount float64) {
	stats.cloneStatsWithFixedAllocations(fixedSubstatCount)
	stats.calculateERBaseline()
}

// Copy to save initial character state with fixed allocations (2 of each substat)
func (stats *SubstatOptimizerDetails) cloneStatsWithFixedAllocations(fixedSubstatCount float64) {
	for i, char := range stats.simcfg.Characters.Profile {
		stats.charProfilesInitial[i] = char.Clone()
		for idxStat, stat := range stats.substatValues {
			if stat == 0 {
				continue
			}
			if core.StatType(idxStat) == core.ER {
				stats.charProfilesInitial[i].Stats[idxStat] += fixedSubstatCount * stat
			} else {
				stats.charProfilesInitial[i].Stats[idxStat] += fixedSubstatCount * stat * stats.charSubstatRarityMod[i]
			}
		}
	}
}

// Add some points into CR/CD to reduce crit variance and have reasonable baseline stats
// Also helps to slightly better evaluate the impact of favonius
// Current concern is that optimization on 2nd stage doesn't perform very well due to messed up rotation
func (stats *SubstatOptimizerDetails) calculateERBaseline() {
	for i, char := range stats.charProfilesInitial {
		stats.charProfilesERBaseline[i] = char.Clone()
		// Need special exception to Raiden due to her burst mechanics
		// TODO: Don't think there's a better solution without an expensive recursive solution to check across all Raiden ER states
		// Practically high ER substat Raiden is always currently unoptimal, so we just set her initial stacks low
		erStack := stats.charSubstatLimits[i][core.ER]
		if char.Base.Key == core.Raiden {
			erStack = 0
		}
		stats.charSubstatFinal[i][core.ER] = erStack

		stats.charProfilesERBaseline[i].Stats[core.ER] += float64(erStack) * stats.substatValues[core.ER]
		stats.charProfilesERBaseline[i].Stats[core.CR] += 4 * stats.substatValues[core.CR] * stats.charSubstatRarityMod[i]
		stats.charProfilesERBaseline[i].Stats[core.CD] += 4 * stats.substatValues[core.CD] * stats.charSubstatRarityMod[i]

		if strings.Contains(char.Weapon.Name, "favonius") {
			stats.calculateERBaselineHandleFav(i)
		}
	}
}

// Current strategy for favonius is to just boost this character's crit values a bit extra for optimal ER calculation purposes
// Then at next step of substat optimization, should naturally see relatively big DPS increases for that character if higher crit matters a lot
// TODO: Do we need a better special case for favonius?
func (stats *SubstatOptimizerDetails) calculateERBaselineHandleFav(i int) {
	stats.charProfilesERBaseline[i].Stats[core.CR] += 4 * stats.substatValues[core.CR] * stats.charSubstatRarityMod[i]
	stats.charWithFavonius[i] = true
}

func InitOptimStats(simcfg core.SimulationConfig, indivLiquidCap int, totalLiquidSubstats int) *SubstatOptimizerDetails {
	s := SubstatOptimizerDetails{}
	s.simcfg = simcfg
	s.indivSubstatLiquidCap = indivLiquidCap
	s.totalLiquidSubstats = totalLiquidSubstats

	// TODO: Will need to update this once artifact keys are introduced, and if more 4* artifact sets are implemented
	s.artifactSets4Star = []string{
		"exile",
		"instructor",
		"theexile",
	}

	s.substatValues = make([]float64, core.EndStatType)
	s.mainstatValues = make([]float64, core.EndStatType)

	// TODO: Is this actually the best way to set these values or am I missing something..?
	s.substatValues[core.ATKP] = 0.0496
	s.substatValues[core.CR] = 0.0331
	s.substatValues[core.CD] = 0.0662
	s.substatValues[core.EM] = 19.82
	s.substatValues[core.ER] = 0.0551
	s.substatValues[core.HPP] = 0.0496
	s.substatValues[core.DEFP] = 0.062
	s.substatValues[core.ATK] = 16.54
	s.substatValues[core.DEF] = 19.68
	s.substatValues[core.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	s.mainstatValues[core.ATKP] = 0.466
	s.mainstatValues[core.CR] = 0.311
	s.mainstatValues[core.CD] = 0.622
	s.mainstatValues[core.EM] = 186.5
	s.mainstatValues[core.ER] = 0.518
	s.mainstatValues[core.HPP] = 0.466
	s.mainstatValues[core.DEFP] = 0.583

	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	s.charRelevantSubstats = map[core.CharKey][]core.StatType{
		core.Albedo:  {core.DEFP},
		core.Hutao:   {core.HPP},
		core.Kokomi:  {core.HPP},
		core.Zhongli: {core.HPP},
		core.Itto:    {core.DEFP},
		core.Yunjin:  {core.DEFP},
		core.Noelle:  {core.DEFP},
		core.Gorou:   {core.DEFP},
	}

	// Final output array that holds [character][substat_count]
	s.charSubstatFinal = make([][]int, len(simcfg.Characters.Profile))
	for i := range simcfg.Characters.Profile {
		s.charSubstatFinal[i] = make([]int, core.EndStatType)
	}

	s.charSubstatLimits = make([][]int, len(simcfg.Characters.Profile))
	s.charSubstatRarityMod = make([]float64, len(simcfg.Characters.Profile))
	s.charProfilesInitial = make([]core.CharacterProfile, len(simcfg.Characters.Profile))

	// Need to make an exception in energy calcs for these characters for optimization purposes
	s.charWithFavonius = make([]bool, len(simcfg.Characters.Profile))
	// Give all characters max ER to set initial state
	s.charProfilesERBaseline = make([]core.CharacterProfile, len(simcfg.Characters.Profile))
	s.charProfilesCopy = make([]core.CharacterProfile, len(simcfg.Characters.Profile))

	return &s
}

// Obtain substat count limits based on main stats and also determine 4* set status
// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
// Most people will have 1 5* artifact which messes things up
// TODO: Check whether taking like an average of the two stat values is good enough?
func (stats *SubstatOptimizerDetails) setStatLimits() bool {
	profileIncludesFourStar := false

	for i, char := range stats.simcfg.Characters.Profile {
		stats.charSubstatLimits[i] = make([]int, core.EndStatType)
		for idxStat, stat := range stats.mainstatValues {
			if stat == 0 {
				continue
			}
			if char.Stats[idxStat] == 0 {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap
			} else {
				stats.charSubstatLimits[i][idxStat] = stats.indivSubstatLiquidCap - (2 * int(math.Round(char.Stats[idxStat]/stats.mainstatValues[idxStat])))
			}
		}

		// Display warning message for 4* sets
		stats.charSubstatRarityMod[i] = 1
		for set := range char.Sets {
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
			sb.WriteString(core.StatTypeString[i])
			sb.WriteString(": ")
			sb.WriteString(fmt.Sprintf("%v", v))
			sb.WriteString(" ")
		}
	}
	return strings.Trim(sb.String(), " ")
}
