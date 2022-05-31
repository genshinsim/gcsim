package core

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
	AttackTagLength
)

type ICDTag int //same ICD tag shares the same counter

const (
	ICDTagNone ICDTag = iota
	ICDTagNormalAttack
	ICDTagExtraAttack
	ICDTagElementalArt
	ICDTagElementalBurst
	ICDTagDash
	ICDTagLisaElectro
	ICDTagYanfeiFire
	ICDTagVentiBurstAnemo
	ICDTagVentiBurstPyro
	ICDTagVentiBurstHydro
	ICDTagVentiBurstCryo
	ICDTagVentiBurstElectro
	ICDTagYelanBurst
	ICDTagMonaWaterDamage
	ICDTagTravelerWakeOfEarth
	ICDTagKleeFireDamage
	ICDTagTartagliaRiptideFlash
	ICDReactionDamageDelim
	ICDTagOverloadDamage
	ICDTagSuperconductDamage
	ICDTagECDamage
	ICDTagShatter
	ICDTagSwirlPyro
	ICDTagSwirlHydro
	ICDTagSwirlCryo
	ICDTagSwirlElectro
	ICDTagLength
)

//group dictate both the sequence and the reset timer
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
	ICDGroupReactionA
	ICDGroupReactionB
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
	30,  //reaction a
	30,  //reaciton b
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
	//reaction a
	{1.0, 1.0},
	//reaction b
	{1.0, 1.0},
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
	//ele A
	{1.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	//ele B
	{1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
}
