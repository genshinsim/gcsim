package core

type AttackEvent struct {
	Info    AttackInfo
	Pattern AttackPattern
	// Timing        AttackTiming
	Cancelled   bool //provide a way to cancel an attack event
	Snapshot    Snapshot
	SourceFrame int            //source frame
	Callbacks   []AttackCBFunc `json:"-"`
}

type AttackCB struct {
	Target      Target
	AttackEvent *AttackEvent
	Damage      float64
	IsCrit      bool
}

type AttackCBFunc func(AttackCB)

// type AttackTiming struct {
// 	SnapshotDelay int
// 	DamageDelay   int
// }

type AttackPattern struct {
	Shape    Shape
	Targets  [TargettableTypeCount]bool
	SelfHarm bool
}

type AttackInfo struct {
	ActorIndex       int    //character this attack belongs to
	DamageSrc        int    //source of this attack (i.e. index of core.Targets); always 0 for player, 1+ for the rest
	Abil             string //name of ability triggering the damage
	AttackTag        AttackTag
	ICDTag           ICDTag
	ICDGroup         ICDGroup
	Element          EleType    //element of ability
	Durability       Durability //durability of aura, 0 if nothing applied
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
	//special flag for sim generated attack
	SourceIsSim bool
	DoNotLog    bool
}

type Snapshot struct {
	CharLvl    int
	ActorEle   EleType
	ExtraIndex int                  //this is currently purely for Kaeya icicle ICD
	Cancelled  bool                 //set to true if this snap should be ignored
	Stats      [EndStatType]float64 //total character stats including from artifact, bonuses, etc...
	BaseAtk    float64              //base attack used in calc
	BaseDef    float64

	SourceFrame int           // frame snapshot was generated at
	Logs        []interface{} // logs for the snapshot
}

type Durability float64

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)
