package substatoptimizer

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/shortcut"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"go.uber.org/zap"
)

type AlgoV1 struct {
	sugarLog   *zap.SugaredLogger
	optionsMap map[string]float64
	cfg        string
	simopt     simulator.Options
	simcfg     *ast.ActionList
	stats      *OptimStats
}

func InitAlgoV1(cfg string, simopt simulator.Options, simcfg *ast.ActionList, optionsMap map[string]float64, sugarLog *zap.SugaredLogger) *AlgoV1 {
	a := AlgoV1{}

	a.cfg = cfg
	a.simopt = simopt
	a.simcfg = simcfg
	a.optionsMap = optionsMap
	a.sugarLog = sugarLog

	return &a
}

// Substat Optimization strategy is very simplistic right now:
// This is not fully optimal - see other comments in code
// 1) User sets team, weapons, artifact sets/main stats, and rotation
// 2) Given those, for each character, sim picks ER substat value that functionally maximizes DPS Mean/SD,
// subject to a penalty on high ER values
//    - Strategy is to just do a dumb grid search over ER substat values for each character
//    - ER substat values are set in increments of 2 to make the search easier
// 3) Given ER values, we then optimize the other substats by doing a "gradient descent" (but not really) method
func (a *AlgoV1) OptimizeSubstats() {
	// Fix iterations at 350 for performance
	// TODO: Seems to be a roughly good number at KQM standards
	a.simcfg.Settings.Iterations = int(a.optionsMap["sim_iter"])

	a.stats = InitOptimStats(a.simcfg, int(a.optionsMap["indiv_liquid_cap"]), int(a.optionsMap["total_liquid_substats"]), int(a.optionsMap["fixed_substats_count"]))

	fourStarFound := a.stats.setStatLimits()
	if fourStarFound {
		a.sugarLog.Warn("Warning: 4* artifact set detected. Optimizer currently assumes that ER substats take 5* values, and all other substats take 4* values.")
	}

	a.stats.setInitialSubstats(a.stats.fixedSubstatCount)
	a.sugarLog.Info("Starting ER Optimization...")

	a.stats.setCharProfilesCopy(a.stats.charProfilesERBaseline)

	// Tolerance cutoffs for mean and SD from initial state
	// Initial state is used rather than checking across each iteration due to noise
	// TODO: May want to adjust further?
	tolMean := a.optionsMap["tol_mean"]
	tolSD := a.optionsMap["tol_sd"]
	runner := runSimWithConfigGenerator(a.cfg, a.simopt)

	debugLogs := a.stats.optimizeERSubstats(tolMean, tolSD, runner)
	for _, debugLog := range debugLogs {
		a.sugarLog.Info(debugLog)
	}

	debugLogs = a.stats.optimizeNonERSubstats(runner)
	for _, debugLog := range debugLogs {
		a.sugarLog.Info(debugLog)
	}

	a.sugarLog.Info("Final config substat strings:")
}

// Final output
// This doesn't take much time relatively speaking, so just always do the processing...
func (a *AlgoV1) PrettyPrintOptimizerOutput(output string, re *OptimRegex, statsFinal *OptimStats) string {
	charNames := make(map[keys.Char]string)

	for _, match := range re.GetCharNames.FindAllStringSubmatch(output, -1) {
		charKey := shortcut.CharNameToKey[match[1]]
		charNames[charKey] = match[1]
	}

	for idxChar, char := range statsFinal.charProfilesInitial {
		finalString := fmt.Sprintf("%v add stats", charNames[char.Base.Key])

		for idxSubstat, value := range statsFinal.substatValues {
			if value <= 0 {
				continue
			}
			finalString += fmt.Sprintf(" %v=%.6g", attributes.StatTypeString[idxSubstat], value*float64(statsFinal.fixedSubstatCount+statsFinal.charSubstatFinal[idxChar][idxSubstat]))
		}

		fmt.Println(finalString + ";")

		output = ReplaceSimOutputForChar(charNames[char.Base.Key], output, finalString)
	}

	return output
}
