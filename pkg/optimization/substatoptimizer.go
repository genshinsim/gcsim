package optimization

import (
	"errors"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulator"
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
		optionsMap, err := parseOptimizerCfg(additionalOptions, optionsMap)
		sugarLog = newLogger(optionsMap["verbose"] == 1)
		if err != nil {
			sugarLog.Panic(err.Error())
		}
	} else {
		sugarLog = newLogger(optionsMap["verbose"] == 1)
	}

	// Parse config
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	clean, err := removeSubstatLines(cfg)
	if errors.Is(err, errInvalidStats) {
		sugarLog.Panic("Error: Could not identify valid main artifact stat rows for all characters based on flower HP values.\n5* flowers must have 4780 HP, and 4* flowers must have 3571 HP.")
		os.Exit(1)
	}

	if err != nil {
		sugarLog.Warn(err.Error())
	}

	parser := ast.New(clean)
	simcfg, gcsl, err := parser.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	optimizer := NewSubstatOptimizer(optionsMap, sugarLog)
	optimizer.Run(cfg, simopt, simcfg, gcsl)
	output := optimizer.PrettyPrint(clean, optimizer.details)

	// Sticks optimized substat string into config and output
	if simopt.ResultSaveToPath != "" {
		output = strings.TrimSpace(output) + "\n"
		// try creating file to write to
		err = os.WriteFile(simopt.ResultSaveToPath, []byte(output), 0o644)
		if err != nil {
			log.Panic(err)
		}
		sugarLog.Infof("Saved to the following location: %v", simopt.ResultSaveToPath)
	}
}
