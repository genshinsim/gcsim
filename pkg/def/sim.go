package def

import "math/rand"

type Sim interface {
	//sim controls
	SwapCD() int
	Stam() float64
	Frame() int //current frame
	Flags() Flags
	SetCustomFlag(key string, val int)
	GetCustomFlag(key string) (int, bool)

	//character related
	CharByPos(ind int) (Character, bool)
	CharByName(name string) (Character, bool)
	DistributeParticle(p Particle)
	ActiveCharIndex() int
	ActiveDuration() int
	Characters() []Character

	//damage related
	ApplyDamage(ds *Snapshot)

	//target related
	TargetHasDebuff(debuff string, param int) bool
	TargetHasElement(ele EleType, param int) bool
	Targets() []Target
	AddOnAttackWillLand(f func(t Target, ds *Snapshot), key string)
	OnAttackWillLand(t Target, ds *Snapshot)
	AddOnAttackLanded(f func(t Target, ds *Snapshot, dmg float64, crit bool), key string)
	OnAttackLanded(t Target, ds *Snapshot, dmg float64, crit bool) //basically after damage
	//these are on reaction damage about to happen
	AddOnAmpReaction(f func(t Target, ds *Snapshot), key string)
	OnAmpReaction(t Target, ds *Snapshot)
	AddOnTransReaction(f func(t Target, ds *Snapshot), key string)
	OnTransReaction(t Target, ds *Snapshot)
	// ReactionBonus(ds *Snapshot) float64
	// AddReactionBonus(f func(ds *Snapshot) float64, key string)
	OnReaction(t Target, ds *Snapshot)
	AddOnReaction(f func(t Target, ds *Snapshot), key string)
	OnTargetDefeated(t Target)
	AddOnTargetDefeated(f func(t Target), key string)

	//initial
	AddInitHook(f func())

	//status
	AddStatus(key string, dur int)
	Status(key string) int //return how many more frames status will last
	DeleteStatus(key string)

	//shields
	AddShield(shd Shield)
	IsShielded() bool
	GetShield(t ShieldType) Shield
	AddShieldBonus(f func() float64)

	//healing
	HealActive(amt float64)
	HealAll(amt float64)
	HealAllPercent(percent float64)
	AddIncHealBonus(f func() float64)

	//constructs
	NewConstruct(c Construct, refresh bool)
	NewNoLimitCons(c Construct, refresh bool)
	ConstructCount() int
	ConstructCountType(t GeoConstructType) int
	Destroy(key int) bool
	HasConstruct(key int) bool

	AddOnHurt(f func(s Sim))

	//hooks
	AddEventHook(f func(s Sim) bool, key string, hook EventHookType)

	//other
	Rand() *rand.Rand
}

type Flags struct {
	HPMode              bool
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
	ShieldBell
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
	PostDashHook
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
	"POST_DASH",
}

func (e EventHookType) String() string {
	return eventHookTypeString[e]
}
