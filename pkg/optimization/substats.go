package optimization

import (
	"fmt"
	"math"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/shortcut"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"go.uber.org/zap"
)

type SubstatOptimizer struct {
	logger     *zap.SugaredLogger
	optionsMap map[string]float64
	details    *SubstatOptimizerDetails
	verbose    bool
}

func NewSubstatOptimizer(optionsMap map[string]float64, sugarLog *zap.SugaredLogger, verbose bool) *SubstatOptimizer {
	o := SubstatOptimizer{}
	o.optionsMap = optionsMap
	o.logger = sugarLog
	o.verbose = verbose
	return &o
}

// Substat Optimization strategy is very simplistic right now:
// This is not fully optimal - see other comments in code
// 1) User sets team, weapons, artifact sets/main stats, and rotation
// 2) Given those, for each character, sim picks ER substat value that functionally maximizes DPS Mean/SD,
// subject to a penalty on high ER values
//   - Strategy is to just do a dumb grid search over ER substat values for each character
//   - ER substat values are set in increments of 2 to make the search easier
//
// 3) Given ER values, we then optimize the other substats by doing a "gradient descent" (but not really) method
func (o *SubstatOptimizer) Run(cfg string, simopt simulator.Options, simcfg *info.ActionList, gcsl ast.Node) {
	simcfg.Settings.Iterations = int(o.optionsMap["sim_iter"])
	// disable stats collection since optimizer has no use for it
	simcfg.Settings.CollectStats = []string{""}

	o.details = NewSubstatOptimizerDetails(
		o,
		cfg,
		simopt,
		simcfg,
		gcsl,
		int(o.optionsMap["indiv_liquid_cap"]),
		int(o.optionsMap["total_liquid_substats"]),
		int(o.optionsMap["fixed_substats_count"]),
	)

	fourStarFound := o.details.setStatLimits()
	if fourStarFound {
		o.logger.Warn("Warning: 4* artifact set detected. Optimizer currently assumes that ER substats take 5* values, and all other substats take 4* values.")
	}

	o.details.setInitialSubstats(o.details.fixedSubstatCount)
	o.logger.Info("Starting ER Optimization...")

	for i := range o.details.charProfilesERBaseline {
		o.details.charProfilesCopy[i] = o.details.charProfilesERBaseline[i].Clone()
	}

	// TODO: Maybe add a configuration to only calculate ER?
	o.details.optimizeERSubstats()

	o.logger.Info("Calculating optimal DMG substat distribution...")
	debugLogs := o.details.optimizeNonERSubstats()
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
	}

	if o.optionsMap["fine_tune"] != 0 {
		o.logger.Info("Fine tuning optimal ER vs DMG substat distribution...")
		debugLogs = o.details.optimizeERAndDMGSubstats()
		for _, debugLog := range debugLogs {
			o.logger.Info(debugLog)
		}
	}
}

// Final output
// This doesn't take much time relatively speaking, so just always do the processing...
func (o *SubstatOptimizer) PrettyPrint(output string, statsFinal *SubstatOptimizerDetails) string {
	charNames := make(map[keys.Char]string)
	o.logger.Info("Final config substat strings:")

	for _, match := range regexpLineCharname.FindAllStringSubmatch(output, -1) {
		charKey := shortcut.CharNameToKey[match[1]]
		charNames[charKey] = match[1]
	}

	for idxChar := range statsFinal.charProfilesInitial {
		finalString := fmt.Sprintf("%v add stats", charNames[statsFinal.charProfilesInitial[idxChar].Base.Key])

		for idxSubstat, value := range statsFinal.substatValues {
			if value <= 0 {
				continue
			}
			finalString += fmt.Sprintf(" %v=%.6g", attributes.StatTypeString[idxSubstat], value*float64(statsFinal.fixedSubstatCount+statsFinal.charSubstatFinal[idxChar][idxSubstat]))
		}

		fmt.Println(finalString + ";")

		output = replaceSimOutputForChar(charNames[statsFinal.charProfilesInitial[idxChar].Base.Key], output, finalString)
	}

	return output
}

func NewSubstatOptimizerDetails(
	optimizer *SubstatOptimizer,
	cfg string,
	simopt simulator.Options,
	simcfg *info.ActionList,
	gcsl ast.Node,
	indivLiquidCap int,
	totalLiquidSubstats int,
	fixedSubstatCount int,
) *SubstatOptimizerDetails {
	s := SubstatOptimizerDetails{}
	s.optimizer = optimizer
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

	s.charRelevantSubstats = make([][]attributes.Stat, len(simcfg.Characters))
	for i := range simcfg.Characters {
		// ER is omitted because there is a dedicated ER step.
		s.charRelevantSubstats[i] = []attributes.Stat{
			attributes.HPP,
			attributes.HP,
			attributes.DEFP,
			attributes.DEF,
			attributes.ATKP,
			attributes.ATK,
			attributes.CR,
			attributes.CD,
			attributes.EM,
		}
	}

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
		for idxStat, stat := range stats.substatValues {
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
	sb.WriteString("Liquid Substat Counts: ")
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
