// Package shield provide a handler to keep track of shields and
// add shields etc...
package shield

import "github.com/genshinsim/gcsim/pkg/core/attributes"

type ShieldType int

const (
	ShieldCrystallize ShieldType = iota //lasts 15 seconds
	ShieldNoelleSkill
	ShieldNoelleA1
	ShieldZhongliJadeShield
	ShieldDionaSkill
	ShieldBeidouThunderShield
	ShieldXinyanSkill
	ShieldXinyanC2
	ShieldKaeyaC4
	ShieldYanfeiC4
	ShieldBell
	ShieldYunjinSkill
	ShieldThomaSkill
	ShieldCandaceSkill
	ShieldLaylaSkill
	ShieldBaizhuBurst
	ShieldKiraraSkill
	EndShieldType
)

type Shield interface {
	ShieldOwner() int
	Key() int
	Type() ShieldType
	ShieldStrength(ele attributes.Element, bonus float64) float64
	OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) //return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	CurrentHP() float64
	Element() attributes.Element
	Desc() string
}
