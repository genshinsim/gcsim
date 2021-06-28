package def

import "math/rand"

type Sim interface {
	//sim controls
	SwapCD() int
	Stam() float64
	Frame() int //current frame
	Flags() Flags

	//character related
	CharByPos(ind int) (Character, bool)
	CharByName(name string) (Character, bool)
	DistributeParticle(p Particle)

	//damage related
	ApplyDamage(ds *Snapshot)

	//target related
	TargetHasDebuff(debuff string, param int) bool
	TargetHasElement(ele EleType, param int) bool
	OnAttackLanded(t Target) //basically after damage
	ReactionBonus() float64

	//status
	Status(key string) int //return how many more frames status will last

	//shields
	AddShield(shd Shield)
	IsShielded() bool
	GetShield(t ShieldType) Shield

	//other
	Rand() *rand.Rand
}

type Flags struct {
	ChildeActive        bool
	AmpReactionDidOccur bool
	AmpReactionType     ReactionType
	NextAttackMVMult    float64 // melt vape multiplier
	// ReactionDamageTriggered bool
	Custom map[string]int
}

type ShieldType int

const (
	ShieldCrystallize ShieldType = iota //lasts 15 seconds
	ShieldNoelleSkill
	ShieldNoelleA2
	ShieldZhongliJadeShield
	ShieldDionaSkill
	ShieldBeidouThunderShield
	ShieldXinyanSkill
	ShieldXinyanC2
	ShieldKaeyaC4
	ShieldYanfeiC4
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

const (
	MaxTeamPlayerCount int = 4
)
