package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
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
	ActorIndex       int               // character this attack belongs to
	DamageSrc        targets.TargetKey // source of this attack; should be a unique key identifying the target
	Abil             string            // name of ability triggering the damage
	AttackTag        attacks.AttackTag
	PoiseDMG         float64 // only needed on blunt attacks for frozen consumption before shatter for now
	ICDTag           attacks.ICDTag
	ICDGroup         attacks.ICDGroup
	Element          attributes.Element   // element of ability
	Durability       reactions.Durability // durability of aura, 0 if nothing applied
	NoImpulse        bool
	HitWeakPoint     bool
	Mult             float64 // ability multiplier. could set to 0 from initial Mona dmg
	StrikeType       attacks.StrikeType
	UseDef           bool    // we use this instead of flatdmg to make sure stat snapshotting works properly
	FlatDmg          float64 // flat dmg;
	IgnoreDefPercent float64 // by default this value is 0; if = 1 then the attack will ignore defense; raiden c2 should be set to 0.6 (i.e. ignore 60%)
	IgnoreInfusion   bool
	// amp info
	Amped   bool                   // new flag used by new reaction system
	AmpMult float64                // amplier
	AmpType reactions.ReactionType // melt or vape i guess
	// catalyze info
	Catalyzed     bool
	CatalyzedType reactions.ReactionType
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
	CharLvl    int
	ActorEle   attributes.Element
	ExtraIndex int                             // this is currently purely for Kaeya icicle ICD
	Cancelled  bool                            // set to true if this snap should be ignored
	Stats      [attributes.EndStatType]float64 // total character stats including from artifact, bonuses, etc...
	BaseAtk    float64                         // base attack used in calc
	BaseDef    float64
	BaseHP     float64

	SourceFrame int           // frame snapshot was generated at
	Logs        []interface{} // logs for the snapshot
}
