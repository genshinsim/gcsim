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
	ICDGroupYaeCharged
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
	ICDGroupYaoyaoRadishSkill
	ICDGroupYaoyaoRadishBurst
	ICDGroupBaizhuC2
	ICDGroupLength
)

var ICDGroupResetTimer []int
var ICDGroupEleApplicationSequence [][]float64
var ICDGroupDamageSequence [][]float64

func init() {
	ICDGroupResetTimer = make([]int, ICDGroupLength)
	ICDGroupResetTimer[ICDGroupDefault] = 150
	ICDGroupResetTimer[ICDGroupAmber] = 60
	ICDGroupResetTimer[ICDGroupVenti] = 60
	ICDGroupResetTimer[ICDGroupFischl] = 300
	ICDGroupResetTimer[ICDGroupDiluc] = 300
	ICDGroupResetTimer[ICDGroupPoleExtraAttack] = 30
	ICDGroupResetTimer[ICDGroupXiaoDash] = 6
	ICDGroupResetTimer[ICDGroupYaeCharged] = 30
	ICDGroupResetTimer[ICDGroupYelanBreakthrough] = 18
	ICDGroupResetTimer[ICDGroupYelanBurst] = 120
	ICDGroupResetTimer[ICDGroupColleiBurst] = 180
	ICDGroupResetTimer[ICDGroupTighnari] = 150
	ICDGroupResetTimer[ICDGroupCynoBolt] = 150
	ICDGroupResetTimer[ICDGroupDoriBurst] = 180
	ICDGroupResetTimer[ICDGroupNilou] = 114
	ICDGroupResetTimer[ICDGroupReactionA] = 30
	ICDGroupResetTimer[ICDGroupReactionB] = 30
	ICDGroupResetTimer[ICDGroupBurning] = 120
	ICDGroupResetTimer[ICDGroupNahidaSkill] = 60
	ICDGroupResetTimer[ICDGroupLayla] = 180
	ICDGroupResetTimer[ICDGroupWandererC6] = 120
	ICDGroupResetTimer[ICDGroupWandererA4] = 60
	ICDGroupResetTimer[ICDGroupAlhaithamProjectionAttack] = 720
	ICDGroupResetTimer[ICDGroupAlhaithamExtraAttack] = 120
	ICDGroupResetTimer[ICDGroupYaoyaoRadishSkill] = 150
	ICDGroupResetTimer[ICDGroupYaoyaoRadishBurst] = 90
	ICDGroupResetTimer[ICDGroupBaizhuC2] = 240

	ICDGroupEleApplicationSequence = make([][]float64, ICDGroupLength)
	ICDGroupEleApplicationSequence[ICDGroupDefault] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupAmber] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupVenti] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupFischl] = []float64{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupDiluc] = []float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupPoleExtraAttack] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupXiaoDash] = []float64{1, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupYaeCharged] = []float64{1, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupYelanBreakthrough] = []float64{1, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupYelanBurst] = []float64{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupColleiBurst] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupTighnari] = []float64{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupCynoBolt] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupDoriBurst] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupNilou] = []float64{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupReactionA] = []float64{1, 1}
	ICDGroupEleApplicationSequence[ICDGroupReactionB] = []float64{1, 1}
	ICDGroupEleApplicationSequence[ICDGroupBurning] = []float64{1, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupNahidaSkill] = []float64{1.5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupLayla] = []float64{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupWandererC6] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupWandererA4] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupAlhaithamProjectionAttack] = []float64{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0}
	ICDGroupEleApplicationSequence[ICDGroupAlhaithamExtraAttack] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupYaoyaoRadishSkill] = []float64{1, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupYaoyaoRadishBurst] = []float64{1, 0, 0, 0, 0, 0}
	ICDGroupEleApplicationSequence[ICDGroupBaizhuC2] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	ICDGroupDamageSequence = make([][]float64, ICDGroupLength)
	ICDGroupDamageSequence[ICDGroupDefault] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupAmber] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupVenti] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupFischl] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupDiluc] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupPoleExtraAttack] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupXiaoDash] = []float64{1, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupYaeCharged] = []float64{1, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupYelanBreakthrough] = []float64{1, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupYelanBurst] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupColleiBurst] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupTighnari] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupCynoBolt] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupDoriBurst] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupNilou] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupReactionA] = []float64{1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	ICDGroupDamageSequence[ICDGroupReactionB] = []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//actual data: {1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0}
	//however there seems to be no limit to the amount of burning dmg a target can take
	ICDGroupDamageSequence[ICDGroupBurning] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupNahidaSkill] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupLayla] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupWandererC6] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupWandererA4] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupAlhaithamProjectionAttack] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupAlhaithamExtraAttack] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupYaoyaoRadishSkill] = []float64{1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupYaoyaoRadishBurst] = []float64{1, 1, 1, 1, 1, 1}
	ICDGroupDamageSequence[ICDGroupBaizhuC2] = []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
}
