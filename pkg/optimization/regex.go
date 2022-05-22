package optimization

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	REGEXP_LINE_MAINSTAT = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
	REGEXP_LINE_SUBSTAT  = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\b[^;]*;.*\n`)
	REGEXP_LINE_CHARNAME = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	REGEXP_LINE_OPTIONS  = regexp.MustCompile(`([a-z_]+)=([0-9.]+)`)
)

type OptimRegex struct {
	Mainstats      *regexp.Regexp
	GetCharNames   *regexp.Regexp
	Substats       *regexp.Regexp
	Options        *regexp.Regexp
	InsertLocation *regexp.Regexp
}

var InvalidStats = errors.New("Error: Could not identify valid main artifact stat rows for all characters based on flower HP values.\n5* flowers must have 4780 HP, and 4* flowers must have 3571 HP.")

func ReplaceSimOutputForChar(char, src, out string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^(%v\s+add\s+stats\b.*)$`, char))
	return re.ReplaceAllString(src, fmt.Sprintf("$1\n%v;", out))
}

func RemoveSubstatLines(cfg string) (string, error) {
	if len(REGEXP_LINE_MAINSTAT.FindAllString(cfg, -1)) != len(REGEXP_LINE_CHARNAME.FindAllString(cfg, -1)) {
		return "", InvalidStats
	}

	clean := cfg
	errorPrinted := false
	for _, match := range REGEXP_LINE_SUBSTAT.FindAllString(cfg, -1) {
		if REGEXP_LINE_MAINSTAT.MatchString(match) {
			continue
		}

		clean = strings.Replace(clean, match, "", -1)
		errorPrinted = true
	}

	if errorPrinted {
		return clean, errors.New("Warning: Config found to have existing substat information. Ignoring...")
	}

	return clean, nil
}

func ParseOptimizerCfg(additionalOptions string, optionsMap map[string]float64) (map[string]float64, error) {
	parsedOptions := REGEXP_LINE_OPTIONS.FindAllStringSubmatch(additionalOptions, -1)
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
