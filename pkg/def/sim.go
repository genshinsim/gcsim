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
