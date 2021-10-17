package core

type Snapshot struct {
	CharLvl    int
	ActorEle   EleType
	Actor      string //name of the character triggering the damage
	ActorIndex int
	ExtraIndex int  //this is currently purely for Kaeya icicle ICD
	Cancelled  bool //set to true if this snap should be ignored

	DamageSrc int //this is used for the purpose of calculating self harm
	SelfHarm  bool
	Targets   int //if TargetAll then not single target, resolve hitbox; other = target index

	SourceFrame     int
	AnimationFrames int //really only for amos bow...

	Abil        string      //name of ability triggering the damage
	WeaponClass WeaponClass //b.c. Gladiators...
	AttackTag   AttackTag
	ICDTag      ICDTag
	ICDGroup    ICDGroup
	ImpulseLvl  int

	CritHits     []bool
	HitWeakPoint bool
	Mult         float64 //ability multiplier. could set to 0 from initial Mona dmg
	StrikeType   StrikeType
	Element      EleType    //element of ability
	Durability   Durability //durability of aura, 0 if nothing applied

	UseDef  bool    //default false
	FlatDmg float64 //flat dmg; so far only zhongli

	Stats []float64 //total character stats including from artifact, bonuses, etc...

	BaseAtk float64 //base attack used in calc
	BaseDef float64 //base def used in calc
	//DmgBonus float64   //total damage bonus, including appropriate ele%, etc..
	DefAdj float64 //attack specific def shred (raiden c2)

	//reaction flags
	IsReactionDamage bool
	IsReaction       bool
	ReactionType     ReactionType
	IsMeltVape       bool
	ReactMult        float64 //reaction multiplier for melt/vape
	ReactBonus       float64 //reaction bonus %+ such as witch; should be 0 and only affected by hooks

	//callbacks
	OnHitCallback func(t Target)
}

type Durability float64

func (s *Snapshot) Clone() Snapshot {
	c := *s
	c.Stats = make([]float64, len(s.Stats))
	copy(c.Stats, s.Stats)
	return c
}

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

const (
	TargetPlayer int = -2
	TargetAll    int = -1
)
