package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type Character struct {
	Core  *core.Core
	Index int

	//Character Profile
	Base     CharacterBase
	Weapon   WeaponProfile
	Talents  TalentProfile
	SkillCon int
	BurstCon int
	CharZone ZoneType

	//current status
	Energy    float64
	EnergyMax float64
	HPCurrent float64

	//normal attack counter
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//stats related
	Stats [attributes.EndStat]float64

	//Tags
	Tags map[string]int
}
