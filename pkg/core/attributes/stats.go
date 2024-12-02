package attributes

import (
	"strconv"
	"strings"
)

type (
	Stat  int
	Stats [EndStatType]float64
)

// stat types
const (
	NoStat Stat = iota
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
	DendroP
	PhyP
	AtkSpd
	DmgP
	DelimBaseStat
	BaseHP
	BaseATK
	BaseDEF
	// delim
	EndStatType
)

func (s Stat) String() string {
	return StatTypeString[s]
}

func (s Stats) MaxHP() float64    { return s[BaseHP]*(1+s[HPP]) + s[HP] }
func (s Stats) TotalATK() float64 { return s[BaseATK]*(1+s[ATKP]) + s[ATK] }
func (s Stats) TotalDEF() float64 { return s[BaseDEF]*(1+s[DEFP]) + s[DEF] }

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
		if v == 0 {
			continue
		}
		sb.WriteString(StatTypeString[i])
		sb.WriteString(": ")
		sb.WriteString(strconv.FormatFloat(v, 'f', 2, 32))
		r = append(r, sb.String())
		sb.Reset()
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
	"dendro%",
	"phys%",
	"atkspd%",
	"dmg%",
	"",
	"base_hp",
	"base_atk",
	"base_def",
}

func StrToStatType(s string) Stat {
	for i, v := range StatTypeString {
		if v == s {
			return Stat(i)
		}
	}
	return -1
}
