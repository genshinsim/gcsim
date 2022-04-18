package substatoptimizer

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

// Additional runtime option to optimize substats according to KQM standards
func RunSubstatOptim(simopt simulator.Options, verbose bool, additionalOptions string) {
	// Each optimizer run should not be saving anything out for the GZIP
	simopt.GZIPResult = false

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

	re := InitRegex()

	// Parse and set all special sim options
	var sugarLog *zap.SugaredLogger
	if additionalOptions != "" {
		optionsMap, err := re.scrubAdditionalOptimizerCfg(additionalOptions, optionsMap)
		sugarLog = InitLogger(optionsMap["verbose"] == 1)
		if err != nil {
			sugarLog.Panic(err.Error())
		}
	} else {
		sugarLog = InitLogger(optionsMap["verbose"] == 1)
	}

	// Parse config
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	srcCleaned, err := re.scrubSimCfg(cfg)
	if err == RegexParseFatalErr {
		sugarLog.Panic(err.Error())
		os.Exit(1)
	}

	if err != nil {
		sugarLog.Warn(err.Error())
	}

	parser := ast.New(srcCleaned)
	simcfg, err := parser.Parse()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	algo := InitAlgoV1(cfg, simopt, simcfg, optionsMap, sugarLog)
	algo.OptimizeSubstats()
	output := algo.PrettyPrintOptimizerOutput(srcCleaned, re, algo.stats)

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

type simRunner func(config *ast.ActionList) result.Summary

func runSimWithConfigGenerator(cfg string, simopt simulator.Options) simRunner {
	return func(simcfg *ast.ActionList) result.Summary {
		return runSimWithConfig(cfg, simcfg, simopt)
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
