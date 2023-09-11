package optimization

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexpLineMainstat = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
	regexpLineSubstat  = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\b[^;]*;.*\n`)
	regexpLineCharname = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	regexpLineOptions  = regexp.MustCompile(`([a-z_]+)=([0-9.]+)`)
)

var errInvalidStats = errors.New("unidentifiable main stats")

func replaceSimOutputForChar(char, src, out string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^(%v\s+add\s+stats\b.*)$`, char))
	return re.ReplaceAllString(src, fmt.Sprintf("$1\n%v;", out))
}

func removeSubstatLines(cfg string) (string, error) {
	if len(regexpLineMainstat.FindAllString(cfg, -1)) != len(regexpLineCharname.FindAllString(cfg, -1)) {
		return "", errInvalidStats
	}

	clean := cfg
	errorPrinted := false
	for _, match := range regexpLineSubstat.FindAllString(cfg, -1) {
		if regexpLineMainstat.MatchString(match) {
			continue
		}

		clean = strings.Replace(clean, match, "", -1)
		errorPrinted = true
	}

	if errorPrinted {
		return clean, errors.New("warning: Config found to have existing substat information. Ignoring")
	}

	return clean, nil
}

func parseOptimizerCfg(additionalOptions string, optionsMap map[string]float64) (map[string]float64, error) {
	parsedOptions := regexpLineOptions.FindAllStringSubmatch(additionalOptions, -1)
	for _, val := range parsedOptions {
		if _, ok := optionsMap[val[1]]; ok {
			optionsMap[val[1]], _ = strconv.ParseFloat(val[2], 64)
		} else {
			err := fmt.Errorf("invalid substat optimization option found: %v: %v", val[1], val[2])
			return optionsMap, err
		}
	}

	return optionsMap, nil
}
