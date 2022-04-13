package substatoptimizer

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/shortcut"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

// Additional runtime option to optimize substats according to KQM standards
func RunSubstatOptim(simopt simulator.Options, verbose bool, additionalOptions string) {
	// Substat Optimization strategy is very simplistic right now:
	// This is not fully optimal - see other comments in code
	// 1) User sets team, weapons, artifact sets/main stats, and rotation
	// 2) Given those, for each character, sim picks ER substat value that functionally maximizes DPS Mean/SD,
	// subject to a penalty on high ER values
	//    - Strategy is to just do a dumb grid search over ER substat values for each character
	//    - ER substat values are set in increments of 2 to make the search easier
	// 3) Given ER values, we then optimize the other substats by doing a "gradient descent" (but not really) method
	sugarLog := InitLogger(verbose)

	// Each optimizer run should not be saving anything out for the GZIP
	simopt.GZIPResult = false

	// Parse config
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	re := InitRegex()
	srcCleaned := re.scrubSimCfg(cfg, sugarLog)

	parser := ast.New(srcCleaned)
	simcfg, err := parser.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	optionsMap := map[string]float64{
		"total_liquid_substats": 20,
		"indiv_liquid_cap":      10,
		"fixed_substats_count":  2,
		"sim_iter":              350,
		"tol_mean":              0.015,
		"tol_sd":                0.33,
	}

	// Parse and set all special sim options
	if additionalOptions != "" {
		re.scrubAdditionalOptimizerCfg(additionalOptions, optionsMap, sugarLog)
	}

	// Fix iterations at 350 for performance
	// TODO: Seems to be a roughly good number at KQM standards
	simcfg.Settings.Iterations = int(optionsMap["sim_iter"])

	stats := InitOptimStats(simcfg, int(optionsMap["indiv_liquid_cap"]), int(optionsMap["fixed_substats_count"]))
	fourStarFound := stats.setStatLimits()
	if fourStarFound {
		sugarLog.Warn("Warning: 4* artifact set detected. Optimizer currently assumes that ER substats take 5* values, and all other substats take 4* values.")
	}

	// Copy to save initial character state with fixed allocations (2 of each substat)
	charProfilesInitial := make([]character.CharacterProfile, len(simcfg.Characters))

	for i, char := range simcfg.Characters {
		charProfilesInitial[i] = char.Clone()
		for idxStat, stat := range stats.substatValues {
			if stat == 0 {
				continue
			}
			if attributes.Stat(idxStat) == attributes.ER {
				charProfilesInitial[i].Stats[idxStat] += float64(stats.fixedSubstatCount) * stat
			} else {
				charProfilesInitial[i].Stats[idxStat] += float64(stats.fixedSubstatCount) * stat * stats.charSubstatRarityMod[i]
			}
		}
	}

	// Need to special case these characters for optimization purposes
	charWithFavonius := make([]bool, len(simcfg.Characters))

	// Give all characters max ER to set initial state
	charProfilesERBaseline := make([]character.CharacterProfile, len(simcfg.Characters))

	sugarLog.Info("Starting ER Optimization...")

	// Add some points into CR/CD to reduce crit variance and have reasonable baseline stats
	// Also helps to slightly better evaluate the impact of favonius
	// TODO: Do we need a better special case for favonius?
	// Current concern is that optimization on 2nd stage doesn't perform very well due to messed up rotation
	for i, char := range charProfilesInitial {
		charProfilesERBaseline[i] = char.Clone()
		// Need special exception to Raiden due to her burst mechanics
		// TODO: Don't think there's a better solution without an expensive recursive solution to check across all Raiden ER states
		// Practically high ER substat Raiden is always currently unoptimal, so we just set her initial stacks low
		erStack := stats.charSubstatLimits[i][attributes.ER]
		if char.Base.Key == keys.Raiden {
			erStack = 0
		}
		stats.charSubstatFinal[i][attributes.ER] = erStack

		charProfilesERBaseline[i].Stats[attributes.ER] += float64(erStack) * stats.substatValues[attributes.ER]
		charProfilesERBaseline[i].Stats[attributes.CR] += 4 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[i]
		charProfilesERBaseline[i].Stats[attributes.CD] += 4 * stats.substatValues[attributes.CD] * stats.charSubstatRarityMod[i]

		// Current strategy for favonius is to just boost this character's crit values a bit extra for optimal ER calculation purposes
		// Then at next step of substat optimization, should naturally see relatively big DPS increases for that character if higher crit matters a lot
		if strings.Contains(char.Weapon.Name, "favonius") {
			charProfilesERBaseline[i].Stats[attributes.CR] += 4 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[i]
			charWithFavonius[i] = true
		}
	}

	// Find optimal ER cutoffs for each character
	// For each character, do grid search to find optimal ER values
	// TODO: Can maybe replace with some kind of gradient descent for speed improvements/allow for 1 ER substat moves?
	// When I tried before, it was hard to define a good step size and penalty on high ER substats that generally worked well
	// At least this version works semi-reliably...
	charProfilesCopy := make([]character.CharacterProfile, len(simcfg.Characters))
	for i, char := range charProfilesERBaseline {
		charProfilesCopy[i] = char.Clone()
	}

	tolMean := optionsMap["tol_mean"]
	tolSD := optionsMap["tol_sd"]
	// Interior loop of the ER optimization - takes in the parsed character index and a character profile
	// No direct output - changes state of local variables instead
	findOptimalERforChar := func(idxChar int, char character.CharacterProfile) {
		var initialMean float64
		var initialSD float64
		sugarLog.Debugf("%v", char.Base.Key)
		for erStack := 0; erStack <= stats.indivSubstatLiquidCap; erStack += 2 {
			charProfilesCopy[idxChar] = char.Clone()
			charProfilesCopy[idxChar].Stats[attributes.ER] -= float64(erStack) * stats.substatValues[attributes.ER]

			simcfg.Characters = charProfilesCopy

			result := runSimWithConfig(cfg, simcfg, simopt)
			sugarLog.Debugf("%v: %v (%v)", stats.charSubstatFinal[idxChar][attributes.ER]-erStack, result.DPS.Mean, result.DPS.SD)

			if erStack == 0 {
				initialMean = result.DPS.Mean
				initialSD = result.DPS.SD
			}

			condition := result.DPS.Mean/initialMean-1 < -tolMean || result.DPS.SD/initialSD-1 > tolSD
			// For Raiden, we can't use DPS directly as a measure since she scales off of her own ER
			// Instead we ONLY use the SD tolerance as big jumps indicate the rotation is becoming more unstable
			if char.Base.Key == keys.Raiden {
				condition = result.DPS.SD/initialSD-1 > tolSD
			}

			// If differences exceed tolerances, then immediately break
			if condition {
				// Reset character stats
				charProfilesCopy[idxChar] = char.Clone()
				// Save ER value - optimal value is the value immediately prior, so we subtract 2
				stats.charSubstatFinal[idxChar][attributes.ER] -= (erStack - 2)
				break
			}

			// Reached minimum possible ER stacks, so optimal is the minimum amount of ER stacks
			if stats.charSubstatFinal[idxChar][attributes.ER]-erStack == 0 {
				// Reset character stats
				charProfilesCopy[idxChar] = char.Clone()
				stats.charSubstatFinal[idxChar][attributes.ER] -= erStack
				break
			}
		}
	}

	// Tolerance cutoffs for mean and SD from initial state
	// Initial state is used rather than checking across each iteration due to noise
	// TODO: May want to adjust further?
	for idxChar, char := range charProfilesERBaseline {
		findOptimalERforChar(idxChar, char)
	}

	// Need a separate optimization routine for strong battery characters (currently Raiden only, maybe EMC?)
	// Need to set all other character's ER substats at final value, then see added benefit from ER for global battery chars
	for i, char := range charProfilesERBaseline {
		charProfilesERBaseline[i].Stats[attributes.ER] = charProfilesInitial[i].Stats[attributes.ER]

		if char.Base.Key == keys.Raiden {
			stats.charSubstatFinal[i][attributes.ER] = stats.indivSubstatLiquidCap
		}

		charProfilesERBaseline[i].Stats[attributes.ER] += float64(stats.charSubstatFinal[i][attributes.ER]) * stats.substatValues[attributes.ER]
	}
	for idxChar, char := range charProfilesERBaseline {
		if char.Base.Key != keys.Raiden {
			continue
		}
		sugarLog.Info("Raiden found in team comp - running secondary optimization routine...")
		findOptimalERforChar(idxChar, char)
	}

	// Fix ER at previously found values then optimize all other substats
	sugarLog.Info("Optimized ER Liquid Substats by character:")
	printVal := ""
	for i, char := range charProfilesInitial {
		printVal += fmt.Sprintf("%v: %.4g, ", char.Base.Key.String(), float64(stats.charSubstatFinal[i][attributes.ER])*stats.substatValues[attributes.ER])
	}
	sugarLog.Info(printVal)

	// Calculate per-character per-substat "gradients" at initial state using finite differences
	// In practical evaluations, adding small numbers of substats (<10) can be VERY noisy
	// Therefore, "gradient" evaluations are done in groups of 10 substats
	// Allocation strategy is to just max substats according to highest gradient to lowest
	// TODO: Probably want to refactor to potentially run gradient step at least twice:
	// once initially then another at 10 assigned liquid substats
	// Fine grained evaluations are too expensive time wise, but can perhaps add in an option for people who want to sit around for a while
	sugarLog.Info("Calculating optimal substat distribution...")

	// Get initial DPS value
	simcfg.Characters = charProfilesCopy
	initialResult := runSimWithConfig(cfg, simcfg, simopt)
	initialMean := initialResult.DPS.Mean
	sugarLog.Debug(initialMean)

	for idxChar, char := range charProfilesCopy {
		sugarLog.Info(char.Base.Key)

		// Reset favonius char crit rate
		if charWithFavonius[idxChar] {
			charProfilesCopy[idxChar].Stats[attributes.CR] -= 8 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
		}

		// Get relevant substats, and add additional ones for special characters if needed
		relevantSubstats := []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD, attributes.EM}
		// RIP crystallize...
		if keys.CharKeyToEle[char.Base.Key] == attributes.Geo {
			relevantSubstats = []attributes.Stat{attributes.ATKP, attributes.CR, attributes.CD}
		}

		addlSubstats := stats.charRelevantSubstats[char.Base.Key]
		if len(addlSubstats) > 0 {
			relevantSubstats = append(relevantSubstats, addlSubstats...)
		}

		substatGradients := make([]float64, len(relevantSubstats))

		// Build "gradient" by substat
		for idxSubstat, substat := range relevantSubstats {
			charProfilesCopy[idxChar].Stats[substat] += 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]

			simcfg.Characters = charProfilesCopy
			substatEvalResult := runSimWithConfig(cfg, simcfg, simopt)
			// sugarLog.Debugf("%v: %v (%v)", substat.String(), substatEvalResult.DPS.Mean, substatEvalResult.DPS.SD)

			substatGradients[idxSubstat] = substatEvalResult.DPS.Mean - initialMean

			// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
			// might give them "too much" cr (= max out liquid cr subs) but that's probably not a big deal
			if charWithFavonius[idxChar] && substat == attributes.CR {
				substatGradients[idxSubstat] += 1000
			}

			charProfilesCopy[idxChar].Stats[substat] -= 10 * stats.substatValues[substat] * stats.charSubstatRarityMod[idxChar]
		}

		// Allocate substats
		sorted := NewSlice(substatGradients...)
		sort.Sort(sort.Reverse(sorted))

		printVal = ""
		for i, idxSorted := range sorted.idx {
			printVal += fmt.Sprintf("%v: %5.5g, ", relevantSubstats[idxSorted], sorted.slice[i])
		}
		sugarLog.Debug(printVal)

		// Assigns substats and returns the remaining global limit and individual substat limit
		assignSubstats := func(substat attributes.Stat, amt int) (int, int) {
			totalSubstatCount := 0
			for _, val := range stats.charSubstatFinal[idxChar] {
				totalSubstatCount += val
			}

			baseLiquidSubstats := int(optionsMap["total_liquid_substats"])
			for set, count := range char.Sets {
				for _, setfourstar := range stats.artifactSets4Star {
					if set == setfourstar {
						baseLiquidSubstats -= 2 * count
					}
				}
			}

			remainingLiquidSubstats := baseLiquidSubstats - totalSubstatCount
			// Minimum of individual limit, global limit, desired amount
			amtToAdd := minInt(stats.charSubstatLimits[idxChar][substat]-stats.charSubstatFinal[idxChar][substat], remainingLiquidSubstats, amt)
			stats.charSubstatFinal[idxChar][substat] += amtToAdd

			return remainingLiquidSubstats - amtToAdd, stats.charSubstatLimits[idxChar][substat] - stats.charSubstatFinal[idxChar][substat]
		}

		for idxGrad, idxSubstat := range sorted.idx {
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
				crCDSubstatRatio = (gradCR / gradCD)
			case attributes.CD:
				gradCD := sorted.slice[idxGrad]
				gradCR := 0.0
				for i, idxSubstatTemp := range sorted.idx {
					if relevantSubstats[idxSubstatTemp] == attributes.CR {
						gradCR = sorted.slice[i]
					}
				}
				crCDSubstatRatio = (gradCR / gradCD)
			default:
				gradStat = sorted.slice[idxGrad]
			}

			// If DPS change is really low, then it's usually better to just toss a few extra points into ER for stability
			if gradStat < 50 && crCDSubstatRatio == 0 {
				assignSubstats(attributes.ER, 4)
				sugarLog.Info("Low damage contribution from substats - adding some points to ER instead")
			}

			var globalLimit int
			var crLimit int
			var cdLimit int
			if crCDSubstatRatio > 0 {
				globalLimit, crLimit = assignSubstats(attributes.CR, 0)
				_, cdLimit = assignSubstats(attributes.CD, 0)

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
						globalLimit, crLimit = assignSubstats(attributes.CR, amtCR)
						currentStat = attributes.CD
					} else if currentStat == attributes.CD {
						globalLimit, cdLimit = assignSubstats(attributes.CD, amtCD)
						currentStat = attributes.CR
					}
					iteration += 1
				}
			} else {
				globalLimit, _ = assignSubstats(substatToMax, stats.indivSubstatLiquidCap + stats.fixedSubstatCount)
			}
			if globalLimit == 0 {
				break
			}
		}

		sugarLog.Info("Final Liquid Substat Counts: ", PrettyPrintStatsCounts(stats.charSubstatFinal[idxChar]))

		// Reset favonius char crit rate... again
		if charWithFavonius[idxChar] {
			charProfilesCopy[idxChar].Stats[attributes.CR] += 8 * stats.substatValues[attributes.CR] * stats.charSubstatRarityMod[idxChar]
		}
	}

	sugarLog.Info("Final config substat strings:")
	// Final output
	// This doesn't take much time relatively speaking, so just always do the processing...
	output := srcCleaned
	charNames := make(map[keys.Char]string)
	for _, match := range re.GetCharNames.FindAllStringSubmatch(output, -1) {
		charKey := shortcut.CharNameToKey[match[1]]
		charNames[charKey] = match[1]
	}

	for idxChar, char := range charProfilesInitial {
		finalString := fmt.Sprintf("%v add stats", charNames[char.Base.Key])

		for idxSubstat, value := range stats.substatValues {
			if value <= 0 {
				continue
			}
			finalString += fmt.Sprintf(" %v=%.6g", attributes.StatTypeString[idxSubstat], value*float64(stats.fixedSubstatCount+stats.charSubstatFinal[idxChar][idxSubstat]))
		}

		fmt.Println(finalString + ";")

		reInsertLocation := regexp.MustCompile(fmt.Sprintf(`(?m)^(%v\s+add\s+stats\b.*)$`, charNames[char.Base.Key]))
		output = reInsertLocation.ReplaceAllString(output, fmt.Sprintf("$1\n%v;", finalString))
	}

	// Sticks optimized substat string into config and output
	if simopt.ResultSaveToPath != "" {
		output = strings.TrimSpace(output) + "\n"
		//try creating file to write to
		err := os.WriteFile(simopt.ResultSaveToPath, []byte(output), 0644)
		if err != nil {
			log.Panic(err)
		}
		sugarLog.Infof("Saved to the following location: %v", simopt.ResultSaveToPath)
	}
}

// Just runs the sim with specified settings
func runSimWithConfig(cfg string, simcfg *ast.ActionList, simopt simulator.Options) result.Summary {
	result, err := simulator.RunWithConfig(cfg, simcfg, simopt)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return result
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

// Gets the minimum of a slice of integers
func minInt(vars ...int) int {
	min := vars[0]

	for _, val := range vars {
		if min > val {
			min = val
		}
	}
	return min
}

// Thin wrapper around sort Slice to retrieve the sorted indices as well
type Slice struct {
	slice sort.Float64Slice
	idx   []int
}

func (s Slice) Len() int {
	return len(s.slice)
}

func (s Slice) Less(i, j int) bool {
	return s.slice[i] < s.slice[j]
}

func (s Slice) Swap(i, j int) {
	s.slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n ...float64) *Slice {
	s := &Slice{
		slice: sort.Float64Slice(n),
		idx:   make([]int, len(n)),
	}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
