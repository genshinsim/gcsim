package core

import (
	"strconv"
	"strings"
)

type StatType int

//stat types
const (
	NoStat StatType = iota
	DEFP
	DEF
	HP
	HPP
	ATK
	ATKP
	ER
	EM
	CR
	CD
	Heal
	PyroP
	HydroP
	CryoP
	ElectroP
	AnemoP
	GeoP
	PhyP
	DendroP
	AtkSpd
	DmgP
	//delim
	EndStatType
)

func (s StatType) String() string {
	return StatTypeString[s]
}

func PrettyPrintStats(stats []float64) string {
	var sb strings.Builder
	for i, v := range stats {
		if v > 0 {
			sb.WriteString(StatTypeString[i])
			sb.WriteString(": ")
			sb.WriteString(strconv.FormatFloat(v, 'f', 2, 32))
			sb.WriteString(" ")
		}
	}
	return strings.Trim(sb.String(), " ")
}

func PrettyPrintStatsSlice(stats []float64) []string {
	r := make([]string, 0)
	var sb strings.Builder
	for i, v := range stats {
		if v > 0 {
			sb.WriteString(StatTypeString[i])
			sb.WriteString(": ")
			sb.WriteString(strconv.FormatFloat(v, 'f', 2, 32))
			r = append(r, sb.String())
			sb.Reset()
		}
	}

	return r
}

var StatTypeString = [...]string{
	"n/a",
	"def%",
	"def",
	"hp",
	"hp%",
	"atk",
	"atk%",
	"er",
	"em",
	"cr",
	"cd",
	"heal",
	"pyro%",
	"hydro%",
	"cryo%",
	"electro%",
	"anemo%",
	"geo%",
	"ele%",
	"phys%",
	"dendro%",
	"atkspd%",
	"dmg%",
}

func StrToStatType(s string) StatType {
	for i, v := range StatTypeString {
		if v == s {
			return StatType(i)
		}
	}
	return -1
}
