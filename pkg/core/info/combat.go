package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type AttackEvent struct {
	Info    AttackInfo
	Pattern AttackPattern
	// Timing        AttackTiming
	Snapshot    Snapshot
	SourceFrame int            // source frame
	Callbacks   []AttackCBFunc `json:"-"`
	Reacted     bool           // true if a reaction already took place - for purpose of attach/refill
}

type AttackCB struct {
	Target      Target
	AttackEvent *AttackEvent
	Damage      float64
	IsCrit      bool
}

type AttackCBFunc func(AttackCB)

type AttackInfo struct {
	ActorIndex       int       // character this attack belongs to
	DamageSrc        TargetKey // source of this attack; should be a unique key identifying the target
	Abil             string    // name of ability triggering the damage
	AttackTag        attacks.AttackTag
	AdditionalTags   []attacks.AdditionalTag
	PoiseDMG         float64 // only needed on blunt attacks for frozen consumption before shatter for now
	ICDTag           attacks.ICDTag
	ICDGroup         attacks.ICDGroup
	Element          attributes.Element // element of ability
	Durability       Durability         // durability of aura, 0 if nothing applied
	NoImpulse        bool
	HitWeakPoint     bool
	Mult             float64 // ability multiplier. could set to 0 from initial Mona dmg
	StrikeType       attacks.StrikeType
	UseDef           bool    // we use this instead of flatdmg to make sure stat snapshotting works properly
	UseHP            bool    // we use this instead of flatdmg to make sure stat snapshotting works properly
	FlatDmg          float64 // flat dmg;
	IgnoreDefPercent float64 // by default this value is 0; if = 1 then the attack will ignore defense; raiden c2 should be set to 0.6 (i.e. ignore 60%)
	IgnoreInfusion   bool
	// amp info
	Amped   bool         // new flag used by new reaction system
	AmpMult float64      // amplier
	AmpType ReactionType // melt or vape i guess
	// catalyze info
	Catalyzed     bool
	CatalyzedType ReactionType
	// special flag for sim generated attack
	SourceIsSim bool
	DoNotLog    bool
	// hitlag stuff
	HitlagHaltFrames     float64 // this is the number of frames to pause by
	HitlagFactor         float64 // this is factor to slow clock by
	CanBeDefenseHalted   bool    // for whacking ruin gaurds
	IsDeployable         bool    // if this is true, then hitlag does not affect owner
	HitlagOnHeadshotOnly bool    // if this is true, will only apply if HitWeakpoint is also true
}

type Snapshot struct {
	Stats   attributes.Stats // total character stats including from artifact, bonuses, etc...
	CharLvl int

	SourceFrame int   // frame snapshot was generated at
	Logs        []any // logs for the snapshot
}

type Target interface {
	Key() TargetKey        // unique key for the target
	SetKey(k TargetKey)    // update key
	Type() TargettableType // type of target
	Shape() Shape          // info.Shape of target
	Pos() Point            // center of target
	SetPos(p Point)        // move target
	IsAlive() bool
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)
	HandleAttack(*AttackEvent) float64
	AttackWillLand(a AttackPattern) (bool, string) // hurtbox collides with AttackPattern
	IsWithinArea(a AttackPattern) bool             // center is in AttackPattern
	Tick()                                         // called every tick
	Kill()
	// for collision check
	CollidableWith(TargettableType) bool
	CollidedWith(t Target)
	WillCollide(Shape) bool
	// direction related
	Direction() Point                  // returns viewing direction as a info.Point
	SetDirection(trg Point)            // calculates viewing direction relative to default direction (0, 1)
	SetDirectionToClosestEnemy()       // looks for closest enemy
	CalcTempDirection(trg Point) Point // used for stuff like Bow CA
}

type TargetWithAura interface {
	Target
	AuraContains(e ...attributes.Element) bool
}

type AttackPattern struct {
	Shape       Shape
	SkipTargets [TargettableTypeCount]bool
	IgnoredKeys []TargetKey
}
