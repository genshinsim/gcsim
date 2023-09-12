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
	simcfg                 *info.ActionList
	gcsl                   ast.Node
	simopt                 simulator.Options
	cfg                    string
	fixedSubstatCount      int
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
func (stats *SubstatOptimizerDetails) optimizeNonERSubstats() []string {
	var (
		opDebug   []string
		charDebug []string
	)

	stats.simcfg.Characters = stats.charProfilesCopy

	// Get initial DPS value
	initialResult, _ := simulator.RunWithConfig(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now())
	initialMean := *initialResult.Statistics.DPS.Mean

	opDebug = append(opDebug, "Calculating optimal substat distribution...")

	for idxChar := range stats.charProfilesCopy {
		charDebug = stats.optimizeNonErSubstatsForChar(idxChar, stats.charProfilesCopy[idxChar], initialMean)
		opDebug = append(opDebug, charDebug...)
	}

	return opDebug
}

func (stats *SubstatOptimizerDetails) optimizeNonErSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
	initialMean float64,
) []string {
	var opDebug []string
	opDebug = append(opDebug, fmt.Sprintf("%v", char.Base.Key))

	// Reset favonius char crit rate
	if stats.charWithFavonius[idxChar] {
		stats.charProfilesCopy[idxChar].Stats[attributes.CR] -= 8 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
	}

	relevantSubstats := stats.getNonErSubstatsToOptimizeForChar(char)

	addlSubstats := stats.charRelevantSubstats[char.Base.Key]
	if len(addlSubstats) > 0 {
		relevantSubstats = append(relevantSubstats, addlSubstats...)
	}

	substatGradients := stats.calculateSubstatGradientsForChar(idxChar, relevantSubstats, initialMean)

	allocDebug := stats.allocateSubstatGradientsForChar(idxChar, char, substatGradients, relevantSubstats)
	opDebug = append(opDebug, allocDebug...)

	return opDebug
}

func (stats *SubstatOptimizerDetails) allocateSubstatGradientsForChar(
	idxChar int,
	char info.CharacterProfile,
	substatGradient []float64,
	relevantSubstats []attributes.Stat,
) []string {
	var opDebug []string

	sorted := newSlice(substatGradient...)
	sort.Sort(sort.Reverse(sorted))

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
		stats.charProfilesCopy[idxChar].Stats[attributes.CR] += 8 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
	}
}

func (stats *SubstatOptimizerDetails) allocateSubstatGradientForChar(
	idxChar int,
	char info.CharacterProfile,
	sorted *Slice,
	idxGrad int,
	idxSubstat int,
	relevantSubstats []attributes.Stat,
) []string {
	var opDebug []string

	substatToMax := relevantSubstats[idxSubstat]

	// TODO: Improve this by adding a mix of CR/CD substats based on the ratio of gradient increase from CR/CD
	// If CR/CD is one of the selected substats, then adding them in a mix is generally most optimal
	// Use the ratio between gradient values to determine mix %
	// Need manual override here since gradient method from init does not always find this result
	var crCDSubstatRatio float64
	var gradStat float64
	switch substatToMax {
	case attributes.CR:
		gradCR := sorted.slice[idxGrad]
		gradCD := 0.0
		for i, idxSubstatTemp := range sorted.idx {
			if relevantSubstats[idxSubstatTemp] == attributes.CD {
				gradCD = sorted.slice[i]
			}
		}
		crCDSubstatRatio = gradCR / gradCD
	case attributes.CD:
		gradCD := sorted.slice[idxGrad]
		gradCR := 0.0
		for i, idxSubstatTemp := range sorted.idx {
			if relevantSubstats[idxSubstatTemp] == attributes.CR {
				gradCR = sorted.slice[i]
			}
		}
		crCDSubstatRatio = gradCR / gradCD
	default:
		gradStat = sorted.slice[idxGrad]
	}

	// If DPS change is really low, then it's usually better to just toss a few extra points into ER for stability
	if gradStat < 50 && crCDSubstatRatio == 0 {
		stats.assignSubstatsForChar(idxChar, char, attributes.ER, 4)
		opDebug = append(opDebug, "Low damage contribution from substats - adding some points to ER instead")
	}

	handleCRCD := func() {
		if crCDSubstatRatio <= 0 {
			stats.assignSubstatsForChar(idxChar, char, substatToMax, stats.indivSubstatLiquidCap+stats.fixedSubstatCount)
			return
		}

		globalLimit, crLimit := stats.assignSubstatsForChar(idxChar, char, attributes.CR, 0)
		_, cdLimit := stats.assignSubstatsForChar(idxChar, char, attributes.CD, 0)

		// Continually add CR/CD to try to align CR/CD ratio to ratio of gradients until we hit a limit
		var currentRatio float64
		var amtCR int
		var amtCD int
		currentStat := attributes.CR
		// Debug to avoid runaway loops...
		var iteration int
		// Want this to continue until either global cap is reached, or we can neither add CR/CD
		for globalLimit > 0 && (crLimit > 0 || cdLimit > 0) && iteration < 100 {
			if stats.charSubstatFinal[idxChar][attributes.CD] == 0 {
				currentRatio = float64(stats.charSubstatFinal[idxChar][attributes.CR])
			} else {
				currentRatio = float64(stats.charSubstatFinal[idxChar][attributes.CR]) / float64(stats.charSubstatFinal[idxChar][attributes.CD])
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
				amtCD = stats.indivSubstatLiquidCap
			}
			if cdLimit == 0 {
				amtCR = stats.indivSubstatLiquidCap
			}

			if currentStat == attributes.CR {
				globalLimit, crLimit = stats.assignSubstatsForChar(idxChar, char, attributes.CR, amtCR)
				currentStat = attributes.CD
			} else if currentStat == attributes.CD {
				globalLimit, cdLimit = stats.assignSubstatsForChar(idxChar, char, attributes.CD, amtCD)
				currentStat = attributes.CR
			}
			iteration += 1
		}
	}

	handleCRCD()
	return opDebug
}

