package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type AttackEvent struct {
	Info    AttackInfo
	Pattern AttackPattern
	// Timing        AttackTiming
	Snapshot    Snapshot
	SourceFrame int            //source frame
	Callbacks   []AttackCBFunc `json:"-"`
	Reacted     bool           // true if a reaction already took place - for purpose of attach/refill
	OnICD       bool           // set this to true if on ICD so we don't accidentally increment counter twice; icd check only happens once
}

type AttackCB struct {
	Target      Target
	AttackEvent *AttackEvent
	Damage      float64
	IsCrit      bool
}

type AttackCBFunc func(AttackCB)

type AttackInfo struct {
	ActorIndex       int       //character this attack belongs to
	DamageSrc        TargetKey //source of this attack; should be a unique key identifying the target
	Abil             string    //name of ability triggering the damage
	AttackTag        AttackTag
	ICDTag           ICDTag
	ICDGroup         ICDGroup
	Element          attributes.Element //element of ability
	Durability       Durability         //durability of aura, 0 if nothing applied
	NoImpulse        bool
	HitWeakPoint     bool
	Mult             float64 //ability multiplier. could set to 0 from initial Mona dmg
	StrikeType       StrikeType
	UseDef           bool    //we use this instead of flatdmg to make sure stat snapshotting works properly
	FlatDmg          float64 //flat dmg;
	IgnoreDefPercent float64 //by default this value is 0; if = 1 then the attack will ignore defense; raiden c2 should be set to 0.6 (i.e. ignore 60%)
	IgnoreInfusion   bool
	//amp info
	Amped   bool         //new flag used by new reaction system
	AmpMult float64      //amplier
	AmpType ReactionType //melt or vape i guess
	// catalyze info
	Catalyzed     bool
	CatalyzedType ReactionType
	//special flag for sim generated attack
	SourceIsSim bool
	DoNotLog    bool
	//hitlag stuff
	HitlagHaltFrames     float64 //this is the number of frames to pause by
	HitlagFactor         float64 //this is factor to slow clock by
	CanBeDefenseHalted   bool    //for whacking ruin gaurds
	IsDeployable         bool    //if this is true, then hitlag does not affect owner
	HitlagOnHeadshotOnly bool    //if this is true, will only apply if HitWeakpoint is also true
}

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

type Snapshot struct {
	CharLvl    int
	ActorEle   attributes.Element
	ExtraIndex int                             //this is currently purely for Kaeya icicle ICD
	Cancelled  bool                            //set to true if this snap should be ignored
	Stats      [attributes.EndStatType]float64 //total character stats including from artifact, bonuses, etc...
	BaseAtk    float64                         //base attack used in calc
	BaseDef    float64
	BaseHP     float64

	SourceFrame int           // frame snapshot was generated at
	Logs        []interface{} // logs for the snapshot
}
