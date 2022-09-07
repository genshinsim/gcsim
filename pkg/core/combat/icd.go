package combat

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

type ICDTag int //same ICD tag shares the same counter

const (
	ICDTagNone ICDTag = iota
	ICDTagNormalAttack
	ICDTagExtraAttack
	ICDTagElementalArt
	ICDTagElementalArtAnemo
	ICDTagElementalArtPyro
	ICDTagElementalArtHydro
	ICDTagElementalArtCryo
	ICDTagElementalArtElectro
	ICDTagElementalBurst
	ICDTagElementalBurstAnemo
	ICDTagElementalBurstPyro
	ICDTagElementalBurstHydro
	ICDTagElementalBurstCryo
	ICDTagElementalBurstElectro
	ICDTagDash
	ICDTagLisaElectro
	ICDTagYanfeiFire
	ICDTagYelanBreakthrough
	ICDTagYelanBurst
	ICDTagCynoBolt
	ICDTagMonaWaterDamage
	ICDTagTravelerWakeOfEarth
	ICDTagKleeFireDamage
	ICDTagTartagliaRiptideFlash
	ICDTagColleiSprout
	ICDReactionDamageDelim
	ICDTagOverloadDamage
	ICDTagSuperconductDamage
	ICDTagECDamage
	ICDTagShatter
	ICDTagSwirlPyro
	ICDTagSwirlHydro
	ICDTagSwirlCryo
	ICDTagSwirlElectro
	ICDTagBurningDamage
	ICDTagBloomDamage
	ICDTagBurgeonDamage
	ICDTagHyperbloomDamage
	ICDTagLength
)

// group dictate both the sequence and the reset timer
type ICDGroup int //same ICD group shares the same timer

const (
	ICDGroupDefault ICDGroup = iota
	ICDGroupAmber
	ICDGroupVenti
	ICDGroupFischl
	ICDGroupDiluc
	ICDGroupPole
	ICDGroupXiaoDash
	ICDGroupYelanBreakthrough
	ICDGroupYelanBurst
	ICDGroupColleiBurst
	ICDGroupTighnari
	ICDGroupCynoBolt
	ICDGroupReactionA
	ICDGroupReactionB
	ICDGroupBurning
	ICDGroupLength
)

var ICDGroupResetTimer = []int{
	150, //default
	60,  //amber
	60,  //venti
	300, //fischl
	300, //diluc
	30,  //pole
	6,   //xiao dash
	18,  //yelan pew pew
	120, //yelan burst
	180, //collei burst
	150, //tighnari
	150, //cyno skill bolts TODO:verify reset timer
	30,  //reaction a
	30,  //reaciton b
	120, //burning
}

var ICDGroupEleApplicationSequence = [][]int{
	//default tag
	{1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	//amber tag
	{1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	//venti tag
	{1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	//fischl
	{1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
	//diluc
	{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0},
	//pole
	{1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	//xiao dash
	{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//yelan pew pew
	{1.0, 0.0, 0.0, 0.0},
	//yelan burst
	{1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	//collei burst
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//tighnari
	{1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0},
	//cyno skill bolt
	{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//reaction a
	{1.0, 1.0},
	//reaction b
	{1.0, 1.0},
	//burning
	{1.0, 0, 0, 0, 0, 0, 0, 0},
}

var ICDGroupDamageSequence = [][]float64{
	//default
	{1, 1, 1, 1, 1},
	//amber
	{1, 1, 1, 1, 1},
	//venti
	{1, 1, 1, 1, 1},
	//fischl
	{1, 1, 1, 1, 1},
	//diluc
	{1, 1, 1, 1, 1},
	//pole
	{1, 1, 1, 1, 1},
	//xiao
	{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//yelan pew pew
	{1.0, 0.0, 0.0, 0.0},
	//yelan burst
	{1, 1, 1, 1, 1},
	//collei burst
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//tighnari
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//cyno skill bolt
	{1, 1, 1, 1, 1},
	//ele A
	{1.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//ele B
	{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//burning
	//actual data: {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//however there seems to be no limit to the amount of burning dmg a target can take
	{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
}
