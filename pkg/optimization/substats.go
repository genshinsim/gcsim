package optimization

import (
	"fmt"
	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/genshinsim/gcsim/pkg/core"
	"go.uber.org/zap"
)

type SubstatOptimizer struct {
	logger     *zap.SugaredLogger
	optionsMap map[string]float64
	Details    *SubstatOptimizerDetails
}

func NewSubstatOptimizer(optionsMap map[string]float64, sugarLog *zap.SugaredLogger) *SubstatOptimizer {
	o := SubstatOptimizer{}
	o.optionsMap = optionsMap
	o.logger = sugarLog

	return &o
}

// Substat Optimization strategy is very simplistic right now:
// This is not fully optimal - see other comments in code
// 1) User sets team, weapons, artifact sets/main stats, and rotation
// 2) Given those, for each character, sim picks ER substat value that functionally maximizes DPS Mean/SD,
// subject to a penalty on high ER values
//    - Strategy is to just do a dumb grid search over ER substat values for each character
//    - ER substat values are set in increments of 2 to make the search easier
// 3) Given ER values, we then optimize the other substats by doing a "gradient descent" (but not really) method
func (o *SubstatOptimizer) Run(cfg string, simopt simulator.Options, simcfg core.SimulationConfig) {

	simcfg.Settings.Iterations = int(o.optionsMap["sim_iter"])
	o.Details = NewSubstatOptimizerDetails(cfg, simopt, simcfg, int(o.optionsMap["indiv_liquid_cap"]), int(o.optionsMap["total_liquid_substats"]))

	fourStarFound := o.Details.setStatLimits()
	if fourStarFound {
		o.logger.Warn("Warning: 4* artifact set detected. Optimizer currently assumes that ER substats take 5* values, and all other substats take 4* values.")
	}

	fixedSubstatCount := o.optionsMap["fixed_substats_count"]
	o.Details.setInitialSubstats(fixedSubstatCount)
	o.logger.Info("Starting ER Optimization...")

	o.Details.setCharProfilesCopy(o.Details.charProfilesERBaseline)

	// Tolerance cutoffs for mean and SD from initial state
	// Initial state is used rather than checking across each iteration due to noise
	// TODO: May want to adjust further?
	tolMean := o.optionsMap["tol_mean"]
	tolSD := o.optionsMap["tol_sd"]

	debugLogs := o.Details.optimizeERSubstats(tolMean, tolSD)
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
	}

	debugLogs = o.Details.optimizeNonERSubstats()
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
	}
}

// Final output
// This doesn't take much time relatively speaking, so just always do the processing...
func (o *SubstatOptimizer) PrettyPrint(output string, statsFinal *SubstatOptimizerDetails) string {
	charNames := make(map[core.CharKey]string)
	o.logger.Info("Final config substat strings:")

	for _, match := range REGEXP_LINE_CHARNAME.FindAllStringSubmatch(output, -1) {
		charKey := core.CharNameToKey[match[1]]
		charNames[charKey] = match[1]
	}

	for idxChar, char := range statsFinal.charProfilesInitial {
		finalString := fmt.Sprintf("%v add stats", charNames[char.Base.Key])

		for idxSubstat, value := range statsFinal.substatValues {
			if value <= 0 {
				continue
			}
			finalString += fmt.Sprintf(" %v=%.6g", core.StatTypeString[idxSubstat], value*float64(2+statsFinal.charSubstatFinal[idxChar][idxSubstat]))
		}

		fmt.Println(finalString + ";")

		output = ReplaceSimOutputForChar(charNames[char.Base.Key], output, finalString)
	}

	return output
}
