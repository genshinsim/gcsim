package core

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
	EndShieldType
)

type Shield interface {
	Key() int
	Type() ShieldType
	OnDamage(dmg float64, ele EleType, bonus float64) (float64, bool) //return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	CurrentHP() float64
	Element() EleType
	Desc() string
}