// Assigns substats and returns the remaining global limit and individual substat limit
func (stats *SubstatOptimizerDetails) assignSubstatsForChar(
	idxChar int,
	char info.CharacterProfile,
	substat attributes.Stat,
	amt int,
) (int, int) {
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
	amtToAdd := minInt(
		stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat],
		remainingLiquidSubstats,
		amt,
	)
	stats.charSubstatFinal[idxChar][substat] += amtToAdd

	return remainingLiquidSubstats - amtToAdd, stats.charSubstatLimits[idxChar][substat] - stats.charSubstatFinal[idxChar][substat]
}

func (stats *SubstatOptimizerDetails) calculateSubstatGradientsForChar(
	idxChar int,
	relevantSubstats []attributes.Stat,
	initialMean float64,
) []float64 {
	substatGradients := make([]float64, len(relevantSubstats))

	// Build "gradient" by substat
	for idxSubstat, substat := range relevantSubstats {
		stats.charProfilesCopy[idxChar].Stats[substat] += 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

		stats.simcfg.Characters = stats.charProfilesCopy
		substatEvalResult, _ := simulator.RunWithConfig(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now())

		substatGradients[idxSubstat] = *substatEvalResult.Statistics.DPS.Mean - initialMean

		// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
		// might give them "too much" cr (= max out liquid cr subs) but that's probably not a big deal
		if stats.charWithFavonius[idxChar] && substat == attributes.CR {
			substatGradients[idxSubstat] += 1000
		}
		stats.charProfilesCopy[idxChar].Stats[substat] -= 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
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
// For each character, do grid search to find optimal ER values
// TODO: Can maybe replace with some kind of gradient descent for speed improvements/allow for 1 ER substat moves?
// When I tried before, it was hard to define a good step size and penalty on high ER substats that generally worked well
// At least this version works semi-reliably...
func (stats *SubstatOptimizerDetails) optimizeERSubstats(tolMean, tolSD float64) []string {
	var opDebug []string

	for idxChar := range stats.charProfilesERBaseline {
		stats.findOptimalERforChar(idxChar, stats.charProfilesERBaseline[idxChar], tolMean, tolSD)
	}

	// Need a separate optimization routine for strong battery characters (currently Raiden only, maybe EMC?)
	// Need to set all other character's ER substats at final value, then see added benefit from ER for global battery chars
	for i := range stats.charProfilesERBaseline {
		stats.charProfilesERBaseline[i].Stats[attributes.ER] = stats.charProfilesInitial[i].Stats[attributes.ER]

		if stats.charProfilesERBaseline[i].Base.Key == keys.Raiden {
			stats.charSubstatFinal[i][attributes.ER] = stats.indivSubstatLiquidCap
		}

		stats.charProfilesERBaseline[i].Stats[attributes.ER] += float64(
			stats.charSubstatFinal[i][attributes.ER],
		) * stats.substatValues[attributes.ER]
	}

	for i := range stats.charProfilesERBaseline {
		if stats.charProfilesERBaseline[i].Base.Key != keys.Raiden {
			continue
		}
		opDebug = append(opDebug, "Raiden found in team comp - running secondary optimization routine...")
		stats.findOptimalERforChar(i, stats.charProfilesERBaseline[i], tolMean, tolSD)
	}

	// Fix ER at previously found values then optimize all other substats
	opDebug = append(opDebug, "Optimized ER Liquid Substats by character:")
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

func (stats *SubstatOptimizerDetails) findOptimalERforChar(
	idxChar int,
	char info.CharacterProfile,
	tolMean float64,
	tolSD float64,
) {
	var initialMean float64
	var initialSD float64

	for erStack := 0; erStack <= stats.indivSubstatLiquidCap; erStack += 2 {
		stats.charProfilesCopy[idxChar] = char.Clone()
		stats.charProfilesCopy[idxChar].Stats[attributes.ER] -= float64(erStack) * stats.substatValues[attributes.ER]

		stats.simcfg.Characters = stats.charProfilesCopy

		result, _ := simulator.RunWithConfig(context.TODO(), stats.cfg, stats.simcfg, stats.gcsl, stats.simopt, time.Now())

		if erStack == 0 {
			initialMean = *result.Statistics.DPS.Mean
			initialSD = *result.Statistics.DPS.SD
		}

		condition := *result.Statistics.DPS.Mean/initialMean-1 < -tolMean || *result.Statistics.DPS.SD/initialSD-1 > tolSD
		// For Raiden, we can't use DPS directly as a measure since she scales off of her own ER
		// Instead we ONLY use the SD tolerance as big jumps indicate the rotation is becoming more unstable
		if char.Base.Key == keys.Raiden {
			condition = *result.Statistics.DPS.SD/initialSD-1 > tolSD
		}

		// If differences exceed tolerances, then immediately break
		if condition {
			// Reset character stats
			stats.charProfilesCopy[idxChar] = char.Clone()
			// Save ER value - optimal value is the value immediately prior, so we subtract 2
			stats.charSubstatFinal[idxChar][attributes.ER] -= erStack - 2
			break
		}

		// Reached minimum possible ER stacks, so optimal is the minimum amount of ER stacks
		if stats.charSubstatFinal[idxChar][attributes.ER]-erStack == 0 {
			// Reset character stats
			stats.charProfilesCopy[idxChar] = char.Clone()
			stats.charSubstatFinal[idxChar][attributes.ER] -= erStack
			break
		}
	}
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
			erStack = 0
		}
		stats.charSubstatFinal[i][attributes.ER] = erStack

		stats.charProfilesERBaseline[i].Stats[attributes.ER] += float64(erStack) * stats.substatValues[attributes.ER]
		stats.charProfilesERBaseline[i].Stats[attributes.CR] += 4 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[i]
		stats.charProfilesERBaseline[i].Stats[attributes.CD] += 4 * stats.substatValues[attributes.CD] * stats.charSubstatRarityMod[i]

		if strings.Contains(stats.charProfilesInitial[i].Weapon.Name, "favonius") {
			stats.calculateERBaselineHandleFav(i)
		}
	}
}

// Current strategy for favonius is to just boost this character's crit values a bit extra for optimal ER calculation purposes
// Then at next step of substat optimization, should naturally see relatively big DPS increases for that character if higher crit matters a lot
// TODO: Do we need a better special case for favonius?
func (stats *SubstatOptimizerDetails) calculateERBaselineHandleFav(i int) {
	stats.charProfilesERBaseline[i].Stats[attributes.CR] += 4 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[i]
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
		keys.Albedo:  {attributes.DEFP},
		keys.Hutao:   {attributes.HPP},
		keys.Kokomi:  {attributes.HPP},
		keys.Zhongli: {attributes.HPP},
		keys.Itto:    {attributes.DEFP},
		keys.Yunjin:  {attributes.DEFP},
		keys.Noelle:  {attributes.DEFP},
		keys.Gorou:   {attributes.DEFP},
		keys.Yelan:   {attributes.HPP},
		keys.Candace: {attributes.HPP},
		keys.Nilou:   {attributes.HPP},
		keys.Layla:   {attributes.HPP},
	}

	// Final output array that holds [character][substat_count]
	s.charSubstatFinal = make([][]int, len(simcfg.Characters))
	for i := range simcfg.Characters {
		s.charSubstatFinal[i] = make([]int, attributes.EndStatType)
	}

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
