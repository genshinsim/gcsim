package substatoptimizer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type OptimRegex struct {
	Mainstats      *regexp.Regexp
	GetCharNames   *regexp.Regexp
	Substats       *regexp.Regexp
	Options        *regexp.Regexp
	InsertLocation *regexp.Regexp
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

func ReplaceSimOutputForChar(charName, src, finalOutput string) string {
	reInsertLocation := regexp.MustCompile(fmt.Sprintf(`(?m)^(%v\s+add\s+stats\b.*)$`, charName))
	return reInsertLocation.ReplaceAllString(src, fmt.Sprintf("$1\n%v;", finalOutput))
}

func (re *OptimRegex) scrubSimCfg(cfg string) (string, error) {
	if len(re.Mainstats.FindAllString(cfg, -1)) != len(re.GetCharNames.FindAllString(cfg, -1)) {
		msg := "Error: Could not identify valid main artifact stat rows for all characters based on flower HP values.\n"
		msg += "5* flowers must have 4780 HP, and 4* flowers must have 3571 HP."

		return "", NewFatalErr(msg)
	}

	srcCleaned := cfg
	errorPrinted := false
	for _, match := range re.Substats.FindAllString(cfg, -1) {
		if re.Mainstats.MatchString(match) {
			continue
		}
		if !errorPrinted {
			errorPrinted = true
		}
		srcCleaned = strings.Replace(srcCleaned, match, "", -1)
	}

	if errorPrinted {
		return srcCleaned, errors.New("Warning: Config found to have existing substat information. Ignoring...")
	}

	return srcCleaned, nil
}

func (re *OptimRegex) scrubAdditionalOptimizerCfg(additionalOptions string, optionsMap map[string]float64) (map[string]float64, error) {
	parsedOptions := re.Options.FindAllStringSubmatch(additionalOptions, -1)
	for _, val := range parsedOptions {
		if _, ok := optionsMap[val[1]]; ok {
			optionsMap[val[1]], _ = strconv.ParseFloat(val[2], 64)
		} else {
			err := errors.New(fmt.Sprintf("Invalid substat optimization option found: %v: %v", val[1], val[2]))
			return optionsMap, err
		}
	}

	return optionsMap, nil
}
