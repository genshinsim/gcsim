package substatoptimizer

import (
	"github.com/genshinsim/gcsim/pkg/optimization"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"

	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/genshinsim/gcsim/pkg/parse"
)

// Additional runtime option to optimize substats according to KQM standards
func RunSubstatOptim(simopt simulator.Options, verbose bool, additionalOptions string) {
	// Each optimizer run should not be saving anything out for the GZIP
	simopt.GZIPResult = false

	// Fix iterations at 350 for performance
	// TODO: Seems to be a roughly good number at KQM standards
	optionsMap := map[string]float64{
		"total_liquid_substats": 20,
		"indiv_liquid_cap":      10,
		"fixed_substats_count":  2,
		"sim_iter":              350,
		"tol_mean":              0.015,
		"tol_sd":                0.33,
		"verbose":               0,
	}

	if verbose {
		optionsMap["verbose"] = 1
	}

	// Parse and set all special sim options
	var sugarLog *zap.SugaredLogger
	if additionalOptions != "" {
		optionsMap, err := optimization.ParseOptimizerCfg(additionalOptions, optionsMap)
		sugarLog = NewLogger(optionsMap["verbose"] == 1)
		if err != nil {
			sugarLog.Panic(err.Error())
		}
	} else {
		sugarLog = NewLogger(optionsMap["verbose"] == 1)
	}

	// Parse config
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	clean, err := optimization.RemoveSubstatLines(cfg)
	if err == optimization.InvalidStats {
		sugarLog.Panic(err.Error())
		os.Exit(1)
	}

	if err != nil {
		sugarLog.Warn(err.Error())
	}

	parser := parse.New("single", clean)
	simcfg, err := parser.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	optimizer := optimization.NewSubstatOptimizer(optionsMap, sugarLog)
	optimizer.Run(cfg, simopt, simcfg)
	output := optimizer.PrettyPrint(clean, optimizer.Details)

	// Sticks optimized substat string into config and output
	if simopt.ResultSaveToPath != "" {
		output = strings.TrimSpace(output) + "\n"
		//try creating file to write to
		err = os.WriteFile(simopt.ResultSaveToPath, []byte(output), 0644)
		if err != nil {
			log.Panic(err)
		}
		sugarLog.Infof("Saved to the following location: %v", simopt.ResultSaveToPath)
	}
}
