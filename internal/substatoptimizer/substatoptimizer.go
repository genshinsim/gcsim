package substatoptimizer

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"github.com/genshinsim/gcsim/pkg/result"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Only includes damage related substats scaling. Ignores things like HP for Barbara
	charRelevantSubstats = map[core.CharKey][]core.StatType{
		core.Albedo:  {core.DEFP},
		core.Hutao:   {core.HPP},
		core.Kokomi:  {core.HPP},
		core.Zhongli: {core.HPP},
		core.Itto:    {core.DEFP},
		core.Yunjin:  {core.DEFP},
		core.Noelle:  {core.DEFP},
		core.Gorou:   {core.DEFP},
	}

	// TODO: Will need to update this once artifact keys are introduced, and if more 4* artifact sets are implemented
	artifactSets4Star = []string{
		"exile",
		"instructor",
		"theexile",
	}

	substatValues  = make([]float64, core.EndStatType)
	mainstatValues = make([]float64, core.EndStatType)
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

	// TODO: Is this actually the best way to set these values or am I missing something..?
	substatValues[core.ATKP] = 0.0496
	substatValues[core.CR] = 0.0331
	substatValues[core.CD] = 0.0662
	substatValues[core.EM] = 19.82
	substatValues[core.ER] = 0.0551
	substatValues[core.HPP] = 0.0496
	substatValues[core.DEFP] = 0.062
	substatValues[core.ATK] = 16.54
	substatValues[core.DEF] = 19.68
	substatValues[core.HP] = 253.94

	// Used to try to back out artifact main stats for limits
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	mainstatValues[core.ATKP] = 0.466
	mainstatValues[core.CR] = 0.311
	mainstatValues[core.CD] = 0.622
	mainstatValues[core.EM] = 186.5
	mainstatValues[core.ER] = 0.518
	mainstatValues[core.HPP] = 0.466
	mainstatValues[core.DEFP] = 0.583

	// Each optimizer run should not be saving anything out for the GZIP
	simopt.GZIPResult = false

	// Start logger
	zapcfg := zap.NewDevelopmentConfig()
	zapcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapcfg.EncoderConfig.CallerKey = ""
	zapcfg.EncoderConfig.StacktraceKey = ""
	zapcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if verbose {
		zapcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	logger, _ := zapcfg.Build()
	defer logger.Sync()
	sugarLog := logger.Sugar()

	// Parse config
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	// Regex to identify main stats based on flower. Check that characters all have one that we can recognize
	var reMainstats = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
	var reGetCharNames = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	if len(reMainstats.FindAllString(cfg, -1)) != len(reGetCharNames.FindAllString(cfg, -1)) {
		sugarLog.Error("Error: Could not identify valid main artifact stat rows for all characters based on flower HP values.")
		sugarLog.Error("5* flowers must have 4780 HP, and 4* flowers must have 3571 HP.")
		os.Exit(1)
	}

	// Regex to remove stat rows that do not look like mainstat rows from the config
	var reSubstats = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\b[^;]*;.*\n`)
	srcCleaned := string(cfg)
	errorPrinted := false
	for _, match := range reSubstats.FindAllString(cfg, -1) {
		if reMainstats.MatchString(string(match)) {
			continue
		}
		if !errorPrinted {
			sugarLog.Warn("Warning: Config found to have existing substat information. Ignoring...")
			errorPrinted = true
		}
		srcCleaned = strings.Replace(srcCleaned, string(match), "", -1)
	}

	parser := parse.New("single", srcCleaned)
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
		reOptions := regexp.MustCompile(`([a-z_]+)=([0-9.]+)`)
		parsedOptions := reOptions.FindAllStringSubmatch(additionalOptions, -1)
		for _, val := range parsedOptions {
			if _, ok := optionsMap[val[1]]; ok {
				optionsMap[val[1]], _ = strconv.ParseFloat(val[2], 64)
			} else {
				sugarLog.Panic("Invalid substat optimization option found: %v", val[1], val[2])
			}
		}
	}

	// Fix iterations at 350 for performance
	// TODO: Seems to be a roughly good number at KQM standards
	simcfg.Settings.Iterations = int(optionsMap["sim_iter"])

	// Final output array that holds [character][substat_count]
	charSubstatFinal := make([][]int, len(simcfg.Characters.Profile))
	for i := range simcfg.Characters.Profile {
		charSubstatFinal[i] = make([]int, core.EndStatType)
	}

	// Obtain substat count limits based on main stats and also determine 4* set status
	// TODO: Not sure how to handle 4* artifact sets... Config can't really identify these instances easily
	// Most people will have 1 5* artifact which messes things up
	// TODO: Check whether taking like an average of the two stat values is good enough?
	indivSubstatLiquidCap := int(optionsMap["indiv_liquid_cap"])
	charSubstatLimits := make([][]int, len(simcfg.Characters.Profile))
	charSubstatRarityMod := make([]float64, len(simcfg.Characters.Profile))
	for i, char := range simcfg.Characters.Profile {
		charSubstatLimits[i] = make([]int, core.EndStatType)
		for idxStat, stat := range mainstatValues {
			if stat == 0 {
				continue
			}
			if char.Stats[idxStat] == 0 {
				charSubstatLimits[i][idxStat] = indivSubstatLiquidCap
			} else {
				charSubstatLimits[i][idxStat] = indivSubstatLiquidCap - (2 * int(math.Round(char.Stats[idxStat]/mainstatValues[idxStat])))
			}
		}

		// Display warning message for 4* sets
		charSubstatRarityMod[i] = 1
		for set := range char.Sets {
			for _, fourStar := range artifactSets4Star {
				if set == fourStar {
					sugarLog.Warn("Warning: 4* artifact set detected. Optimizer currently assumes that ER substats take 5* values, and all other substats take 4* values.")
					charSubstatRarityMod[i] = 0.8
				}
			}
		}
	}

	// Copy to save initial character state with fixed allocations (2 of each substat)
	charProfilesInitial := make([]core.CharacterProfile, len(simcfg.Characters.Profile))

	fixedSubstatCount := optionsMap["fixed_substats_count"]
	for i, char := range simcfg.Characters.Profile {
		charProfilesInitial[i] = char.Clone()
		for idxStat, stat := range substatValues {
			if stat == 0 {
				continue
			}
			if core.StatType(idxStat) == core.ER {
				charProfilesInitial[i].Stats[idxStat] += fixedSubstatCount * stat
			} else {
				charProfilesInitial[i].Stats[idxStat] += fixedSubstatCount * stat * charSubstatRarityMod[i]
			}
		}
	}

	// Need to special case these characters for optimization purposes
	charWithFavonius := make([]bool, len(simcfg.Characters.Profile))

	// Give all characters max ER to set initial state
	charProfilesERBaseline := make([]core.CharacterProfile, len(simcfg.Characters.Profile))

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
		erStack := charSubstatLimits[i][core.ER]
		if char.Base.Key == core.Raiden {
			erStack = 0
		}
		charSubstatFinal[i][core.ER] = erStack

		charProfilesERBaseline[i].Stats[core.ER] += float64(erStack) * substatValues[core.ER]
		charProfilesERBaseline[i].Stats[core.CR] += 4 * substatValues[core.CR] * charSubstatRarityMod[i]
		charProfilesERBaseline[i].Stats[core.CD] += 4 * substatValues[core.CD] * charSubstatRarityMod[i]

		// Current strategy for favonius is to just boost this character's crit values a bit extra for optimal ER calculation purposes
		// Then at next step of substat optimization, should naturally see relatively big DPS increases for that character if higher crit matters a lot
		if strings.Contains(char.Weapon.Name, "favonius") {
			charProfilesERBaseline[i].Stats[core.CR] += 4 * substatValues[core.CR] * charSubstatRarityMod[i]
			charWithFavonius[i] = true
		}
	}

	// Find optimal ER cutoffs for each character
	// For each character, do grid search to find optimal ER values
	// TODO: Can maybe replace with some kind of gradient descent for speed improvements/allow for 1 ER substat moves?
	// When I tried before, it was hard to define a good step size and penalty on high ER substats that generally worked well
	// At least this version works semi-reliably...
	charProfilesCopy := make([]core.CharacterProfile, len(simcfg.Characters.Profile))
	for i, char := range charProfilesERBaseline {
		charProfilesCopy[i] = char.Clone()
	}

	tolMean := optionsMap["tol_mean"]
	tolSD := optionsMap["tol_sd"]
	// Interior loop of the ER optimization - takes in the parsed character index and a character profile
	// No direct output - changes state of local variables instead
	findOptimalERforChar := func(idxChar int, char core.CharacterProfile) {
		var initialMean float64
		var initialSD float64
		sugarLog.Debugf("%v", char.Base.Key)
		for erStack := 0; erStack <= 10; erStack += 2 {
			charProfilesCopy[idxChar] = char.Clone()
			charProfilesCopy[idxChar].Stats[core.ER] -= float64(erStack) * substatValues[core.ER]

			simcfg.Characters.Profile = charProfilesCopy

			result := runSimWithConfig(cfg, simcfg, simopt)
			sugarLog.Debugf("%v: %v (%v)", charSubstatFinal[idxChar][core.ER]-erStack, result.DPS.Mean, result.DPS.SD)

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
				charProfilesCopy[idxChar] = char.Clone()
				// Save ER value - optimal value is the value immediately prior, so we subtract 2
				charSubstatFinal[idxChar][core.ER] -= (erStack - 2)
				break
			}

			// Reached minimum possible ER stacks, so optimal is the minimum amount of ER stacks
			if charSubstatFinal[idxChar][core.ER]-erStack == 0 {
				// Reset character stats
				charProfilesCopy[idxChar] = char.Clone()
				charSubstatFinal[idxChar][core.ER] -= erStack
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
		charProfilesERBaseline[i].Stats[core.ER] = charProfilesInitial[i].Stats[core.ER]

		if char.Base.Key == core.Raiden {
			charSubstatFinal[i][core.ER] = 10
		}

		charProfilesERBaseline[i].Stats[core.ER] += float64(charSubstatFinal[i][core.ER]) * substatValues[core.ER]
	}
	for idxChar, char := range charProfilesERBaseline {
		if char.Base.Key != core.Raiden {
			continue
		}
		sugarLog.Info("Raiden found in team comp - running secondary optimization routine...")
		findOptimalERforChar(idxChar, char)
	}

	// Fix ER at previously found values then optimize all other substats
	sugarLog.Info("Optimized ER Liquid Substats by character:")
	printVal := ""
	for i, char := range charProfilesInitial {
		printVal += fmt.Sprintf("%v: %.4g, ", char.Base.Key.String(), float64(charSubstatFinal[i][core.ER])*substatValues[core.ER])
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
	simcfg.Characters.Profile = charProfilesCopy
	initialResult := runSimWithConfig(cfg, simcfg, simopt)
	initialMean := initialResult.DPS.Mean
	sugarLog.Debug(initialMean)

	for idxChar, char := range charProfilesCopy {
		sugarLog.Info(char.Base.Key)

		// Reset favonius char crit rate
		if charWithFavonius[idxChar] {
			charProfilesCopy[idxChar].Stats[core.CR] -= 8 * substatValues[core.CR] * charSubstatRarityMod[idxChar]
		}

		// Get relevant substats, and add additional ones for special characters if needed
		relevantSubstats := []core.StatType{core.ATKP, core.CR, core.CD, core.EM}
		// RIP crystallize...
		if core.CharKeyToEle[char.Base.Key] == core.Geo {
			relevantSubstats = []core.StatType{core.ATKP, core.CR, core.CD}
		}

		addlSubstats := charRelevantSubstats[char.Base.Key]
		if len(addlSubstats) > 0 {
			relevantSubstats = append(relevantSubstats, addlSubstats...)
		}

		substatGradients := make([]float64, len(relevantSubstats))

		// Build "gradient" by substat
		for idxSubstat, substat := range relevantSubstats {
			charProfilesCopy[idxChar].Stats[substat] += 10 * substatValues[substat] * charSubstatRarityMod[idxChar]

			simcfg.Characters.Profile = charProfilesCopy
			substatEvalResult := runSimWithConfig(cfg, simcfg, simopt)
			// sugarLog.Debugf("%v: %v (%v)", substat.String(), substatEvalResult.DPS.Mean, substatEvalResult.DPS.SD)

			substatGradients[idxSubstat] = substatEvalResult.DPS.Mean - initialMean

			// fixes cases in which fav holders don't get enough crit rate to reliably proc fav (an important example would be fav kazuha)
			// might give them "too much" cr (= max out liquid cr subs) but that's probably not a big deal
			if charWithFavonius[idxChar] && substat == core.CR {
				substatGradients[idxSubstat] += 1000
			}

			charProfilesCopy[idxChar].Stats[substat] -= 10 * substatValues[substat] * charSubstatRarityMod[idxChar]
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
		assignSubstats := func(substat core.StatType, amt int) (int, int) {
			totalSubstatCount := 0
			for _, val := range charSubstatFinal[idxChar] {
				totalSubstatCount += val
			}

			baseLiquidSubstats := int(optionsMap["total_liquid_substats"])
			for set, count := range char.Sets {
				for _, setfourstar := range artifactSets4Star {
					if set == setfourstar {
						baseLiquidSubstats -= 2 * count
					}
				}
			}

			remainingLiquidSubstats := baseLiquidSubstats - totalSubstatCount
			// Minimum of individual limit, global limit, desired amount
			amtToAdd := minInt(charSubstatLimits[idxChar][substat]-charSubstatFinal[idxChar][substat], remainingLiquidSubstats, amt)
			charSubstatFinal[idxChar][substat] += amtToAdd

			return remainingLiquidSubstats - amtToAdd, charSubstatLimits[idxChar][substat] - charSubstatFinal[idxChar][substat]
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
			case core.CR:
				gradCR := sorted.slice[idxGrad]
				gradCD := 0.0
				for i, idxSubstatTemp := range sorted.idx {
					if relevantSubstats[idxSubstatTemp] == core.CD {
						gradCD = sorted.slice[i]
					}
				}
				crCDSubstatRatio = (gradCR / gradCD)
			case core.CD:
				gradCD := sorted.slice[idxGrad]
				gradCR := 0.0
				for i, idxSubstatTemp := range sorted.idx {
					if relevantSubstats[idxSubstatTemp] == core.CR {
						gradCR = sorted.slice[i]
					}
				}
				crCDSubstatRatio = (gradCR / gradCD)
			default:
				gradStat = sorted.slice[idxGrad]
			}

			// If DPS change is really low, then it's usually better to just toss a few extra points into ER for stability
			if gradStat < 50 && crCDSubstatRatio == 0 {
				assignSubstats(core.ER, 4)
				sugarLog.Info("Low damage contribution from substats - adding some points to ER instead")
			}

			var globalLimit int
			var crLimit int
			var cdLimit int
			if crCDSubstatRatio > 0 {
				globalLimit, crLimit = assignSubstats(core.CR, 0)
				_, cdLimit = assignSubstats(core.CD, 0)

				// Continually add CR/CD to try to align CR/CD ratio to ratio of gradients until we hit a limit
				var currentRatio float64
				var amtCR int
				var amtCD int
				currentStat := core.CR
				// Debug to avoid runaway loops...
				var iteration int
				// Want this to continue until either global cap is reached, or we can neither add CR/CD
				for globalLimit > 0 && (crLimit > 0 || cdLimit > 0) && iteration < 100 {
					if charSubstatFinal[idxChar][core.CD] == 0 {
						currentRatio = float64(charSubstatFinal[idxChar][core.CR])
					} else {
						currentRatio = float64(charSubstatFinal[idxChar][core.CR]) / float64(charSubstatFinal[idxChar][core.CD])
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
						globalLimit, crLimit = assignSubstats(core.CR, amtCR)
						currentStat = core.CD
					} else if currentStat == core.CD {
						globalLimit, cdLimit = assignSubstats(core.CD, amtCD)
						currentStat = core.CR
					}
					iteration += 1
				}
			} else {
				globalLimit, _ = assignSubstats(substatToMax, 12)
			}
			if globalLimit == 0 {
				break
			}
		}

		sugarLog.Info("Final Liquid Substat Counts: ", PrettyPrintStatsCounts(charSubstatFinal[idxChar]))

		// Reset favonius char crit rate... again
		if charWithFavonius[idxChar] {
			charProfilesCopy[idxChar].Stats[core.CR] += 8 * substatValues[core.CR] * charSubstatRarityMod[idxChar]
		}
	}

	sugarLog.Info("Final config substat strings:")
	// Final output
	// This doesn't take much time relatively speaking, so just always do the processing...
	output := srcCleaned
	charNames := make(map[core.CharKey]string)
	for _, match := range reGetCharNames.FindAllStringSubmatch(output, -1) {
		charKey := core.CharNameToKey[match[1]]
		charNames[charKey] = match[1]
	}

	for idxChar, char := range charProfilesInitial {
		finalString := fmt.Sprintf("%v add stats", charNames[char.Base.Key])

		for idxSubstat, value := range substatValues {
			if value <= 0 {
				continue
			}
			finalString += fmt.Sprintf(" %v=%.6g", core.StatTypeString[idxSubstat], value*float64(2+charSubstatFinal[idxChar][idxSubstat]))
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
func runSimWithConfig(cfg string, simcfg core.SimulationConfig, simopt simulator.Options) result.Summary {
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
			sb.WriteString(core.StatTypeString[i])
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
