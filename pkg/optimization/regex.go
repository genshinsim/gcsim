package optimization

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	RegexpLineMainstat = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\s+hp=(4780|3571)\b[^;]*;`)
	RegexpLineSubstat  = regexp.MustCompile(`(?m)^[a-z]+\s+add\s+stats\b[^;]*;.*\n`)
	RegexpLineCharname = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	RegexpLineOptions  = regexp.MustCompile(`([a-z_]+)=([0-9.]+)`)
)

type OptimRegex struct {
	Mainstats      *regexp.Regexp
	GetCharNames   *regexp.Regexp
	Substats       *regexp.Regexp
	Options        *regexp.Regexp
	InsertLocation *regexp.Regexp
}

var ErrInvalidStats = errors.New("unidentifiable main stats")

func ReplaceSimOutputForChar(char, src, out string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^(%v\s+add\s+stats\b.*)$`, char))
	return re.ReplaceAllString(src, fmt.Sprintf("$1\n%v;", out))
}

func RemoveSubstatLines(cfg string) (string, error) {
	if len(RegexpLineMainstat.FindAllString(cfg, -1)) != len(RegexpLineCharname.FindAllString(cfg, -1)) {
		return "", ErrInvalidStats
	}

	clean := cfg
	errorPrinted := false
	for _, match := range RegexpLineSubstat.FindAllString(cfg, -1) {
		if RegexpLineMainstat.MatchString(match) {
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
	parsedOptions := RegexpLineOptions.FindAllStringSubmatch(additionalOptions, -1)
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
