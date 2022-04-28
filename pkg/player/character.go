package player

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

type MasterChar struct {
	Player *Player
	Index  int
	Character

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

	//infusion
	Infusion WeaponInfusion

	//Tags
	Tags map[string]int
}

type Character interface {
	Attack(p map[string]int) ActionInfo
	Aimed(p map[string]int) ActionInfo
	ChargeAttack(p map[string]int) ActionInfo
	HighPlungeAttack(p map[string]int) ActionInfo
	LowPlungeAttack(p map[string]int) ActionInfo
	Skill(p map[string]int) ActionInfo
	Burst(p map[string]int) ActionInfo
	Dash(p map[string]int) ActionInfo

	ActionReady(a Action, p map[string]int) bool
	ActionStam(a Action, p map[string]int) float64

	SetCD(a Action, dur int)
	Cooldown(a Action) int
	ResetActionCooldown(a Action)
	ReduceActionCooldown(a Action, v int)
	Charges(Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot
}

type WeaponInfusion struct {
	Key    string
	Ele    attributes.Element
	Tags   []combat.AttackTag
	Expiry int
}

type WeaponClass int

const (
	WeaponClassSword WeaponClass = iota
	WeaponClassClaymore
	WeaponClassSpear
	WeaponClassBow
	WeaponClassCatalyst
	EndWeaponClass
)

var weaponName = []string{
	"sword",
	"claymore",
	"polearm",
	"bow",
	"catalyst",
}

func (w WeaponClass) String() string {
	return weaponName[w]
}
