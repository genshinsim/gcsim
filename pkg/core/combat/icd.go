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
	ICDTagDoriC2
	ICDTagDoriChargingStation
	ICDTagNilouTranquilityAura
	ICDTagWandererC6
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
	ICDTagNahidaSkill
	ICDTagNahidaC6
	ICDTagWandererA4
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
	ICDGroupPoleExtraAttack
	ICDGroupXiaoDash
	ICDGroupYelanBreakthrough
	ICDGroupYelanBurst
	ICDGroupColleiBurst
	ICDGroupTighnari
	ICDGroupCynoBolt
	ICDGroupDoriBurst
	ICDGroupNilou
	ICDGroupReactionA
	ICDGroupReactionB
	ICDGroupBurning
	ICDGroupNahidaSkill
	ICDGroupLayla
	ICDGroupWandererC6
	ICDGroupWandererA4
	ICDGroupAlhaithamProjectionAttack
	ICDGroupAlhaithamExtraAttack //CA
	ICDGroupLength
)

var ICDGroupResetTimer = []int{
	150, //default
	60,  //amber
	60,  //venti
	300, //fischl
	300, //diluc
	30,  //pole extra
	6,   //xiao dash
	18,  //yelan pew pew
	120, //yelan burst
	180, //collei burst
	150, //tighnari
	150, //cyno skill bolts
	180, //dori burst
	114, //nilou
	30,  //reaction a
	30,  //reaciton b
	120, //burning
	60,  //nahida skill
	180, //layla
	120, //wanderer c6
	60,  //wanderer a4
	720, //alhaitham projection
	120, //alhaitham CA
}

var ICDGroupEleApplicationSequence = [][]float64{
	//default tag
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	//amber tag
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	//venti tag
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	//fischl
	{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
	//diluc
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	//pole extra
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	//xiao dash
	{1, 0, 0, 0, 0, 0, 0},
	//yelan pew pew
	{1, 0, 0, 0},
	//yelan burst
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	//collei burst
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//tighnari
	{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0},
	//cyno skill bolt
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//dori burst
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//nilou
	{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
	//reaction a
	{1, 1},
	//reaction b
	{1, 1},
	//burning
	{1, 0, 0, 0, 0, 0, 0, 0},
	//nahida skill
	{1.5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//layla
	{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	//wanderer c6
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//wanderer a4
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//alhaitham projection
	{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
	//alhaitham CA
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

var ICDGroupDamageSequence = [][]float64{
	//default
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//amber
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//venti
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//fischl
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//diluc
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//pole extra
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//xiao
	{1, 0, 0, 0, 0, 0, 0},
	//yelan pew pew
	{1, 0, 0, 0},
	//yelan burst
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//collei burst
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//tighnari
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//cyno skill bolt
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//dori burst
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//nilou
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//ele A
	{1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	//ele B
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//burning
	//actual data: {1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}
	//however there seems to be no limit to the amount of burning dmg a target can take
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//nahida-skill
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//layla
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//wanderer c6
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//wanderer a4
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//alhaitham-projection
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	//alhaitham-CA
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}
