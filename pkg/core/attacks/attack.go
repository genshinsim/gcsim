package attacks

type AttackTag int //attacktag is used instead of actions etc..

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
	ReactionAttackDelim
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
	AttackTagBurgeon
	AttackTagHyperbloom
	AttackTagLength
)
