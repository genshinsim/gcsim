package attacks

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
