package optimization

import (
	"fmt"

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
//   - Strategy is to just do a dumb grid search over ER substat values for each character
//   - ER substat values are set in increments of 2 to make the search easier
//
// 3) Given ER values, we then optimize the other substats by doing a "gradient descent" (but not really) method
func (o *SubstatOptimizer) Run(cfg string, simopt simulator.Options, simcfg *info.ActionList, gcsl ast.Node) {
	simcfg.Settings.Iterations = int(o.optionsMap["sim_iter"])
	// disable stats collection since optimizer has no use for it
	simcfg.Settings.CollectStats = []string{""}

	o.details = NewSubstatOptimizerDetails(
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

	debugLogs := o.details.optimizeERSubstats()
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
	}

	o.logger.Info("Calculating optimal DMG substat distribution...")
	debugLogs = o.details.optimizeNonERSubstats()
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
	}

	o.logger.Info("Fine tuning optimal ER vs DMG substat distribution...")
	debugLogs = o.details.optimizeERAndDMGSubstats()
	for _, debugLog := range debugLogs {
		o.logger.Info(debugLog)
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
