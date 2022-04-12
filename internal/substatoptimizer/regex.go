package substatoptimizer

import (
	"go.uber.org/zap"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type OptimRegex struct {
	Mainstats    *regexp.Regexp
	GetCharNames *regexp.Regexp
	Substats     *regexp.Regexp
	Options      *regexp.Regexp
}

func InitRegex() *OptimRegex {
	re := OptimRegex{}

	// Regex to identify main stats based on flower. Check that characters all have one that we can recognize
	re.Mainstats = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
	re.GetCharNames = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)

	// Regex to remove stat rows that do not look like mainstat rows from the config
	re.Substats = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\b[^;]*;.*\n`)

	re.Options = regexp.MustCompile(`([a-z_]+)=([0-9.]+)`)

	return &re
}

func (re *OptimRegex) scrubSimCfg(cfg string, sugarLog *zap.SugaredLogger) string {
	if len(re.Mainstats.FindAllString(cfg, -1)) != len(re.GetCharNames.FindAllString(cfg, -1)) {
		sugarLog.Error("Error: Could not identify valid main artifact stat rows for all characters based on flower HP values.")
		sugarLog.Error("5* flowers must have 4780 HP, and 4* flowers must have 3571 HP.")
		os.Exit(1)
	}

	srcCleaned := string(cfg)
	errorPrinted := false
	for _, match := range re.Substats.FindAllString(cfg, -1) {
		if re.Mainstats.MatchString(string(match)) {
			continue
		}
		if !errorPrinted {
			sugarLog.Warn("Warning: Config found to have existing substat information. Ignoring...")
			errorPrinted = true
		}
		srcCleaned = strings.Replace(srcCleaned, string(match), "", -1)
	}

	return srcCleaned
}

func (re *OptimRegex) scrubAdditionalOptimizerCfg(additionalOptions string, optionsMap map[string]float64, sugarLog *zap.SugaredLogger) map[string]float64 {
	parsedOptions := re.Options.FindAllStringSubmatch(additionalOptions, -1)
	for _, val := range parsedOptions {
		if _, ok := optionsMap[val[1]]; ok {
			optionsMap[val[1]], _ = strconv.ParseFloat(val[2], 64)
		} else {
			sugarLog.Panic("Invalid substat optimization option found: %v", val[1], val[2])
		}
	}

	return optionsMap
}
