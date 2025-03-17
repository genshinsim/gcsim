// Package shield provide a handler to keep track of shields and
// add shields etc...
package shield

import "github.com/genshinsim/gcsim/pkg/core/attributes"

type Type int

const (
	Crystallize Type = iota // lasts 15 seconds
	NoelleSkill
	NoelleA1
	ZhongliJadeShield
	DionaSkill
	BeidouThunderShield
	BeidouC1
	XinyanSkill
	XinyanC2
	KaeyaC4
	YanfeiC4
	Bell
	YunjinSkill
	ThomaSkill
	CandaceSkill
	LaylaSkill
	BaizhuBurst
	KiraraSkill
	TravelerHydroC4
	SigewinneC2
	LanyanShield
	EndType
)

type Shield interface {
	ShieldOwner() int
	ShieldTarget() int // -1 to apply to all characters, character index otherwise
	Key() int
	Type() Type
	ShieldStrength(ele attributes.Element, bonus float64) float64
	OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) // return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	CurrentHP() float64
	Element() attributes.Element
	Desc() string
}
