package def

type StatType int

//stat types
const (
	DEFP StatType = iota
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
	EleP
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

var StatTypeString = [...]string{
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
