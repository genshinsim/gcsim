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
	AddOnAttackLanded(f func(t Target, ds *Snapshot), key string)
	OnAttackLanded(t Target, ds *Snapshot) //basically after damage
	ReactionBonus() float64

	//status
	AddStatus(key string, dur int)
	Status(key string) int //return how many more frames status will last

	//shields
	AddShield(shd Shield)
	IsShielded() bool
	GetShield(t ShieldType) Shield

	//hooks
	AddEventHook(f func(s Sim) bool, key string, hook EventHookType)

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

const (
	MaxTeamPlayerCount int = 4
)

type EventHookType int

const (
	PreSwapHook EventHookType = iota
	PostSwapHook
	PreBurstHook
	PostBurstHook
	PreSkillHook
	PostSkillHook
	PreAttackHook
	PostAttackHook
	PostShieldHook
	PostParticleHook
	//delim
	EndEventHook
)

var eventHookTypeString = [...]string{
	"PRE_SWAP",
	"POST_SWAP",
	"PRE_BURST",
	"POST_BURST",
	"PRE_SKILL",
	"POST_SKILL",
	"PRE_ATTACK",
	"POST_ATTACK",
	"POST_SHIELD",
	"POST_PARTICLE",
}

func (e EventHookType) String() string {
	return eventHookTypeString[e]
}
