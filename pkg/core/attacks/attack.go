package attacks

type AttackTag int // attacktag is used instead of actions etc..

const (
	AttackTagNone AttackTag = iota
	AttackTagNormal
	AttackTagExtra
	AttackTagPlunge
	AttackTagElementalArt
	AttackTagElementalArtHold
	AttackTagElementalBurst
	AttackTagWeaponSkill
	AttackTagMonaBubbleBreak
	AttackTagNoneStat

	ReactionAttackStartDelim
	AttackTagOverloadDamage
	AttackTagSuperconductDamage
	AttackTagECDamage
	AttackTagShatter
	AttackTagSwirlPyro
	AttackTagSwirlHydro
	AttackTagSwirlCryo
	AttackTagSwirlElectro
	AttackTagBurningDamage
	AttackTagBloom
	AttackTagBountifulCore // special tag for nilou
	AttackTagBurgeon
	AttackTagHyperbloom
	ReactionAttackEndDelim

	LunarReactionStartDelim
	AttackTagReactionLunarCharge
	LunarReactionEndDelim

	DirectLunarReactionStartDelim
	AttackTagDirectLunarCharged
	DirectLunarReactionEndDelim

	AttackTagLength
)

type StrikeType int

const (
	StrikeTypeDefault StrikeType = iota
	StrikeTypePierce
	StrikeTypeBlunt
	StrikeTypeSlash
	StrikeTypeSpear
)

type AdditionalTag int

const (
	AdditionalTagNone AdditionalTag = iota
	AdditionalTagNightsoul
	AdditionalTagKinichCannon
)
